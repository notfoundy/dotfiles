package onepass

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/shell"
)

const DefaultDomain = "my.1password.eu"

// The fields roles/ssh/tasks/ssh_keys.yml reads; keep in sync. The private key
// is left out on purpose: the 1Password agent serves it.
var secretFields = []string{"public_key"}

type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
	Shorthand   string `json:"shorthand"`
}

func (a Account) Label() string {
	if a.Email != "" {
		return fmt.Sprintf("%s (%s)", a.Email, a.URL)
	}
	return a.URL
}

func (a Account) Ref() string {
	if a.Shorthand != "" {
		return a.Shorthand
	}
	if a.URL != "" {
		return a.URL
	}
	return a.AccountUUID
}

type Client struct {
	log *config.Logger
}

func New(log *config.Logger) *Client {
	return &Client{log: log}
}

func Installed() bool { return shell.Exists("op") }

func (c *Client) Authenticated() bool {
	return shell.Succeeds(c.log, shell.Cmd{Name: "op", Args: []string{"vault", "list"}})
}

func (c *Client) Accounts() ([]Account, error) {
	out, err := shell.Capture(c.log, shell.Cmd{Name: "op", Args: []string{"account", "list", "--format=json"}})
	if err != nil {
		if reason := shell.Output(err); reason != "" {
			return nil, fmt.Errorf("listing 1Password accounts: %s", reason)
		}
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

func (c *Client) AddAccount(domain string) error {
	if domain == "" {
		domain = DefaultDomain
	}
	return shell.Interactive(c.log, shell.Cmd{
		Name: "op",
		Args: []string{"account", "add", "--address", domain},
	})
}

// SignIn returns "" with the desktop app integration, which keeps the session itself.
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

// VerifySecrets fails before the playbook would, midway, on a missing item.
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
			break
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("unreadable 1Password items:\n  - %s", strings.Join(missing, "\n  - "))
	}
	return nil
}
