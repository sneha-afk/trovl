package cmd

import (
	"os"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <symlink> [more_symlinks]",
	Short: "Removes a specified symlink while keeping the target file as-is.",
	Long: `Removes symlinks while keeping the target file untouched. Validates any argument passed
	in as truly being a symlink to prevent data loss.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, symlink := range args {
			if err := links.RemoveByPath(State, symlink); err != nil {
				State.Logger.Error("Could not remove symlink", "error", err)
				os.Exit(1)
			}

			if State.Verbose() {
				State.Logger.Info("Successfully removed symlink", "link", symlink)
			}
		}
	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"unlink", "delete"},
	Example: "trovl remove ~/.vimrc (where it is a symlink)",
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
