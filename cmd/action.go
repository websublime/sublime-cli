/*
Copyright © 2022 Websublime.dev organization@websublime.dev

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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/core/clients"
	"github.com/websublime/sublime-cli/utils"
)

type ActionCommand struct {
	Kind    string
	Client  string
	Bucket  string
	Key     string
	BaseUrl string
}

func init() {
	cmd := &ActionCommand{}

	actionCmd := NewActionCmd(cmd)
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().StringVar(&cmd.Kind, "kind", "branch", "Kind of action (branch or tag)")
	actionCmd.Flags().StringVar(&cmd.Client, "client", "supabase", "Client to use to upload to storage")

	actionCmd.Flags().StringVar(&cmd.Bucket, "bucket", "", "Bucket storage name [REQUIRED]")
	actionCmd.MarkFlagRequired("bucket")

	actionCmd.Flags().StringVar(&cmd.Key, "key", "", "Api key [REQUIRED]")
	actionCmd.MarkFlagRequired("key")

	actionCmd.Flags().StringVar(&cmd.BaseUrl, "url", "", "Api base url [REQUIRED]")
	actionCmd.MarkFlagRequired("url")
}

func NewActionCmd(cmdAction *ActionCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "action",
		Short: "Run action",
		Run: func(cmd *cobra.Command, args []string) {
			if cmdAction.Kind == "tag" {
				cmdAction.Tag(cmd)
			}

			if cmdAction.Kind == "branch" {
				cmdAction.Branch((cmd))
			}
		},
	}
}

func (ctx *ActionCommand) Branch(cmd *cobra.Command) {
	sublime := core.GetSublime()
	dir, _ := os.Getwd()
	pkgs := []core.Packages{}

	count, _ := utils.GetCommitsCount(sublime.Root)
	counter, err := strconv.ParseInt(strings.Replace(count, "\n", "", -1), 10, 0)
	if err != nil || counter <= 0 {
		cobra.CheckErr("No commits founded. Please commit first")
	}

	lastCommit, _ := utils.GetLastCommit(sublime.Root)
	// beforeCommit, _ := utils.GetBeforeLastCommit(sublime.Root)
	hash, _ := utils.GetShortCommit(sublime.Root, lastCommit)
	supabase := clients.NewSupabase(ctx.BaseUrl, ctx.Key)

	for key := range sublime.Packages {
		pkgName := sublime.Packages[key].Name
		var output string = ""

		if counter >= 2 {
			output, _ = utils.GetBeforeAndLastDiff(dir, pkgName)
		} else if counter == 1 {
			output, _ = utils.GetBeforeDiff(dir, pkgName)
		}

		counted := strings.Count(output, "\n")

		if counted > 0 {
			pkgs = append(pkgs, sublime.Packages[key])
		}
	}

	for key := range pkgs {
		var libDir string
		if pkgs[key].Type == "lib" {
			libDir = "libs"
		}

		if pkgs[key].Type == "pkg" {
			libDir = "packages"
		}

		libFolder := filepath.Join(sublime.Root, libDir, pkgs[key].Name)
		distFolder := filepath.Join(libFolder, "dist")

		pkgJson := &utils.PackageJson{}
		data, _ := os.ReadFile(filepath.Join(libFolder, "package.json"))
		json.Unmarshal(data, &pkgJson)

		destinationFolder := fmt.Sprintf("packages/%s@%s-%s", pkgs[key].Name, pkgJson.Version, hash)

		supabase.Upload(distFolder, destinationFolder)
	}
}

func (ctx *ActionCommand) Tag(cmd *cobra.Command) {}
