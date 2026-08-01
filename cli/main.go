// Command dotfiles bootstraps and runs the Ansible-managed dotfiles setup.
package main

import (
	"fmt"
	"os"

	"github.com/notfoundy/dotfiles/cli/internal/cmd"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
