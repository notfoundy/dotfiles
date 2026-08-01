// Package run orchestrates a full dotfiles run: prepare the host, resolve what
// to apply, then hand the terminal over to Ansible.
package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/notfoundy/dotfiles/cli/internal/ansible"
	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/onepass"
	"github.com/notfoundy/dotfiles/cli/internal/prompt"
	"github.com/notfoundy/dotfiles/cli/internal/system"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

// Options are the resolved command-line inputs.
type Options struct {
	Repo       string   // explicit checkout; empty means the managed ~/.dotfiles
	Tags       []string // roles requested with -t
	Select     bool     // open the role picker
	SkipUpdate bool     // do not touch git
	Extra      []string // passed through to ansible-playbook
	Version    string
}

// Execute runs the whole flow.
func Execute(opts Options) error {
	log := config.Open()
	defer log.Close()

	log.Printf("=== dotfiles %s started ===", opts.Version)
	defer log.Printf("=== dotfiles run finished ===")

	p := ui.New()

	root, err := repoPath(opts)
	if err != nil {
		return err
	}

	// On a machine that has been set up before, refresh the checkout and
	// resolve the selection up front: a typo in -t then fails instantly,
	// before any package installation or sudo prompt. A fresh machine has no
	// checkout yet, so there it can only happen after the clone.
	var (
		cfg       *config.Config
		selection config.Selection
	)
	ready := isCheckout(root)
	if ready {
		if cfg, selection, err = prepareRun(log, p, opts, root); err != nil {
			return err
		}
	}

	if err := prepareHost(log, p); err != nil {
		return err
	}

	if !ready {
		if cfg, selection, err = prepareRun(log, p, opts, root); err != nil {
			return err
		}
	}
	log.Printf("selected roles: %v (explicit=%t)", selection.Roles, selection.Explicit)

	if err := ansible.InstallGalaxyDeps(log, p, root); err != nil {
		return err
	}

	env, err := setupOnePassword(log, p, root, cfg)
	if err != nil {
		return err
	}

	p.Notice("Applying %d role(s): %s", len(selection.Roles), join(selection.Roles))

	if err := ansible.RunPlaybook(log, ansible.PlaybookOptions{
		Root:      root,
		Selection: selection,
		Extra:     opts.Extra,
		Env:       env,
	}); err != nil {
		return err
	}

	return announceFirstRun(p)
}

func prepareHost(log *config.Logger, p *ui.Printer) error {
	distro, err := system.DetectDistro()
	if err != nil {
		return err
	}

	task := p.Start("Preparing %s", distro.Name)
	if !distro.Supported() {
		task.Fail()
		return fmt.Errorf("%s is not supported — only Fedora is", distro.Name)
	}
	task.Done()

	if err := system.EnsureBasePackages(log, p); err != nil {
		return err
	}
	return system.EnsureSSHKeys(log, p)
}

// repoPath returns the checkout to work from.
//
// With --repo the given directory is used as-is and git is left alone, which is
// what running against a development clone needs. The bash script computed its
// own location and then ignored it, always operating on ~/.dotfiles.
func repoPath(opts Options) (string, error) {
	if opts.Repo != "" {
		return filepath.Abs(opts.Repo)
	}
	return system.DefaultRepoPath()
}

func isCheckout(root string) bool {
	_, err := os.Stat(filepath.Join(root, "main.yml"))
	return err == nil
}

// prepareRun brings the checkout up to date, reads its configuration and
// resolves what to apply.
func prepareRun(log *config.Logger, p *ui.Printer, opts Options, root string) (*config.Config, config.Selection, error) {
	if opts.Repo != "" {
		p.Start("Using checkout at %s", root).Done()
	} else if err := system.EnsureRepo(log, p, root, opts.SkipUpdate); err != nil {
		return nil, config.Selection{}, err
	}

	if !isCheckout(root) {
		return nil, config.Selection{}, fmt.Errorf("%s does not look like a dotfiles checkout: no main.yml", root)
	}

	cfg, err := config.Load(root)
	if err != nil {
		return nil, config.Selection{}, err
	}

	selection, err := resolveSelection(cfg, opts)
	if err != nil {
		return nil, config.Selection{}, err
	}
	return cfg, selection, nil
}

func resolveSelection(cfg *config.Config, opts Options) (config.Selection, error) {
	if !opts.Select {
		return cfg.Resolve(opts.Tags)
	}

	preselected := cfg.DefaultRoles
	if len(opts.Tags) > 0 {
		// -t narrows what the picker starts from; validate it first.
		requested, err := cfg.Resolve(opts.Tags)
		if err != nil {
			return config.Selection{}, err
		}
		preselected = requested.Roles
	}

	chosen, err := prompt.SelectRoles(cfg.Available, preselected)
	if err != nil {
		return config.Selection{}, err
	}
	return cfg.Resolve(chosen)
}

// setupOnePassword makes sure op is installed and a session is available,
// returning the environment additions to pass on to Ansible.
func setupOnePassword(log *config.Logger, p *ui.Printer, root string, cfg *config.Config) ([]string, error) {
	client := onepass.New(log)

	if !onepass.Installed() {
		p.Start("Checking 1Password CLI").Skip("not installed")
		p.Warn("Installing it through Ansible — expect a separate sudo prompt.")
		if err := ansible.RunTag(log, root, "1password"); err != nil {
			return nil, err
		}
		if !onepass.Installed() {
			return nil, fmt.Errorf("1Password CLI still not available after installing the 1password role")
		}
	}

	if client.Authenticated() {
		p.Start("1Password session").Done()
		return nil, verifySecrets(p, client, cfg, nil)
	}

	account, err := chooseAccount(p, client)
	if err != nil {
		return nil, err
	}

	// Announced before signing in: op takes the terminal over to ask for the
	// master password, and the step it belongs to has to be on screen by then.
	task := p.StartInteractive("1Password session")
	session, err := client.SignIn(account)
	if err != nil {
		task.Fail()
		return nil, err
	}
	task.Done()

	var env []string
	if session != "" {
		env = append(env, session)
	}

	return env, verifySecrets(p, client, cfg, env)
}

func chooseAccount(p *ui.Printer, client *onepass.Client) (onepass.Account, error) {
	accounts, err := client.Accounts()
	if err != nil {
		return onepass.Account{}, err
	}

	switch len(accounts) {
	case 1:
		return accounts[0], nil
	case 0:
		domain, err := prompt.AskDomain(onepass.DefaultDomain)
		if err != nil {
			return onepass.Account{}, err
		}
		p.Info("Adding a 1Password account on %s", domain)
		if err := client.AddAccount(domain); err != nil {
			return onepass.Account{}, err
		}

		added, err := client.Accounts()
		if err != nil {
			return onepass.Account{}, err
		}
		if len(added) == 0 {
			return onepass.Account{}, fmt.Errorf("no 1Password account was added")
		}
		return added[0], nil
	default:
		return prompt.SelectAccount(accounts)
	}
}

func verifySecrets(p *ui.Printer, client *onepass.Client, cfg *config.Config, env []string) error {
	if len(cfg.Secrets) == 0 {
		return nil
	}

	task := p.Start("Checking %d 1Password item(s)", len(cfg.Secrets))
	if err := client.VerifySecrets(cfg.Secrets, env); err != nil {
		task.Fail()
		return err
	}
	task.Done()
	return nil
}

// announceFirstRun prints the reboot notice once, tracked by ~/.dotfiles_run.
func announceFirstRun(p *ui.Printer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	marker := filepath.Join(home, ".dotfiles_run")
	if _, err := os.Stat(marker); err == nil {
		return nil
	}

	p.Success("First run complete!")
	p.Notice("Please reboot your computer to complete the setup.")

	f, err := os.Create(marker)
	if err != nil {
		return err
	}
	return f.Close()
}

func join(roles []string) string {
	const max = 6
	if len(roles) <= max {
		return fmt.Sprint(sliceString(roles))
	}
	return fmt.Sprintf("%s and %d more", sliceString(roles[:max]), len(roles)-max)
}

func sliceString(roles []string) string {
	out := ""
	for i, r := range roles {
		if i > 0 {
			out += ", "
		}
		out += r
	}
	return out
}
