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
	"log"
	"os"
	"path/filepath"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/sneha-afk/trovl/internal/models"
	"github.com/spf13/cobra"
)

var (
	useRelative bool
)

// CleanLink defaults to using an absolute filepath, only relative if specified
// Guaranteed that filepath.Clean has been called before returning
func CleanLink(raw string) (string, error) {
	var ret string
	var err error = nil
	if useRelative {
		ret = filepath.Clean(raw)
	} else {
		ret, err = filepath.Abs(raw)
	}
	return ret, err
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a symlink that points to the target file",
	Long: `To add a symlink, specify via:
  trovl add <target> <symlink>
For example, to create a link from a dotfiles/.vimrc to .vimrc at home:
  trovl add ~/dotfiles/.vimrc ~/.vimrc`,
	Run: func(cmd *cobra.Command, args []string) {
		target, err := CleanLink(args[0])
		if err != nil {
			log.Fatalln("[ERROR] Add: invalid path (target)", err)
		}
		symlink, err := CleanLink(args[1])
		if err != nil {
			log.Fatalln("[ERROR] Add: invalid path (symlink)", err)
		}

		targetFile, err := os.Open(target)
		if err != nil {
			log.Fatalln("[ERROR] Add: ", err)
		}
		targetFileInfo, err := targetFile.Stat()
		if err != nil {
			log.Fatalln("[ERROR] Add: could not get target file info: ", err)
		}

		var linkType models.LinkType
		if targetFileInfo.IsDir() {
			linkType = models.LinkDirectory
		} else {
			linkType = models.LinkFile
		}

		targetFile.Close()

		link, err := links.Construct(target, symlink, linkType)
		if err != nil {
			log.Fatalf("[ERROR] Add: could not construct symlink: %v", err)
		}

		if err := links.Add(link); err != nil {
			log.Fatalf("[ERROR] Add: %v	(maybe try running as admin?)\n", err)
		}

		if GlobalState.Verbose {
			log.Printf("[INFO] Add: created link from %v -> %v\n", symlink, target)
		}
	},
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"link", "create"},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVar(&useRelative, "relative", false, "retain relative paths to target")
}
