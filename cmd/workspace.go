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
	Name     string
	Scope    string
	Repo     string
	Username string
	Email    string
}

func init() {
	cmd := &WorkSpaceCommand{}
	workspaceCmd := NewWorkspaceCmd(cmd)
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.Flags().StringVar(&cmd.Name, "name", "", "Workspace folder name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("name")

	workspaceCmd.Flags().StringVar(&cmd.Scope, "scope", "", "Workspace scope name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("scope")

	workspaceCmd.Flags().StringVar(&cmd.Repo, "repo", "", "Github repo shortcut (you/repo) [REQUIRED]")
	workspaceCmd.MarkFlagRequired("repo")

	workspaceCmd.Flags().StringVar(&cmd.Username, "username", "", "Git username [REQUIRED]")
	workspaceCmd.MarkFlagRequired("username")

	workspaceCmd.Flags().StringVar(&cmd.Email, "email", "", "Git email [REQUIRED]")
	workspaceCmd.MarkFlagRequired("email")
}

func NewWorkspaceCmd(cmdWsp *WorkSpaceCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "workspace",
		Short: "Create a workspace project",
		Run: func(cmd *cobra.Command, args []string) {
			cmdWsp.Run(cmd)
			cmdWsp.Workflows()
			cmdWsp.InitGit()
			cmdWsp.InitYarn()
		},
	}
}

func (ctx *WorkSpaceCommand) Run(cmd *cobra.Command) {
	color.Info.Println("🚀 Creating new workspace: ", ctx.Name)

	rootNamespace := strings.Join([]string{ctx.Scope, slug.Make(ctx.Name)}, "/")
	viteNamespace := strings.Join([]string{ctx.Scope, "vite"}, "/")

	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))
	if err := os.Mkdir(workspaceDir, 0755); err != nil {
		color.Error.Printf("Error creating workspace: %s", ctx.Name)
		os.Exit(1)
	}

	gitCmd := exec.Command("git", "clone", "git@github.com:websublime/sublime-workspace-template.git", slug.Make(ctx.Name))
	_, err := gitCmd.Output()
	if err != nil {
		color.Error.Println("Unable to clone workspace template")
		cobra.CheckErr(err)
	}

	color.Info.Println("🛢 Template repo cloned. Initializing config files")

	packageJson, _ := FileTemplates.ReadFile("templates/workspace-package.json")
	vitePackageJson, _ := FileTemplates.ReadFile("templates/vite-package.json")
	tsconfigBaseJson, _ := FileTemplates.ReadFile("templates/tsconfig-base.json")
	changesetConfigJson, _ := FileTemplates.ReadFile("templates/changeset-config.json")
	sublimeConfigJson, _ := FileTemplates.ReadFile("templates/sublime.json")

	pkgJsonFile, _ := os.Create(filepath.Join(workspaceDir, "package.json"))
	pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &utils.PackageJsonVars{
		Namespace: rootNamespace,
		Repo:      ctx.Repo,
		Username:  ctx.Username,
		Email:     ctx.Email,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Package json created and configured!")

	vitePkgJsonFile, _ := os.Create(filepath.Join(workspaceDir, "libs/vite/package.json"))
	vitePkgJsonFile.WriteString(utils.ProcessString(string(vitePackageJson), &utils.ViteJsonVars{
		Namespace: viteNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Vite plugin ready!")

	tsConfigBaseFile, _ := os.Create(filepath.Join(workspaceDir, "tsconfig.base.json"))
	tsConfigBaseFile.WriteString(utils.ProcessString(string(tsconfigBaseJson), &utils.TsConfigJsonVars{
		Namespace: viteNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Tsconfig created and configured!")

	changesetConfigFile, _ := os.Create(filepath.Join(workspaceDir, ".changeset/config.json"))
	changesetConfigFile.WriteString(utils.ProcessString(string(changesetConfigJson), &utils.PackageJsonVars{
		Namespace: rootNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Changeset created and configured!")

	sublimeConfigFile, _ := os.Create(filepath.Join(workspaceDir, ".sublime.json"))
	sublimeConfigFile.WriteString(utils.ProcessString(string(sublimeConfigJson), &utils.SublimeJsonVars{
		Namespace: rootNamespace,
		Name:      slug.Make(ctx.Name),
		Scope:     ctx.Scope,
		Repo:      ctx.Repo,
		Root:      "./",
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Sublime json created and configured!")
}

func (ctx *WorkSpaceCommand) Workflows() {
	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))

	releaseYaml, _ := FileTemplates.ReadFile("templates/workflow-release.yaml")
	featureYaml, _ := FileTemplates.ReadFile("templates/workflow-feature.yaml")
	artifactYaml, _ := FileTemplates.ReadFile("templates/workflow-artifact.yaml")

	releaseYamlFile, _ := os.Create(filepath.Join(workspaceDir, ".github/workflows/release.yaml"))
	releaseYamlFile.WriteString(utils.ProcessString(string(releaseYaml), &utils.ReleaseYamlVars{
		Username: ctx.Username,
		Email:    ctx.Email,
		Scope:    ctx.Scope,
	}, "[[", "]]"))

	color.Info.Println("❤️‍🔥 Github action release created!")

	featureYamlFile, _ := os.Create(filepath.Join(workspaceDir, ".github/workflows/feature.yaml"))
	featureYamlFile.WriteString(utils.ProcessString(string(featureYaml), &utils.EmptyVars{}, "[[", "]]"))

	color.Info.Println("❤️‍🔥 Github action feature created!")

	artifactYamlFile, _ := os.Create(filepath.Join(workspaceDir, ".github/workflows/artifact.yaml"))
	artifactYamlFile.WriteString(utils.ProcessString(string(artifactYaml), &utils.EmptyVars{}, "[[", "]]"))

	color.Info.Println("❤️‍🔥 Github action artifact created!")
}

func (ctx *WorkSpaceCommand) InitGit() {
	color.Info.Println("❤️‍🔥 Init git on workspace")

	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))

	os.RemoveAll(filepath.Join(workspaceDir, ".git"))
	_, err := utils.InitGit(workspaceDir)

	if err != nil {
		color.Error.Println("Git wasn't enabled on", workspaceDir, ". Please do it manually")
	}

}

func (ctx *WorkSpaceCommand) InitYarn() {
	workspaceDir := filepath.Join(core.GetSublime().Root, slug.Make(ctx.Name))

	_, err := utils.YarnInstall(workspaceDir)

	if err != nil {
		color.Error.Println("Yarn wasn't installed on", workspaceDir, ". Please do it manually")
	}

	color.Success.Println("✅ Your app is initialized. Create your first lib or package.")
}
