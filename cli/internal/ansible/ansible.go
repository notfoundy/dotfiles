// Package ansible drives ansible-galaxy and ansible-playbook.
package ansible

import (
	"path/filepath"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/shell"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

// InstallGalaxyDeps installs the collections listed in requirements/common.yml.
func InstallGalaxyDeps(log *config.Logger, p *ui.Printer, root string) error {
	task := p.Start("Updating Ansible Galaxy dependencies")

	cmd := shell.Cmd{
		Name: "ansible-galaxy",
		Args: []string{"install", "-r", filepath.Join(root, "requirements", "common.yml")},
		Dir:  root,
	}
	if _, err := shell.Capture(log, cmd); err != nil {
		task.Fail()
		p.Detail(shell.Output(err))
		return err
	}

	task.Done()
	return nil
}

// PlaybookOptions describes one ansible-playbook invocation.
type PlaybookOptions struct {
	Root      string
	Selection config.Selection
	// Extra is passed through verbatim, e.g. --check or -vvv.
	Extra []string
	// Env holds additional variables for the child, such as the 1Password
	// session token.
	Env []string
}

// RunPlaybook executes main.yml with the terminal handed over to Ansible, so
// its native output behaves normally.
func RunPlaybook(log *config.Logger, opts PlaybookOptions) error {
	return shell.Interactive(log, shell.Cmd{
		Name: "ansible-playbook",
		Args: playbookArgs(opts),
		Dir:  opts.Root,
		Env:  opts.Env,
	})
}

// playbookArgs assembles the ansible-playbook command line.
func playbookArgs(opts PlaybookOptions) []string {
	args := []string{filepath.Join(opts.Root, "main.yml")}

	if opts.Selection.Explicit {
		args = append(args, "-t", strings.Join(opts.Selection.Roles, ","))
	}
	return append(args, opts.Extra...)
}
