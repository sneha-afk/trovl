package cmd

import (
	"os"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a specified symlink while keeping the target file as-is.",
	Long: `To remove a symlink, specify it:
  trovl remove <symlink>
For example, to remove a symbolic link at home to a .vimrc:
  trovl remove ~/.vimrc
The target file will remain untouched, and this command will NOT remove any file that
is not a symbolic link.`,
	Run: func(cmd *cobra.Command, args []string) {
		symlink := args[0]
		if err := links.RemoveByPath(symlink); err != nil {
			State.Logger.Error("Could not remove symlink", "error", err)
			os.Exit(1)
		}

		if State.Verbose() {
			State.Logger.Info("Successfully removed symlink", "link", symlink)
		}
	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"unlink", "delete"},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
