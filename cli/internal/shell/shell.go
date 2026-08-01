// Package shell runs external commands, in one of two modes.
//
// Capture buffers output into the run log and keeps the terminal clean, which
// is what the preparation steps want. Interactive hands the real stdin, stdout
// and stderr over to the child, which is what anything prompting the user needs
// — ansible-playbook's BECOME password, or `op account add --signin`.
//
// Note that sudo prompts work in both modes: sudo reads the password from
// /dev/tty directly rather than from the inherited descriptors.
package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/config"
)

// Cmd describes an external command to run.
type Cmd struct {
	Name string
	Args []string
	Dir  string
	Env  []string // extra variables, appended to the inherited environment
}

func (c Cmd) String() string {
	return strings.TrimSpace(c.Name + " " + strings.Join(c.Args, " "))
}

// Error carries the captured output of a failed command.
type Error struct {
	Cmd    string
	Output string
	Err    error
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %v", e.Cmd, e.Err) }
func (e *Error) Unwrap() error { return e.Err }

func (c Cmd) build() *exec.Cmd {
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Dir = c.Dir
	if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}
	return cmd
}

// Capture runs the command with its output buffered, returning it as a string.
func Capture(log *config.Logger, c Cmd) (string, error) {
	log.Printf("exec: %s", c)

	out, err := c.build().CombinedOutput()
	output := string(out)
	if output != "" {
		log.Printf("output: %s", strings.TrimRight(output, "\n"))
	}
	if err != nil {
		log.Printf("failed: %s: %v", c, err)
		return output, &Error{Cmd: c.String(), Output: output, Err: err}
	}
	return output, nil
}

// Interactive runs the command wired to the real terminal.
func Interactive(log *config.Logger, c Cmd) error {
	log.Printf("exec (interactive): %s", c)

	cmd := c.build()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("failed: %s: %v", c, err)
		return &Error{Cmd: c.String(), Err: err}
	}
	return nil
}

// PromptCapture runs the command with stdin and stderr wired to the terminal
// but stdout buffered and returned.
//
// This is what reading a value out of an interactive tool requires: `op signin`
// writes its password prompt to stderr and reads the answer from the terminal,
// while the session token it emits on stdout must be captured rather than
// printed. It replaces the shell idiom `eval "$(op signin)"`.
func PromptCapture(log *config.Logger, c Cmd) (string, error) {
	log.Printf("exec (prompt): %s", c)

	var stdout strings.Builder
	cmd := c.build()
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("failed: %s: %v", c, err)
		return stdout.String(), &Error{Cmd: c.String(), Err: err}
	}
	return stdout.String(), nil
}

// Succeeds reports whether the command exits zero, discarding its output.
// Used for probes such as `op vault list` or `dnf list --installed`.
func Succeeds(log *config.Logger, c Cmd) bool {
	_, err := Capture(log, c)
	return err == nil
}

// Exists reports whether a binary is on PATH.
func Exists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Output returns the captured output of a failed command, if any.
func Output(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return strings.TrimRight(e.Output, "\n")
	}
	return ""
}
