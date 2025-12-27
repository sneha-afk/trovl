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
	"fmt"
	"log"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/sneha-afk/trovl/internal/manifests"
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
	Use:   "apply",
	Short: "Applies a link list specified by schema.",
	Long:  `Applies a link list specified by schema to bulk add links or fix as needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apply called")

		if 0 < len(args) {
			for _, manifestFilePath := range args {
				m, err := manifests.New(manifestFilePath)
				if err != nil {
					log.Fatalf("[ERROR] Apply: could not read manifest file: %v\n", err)
				}
				m.Apply(GlobalState.Verbose)
			}
		} else {
			for _, path := range defaultFiles {
				ok, err := links.ValidatePath(path)
				if !ok || err != nil {
					continue
				}

				m, err := manifests.New(path)
				if err != nil {
					log.Fatalf("[ERROR] Apply: could not read manifest file: %v\n", err)
				}
				m.Apply(GlobalState.Verbose)
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
