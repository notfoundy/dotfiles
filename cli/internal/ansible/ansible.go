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
// its native output and its BECOME password prompt behave normally.
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
	args = append(args, opts.Extra...)
	if !hasBecomeFlag(opts.Extra) {
		args = append(args, "-K")
	}
	return args
}

// RunTag runs a single tag, used to bootstrap the 1Password role before the
// main run.
func RunTag(log *config.Logger, root, tag string) error {
	return shell.Interactive(log, shell.Cmd{
		Name: "ansible-playbook",
		Args: []string{filepath.Join(root, "main.yml"), "-t", tag, "-K"},
		Dir:  root,
	})
}

// hasBecomeFlag reports whether the caller already dealt with privilege
// escalation, so we do not append a second -K on top of theirs.
func hasBecomeFlag(extra []string) bool {
	for _, arg := range extra {
		switch {
		case arg == "-K", arg == "--ask-become-pass":
			return true
		case strings.HasPrefix(arg, "--become-password-file"):
			return true
		}
	}
	return false
}
