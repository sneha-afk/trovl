/*
Package cmd is the central CobraCLI aspects of trovl.
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/sneha-afk/trovl/internal/state"
	"github.com/spf13/cobra"
)

var (
	cfg   = &state.TrovlOptions{}
	State *state.TrovlState
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "trovl",
	Short: "A cross-platform symlink manager.",
	Long: `trovl is a cross-platform symlink manager that aims to make file management easier and more efficient.
It features configurable paths for files and directories that vary in location depending on the system,
and true-symlinking when possible.
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		State = state.New(cfg)
		slog.SetDefault(State.Logger)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	State = state.DefaultState()
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "have verbose outputs for actions taken")
	rootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, "show debug info")
}
