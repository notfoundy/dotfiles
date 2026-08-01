// Package onepass wraps the 1Password CLI.
//
// It replaces the bash implementation's `eval "$(op signin)"` followed by a
// fixed sequence of cursor-movement escapes to erase the prompts. Here the
// session token is parsed out of op's own output and propagated to
// ansible-playbook through the child environment, so nothing has to be erased.
package onepass

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/shell"
)

// DefaultDomain matches the OP_ACCOUNT_DOMAIN default of the bash script.
const DefaultDomain = "my.1password.eu"

// secretFields are the field references roles/ssh/tasks/ssh_keys.yml reads for
// each key. `op read` only resolves a full op://vault/item/field reference, so
// verifying these is both correct and a dry run of what the playbook will do.
// Keep in sync with that role.
var secretFields = []string{"private_key?ssh-format=openssh", "public_key"}

// Account is one entry of `op account list --format=json`.
type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
	Shorthand   string `json:"shorthand"`
}

// Label is what to show a human choosing between accounts.
func (a Account) Label() string {
	if a.Email != "" {
		return fmt.Sprintf("%s (%s)", a.Email, a.URL)
	}
	return a.URL
}

// Ref is the identifier to hand to `op --account`.
func (a Account) Ref() string {
	if a.Shorthand != "" {
		return a.Shorthand
	}
	if a.URL != "" {
		return a.URL
	}
	return a.AccountUUID
}

// Client talks to the op binary.
type Client struct {
	log *config.Logger
}

// New returns a client bound to a run log.
func New(log *config.Logger) *Client {
	return &Client{log: log}
}

// Installed reports whether the op binary is on PATH.
func Installed() bool { return shell.Exists("op") }

// Authenticated reports whether an existing session is already usable, which
// is the common case on a machine that has been set up before.
func (c *Client) Authenticated() bool {
	return shell.Succeeds(c.log, shell.Cmd{Name: "op", Args: []string{"vault", "list"}})
}

// Accounts lists the accounts configured on this machine.
func (c *Client) Accounts() ([]Account, error) {
	out, err := shell.Capture(c.log, shell.Cmd{Name: "op", Args: []string{"account", "list", "--format=json"}})
	if err != nil {
		return nil, err
	}

	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}

	var accounts []Account
	if err := json.Unmarshal([]byte(out), &accounts); err != nil {
		return nil, fmt.Errorf("parsing `op account list` output: %w", err)
	}
	return accounts, nil
}

// AddAccount registers a new account. It needs the terminal: the user has to
// type their email, secret key and master password.
func (c *Client) AddAccount(domain string) error {
	if domain == "" {
		domain = DefaultDomain
	}
	return shell.Interactive(c.log, shell.Cmd{
		Name: "op",
		Args: []string{"account", "add", "--address", domain},
	})
}

// SignIn authenticates and returns the resulting session as a KEY=VALUE pair
// ready to be appended to a child environment. It returns an empty string when
// op handled the session itself, which is what happens with desktop app
// integration.
func (c *Client) SignIn(account Account) (string, error) {
	args := []string{"signin"}
	if ref := account.Ref(); ref != "" {
		args = append(args, "--account", ref)
	}

	out, err := shell.PromptCapture(c.log, shell.Cmd{Name: "op", Args: args})
	if err != nil {
		return "", err
	}
	return parseSession(out), nil
}

// parseSession extracts OP_SESSION_x=token from the `export OP_SESSION_x="..."`
// line op prints on stdout.
func parseSession(out string) string {
	for line := range strings.SplitSeq(out, "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "export "))

		key, value, found := strings.Cut(line, "=")
		if !found || !strings.HasPrefix(key, "OP_SESSION_") {
			continue
		}

		value = strings.Trim(value, `"'`)
		value = strings.TrimSuffix(value, ";")
		if value == "" {
			continue
		}
		return key + "=" + strings.Trim(value, `"'`)
	}
	return ""
}

// references expands a secret into the op:// references the ssh role reads.
func references(s config.Secret) []string {
	if s.VaultPath == "" {
		return nil
	}

	refs := make([]string, 0, len(secretFields))
	for _, field := range secretFields {
		refs = append(refs, s.VaultPath+"/"+field)
	}
	return refs
}

// VerifySecrets checks that every vault path declared in group_vars/all.yml is
// readable, so a missing or renamed item fails here rather than midway through
// the playbook.
func (c *Client) VerifySecrets(secrets []config.Secret, env []string) error {
	var missing []string

	for _, s := range secrets {
		for _, ref := range references(s) {
			cmd := shell.Cmd{Name: "op", Args: []string{"read", ref}, Env: env}
			_, err := shell.Capture(c.log, cmd)
			if err == nil {
				continue
			}

			reason := shell.Output(err)
			if reason == "" {
				reason = err.Error()
			}
			missing = append(missing, fmt.Sprintf("%s (%s): %s", s.Name, ref, reason))
			// The remaining fields of a key we cannot read add nothing.
			break
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("unreadable 1Password items:\n  - %s", strings.Join(missing, "\n  - "))
	}
	return nil
}
