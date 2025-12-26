/*
Copyright © 2025 Sneha De <55897319+sneha-afk@users.noreply.github.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"
	"runtime"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/sneha-afk/trovl/internal/models"
	"github.com/spf13/cobra"
)

// Paths to resolve when there are no arguments passed in
var defaultFiles = []string{
	".trovl",
	"trovl.json",
	".trovl.json",
}

func ApplyManifest(manifestFilePath string) {
	allPlatforms := make(map[string]struct{})
	allPlatforms["windows"] = struct{}{}
	allPlatforms["linux"] = struct{}{}
	allPlatforms["darwin"] = struct{}{}

	manifestFile, err := os.ReadFile(manifestFilePath)
	if err != nil {
		log.Fatalf("[ERROR] Apply: could not read manifest file: %v\n", err)
	}

	var manifest models.Manifest
	if err := json.Unmarshal(manifestFile, &manifest); err != nil {
		log.Fatalf("[ERROR] Apply: could not unmarshal manifest: %v\n", err)
	}
	manifest.FillDefaults()

	for _, manifestLink := range manifest.Links {
		platformsUsingSpecifiedLink := make(map[string]struct{})
		platformUsingOverrides := make(map[string]string)

		// 1. Separate out platform overrides from the "platforms"
		for _, plat := range manifestLink.Platforms {
			if plat == "all" {
				platformsUsingSpecifiedLink = maps.Clone(allPlatforms)
				break
			}
			platformsUsingSpecifiedLink[plat] = struct{}{}
		}

		for _, po := range manifestLink.PlatformOverrides {
			if _, ok := platformsUsingSpecifiedLink[po.Platform]; !ok {
				delete(platformsUsingSpecifiedLink, po.Platform)
				platformUsingOverrides[po.Platform] = po.Link
			}
		}

		// 2. Detect current OS and then carry out the links
		var linkToUse string
		if _, ok := platformsUsingSpecifiedLink[runtime.GOOS]; ok {
			linkToUse = manifestLink.Link
		} else {
			linkToUse = platformUsingOverrides[runtime.GOOS]
		}

		linkSpec, err := links.Construct(manifestLink.Target, linkToUse, manifestLink.Relative)
		if errors.Is(err, links.ErrDeclinedOverwrite) {
			continue
		}
		if err != nil {
			log.Fatalf("[ERROR] Apply: could not construct link: %v", err)
		}
		if err := links.Add(linkSpec); err != nil {
			log.Fatalf("[ERROR] Apply: could not add link: %v", err)
		}

		if GlobalState.Verbose {
			log.Printf("[INFO] Add: created link from %v -> %v\n", linkToUse, manifestLink.Target)
		}

	}
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies a link list specified by schema.",
	Long:  `Applies a link list specified by schema to bulk add links or fix as needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apply called")

		if 0 < len(args) {
			for _, manifestFilePath := range args {
				ApplyManifest(manifestFilePath)
			}
		} else {
			for _, path := range defaultFiles {
				ok, err := links.ValidatePath(path)
				if !ok || err != nil {
					continue
				}

				ApplyManifest(path)
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
