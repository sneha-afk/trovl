package cmd

import (
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan <manifest_file> [more_manifests]",
	Short: "Describes what will happen during an `apply` without modifying the filesystem",
	Long: `Describes the actions that will happen when a manifest file is applied. This is essentially
an alias for running
  trovl apply --dry-run`,
	Run: func(cmd *cobra.Command, args []string) {
		State.Options.DryRun = true
		State.SetLogLevel()
		applyCmd.Run(cmd, args)
	},
	Example: "trovl plan .trovl",
}

func init() {
	rootCmd.AddCommand(planCmd)
}
