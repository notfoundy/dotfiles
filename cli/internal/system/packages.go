package system

import (
	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/shell"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

// basePackages are what the rest of the run needs before Ansible takes over.
// dnf-plugins-core is required by the seven COPR repositories the roles enable.
var basePackages = []string{"ansible", "git", "dnf-plugins-core"}

// PackageInstalled reports whether dnf knows the package as installed.
func PackageInstalled(log *config.Logger, name string) bool {
	return shell.Succeeds(log, shell.Cmd{Name: "dnf", Args: []string{"list", "--installed", name}})
}

// EnsureBasePackages installs whatever is missing from basePackages.
//
// A full `dnf update` runs first, but only when Ansible itself is absent — that
// is the marker for a machine that has never been set up.
func EnsureBasePackages(log *config.Logger, p *ui.Printer) error {
	if !PackageInstalled(log, "ansible") {
		task := p.Start("Updating system packages")
		if _, err := shell.Capture(log, shell.Cmd{Name: "sudo", Args: []string{"dnf", "update", "-y"}}); err != nil {
			task.Fail()
			p.Detail(shell.Output(err))
			return err
		}
		task.Done()
	}

	for _, pkg := range basePackages {
		if PackageInstalled(log, pkg) {
			continue
		}

		task := p.Start("Installing %s", pkg)
		if _, err := shell.Capture(log, shell.Cmd{Name: "sudo", Args: []string{"dnf", "install", "-y", pkg}}); err != nil {
			task.Fail()
			p.Detail(shell.Output(err))
			return err
		}
		task.Done()
	}
	return nil
}
