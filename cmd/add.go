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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("add called")
		// for _, arg := range args {
		// 	fmt.Println(arg)
		// }

		target, err := CleanLink(args[0])
		if err != nil {
			log.Fatalln("Error: invalid filepath (target): ", target)
		}
		symlink, err := CleanLink(args[1])
		if err != nil {
			log.Fatalln("Error: invalid filepath (symlink): ", symlink)
		}

		targetFile, err := os.Open(target)
		if err != nil {
			log.Fatalln("Error: could not open target file: ", err)
		}
		targetFileInfo, err := targetFile.Stat()
		if err != nil {
			log.Fatalln("Error: could not get target file's info: ", err)
		}

		var linkType models.LinkType
		if targetFileInfo.IsDir() {
			linkType = models.LinkDirectory
		} else {
			linkType = models.LinkFile
		}

		link := models.Link{
			Target:    target,
			LinkMount: symlink,
			Type:      linkType,
		}

		links.Add(link)

		// TODO: verbose mode to print all ops
	},
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"link"},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().BoolVar(&useRelative, "relative", false, "Retain relative path to target")
}
