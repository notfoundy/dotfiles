package cmd

import (
	"github.com/spf13/cobra"

	"github.com/notfoundy/dotfiles/cli/internal/selfupdate"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

func newSelfUpdateCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "self-update",
		Short: "Update the dotfiles binary to the latest release",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			p := ui.New()

			task := p.Start("Checking for a newer release")
			latest, err := selfupdate.LatestVersion()
			if err != nil {
				task.Fail()
				return err
			}

			if !selfupdate.IsNewer(version, latest) {
				task.Skip("already on " + version)
				return nil
			}
			task.Done()

			task = p.Start("Installing %s", latest)
			if err := selfupdate.Apply(); err != nil {
				task.Fail()
				return err
			}
			task.Done()

			p.Success("Updated to %s", latest)
			return nil
		},
	}
}
