package main

import (
	"fmt"
	"os"

	"github.com/notfoundy/dotfiles/cli/internal/cmd"
)

// Set by GoReleaser through -ldflags.
var version = "dev"

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
