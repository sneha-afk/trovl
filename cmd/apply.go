package cmd

import (
	"os"

	"github.com/sneha-afk/trovl/internal/manifests"
	"github.com/sneha-afk/trovl/internal/utils"
	"github.com/spf13/cobra"
)

// Paths to resolve when there are no arguments passed in
var defaultFiles = []string{
	".trovl",
	"trovl.json",
	".trovl.json",
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply <manifest_file> [more_manifests]",
	Short: "Applies a link list specified by schema.",
	Long:  `Applies a link list specified by schema to bulk add links or fix as needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Find one of the default filepaths to apply
		if len(args) <= 0 {
			for _, path := range defaultFiles {
				info, err := utils.GetPathInfo(path)
				if !info.Exists || err != nil {
					continue
				}

				args = append(args, path)
				break
			}
		}

		for _, path := range args {
			m, err := manifests.New(path)
			if err != nil {
				State.Logger.Error("Could not read manifest file", "error", err)
				os.Exit(1)
			}

			if err := m.Apply(State); err != nil {
				State.Logger.Error("Could not apply manifest file", "error", err)
				os.Exit(1)
			}
		}

	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"execute", "run", "do"},
	Example: "trovl apply .trovl",
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().BoolVar(&cfg.OverwriteYes, "overwrite", false, "overwrite any existing symlinks")
	applyCmd.Flags().BoolVar(&cfg.OverwriteNo, "no-overwrite", false, "do not overwrite any existing symlinks")

	applyCmd.MarkFlagsMutuallyExclusive("overwrite", "no-overwrite")
}
