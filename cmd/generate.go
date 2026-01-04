package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a blank manifest file with the current schema.",
	Long: `Generate a blank manifest file with trovl's current schema. By default, this will be generated at the default location of
$XDG_CONFIG_HOME/trovl/manifest.json (see [environment variable usage](../configuration/#environment-variables))`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&cfg.ManifestPath, "path", "p", "", "path to output manifest file")
}
