// Package cmd wires the command-line interface.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/notfoundy/dotfiles/cli/internal/run"
)

// Execute parses the arguments and runs the CLI.
func Execute(version string) error {
	root := newRootCmd(version)
	root.SilenceUsage = true
	root.SilenceErrors = true
	return root.Execute()
}

func newRootCmd(version string) *cobra.Command {
	opts := run.Options{Version: version}

	cmd := &cobra.Command{
		Use:   "dotfiles [flags] [-- ansible-playbook args...]",
		Short: "Set up and update this machine from the dotfiles repo",
		Long: `Prepares the host, refreshes the dotfiles checkout, then applies it with Ansible.

Without arguments every default role is applied, as configured by default_roles
in group_vars/all.yml. Arguments after -- are forwarded to ansible-playbook.`,
		Example: `  dotfiles
  dotfiles -t nvim,go
  dotfiles --select
  dotfiles --repo . -- --check`,
		Version: version,
		Args:    cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			opts.Extra = args
			return run.Execute(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&opts.Tags, "tags", "t", nil, "roles to apply (comma separated); defaults to every default role")
	flags.BoolVarP(&opts.Select, "select", "s", false, "pick the roles interactively")
	flags.StringVar(&opts.Repo, "repo", "", "run against this checkout instead of ~/.dotfiles, leaving git untouched")
	flags.BoolVar(&opts.SkipUpdate, "skip-update", false, "do not pull the dotfiles repository")

	// Everything after the first positional argument belongs to
	// ansible-playbook, so cobra must stop interpreting flags there.
	flags.SetInterspersed(false)

	_ = cmd.RegisterFlagCompletionFunc("tags", completeTags)

	cmd.AddCommand(newSelfUpdateCmd(version))

	return cmd
}
