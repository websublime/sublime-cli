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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/websublime/sublime-cli/api"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type CreateFlags struct {
	Name        string                   `json:"name"`
	Type        utils.PackageType        `json:"type"`
	Template    utils.TemplateType       `json:"template"`
	Description string                   `json:"description"`
	Sublime     models.SublimeViperProps `json:"-"`
	PackageDir  string                   `json:"-"`
	LibTypeDir  string                   `json:"-"`
}

func init() {
	createFlags := &CreateFlags{
		Sublime: models.SublimeViperProps{},
	}
	createCmd := NewCreateCmd(createFlags)

	rootCommand.AddCommand(createCmd)
}

func NewCreateCmd(cmdCreate *CreateFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandCreate,
		Short: utils.MessageCommandCreateShort,
		Long:  utils.MessageCommandCreateLong,
		PreRun: func(cmd *cobra.Command, _ []string) {
			app := core.GetApp()

			err := viper.Unmarshal(&cmdCreate.Sublime)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidWorkspace)
			}

			supabase := api.NewSupabase(utils.ApiUrl, utils.ApiKey, app.Author.Token, "production")
			isUserOrganization, err := supabase.ValidateUserOrganization(app.Author.ID, cmdCreate.Sublime.Organization)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidOrganization)
			}

			if !isUserOrganization {
				utils.ErrorOut(utils.MessageErrorCommandWorkspaceInvalidOrganization, utils.ErrorInvalidOrganization)
			}

			isWorkspaceOrganization, err := supabase.ValidateWorkspaceOrganization(cmdCreate.Sublime.ID, app.OrganizationID)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidOrganization)
			}

			if !isWorkspaceOrganization {
				utils.ErrorOut(utils.MessageErrorCommandWorkspaceInvalidOrganization, utils.ErrorInvalidWorkspace)
			}
		},
		Run: func(cmd *cobra.Command, _ []string) {
			cmdCreate.Run(cmd)
			cmdCreate.CreatePackage()
			cmdCreate.UpdateRepoFiles()
			cmdCreate.YarnLink()
			cmdCreate.CreateCloudPackage()
		},
	}
}

func (ctx *CreateFlags) Run(cmd *cobra.Command) {
	nameContent := models.PromptContent{
		Error: utils.MessageErrorCommandCreateNamePrompt,
		Label: utils.MessageCommandCreateNamePrompt,
		Hide:  false,
	}

	typesContent := models.PromptSelectContent{
		Label: utils.MessageCommandCreateTypePrompt,
		Items: []string{fmt.Sprintf("Package: %s", string(utils.Package)), fmt.Sprintf("Library: %s", string(utils.Library))},
	}

	templateContent := models.PromptSelectContent{
		Label: utils.MessageCommandCreateTemplatePrompt,
		Items: []string{
			fmt.Sprintf("SolidJS: %s", string(utils.Solid)),
			fmt.Sprintf("Lit.dev: %s", string(utils.Lit)),
			fmt.Sprintf("Vue: %s", string(utils.Vue)),
			fmt.Sprintf("Typescript: %s", string(utils.Typescript)),
		},
	}

	name, err := models.PromptGetInput(nameContent, 3)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	idxType, _, err := models.PromptGetSelect(typesContent)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	idxTemplate, _, err := models.PromptGetSelect(templateContent)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	if idxType == 0 {
		ctx.Type = utils.Package
	} else {
		ctx.Type = utils.Library
	}

	if idxTemplate == 0 {
		ctx.Template = utils.Solid
	} else if idxTemplate == 1 {
		ctx.Template = utils.Lit
	} else if idxTemplate == 2 {
		ctx.Template = utils.Vue
	} else {
		ctx.Template = utils.Typescript
	}

	ctx.Name = slug.Make(name)
}

func (ctx *CreateFlags) CreatePackage() {
	config := core.GetConfig()
	config.Progress.SetNumTrackersExpected(5)
	config.Progress.Style().Visibility.Value = false

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandCreateProgressInit, 2)

	ctx.LibTypeDir = "libs"

	if ctx.Type == utils.Package {
		ctx.LibTypeDir = "packages"
	}

	ctx.PackageDir = filepath.Join(config.RootDir, ctx.LibTypeDir, ctx.Name)

	scope := fmt.Sprintf("@%s", ctx.Sublime.Organization)
	libNamespace := strings.Join([]string{scope, ctx.Name}, "/")
	viteRel, err := filepath.Rel(ctx.PackageDir, filepath.Join(config.RootDir, "libs/vite"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorMissingDirectory)
	}

	var templateLink string = ""
	for i := range utils.TemplatesMap {
		if utils.TemplatesMap[i].Template == ctx.Template {
			templateLink = utils.TemplatesMap[i].Link
			break
		}
	}

	if templateLink == "" {
		ctx.CommandError(utils.MessageErrorCommandCreateTemplateInvalid, utils.ErrorInvalidTemplate)
	}

	gitCmd := exec.Command("git", "clone", templateLink, ctx.PackageDir)
	_, err = gitCmd.Output()
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidGit)
	}

	var libPackageJson = "templates/lib-package.json"
	var libTsconfigJson = "templates/tsconfig-lib.json"
	var libViteConfigJson = "templates/vite-config-typescript.json"

	if ctx.Template == utils.Solid {
		libPackageJson = "templates/lib-package-solid.json"
		libTsconfigJson = "templates/tsconfig-lib-solid.json"
		libViteConfigJson = "templates/vite-config-solid.json"
	}

	if ctx.Template == utils.Vue {
		libPackageJson = "templates/lib-package-vue.json"
		libTsconfigJson = "templates/tsconfig-lib-vue.json"
		libViteConfigJson = "templates/vite-config-vue.json"
	}

	if ctx.Template == utils.Lit {
		libViteConfigJson = "templates/vite-config-lit.json"
	}

	packageJson, err := FileTemplates.ReadFile(libPackageJson)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	apiExtractorJson, err := FileTemplates.ReadFile("templates/api-extractor-lib.json")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	tsConfigJson, err := FileTemplates.ReadFile(libTsconfigJson)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	viteConfigJson, err := FileTemplates.ReadFile(libViteConfigJson)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	pkgJsonFile, err := os.Create(filepath.Join(ctx.PackageDir, "package.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}

	_, err = pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &models.PackageJsonFileProps{
		Namespace: libNamespace,
		Repo:      ctx.Sublime.Repo,
		Name:      ctx.Name,
		Scope:     scope,
		Type:      ctx.LibTypeDir,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	apiExtractorFile, err := os.Create(filepath.Join(ctx.PackageDir, "api-extractor.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = apiExtractorFile.WriteString(utils.ProcessString(string(apiExtractorJson), &models.ApiExtractorFileProps{
		Name: ctx.Name,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	tsConfigFile, err := os.Create(filepath.Join(ctx.PackageDir, "tsconfig.json"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	_, err = tsConfigFile.WriteString(utils.ProcessString(string(tsConfigJson), &models.TsConfigJsonFileProps{
		Namespace: libNamespace,
		Vite:      viteRel,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	viteConfigFile, err := os.Create(filepath.Join(ctx.PackageDir, "vite.config.js"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}
	viteConfigFile.WriteString(utils.ProcessString(string(viteConfigJson), &models.ViteJsonFileProps{
		Scope: scope,
		Name:  ctx.Name,
	}, "{{", "}}"))
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidTemplate)
	}

	config.DoneProgress()
}

func (ctx *CreateFlags) UpdateRepoFiles() {
	config := core.GetConfig()
	app := core.GetApp()

	config.AddTracker()

	scope := fmt.Sprintf("@%s", ctx.Sublime.Organization)

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandCreateProgressUpdate, 2)

	update := map[string]interface{}{
		"namespace":    ctx.Sublime.Namespace,
		"name":         ctx.Sublime.Name,
		"repo":         ctx.Sublime.Repo,
		"root":         "./",
		"organization": ctx.Sublime.Organization,
		"id":           ctx.Sublime.ID,
		"description":  ctx.Sublime.Description,
		"packages": append(ctx.Sublime.Packages, models.SublimePackages{
			Name:        ctx.Name,
			Scope:       scope,
			Type:        ctx.Type,
			Description: ctx.Description,
			ID:          "",
		}),
	}

	data, err := json.MarshalIndent(update, "", " ")
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidaIndentation)
	}

	err = os.WriteFile(filepath.Join(config.RootDir, ".sublime.json"), data, 0644)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}

	tsConfigBase, err := app.GetTsconfig()
	if err != nil {
		app.RemoveConfigurationsOnPackageError(ctx.Name)
		ctx.CommandError(err.Error(), utils.ErrorInvalidTypescript)
	}

	tsConfigBase.References = append(tsConfigBase.References, models.TsConfigReferences{
		Path: filepath.Join("./", ctx.LibTypeDir, ctx.Name),
		Name: filepath.Join(scope, ctx.Name),
	})

	tsconfig, err := json.MarshalIndent(tsConfigBase, "", " ")
	if err != nil {
		app.RemoveConfigurationsOnPackageError(ctx.Name)
		ctx.CommandError(err.Error(), utils.ErrorInvalidaIndentation)
	}

	err = os.WriteFile(filepath.Join(config.RootDir, "tsconfig.base.json"), tsconfig, 0644)
	if err != nil {
		app.RemoveConfigurationsOnPackageError(ctx.Name)
		ctx.CommandError(err.Error(), utils.ErrorCreateFile)
	}

	config.DoneProgress()
}

func (ctx *CreateFlags) YarnLink() {
	config := core.GetConfig()
	app := core.GetApp()
	config.AddTracker()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandCreateProgressYarn, 2)

	_, err := utils.YarnInstall(config.RootDir)
	if err != nil {
		app.RemoveConfigurationsOnPackageError(ctx.Name)
		ctx.CommandError(err.Error(), utils.ErrorInvalidYarn)
	}

	config.DoneProgress()
}

func (ctx *CreateFlags) CreateCloudPackage() {
	config := core.GetConfig()
	app := core.GetApp()
	config.AddTracker()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandCreateProgressCloud, 2)

	supabase := api.NewSupabase(utils.ApiUrl, utils.ApiKey, app.Author.Token, "production")
	packages, err := supabase.CreateWorkspacePackage(ctx.Name, ctx.Description, ctx.Type, ctx.Template, ctx.Sublime.ID)
	if err != nil {
		app.RemoveConfigurationsOnPackageError(ctx.Name)
		ctx.CommandError(err.Error(), utils.ErrorInvalidCloudOperation)
	}

	config.UpdateProgress(utils.MessageCommandCreateProgressCloud, 6)
	err = app.UpdatePackage(&packages[0])
	if err != nil {
		app.RemoveConfigurationsOnPackageError(ctx.Name)
		_, _ = supabase.DeletePackageByID(packages[0].ID)
		ctx.CommandError(err.Error(), utils.ErrorInvalidCloudOperation)
	}

	config.UpdateProgress(utils.MessageCommandCreateProgressCloud, 1)
	config.TerminateProgress()
	utils.SuccessOut(utils.MessageCommandCreateSuccess)
}

func (ctx *CreateFlags) CommandError(message string, errorType utils.ErrorType) {
	config := core.GetConfig()

	if ctx.PackageDir != "" {
		os.RemoveAll(ctx.PackageDir)
	}

	config.TerminateErrorProgress(fmt.Sprintf("Error: %s", errorType))
	utils.ErrorOut(message, errorType)
}
