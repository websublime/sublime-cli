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
	"github.com/spf13/viper"
	"github.com/websublime/sublime-cli/api"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type ActionFlags struct {
	Type        string                   `json:"type"`
	Environment string                   `json:"environment"`
	Sublime     models.SublimeViperProps `json:"-"`
	Packages    []models.SublimePackages `json:"-"`
}

func init() {
	actionFlags := &ActionFlags{
		Sublime:  models.SublimeViperProps{},
		Packages: []models.SublimePackages{},
	}
	actionCmd := NewActionCmd(actionFlags)

	rootCommand.AddCommand(actionCmd)

	actionCmd.Flags().StringVar(&actionFlags.Type, utils.CommandFlagActionType, "branch", "Type of action (branch or tag)")
	actionCmd.Flags().StringVar(&actionFlags.Environment, utils.CommandFlagActionEnv, "develop", "Environment")
}

func NewActionCmd(cmdAction *ActionFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandAction,
		Short: utils.MessageCommandActionShort,
		Long:  utils.MessageCommandActionLong,
		PreRun: func(cmd *cobra.Command, _ []string) {
			isCiEnv := viper.GetBool("CI")

			err := viper.Unmarshal(&cmdAction.Sublime)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidEnvironment)
			}

			if !isCiEnv {
				utils.ErrorOut(utils.MessageErrorCommandActionEnv, utils.ErrorInvalidEnvironment)
			}
		},
		Run: func(cmd *cobra.Command, _ []string) {
			cmdAction.Run(cmd)
		},
	}
}

func getBranchDiffPackages(dir string, pkgs []models.SublimePackages) []models.SublimePackages {
	packages := []models.SublimePackages{}

	changedList, err := utils.GetBranchList(dir)
	if err != nil {
		utils.WarningOut(err.Error())
		os.Exit(0)
	}
	list := strings.Fields(changedList)

	for _, pkg := range pkgs {
		// founded := strings.Contains(changedList, pkg.Name)
		founded := utils.Present(list, pkg.Name)

		if founded {
			packages = append(packages, pkg)
		}
	}

	return packages
}

func (ctx *ActionFlags) Run(cmd *cobra.Command) {
	config := core.GetConfig()

	types := utils.GitType(ctx.Type)
	count, err := utils.GetCommitsCount(config.RootDir)
	counter, err := strconv.ParseInt(count, 10, 0)
	if err != nil || counter <= 0 {
		utils.WarningOut(utils.MessageErrorCommandActionNoCommits)
		os.Exit(0)
	}

	switch types {
	case utils.Branch:
		ctx.Packages = getBranchDiffPackages(config.RootDir, ctx.Sublime.Packages)

		if len(ctx.Packages) <= 0 {
			utils.WarningOut(utils.MessageCommandActionNoPackages)
			os.Exit(0)
		} else {
			utils.SuccessOut(fmt.Sprintf(utils.MessageCommandActionFoundPackages, len(ctx.Packages)))
		}

		ctx.DeployArtifacts()
	case utils.Tag:
		ctx.Packages = ctx.Sublime.Packages

		if len(ctx.Packages) <= 0 {
			utils.WarningOut(utils.MessageCommandActionNoPackages)
			os.Exit(0)
		} else {
			utils.SuccessOut(fmt.Sprintf(utils.MessageCommandActionFoundPackages, len(ctx.Packages)))
		}

		ctx.DeployArtifacts()
		ctx.UpdatePackageVersion()
	default:
		utils.WarningOut(utils.MessageCommandActionTypeUnknown)
		os.Exit(0)
	}
}

func (ctx *ActionFlags) DeployArtifacts() {
	config := core.GetConfig()
	env := utils.EnvType(ctx.Environment)
	supabase := api.NewSupabase(utils.ApiUrl, utils.ApiSecret, utils.ApiSecret, string(env))
	scope := fmt.Sprintf("@%s", ctx.Sublime.Organization)
	isBranch := utils.GitType(ctx.Type) == utils.Branch

	for _, pkg := range ctx.Packages {
		var libDir string = "libs"
		isPackage := pkg.Type == utils.Package

		if isPackage {
			libDir = "packages"
		}

		packageDir := filepath.Join(config.RootDir, libDir, pkg.Name)
		packageDistDir := filepath.Join(packageDir, "dist")

		packageJson := models.PackageJson{}
		data, err := os.ReadFile(filepath.Join(packageDir, "package.json"))
		if err != nil {
			utils.WarningOut(utils.MessageErrorReadFile)
			os.Exit(0)
		}

		err = json.Unmarshal(data, &packageJson)
		if err != nil {
			utils.WarningOut(utils.MessageErrorParseFile)
			os.Exit(0)
		}

		var destinationFolder = ""
		if isBranch {
			destinationFolder = fmt.Sprintf("%s@%s-SNAPSHOT", pkg.Name, packageJson.Version)
		} else {
			destinationFolder = fmt.Sprintf("%s@%s", pkg.Name, packageJson.Version)
		}

		distFiles, err := utils.PathWalk(packageDistDir)
		if err != nil {
			utils.WarningOut(err.Error())
			os.Exit(0)
		}

		for _, file := range distFiles {
			upload, err := supabase.Upload(ctx.Sublime.Organization, file, destinationFolder)
			if err != nil {
				utils.WarningOut(err.Error())
				os.Exit(0)
			}

			utils.InfoOut(fmt.Sprintf(utils.MessageCommandActionUploadFile, ctx.Sublime.Organization, upload.Key))
		}

		manifestBaseLink := fmt.Sprintf("%s/%s/object/public/%s/%s", utils.ApiUrl, api.StorageEndpoint, ctx.Sublime.Organization, destinationFolder)
		manifestJson, _ := FileTemplates.ReadFile("templates/manifest.json")
		manifestFile := core.CreateManifest(manifestJson, core.Manifest{
			Name:    pkg.Name,
			Scope:   scope,
			Repo:    ctx.Sublime.Repo,
			Version: packageJson.Version,
			Scripts: &core.ManifestScripts{
				Main: fmt.Sprintf("%s/%s", manifestBaseLink, filepath.Base(packageJson.Main)),
				Esm:  fmt.Sprintf("%s/%s", manifestBaseLink, filepath.Base(packageJson.Module)),
			},
			Styles: make([]string, 0),
			Docs:   fmt.Sprintf("https://websublime.dev/organization/%s/%s", ctx.Sublime.Organization, pkg.Name),
		})

		var manifestDestination = ""
		if isBranch {
			manifestDestination = fmt.Sprintf("manifests/%s/%s", packageJson.Name, string(env))
		} else {
			manifestDestination = packageJson.Name
		}

		manifest, err := supabase.Upload(ctx.Sublime.Organization, manifestFile.Name(), manifestDestination)
		if err != nil {
			utils.WarningOut(err.Error())
			os.Exit(0)
		}

		utils.InfoOut(fmt.Sprintf(utils.MessageCommandActionUploadFile, ctx.Sublime.Organization, manifest.Key))

		os.Remove(manifestFile.Name())

		utils.SuccessOut(utils.MessageCommandActionArtifact)
	}
}

func (ctx *ActionFlags) UpdatePackageVersion() {
	config := core.GetConfig()
	supabase := api.NewSupabase(utils.ApiUrl, utils.ApiSecret, utils.ApiSecret, "production")

	for _, pkg := range ctx.Packages {
		var libDir string = "libs"
		isPackage := pkg.Type == utils.Package

		if isPackage {
			libDir = "packages"
		}

		packageDir := filepath.Join(config.RootDir, libDir, pkg.Name)

		packageJson := models.PackageJson{}
		data, err := os.ReadFile(filepath.Join(packageDir, "package.json"))
		if err != nil {
			utils.WarningOut(utils.MessageErrorReadFile)
			os.Exit(0)
		}

		err = json.Unmarshal(data, &packageJson)
		if err != nil {
			utils.WarningOut(utils.MessageErrorParseFile)
			os.Exit(0)
		}

		_, err = supabase.UpdateWorkspacePackageVersion(pkg.ID, packageJson.Version)
		if err != nil {
			utils.WarningOut(err.Error())
		}

		utils.SuccessOut(fmt.Sprintf(utils.MessageCommandActionVersionUpdate, pkg.Name, packageJson.Version))
	}
}
