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
	"github.com/websublime/sublime-cli/core/clients"
	"github.com/websublime/sublime-cli/utils"
)

type WorkSpaceCommand struct {
	Name         string
	Repo         string
	Username     string
	Email        string
	Organization string
	Description  string
	WorkspaceDir string
}

// Init workspace command and flags
func init() {
	cmd := &WorkSpaceCommand{}
	workspaceCmd := NewWorkspaceCmd(cmd)
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.Flags().StringVar(&cmd.Name, "name", "", "Workspace folder name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("name")

	workspaceCmd.Flags().StringVar(&cmd.Repo, "repo", "", "Github repo shortcut (you/repo) [REQUIRED]")
	workspaceCmd.MarkFlagRequired("repo")

	workspaceCmd.Flags().StringVar(&cmd.Organization, "organization", "", "Github organization name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("organization")

	workspaceCmd.Flags().StringVar(&cmd.Description, "description", "", "Workspace description")
	workspaceCmd.Flags().StringVar(&cmd.Username, "username", "", "Git username")
	workspaceCmd.Flags().StringVar(&cmd.Email, "email", "", "Git email")
}

// Create new workspace command. Run sub tasks and perform
// a pre run to test if current loggedin user belongs to the
// organization that is trying to create the workspace
func NewWorkspaceCmd(cmdWsp *WorkSpaceCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "workspace",
		Short: "Create a workspace project",
		Long:  "Workspace command will create a monorepo directory with all configurations needeed to initiate. Also it will be created on the cloud platform",
		PreRun: func(cmd *cobra.Command, _ []string) {
			sublime := core.GetSublime()
			supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, sublime.Author.Token, "production")
			response, err := supabase.GetUserOrganizations()
			if err != nil {
				cmdWsp.ErrorOut(err, "User not found in this organization")
			}

			orgs := []core.Organization{}

			err = json.Unmarshal([]byte(response), &orgs)
			if err != nil {
				cmdWsp.ErrorOut(err, "Unable to parse organization")
			}

			organization, err := cmd.Flags().GetString("organization")
			if err != nil {
				cmdWsp.ErrorOut(err, "Organization parameter invalid")
			}

			var isUserOrganization bool = false
			for i := range orgs {
				if orgs[i].Name == organization {
					isUserOrganization = true
					sublime.ID = orgs[i].ID
					break
				}
			}

			if !isUserOrganization {
				cmdWsp.ErrorOut(errors.New("Invalid organization user"), "User is not valid to this organization")
			}
		},
		Run: func(cmd *cobra.Command, _ []string) {
			cmdWsp.Run(cmd)
			cmdWsp.Workflows()
			cmdWsp.InitGit()
			cmdWsp.InitYarn()
			cmdWsp.BuildVitePlugin()
			cmdWsp.CreateCloudWorkspace()
		},
	}
}

// Run initial steps of configs and clone template repo
func (ctx *WorkSpaceCommand) Run(cmd *cobra.Command) {
	color.Info.Println("🚀 Creating new workspace: ", ctx.Name)
	sublime := core.GetSublime()

	if strings.HasPrefix(ctx.Organization, "@") {
		ctx.ErrorOut(errors.New("Invalid organization name"), "Please provide a valid github organization name")
	}

	rootNamespace := strings.Join([]string{fmt.Sprintf("@%s", ctx.Organization), slug.Make(ctx.Name)}, "/")
	viteNamespace := strings.Join([]string{fmt.Sprintf("@%s", ctx.Organization), "vite"}, "/")

	var username string
	var email string

	if ctx.Username != "" {
		username = ctx.Username
	} else {
		username = sublime.Author.Username
	}

	if ctx.Email != "" {
		email = ctx.Email
	} else {
		email = sublime.Author.Email
	}

	ctx.WorkspaceDir = filepath.Join(sublime.Root, slug.Make(ctx.Name))

	if err := os.Mkdir(ctx.WorkspaceDir, 0755); err != nil {
		ctx.WorkspaceDir = ""
		ctx.ErrorOut(err, fmt.Sprintf("Error creating workspace: %s", ctx.Name))
	}

	gitCmd := exec.Command("git", "clone", "git@github.com:websublime/sublime-workspace-template.git", ctx.WorkspaceDir)
	_, err := gitCmd.Output()
	if err != nil {
		ctx.ErrorOut(err, "Unable to clone workspace template")
	}

	color.Info.Println("🛢 Template repo cloned. Initializing config files")

	packageJson, err := FileTemplates.ReadFile("templates/workspace-package.json")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read package.json template file")
	}
	vitePackageJson, err := FileTemplates.ReadFile("templates/vite-package.json")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read vite.json template file")
	}
	tsconfigBaseJson, err := FileTemplates.ReadFile("templates/tsconfig-base.json")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read tsconfig-base.json template file")
	}
	changesetConfigJson, err := FileTemplates.ReadFile("templates/changeset-config.json")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read changeset.json template file")
	}
	sublimeConfigJson, err := FileTemplates.ReadFile("templates/sublime.json")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read sublime.json template file")
	}
	readmeConfigJson, err := FileTemplates.ReadFile("templates/readme.md")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read readme.md template file")
	}

	pkgJsonFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "package.json"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create package.json file")
	}
	pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &utils.PackageJsonVars{
		Namespace: rootNamespace,
		Repo:      ctx.Repo,
		Username:  username,
		Email:     email,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Package json created and configured!")

	vitePkgJsonFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "libs/vite/package.json"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create libs/vite/package.json file")
	}
	vitePkgJsonFile.WriteString(utils.ProcessString(string(vitePackageJson), &utils.ViteJsonVars{
		Namespace: viteNamespace,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Vite plugin ready!")

	tsConfigBaseFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "tsconfig.base.json"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create tsconfig.base.json file")
	}
	tsConfigBaseFile.WriteString(utils.ProcessString(string(tsconfigBaseJson), &utils.TsConfigJsonVars{
		Namespace: viteNamespace,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Tsconfig created and configured!")

	changesetConfigFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".changeset/config.json"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create .changeset/config.json file")
	}
	changesetConfigFile.WriteString(utils.ProcessString(string(changesetConfigJson), &utils.PackageJsonVars{
		Namespace: ctx.Repo,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Changeset created and configured!")

	sublimeConfigFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".sublime.json"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create .sublime.json file")
	}
	sublimeConfigFile.WriteString(utils.ProcessString(string(sublimeConfigJson), &utils.SublimeJsonVars{
		Namespace:    rootNamespace,
		Name:         slug.Make(ctx.Name),
		Repo:         ctx.Repo,
		Root:         "./",
		Organization: ctx.Organization,
		ID:           sublime.ID,
		Description:  ctx.Description,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Sublime json created and configured!")

	readmeFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "README.md"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create README.md file")
	}
	readmeFile.WriteString(utils.ProcessString(string(readmeConfigJson), &utils.ReadmeVars{
		Name:         ctx.Name,
		Repo:         ctx.Repo,
		Organization: ctx.Organization,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Readme file created!")
}

// Creat github actions that will publish artifacts to websublime cloud service
func (ctx *WorkSpaceCommand) Workflows() {
	color.Info.Println("🚀 Creating workflows ")

	releaseYaml, err := FileTemplates.ReadFile("templates/workflow-release.yaml")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read workflow-release.yaml template file")
	}
	featureYaml, err := FileTemplates.ReadFile("templates/workflow-feature.yaml")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read workflow-feature.yaml template file")
	}
	artifactYaml, err := FileTemplates.ReadFile("templates/workflow-artifact.yaml")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read workflow-artifact.yaml template file")
	}

	releaseYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/release.yaml"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create .github/workflows/release.yaml file")
	}
	releaseYamlFile.WriteString(utils.ProcessString(string(releaseYaml), &utils.ReleaseYamlVars{
		Username: ctx.Username,
		Email:    ctx.Email,
		Scope:    fmt.Sprintf("@%s", ctx.Organization),
	}, "[[", "]]"))

	color.Success.Println("❤️‍🔥 Github action release created!")

	featureYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/feature.yaml"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create .github/workflows/feature.yaml file")
	}
	featureYamlFile.WriteString(utils.ProcessString(string(featureYaml), &utils.ArtifactsVars{
		Version: Version,
	}, "[[", "]]"))

	color.Success.Println("❤️‍🔥 Github action feature created!")

	artifactYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/artifact.yaml"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create .github/workflows/artifact.yaml file")
	}
	artifactYamlFile.WriteString(utils.ProcessString(string(artifactYaml), &utils.ArtifactsVars{
		Version: Version,
	}, "[[", "]]"))

	color.Success.Println("❤️‍🔥 Github action artifact created!")
}

// Init git on workspace directory
func (ctx *WorkSpaceCommand) InitGit() {
	color.Info.Println("❤️‍🔥 Init git on workspace")

	os.RemoveAll(filepath.Join(ctx.WorkspaceDir, ".git"))
	_, err := utils.InitGit(ctx.WorkspaceDir)
	if err != nil {
		ctx.ErrorOut(err, fmt.Sprintf("Git wasn't enabled on: %s", ctx.WorkspaceDir))
	}

	color.Success.Println("❤️‍🔥 Workspace git init successful!")
}

// Install all packages thru Yarn
func (ctx *WorkSpaceCommand) InitYarn() {
	color.Info.Println("❤️‍🔥 Init yarn install on workspace")

	_, err := utils.YarnInstall(ctx.WorkspaceDir)
	if err != nil {
		ctx.ErrorOut(err, fmt.Sprintf("Yarn wasn't installed on: %s", ctx.WorkspaceDir))
	}

	color.Success.Println("❤️‍🔥 Yarn installed.")
}

// Install all packages thru Yarn
func (ctx *WorkSpaceCommand) BuildVitePlugin() {
	color.Info.Println("❤️‍🔥 Building vite plugin")

	_, err := utils.YarnBuild(filepath.Join(ctx.WorkspaceDir, "libs/vite/"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to build vite plugin")
	}

	color.Success.Println("❤️‍🔥 Vite plugin build with success")
}

// Post the new workspace on organization
func (ctx *WorkSpaceCommand) CreateCloudWorkspace() {
	color.Info.Println("❤️‍🔥 Creating workspace on cloud platform")
	sublime := core.GetSublime()

	supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, sublime.Author.Token, "production")
	_, err := supabase.CreateOrganizationWorkspace(slug.Make(ctx.Name), ctx.Repo, ctx.Description, sublime.ID)
	if err != nil {
		ctx.ErrorOut(err, "Unable to create workspace on cloud")
	}

	color.Success.Println("❤️‍🔥 Workspace created!")
	color.Success.Println("✅ Your app is initialized. Create your first lib or package.")
}

// Prints out errors and delete workspace directory
func (ctx *WorkSpaceCommand) ErrorOut(err error, msg string) {
	if ctx.WorkspaceDir != "" {
		os.RemoveAll(ctx.WorkspaceDir)
	}

	color.Error.Println(msg, err)
	cobra.CheckErr(err)
}
