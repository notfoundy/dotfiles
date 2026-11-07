package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/config"
)

type Cmd struct {
	Name string
	Args []string
	Dir  string
	Env  []string // appended to the inherited environment
}

func (c Cmd) String() string {
	return strings.TrimSpace(c.Name + " " + strings.Join(c.Args, " "))
}

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

// PromptCapture is for `op signin`: it prompts on stderr and prints the session on stdout.
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

func Succeeds(log *config.Logger, c Cmd) bool {
	_, err := Capture(log, c)
	return err == nil
}

func Exists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func Output(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return strings.TrimRight(e.Output, "\n")
	}
	return ""
}
