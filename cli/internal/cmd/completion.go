package cmd

import (
	"github.com/spf13/cobra"

	"github.com/notfoundy/dotfiles/cli/internal/config"
	"github.com/notfoundy/dotfiles/cli/internal/system"
)

// completeTags offers the roles actually present in the checkout, which is the
// same list -t is validated against.
func completeTags(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	root := cmd.Flag("repo").Value.String()
	if root == "" {
		var err error
		if root, err = system.DefaultRepoPath(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}

	cfg, err := config.Load(root)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return cfg.Available, cobra.ShellCompDirectiveNoFileComp
}
