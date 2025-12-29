package cmd

import (
	"os"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a symlink that points to the target file",
	Long: `To add a symlink, specify via:
  trovl add <target> <symlink>
For example, to create a link from a dotfiles/.vimrc to .vimrc at home:
  trovl add ~/dotfiles/.vimrc ~/.vimrc`,
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		symlink := args[1]

		link, err := links.Construct(State, target, symlink)
		if err != nil {
			State.Logger.Error("Could not construct symlink", "error", err)
			os.Exit(1)
		}

		if err := links.Add(&link); err != nil {
			State.Logger.Error("Failed to create link (maybe try running as admin?)", "error", err)
			os.Exit(1)
		}

		State.Logger.Info("Successfully added symlink", "link", symlink, "target", target)
	},
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"link", "create"},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVar(&cfg.UseRelative, "relative", false, "retain relative paths to target")
	addCmd.Flags().BoolVar(&cfg.OverwriteYes, "overwrite", false, "overwrite any existing symlinks")
	addCmd.Flags().BoolVar(&cfg.OverwriteNo, "no-overwrite", false, "do not overwrite any existing symlinks")

	addCmd.MarkFlagsMutuallyExclusive("overwrite", "no-overwrite")
}
