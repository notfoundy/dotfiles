package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/shell"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

// RepoURL is where a fresh machine gets the dotfiles from.
const RepoURL = "https://github.com/notfoundy/dotfiles.git"

// DefaultRepoPath is the managed checkout, ~/.dotfiles.
func DefaultRepoPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dotfiles"), nil
}

// EnsureRepo makes sure path holds an up-to-date checkout, cloning it if
// needed. Updating is skipped when skipUpdate is set.
func EnsureRepo(log *config.Logger, p *ui.Printer, path string, skipUpdate bool) error {
	if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		task := p.Start("Cloning dotfiles repository")
		cmd := shell.Cmd{Name: "git", Args: []string{"clone", "--quiet", RepoURL, path}}
		if _, err := shell.Capture(log, cmd); err != nil {
			task.Fail()
			p.Detail(shell.Output(err))
			return err
		}
		task.Done()
		return nil
	}

	task := p.Start("Updating dotfiles repository")
	if skipUpdate {
		task.Skip("--skip-update")
		return nil
	}
	if err := update(log, path); err != nil {
		task.Fail()
		p.Detail(shell.Output(err))
		return err
	}
	task.Done()
	return nil
}

// update fetches and fast-forwards.
//
// The bash version ran a bare `git pull`, which fails opaquely when the local
// checkout has diverged or carries uncommitted work. Splitting fetch from a
// strict fast-forward lets us say exactly what is wrong.
func update(log *config.Logger, path string) error {
	git := func(args ...string) (string, error) {
		return shell.Capture(log, shell.Cmd{Name: "git", Args: args, Dir: path})
	}

	if _, err := git("fetch", "--quiet", "--prune"); err != nil {
		return err
	}

	dirty, err := git("status", "--porcelain")
	if err != nil {
		return err
	}
	if strings.TrimSpace(dirty) != "" {
		return fmt.Errorf(
			"%s has uncommitted changes; commit, stash or discard them, or pass --skip-update", path)
	}

	if _, err := git("merge", "--ff-only", "--quiet", "@{u}"); err != nil {
		return fmt.Errorf(
			"%s has diverged from its upstream and cannot fast-forward; reconcile it manually, or pass --skip-update", path)
	}
	return nil
}
