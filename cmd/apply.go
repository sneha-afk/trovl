package cmd

import (
	"os"
	"path/filepath"

	"github.com/sneha-afk/trovl/internal/manifests"
	"github.com/sneha-afk/trovl/internal/utils"
	"github.com/spf13/cobra"
)

var defaultFile = "manifest.json"

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply <manifest_file> [more_manifests]",
	Short: "Applies a manifest specified by schema (default: $XDG_CONFIG_HOME/trovl/manifest.json)",
	Long: `Applies a manifest specified by schema to bulk add or fix links as needed.

By default, trovl looks for a manifest in $XDG_CONFIG_HOME/trovl/manifest.json (typically ~/.config/trovl/manifest.json). If $XDG_CONFIG_HOME
is not set, trovl then checks ~/.config/trovl/manifest.json (on all OSes). If any manifest is specified into the command, the default
manifest file is not applied (i.e, this process happens only upon trovl apply)

When backing up a file that would be overwritten by this new symlink, trovl always uses $XDG_CACHE_HOME first, before
falling back to OS defaults. See [trovl's use of environment variables](../configuration/#environment-variables) to learn more.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Find one of the default filepaths to apply
		if len(args) <= 0 {
			configDir, err := utils.GetConfigDir()
			if err != nil {
				State.Logger.Error("Could not read config directory", "error", err)
			}
			path := filepath.Join(configDir, defaultFile)

			defaultManifest, err := manifests.New(path)
			if err != nil {
				if err == os.ErrNotExist {
					State.Logger.Error("Did not find manifest at the default location", "defaultLocation", path)
				} else {
					State.Logger.Error("Error reading the default manifest", "error", err)
				}
				cmd.Help()
				os.Exit(1)
			}

			if err := defaultManifest.Apply(State); err != nil {
				State.Logger.Error("Could not apply default manifest file", "error", err)
				os.Exit(1)
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
	Aliases: []string{"exec", "run", "do"},
	Example: "trovl apply .trovl",
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().BoolVar(&cfg.OverwriteYes, "overwrite", false, "overwrite any existing symlinks")
	applyCmd.Flags().BoolVar(&cfg.OverwriteNo, "no-overwrite", false, "do not overwrite any existing symlinks")
	applyCmd.Flags().BoolVar(&cfg.BackupYes, "backup", false, "backup existing single files if a symlink would overwrite it")
	applyCmd.Flags().BoolVar(&cfg.BackupYes, "no-backup", false, "do not backup existing files and abandon symlink creation")

	applyCmd.MarkFlagsMutuallyExclusive("overwrite", "no-overwrite")
	applyCmd.MarkFlagsMutuallyExclusive("backup", "no-backup")
}
