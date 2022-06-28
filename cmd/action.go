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

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/core/clients"
	"github.com/websublime/sublime-cli/utils"
)

type ActionCommand struct {
	Kind        string
	Client      string
	Bucket      string
	Key         string
	BaseUrl     string
	Environment string
}

func init() {
	cmd := &ActionCommand{}

	actionCmd := NewActionCmd(cmd)
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().StringVar(&cmd.Kind, "kind", "branch", "Kind of action (branch or tag)")
	actionCmd.Flags().StringVar(&cmd.Client, "client", "supabase", "Client to use to upload to storage (supabase, github)")
	actionCmd.Flags().StringVar(&cmd.Environment, "env", "develop", "Environment")

	actionCmd.Flags().StringVar(&cmd.Bucket, "bucket", "", "Bucket storage name [REQUIRED]")
	actionCmd.MarkFlagRequired("bucket")

	actionCmd.Flags().StringVar(&cmd.Key, "key", "", "Api key [REQUIRED]")
	actionCmd.MarkFlagRequired("key")

	actionCmd.Flags().StringVar(&cmd.BaseUrl, "url", "", "Api base url [REQUIRED]")
	actionCmd.MarkFlagRequired("url")
}

func getSublimeBranchPackages(commitsCount int64) []core.Packages {
	sublime := core.GetSublime()
	pkgs := []core.Packages{}

	for key := range sublime.Packages {
		pkgName := sublime.Packages[key].Name
		var output string = ""

		if commitsCount >= 2 {
			output, _ = utils.GetBeforeAndLastDiff(sublime.Root)
		} else if commitsCount == 1 {
			output, _ = utils.GetBeforeDiff(sublime.Root)
		}

		founded := strings.Contains(output, pkgName)

		if founded {
			pkgs = append(pkgs, sublime.Packages[key])
		}
	}

	return pkgs
}

func getSublimeTagPackages() []core.Packages {
	sublime := core.GetSublime()
	pkgs := []core.Packages{}

	for key := range sublime.Packages {
		pkgs = append(pkgs, sublime.Packages[key])
	}

	return pkgs
}

func NewActionCmd(cmdAction *ActionCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "action",
		Short: "Create artifacts on github actions",
		Run: func(cmd *cobra.Command, args []string) {
			cmdAction.ReleaseArtifact(cmd)
		},
	}
}

func (ctx *ActionCommand) ReleaseArtifact(cmd *cobra.Command) {
	color.Info.Println("🥼 Starting Feature Artifacts creation")
	sublime := core.GetSublime()

	var pkgs = []core.Packages{}

	count, _ := utils.GetCommitsCount(sublime.Root)
	counter, err := strconv.ParseInt(count, 10, 0)
	if err != nil || counter <= 0 {
		cobra.CheckErr("No commits founded. Please commit first")
	}

	color.Info.Println("🥼 Commits counted: ", counter)

	supabase := clients.NewSupabase(ctx.BaseUrl, ctx.Key, ctx.Environment)
	github := clients.NewGithub(fmt.Sprintf("https://api.github.com/repos/%s/contents", ctx.BaseUrl), ctx.Key, ctx.Environment)

	if ctx.Kind == "branch" {
		pkgs = getSublimeBranchPackages(counter)
	} else {
		pkgs = getSublimeTagPackages()
	}

	if len(pkgs) <= 0 {
		color.Info.Println("🥼 No packages founded to build artifacts")
	} else {
		color.Info.Println("🥼 Founded", len(pkgs), "package to build artifacts")
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

		color.Info.Println("🥼 Starting creating artifact for:", pkgJson.Name)

		var destinationFolder = ""
		if ctx.Kind == "branch" {
			destinationFolder = fmt.Sprintf("%s@%s-SNAPSHOT", pkgs[key].Name, pkgJson.Version)
		} else {
			destinationFolder = fmt.Sprintf("%s@%s", pkgs[key].Name, pkgJson.Version)
		}

		fileList, err := utils.PathWalk(distFolder)
		if err != nil {
			panic(err)
		}

		color.Info.Println("🥼 Founded", len(fileList), "files to be upload to assets bucket")

		for idx := range fileList {
			if ctx.Client == "supabase" {
				supabase.Upload(ctx.Bucket, fileList[idx], destinationFolder)
			} else if ctx.Client == "github" {
				github.Upload("assets", fileList[idx], destinationFolder)
			}
		}

		color.Info.Println("🥼 Files uploaded to bucket. Starting manifest creation")

		var manifestBaseLink = ""
		if ctx.Client == "supabase" {
			manifestBaseLink = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", ctx.BaseUrl, ctx.Bucket, destinationFolder)
		} else if ctx.Client == "github" {
			manifestBaseLink = fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s/%s", ctx.BaseUrl, "assets", destinationFolder)
		}

		manifestJson, _ := FileTemplates.ReadFile("templates/manifest.json")
		manifestFile := core.CreateManifest(manifestJson, core.Manifest{
			Name:    pkgs[key].Name,
			Scope:   sublime.Scope,
			Repo:    sublime.Repo,
			Version: pkgJson.Version,
			Scripts: &core.ManifestScripts{
				Main: fmt.Sprintf("%s/%s", manifestBaseLink, filepath.Base(pkgJson.Main)),
				Esm:  fmt.Sprintf("%s/%s", manifestBaseLink, filepath.Base(pkgJson.Module)),
			},
			Styles: make([]string, 0),
			Docs:   "",
		})

		var manifestDestination = ""
		if ctx.Kind == "branch" {
			manifestDestination = fmt.Sprintf("%s/%s", pkgJson.Name, ctx.Environment)
		} else {
			manifestDestination = pkgJson.Name
		}

		var manifestLink = ""
		if ctx.Client == "supabase" {
			supabase.Upload("manifests", manifestFile.Name(), manifestDestination)
			manifestLink = fmt.Sprintf("%s/storage/v1/object/public/%s/%s/%s", ctx.BaseUrl, "manifests", manifestDestination, filepath.Base(manifestFile.Name()))
		} else if ctx.Client == "github" {
			github.Upload("manifests", manifestFile.Name(), manifestDestination)
			manifestLink = fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s/%s/%s", ctx.BaseUrl, "manifests", manifestDestination, filepath.Base(manifestFile.Name()))
		}

		os.Remove(manifestFile.Name())

		color.Info.Println("🥼 Manifest uploaded to:", manifestLink)
	}
}
