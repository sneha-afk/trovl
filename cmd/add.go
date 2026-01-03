package cmd

import (
	"os"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <target> <symlink> [target2, symlink2, ...]",
	Short: "Adds a symlink that points to the target file",
	Long:  `When possible, add a true symlink (as in, not a junction or hard link) to a target file.`,
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < len(args); i += 2 {
			target := args[i]
			symlink := args[i+1]

			if err := links.Add(State, target, symlink); err != nil {
				State.Logger.Error("Failed to create link (maybe try running as admin?)", "error", err)
				os.Exit(1)
			}

			if !State.Options.DryRun {
				State.Logger.Info("Successfully added symlink", "link", symlink, "target", target)
			}
		}
	},
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"link", "create", "new"},
	Example: "trovl add ~/dotfiles/.vimrc ~/.vimrc",
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVar(&cfg.UseRelative, "relative", false, "retain relative paths to target")
	addCmd.Flags().BoolVar(&cfg.OverwriteYes, "overwrite", false, "overwrite any existing symlinks")
	addCmd.Flags().BoolVar(&cfg.OverwriteNo, "no-overwrite", false, "do not overwrite any existing symlinks")

	addCmd.MarkFlagsMutuallyExclusive("overwrite", "no-overwrite")
}
