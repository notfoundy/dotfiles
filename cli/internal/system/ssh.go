package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/shell"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

// EnsureSSHKeys generates an ed25519 key pair and self-authorizes it, unless
// ~/.ssh/authorized_keys already exists.
func EnsureSSHKeys(log *config.Logger, p *ui.Printer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshDir := filepath.Join(home, ".ssh")
	authorized := filepath.Join(sshDir, "authorized_keys")
	if _, err := os.Stat(authorized); err == nil {
		return nil
	}

	task := p.Start("Generating SSH keys")

	if err := os.MkdirAll(sshDir, 0o700); err != nil {
		task.Fail()
		return err
	}
	if err := os.Chmod(sshDir, 0o700); err != nil {
		task.Fail()
		return err
	}

	key := filepath.Join(sshDir, "id_ed25519")
	if _, err := os.Stat(key); err != nil {
		comment := fmt.Sprintf("%s@%s", username(), hostname())
		cmd := shell.Cmd{Name: "ssh-keygen", Args: []string{"-t", "ed25519", "-f", key, "-N", "", "-C", comment}}
		if _, err := shell.Capture(log, cmd); err != nil {
			task.Fail()
			p.Detail(shell.Output(err))
			return err
		}
	}

	pub, err := os.ReadFile(key + ".pub")
	if err != nil {
		task.Fail()
		return err
	}
	if err := appendLine(authorized, pub); err != nil {
		task.Fail()
		return err
	}

	task.Done()
	return nil
}

func appendLine(path string, content []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.Write(content)
	return err
}

func username() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return "user"
}

func hostname() string {
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	return "localhost"
}
