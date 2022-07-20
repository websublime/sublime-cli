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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/utils"
)

type WorkSpaceCommand struct {
	Name         string
	Repo         string
	Username     string
	Email        string
	Organization string
}

func init() {
	cmd := &WorkSpaceCommand{}
	workspaceCmd := NewWorkspaceCmd(cmd)
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.Flags().StringVar(&cmd.Name, "name", "", "Workspace folder name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("name")

	workspaceCmd.Flags().StringVar(&cmd.Repo, "repo", "", "Github repo shortcut (you/repo) [REQUIRED]")
	workspaceCmd.MarkFlagRequired("repo")

	workspaceCmd.Flags().StringVar(&cmd.Repo, "organization", "", "Github organization name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("organization")

	workspaceCmd.Flags().StringVar(&cmd.Username, "username", "", "Git username")
	workspaceCmd.Flags().StringVar(&cmd.Email, "email", "", "Git email")
}

func NewWorkspaceCmd(cmdWsp *WorkSpaceCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "workspace",
		Short: "Create a workspace project",
		Run: func(cmd *cobra.Command, _ []string) {
			cmdWsp.Run(cmd)
			cmdWsp.Workflows()
			cmdWsp.InitGit()
			cmdWsp.InitYarn()
		},
	}
}

func (ctx *WorkSpaceCommand) Run(cmd *cobra.Command) {
	color.Info.Println("🚀 Creating new workspace: ", ctx.Name)

	if strings.HasPrefix(ctx.Organization, "@") {
		color.Error.Println("Please provide a valid github organization name")
		cobra.CheckErr(errors.New("Invalid organization name"))
	}

	rootNamespace := strings.Join([]string{fmt.Sprintf("@%s", ctx.Organization), slug.Make(ctx.Name)}, "/")
	viteNamespace := strings.Join([]string{ctx.Organization, "vite"}, "/")

	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))
	if err := os.Mkdir(workspaceDir, 0755); err != nil {
		color.Error.Printf("Error creating workspace: %s", ctx.Name)
		cobra.CheckErr(err)
	}

	gitCmd := exec.Command("git", "clone", "git@github.com:websublime/sublime-workspace-template.git", slug.Make(ctx.Name))
	_, err := gitCmd.Output()
	if err != nil {
		color.Error.Println("Unable to clone workspace template")
		cobra.CheckErr(err)
	}

	color.Info.Println("🛢 Template repo cloned. Initializing config files")

	packageJson, err := FileTemplates.ReadFile("templates/workspace-package.json")
	if err != nil {
		color.Error.Println("Unable to read package.json template file")
		cobra.CheckErr(err)
	}
	vitePackageJson, err := FileTemplates.ReadFile("templates/vite-package.json")
	if err != nil {
		color.Error.Println("Unable to read vite.json template file")
		cobra.CheckErr(err)
	}
	tsconfigBaseJson, err := FileTemplates.ReadFile("templates/tsconfig-base.json")
	if err != nil {
		color.Error.Println("Unable to read tsconfig-base.json template file")
		cobra.CheckErr(err)
	}
	changesetConfigJson, err := FileTemplates.ReadFile("templates/changeset-config.json")
	if err != nil {
		color.Error.Println("Unable to read changeset.json template file")
		cobra.CheckErr(err)
	}
	sublimeConfigJson, err := FileTemplates.ReadFile("templates/sublime.json")
	if err != nil {
		color.Error.Println("Unable to read sublime.json template file")
		cobra.CheckErr(err)
	}

	pkgJsonFile, err := os.Create(filepath.Join(workspaceDir, "package.json"))
	if err != nil {
		color.Error.Println("Unable to create package.json file")
		cobra.CheckErr(err)
	}
	pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &utils.PackageJsonVars{
		Namespace: rootNamespace,
		Repo:      ctx.Repo,
		Username:  ctx.Username,
		Email:     ctx.Email,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Package json created and configured!")

	vitePkgJsonFile, err := os.Create(filepath.Join(workspaceDir, "libs/vite/package.json"))
	if err != nil {
		color.Error.Println("Unable to create libs/vite/package.json file")
		cobra.CheckErr(err)
	}
	vitePkgJsonFile.WriteString(utils.ProcessString(string(vitePackageJson), &utils.ViteJsonVars{
		Namespace: viteNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Vite plugin ready!")

	tsConfigBaseFile, err := os.Create(filepath.Join(workspaceDir, "tsconfig.base.json"))
	if err != nil {
		color.Error.Println("Unable to create tsconfig.base.json file")
		cobra.CheckErr(err)
	}
	tsConfigBaseFile.WriteString(utils.ProcessString(string(tsconfigBaseJson), &utils.TsConfigJsonVars{
		Namespace: viteNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Tsconfig created and configured!")

	changesetConfigFile, err := os.Create(filepath.Join(workspaceDir, ".changeset/config.json"))
	if err != nil {
		color.Error.Println("Unable to create .changeset/config.json file")
		cobra.CheckErr(err)
	}
	changesetConfigFile.WriteString(utils.ProcessString(string(changesetConfigJson), &utils.PackageJsonVars{
		Namespace: ctx.Repo,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Changeset created and configured!")

	sublimeConfigFile, err := os.Create(filepath.Join(workspaceDir, ".sublime.json"))
	if err != nil {
		color.Error.Println("Unable to create .sublime.json file")
		cobra.CheckErr(err)
	}
	sublimeConfigFile.WriteString(utils.ProcessString(string(sublimeConfigJson), &utils.SublimeJsonVars{
		Namespace:    rootNamespace,
		Name:         slug.Make(ctx.Name),
		Repo:         ctx.Repo,
		Root:         "./",
		Organization: ctx.Organization,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Sublime json created and configured!")
}

// @deprecated
// Action will be deploy to websublime cloud service
func (ctx *WorkSpaceCommand) Workflows() {
	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))

	releaseYaml, err := FileTemplates.ReadFile("templates/workflow-release.yaml")
	if err != nil {
		color.Error.Println("Unable to read workflow-release.yaml template file")
		cobra.CheckErr(err)
	}
	featureYaml, err := FileTemplates.ReadFile("templates/workflow-feature.yaml")
	if err != nil {
		color.Error.Println("Unable to read workflow-feature.yaml template file")
		cobra.CheckErr(err)
	}
	artifactYaml, err := FileTemplates.ReadFile("templates/workflow-artifact.yaml")
	if err != nil {
		color.Error.Println("Unable to read workflow-artifact.yaml template file")
		cobra.CheckErr(err)
	}

	releaseYamlFile, err := os.Create(filepath.Join(workspaceDir, ".github/workflows/release.yaml"))
	if err != nil {
		color.Error.Println("Unable to create .github/workflows/release.yaml file")
		cobra.CheckErr(err)
	}
	releaseYamlFile.WriteString(utils.ProcessString(string(releaseYaml), &utils.ReleaseYamlVars{
		Username: ctx.Username,
		Email:    ctx.Email,
		Scope:    fmt.Sprintf("@%s", ctx.Organization),
	}, "[[", "]]"))

	color.Info.Println("❤️‍🔥 Github action release created!")

	featureYamlFile, err := os.Create(filepath.Join(workspaceDir, ".github/workflows/feature.yaml"))
	if err != nil {
		color.Error.Println("Unable to create .github/workflows/feature.yaml file")
		cobra.CheckErr(err)
	}
	featureYamlFile.WriteString(utils.ProcessString(string(featureYaml), &utils.ArtifactsVars{
		Version: Version,
	}, "[[", "]]"))

	color.Info.Println("❤️‍🔥 Github action feature created!")

	artifactYamlFile, err := os.Create(filepath.Join(workspaceDir, ".github/workflows/artifact.yaml"))
	if err != nil {
		color.Error.Println("Unable to create .github/workflows/artifact.yaml file")
		cobra.CheckErr(err)
	}
	artifactYamlFile.WriteString(utils.ProcessString(string(artifactYaml), &utils.ArtifactsVars{
		Version: Version,
	}, "[[", "]]"))

	color.Info.Println("❤️‍🔥 Github action artifact created!")
}

func (ctx *WorkSpaceCommand) InitGit() {
	color.Info.Println("❤️‍🔥 Init git on workspace")

	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))

	os.RemoveAll(filepath.Join(workspaceDir, ".git"))
	_, err := utils.InitGit(workspaceDir)

	if err != nil {
		color.Error.Println("Git wasn't enabled on", workspaceDir, ". Please do it manually")
		cobra.CheckErr(err)
	}

}

func (ctx *WorkSpaceCommand) InitYarn() {
	color.Info.Println("❤️‍🔥 Init yarn install on workspace")

	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))

	_, err := utils.YarnInstall(workspaceDir)

	if err != nil {
		color.Error.Println("Yarn wasn't installed on", workspaceDir, ". Please do it manually")
		color.Error.Println("Yarn error:", err.Error())
		cobra.CheckErr(err)
	}

	color.Success.Println("✅ Your app is initialized. Create your first lib or package.")
}
