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
	actionCmd.Flags().StringVar(&cmd.Environment, "env", "develop", "Environment")
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
		Run: func(cmd *cobra.Command, _ []string) {
			cmdAction.ReleaseArtifact(cmd)
		},
	}
}

func (ctx *ActionCommand) ReleaseArtifact(cmd *cobra.Command) {
	color.Info.Println("🥼 Starting Feature Artifacts creation")
	sublime := core.GetSublime()
	scope := fmt.Sprintf("@%s", sublime.Organization)

	var pkgs = []core.Packages{}

	count, _ := utils.GetCommitsCount(sublime.Root)
	counter, err := strconv.ParseInt(count, 10, 0)
	if err != nil || counter <= 0 {
		cobra.CheckErr("No commits founded. Please commit first")
	}

	color.Info.Println("🥼 Commits counted: ", counter)

	supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, utils.ApiSecret, ctx.Environment)

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

		pkgJson := &core.PackageJson{}
		data, _ := os.ReadFile(filepath.Join(libFolder, "package.json"))
		json.Unmarshal(data, &pkgJson)

		color.Info.Println("🥼 Starting creating artifact for:", pkgJson.Name)

		var destinationFolder = ""
		if ctx.Kind == "branch" {
			destinationFolder = fmt.Sprintf("%s/%s@%s-SNAPSHOT", sublime.Namespace, pkgs[key].Name, pkgJson.Version)
		} else {
			destinationFolder = fmt.Sprintf("%s/%s@%s", sublime.Namespace, pkgs[key].Name, pkgJson.Version)
		}

		fileList, err := utils.PathWalk(distFolder)
		if err != nil {
			panic(err)
		}

		color.Info.Println("🥼 Founded", len(fileList), "files to be upload to assets bucket")

		for idx := range fileList {
			supabase.Upload("assets", fileList[idx], destinationFolder)
		}

		color.Info.Println("🥼 Files uploaded to bucket. Starting manifest creation")

		manifestBaseLink := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", utils.ApiUrl, "assets", destinationFolder)

		manifestJson, _ := FileTemplates.ReadFile("templates/manifest.json")
		manifestFile := core.CreateManifest(manifestJson, core.Manifest{
			Name:    pkgs[key].Name,
			Scope:   scope,
			Repo:    sublime.Repo,
			Version: pkgJson.Version,
			Scripts: &core.ManifestScripts{
				Main: fmt.Sprintf("%s/%s", manifestBaseLink, filepath.Base(pkgJson.Main)),
				Esm:  fmt.Sprintf("%s/%s", manifestBaseLink, filepath.Base(pkgJson.Module)),
			},
			Styles: make([]string, 0),
			Docs:   fmt.Sprintf("https://websublime.dev/%s/%s", sublime.Namespace, pkgs[key].Name),
		})

		var manifestDestination = ""
		if ctx.Kind == "branch" {
			manifestDestination = fmt.Sprintf("%s/%s", pkgJson.Name, ctx.Environment)
		} else {
			manifestDestination = pkgJson.Name
		}

		supabase.Upload("manifests", manifestFile.Name(), manifestDestination)
		manifestLink := fmt.Sprintf("%s/storage/v1/object/public/%s/%s/%s", utils.ApiUrl, "manifests", manifestDestination, filepath.Base(manifestFile.Name()))

		os.Remove(manifestFile.Name())

		color.Success.Println("🥼 Manifest uploaded to:", manifestLink)
	}
}
