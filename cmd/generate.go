package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sneha-afk/trovl/internal/manifests"
	"github.com/sneha-afk/trovl/internal/utils"
	"github.com/spf13/cobra"
)

func generate(path string) {
	blankManifest := manifests.Manifest{}
	blankManifest.Links = append(blankManifest.Links, manifests.ManifestLink{
		Target:    "example_target",
		Link:      "example_symlink",
		Platforms: []string{"all"},
		Relative:  false,
		PlatformOverrides: map[string]manifests.PlatformOverride{
			"linux": {Link: "example_override"},
		},
	})
	blankManifest.FillDefaults()

	out := struct {
		Schema string `json:"$schema"`
		manifests.Manifest
	}{
		Schema:   "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
		Manifest: blankManifest,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		State.Logger.Error("Could not marshal manifest", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		State.Logger.Error("Could not create parent directory", "dir", filepath.Dir(path), "error", err)
		os.Exit(1)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		State.Logger.Error("Could not write manifest file", "path", path, "error", err)
		os.Exit(1)
	}

	State.Logger.Info("Generated manifest", "path", path)
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [optional path(s)]",
	Short: "Generate a blank manifest file with the current schema (default: `$XDG_CONFIG_HOME/trovl/manifest.json`).",
	Long: `Generate a blank manifest file with trovl's current schema. By default, this will be generated at the default location of ` +
		"`$XDG_CONFIG_HOME/trovl/manifest.json` (see [environment variable usage](/trovl/configuration/#environment-variables))",
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if 0 < len(args) {
			for _, arg := range args {
				path, err := utils.CleanPath(arg, true)
				if err != nil {
					State.Logger.Error("Could not clean up argument path", "error", err)
					os.Exit(1)
				}
				generate(path)
			}
		} else {
			configDir, err := utils.GetConfigDir()
			if err != nil {
				State.Logger.Error("Could not read config directory", "error", err)
				os.Exit(1)
			}
			path = filepath.Join(configDir, defaultFile)
			generate(path)
		}
	},
	Aliases: []string{"gen"},
	Example: `trovl generate     # Default location
trovl generate here.json`,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
