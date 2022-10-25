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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/api"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type CreateWorkspace struct {
	Name         string `json:"name"`
	Repo         string `json:"repo"`
	Organization string `json:"organization"`
	Description  string `json:"description"`
	WorkspaceDir string `json:"-"`
}

func init() {
	createWorkspace := &CreateWorkspace{}
	workspaceCmd := NewWorkspaceCmd(createWorkspace)

	workspaceCmd.Flags().StringVar(&createWorkspace.Organization, utils.CommandFlagWorkspaceOrganization, "", utils.MessageCommandWorkspaceOrganization)
	workspaceCmd.MarkFlagRequired(utils.CommandFlagWorkspaceOrganization)

	rootCommand.AddCommand(workspaceCmd)
}

func NewWorkspaceCmd(cmdWorkspace *CreateWorkspace) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandWorkspace,
		Short: utils.MessageCommandWorkspaceShort,
		Long:  utils.MessageCommandWorkspaceLong,
		PreRun: func(cmd *cobra.Command, _ []string) {
			app := core.GetApp()
			organization, err := cmd.Flags().GetString(utils.CommandFlagWorkspaceOrganization)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidFlag)
			}

			if strings.HasPrefix(organization, "@") {
				utils.ErrorOut(utils.MessageErrorCommandWorkspaceInvalidNamespace, utils.ErrorInvalidFlag)
			}

			supabase := api.NewSupabase(utils.ApiUrl, utils.ApiKey, app.Author.Token, "production")
			isUserOrganization, err := supabase.ValidateUserOrganization(app.Author.ID, organization)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidOrganization)
			}

			if !isUserOrganization {
				utils.ErrorOut(utils.MessageErrorCommandWorkspaceInvalidOrganization, utils.ErrorInvalidOrganization)
			}

		},
		Run: func(cmd *cobra.Command, _ []string) {
			cmdWorkspace.Run(cmd)
			cmdWorkspace.CreateWorkTree(cmd)
			cmdWorkspace.Workflows()
			cmdWorkspace.InitGit()
			cmdWorkspace.InitYarn()
			cmdWorkspace.BuildVitePlugin()
			cmdWorkspace.CreateCloudWorkspace()
		},
	}
}

func (ctx *CreateWorkspace) Run(cmd *cobra.Command) {
	nameContent := models.PromptContent{
		Error: utils.MessageErrorCommandWorkspaceNamePrompt,
		Label: utils.MessageCommandWorkspaceNamePrompt,
		Hide:  false,
	}

	repoContent := models.PromptContent{
		Error: utils.MessageErrorCommandWorkspaceRepoPrompt,
		Label: utils.MessageCommandWorkspaceRepoPrompt,
		Hide:  false,
	}

	descriptionContent := models.PromptContent{
		Error: utils.MessageErrorCommandWorkspaceDescriptionPrompt,
		Label: utils.MessageCommandWorkspaceDescriptionPrompt,
		Hide:  false,
	}

	name, err := models.PromptGetInput(nameContent, 3)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	repo, err := models.PromptGetInput(repoContent, 3)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	description, err := models.PromptGetInput(descriptionContent, 3)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	ctx.Description = description
	ctx.Name = name
	ctx.Repo = repo
}

func (ctx *CreateWorkspace) CreateWorkTree(cmd *cobra.Command) {
	config := core.GetConfig()
	config.Progress.SetNumTrackersExpected(6)
	config.Progress.Style().Visibility.Value = false
	app := core.GetApp()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressInit, 2)

	rootNamespace := strings.Join([]string{fmt.Sprintf("@%s", ctx.Organization), slug.Make(ctx.Name)}, "/")
	viteNamespace := strings.Join([]string{fmt.Sprintf("@%s", ctx.Organization), "vite"}, "/")

	ctx.WorkspaceDir = filepath.Join(config.RootDir, slug.Make(ctx.Name))

	if err := os.Mkdir(ctx.WorkspaceDir, 0755); err != nil {
		ctx.CommandError(utils.MessageErrorCommandWorkspaceInvalidDirectory, utils.ErrorCreateDirectory)
	}

	gitCmd := exec.Command("git", "clone", "git@github.com:websublime/sublime-workspace-template.git", ctx.WorkspaceDir)
	_, err := gitCmd.Output()
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidGit)
	}

	packageJson, err := FileTemplates.ReadFile("templates/workspace-package.json")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	vitePackageJson, err := FileTemplates.ReadFile("templates/vite-package.json")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	tsconfigBaseJson, err := FileTemplates.ReadFile("templates/tsconfig-base.json")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	changesetConfigJson, err := FileTemplates.ReadFile("templates/changeset-config.json")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	sublimeConfigJson, err := FileTemplates.ReadFile("templates/sublime.json")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	readmeConfigJson, err := FileTemplates.ReadFile("templates/readme.md")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	pkgJsonFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "package.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &models.PackageJsonFileProps{
		Namespace: rootNamespace,
		Repo:      ctx.Repo,
		Username:  app.Author.Username,
		Email:     app.Author.Email,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	vitePkgJsonFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "libs/vite/package.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = vitePkgJsonFile.WriteString(utils.ProcessString(string(vitePackageJson), &models.ViteJsonFileProps{
		Namespace: viteNamespace,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	tsConfigBaseFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "tsconfig.base.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = tsConfigBaseFile.WriteString(utils.ProcessString(string(tsconfigBaseJson), &models.TsConfigJsonFileProps{
		Namespace: viteNamespace,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	changesetConfigFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".changeset/config.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = changesetConfigFile.WriteString(utils.ProcessString(string(changesetConfigJson), &models.PackageJsonFileProps{
		Namespace: ctx.Repo,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	sublimeConfigFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".sublime.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = sublimeConfigFile.WriteString(utils.ProcessString(string(sublimeConfigJson), &models.SublimeJsonFileProps{
		Namespace:    rootNamespace,
		Name:         slug.Make(ctx.Name),
		Repo:         ctx.Repo,
		Root:         "./",
		Organization: ctx.Organization,
		ID:           "",
		Description:  ctx.Description,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	readmeFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, "README.md"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = readmeFile.WriteString(utils.ProcessString(string(readmeConfigJson), &models.ReadmeFileProps{
		Name:         ctx.Name,
		Repo:         ctx.Repo,
		Organization: ctx.Organization,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	config.DoneProgress()
}

func (ctx *CreateWorkspace) Workflows() {
	config := core.GetConfig()
	config.AddTracker()

	app := core.GetApp()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressWorkflows, 2)
	releaseYaml, err := FileTemplates.ReadFile("templates/workflow-release.yaml")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	featureYaml, err := FileTemplates.ReadFile("templates/workflow-feature.yaml")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	artifactYaml, err := FileTemplates.ReadFile("templates/workflow-artifact.yaml")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}
	snapshotYaml, err := FileTemplates.ReadFile("templates/workflow-snapshot.yaml")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	releaseYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/release.yaml"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = releaseYamlFile.WriteString(utils.ProcessString(string(releaseYaml), &models.ReleaseYamlFileProps{
		Username: app.Author.Username,
		Email:    app.Author.Email,
		Scope:    fmt.Sprintf("@%s", ctx.Organization),
	}, "[[", "]]"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	featureYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/feature.yaml"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = featureYamlFile.WriteString(utils.ProcessString(string(featureYaml), &models.ArtifactsYamlFileProps{
		Version: Version,
	}, "[[", "]]"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	artifactYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/artifact.yaml"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = artifactYamlFile.WriteString(utils.ProcessString(string(artifactYaml), &models.ArtifactsYamlFileProps{
		Version: Version,
	}, "[[", "]]"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	snapshotYamlFile, err := os.Create(filepath.Join(ctx.WorkspaceDir, ".github/workflows/snapshot.yaml"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = snapshotYamlFile.WriteString(utils.ProcessString(string(snapshotYaml), &models.SnapshotsYamlFileProps{
		Version: Version,
	}, "[[", "]]"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	config.DoneProgress()
}

func (ctx *CreateWorkspace) InitGit() {
	config := core.GetConfig()
	config.AddTracker()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressGit, 2)
	_ = os.RemoveAll(filepath.Join(ctx.WorkspaceDir, ".git"))
	_, err := utils.InitGit(ctx.WorkspaceDir)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidGit)
	}

	config.DoneProgress()
}

func (ctx *CreateWorkspace) InitYarn() {
	config := core.GetConfig()
	config.AddTracker()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressYarn, 2)
	_, err := utils.YarnInstall(ctx.WorkspaceDir)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidYarn)
	}

	config.DoneProgress()
}

func (ctx *CreateWorkspace) BuildVitePlugin() {
	config := core.GetConfig()
	config.AddTracker()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressVite, 2)
	_, err := utils.YarnBuild(ctx.WorkspaceDir)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidBuild)
	}

	config.DoneProgress()
}

func (ctx *CreateWorkspace) CreateCloudWorkspace() {
	config := core.GetConfig()
	config.AddTracker()
	app := core.GetApp()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressCloud, 2)
	supabase := api.NewSupabase(utils.ApiUrl, utils.ApiKey, app.Author.Token, "production")
	workspaces, err := supabase.CreateOrganizationWorkspace(ctx.Name, ctx.Repo, ctx.Description, app.OrganizationID)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorInvalidCloudOperation)
	}

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressCloud, 6)
	err = app.UpdateWorkspace(&workspaces[0])
	if err != nil {
		_, _ = supabase.DeleteWorkspaceByID(workspaces[0].ID)
		utils.ErrorOut(err.Error(), utils.ErrorInvalidWorkspace)
	}

	config.UpdateProgress(utils.MessageCommandWorkspaceProgressCloud, 1)
	config.TerminateProgress()
	utils.SuccessOut(utils.MessageCommandWorkspaceSuccess)
}

func (ctx *CreateWorkspace) CommandError(message string, errorType utils.ErrorType) {
	config := core.GetConfig()

	if ctx.WorkspaceDir != "" {
		os.RemoveAll(ctx.WorkspaceDir)
	}

	config.TerminateErrorProgress(fmt.Sprintf("Error: %s", errorType))
	utils.ErrorOut(message, errorType)
}
