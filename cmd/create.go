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

type CreateCommand struct {
	Name        string
	Type        string
	Template    string
	Description string
	PackageDir  string
	WorkspaceID string
	Templates   []LibTemplate
}

type LibTemplate struct {
	Type string
	Link string
}

func init() {
	templates := []LibTemplate{
		{
			Type: "vue",
			Link: "git@github.com:websublime/sublime-vue-template.git",
		},
		{
			Type: "lit",
			Link: "git@github.com:websublime/sublime-lit-template.git",
		},
		{
			Type: "solid",
			Link: "git@github.com:websublime/sublime-solid-template.git",
		},
		{
			Type: "typescript",
			Link: "git@github.com:websublime/sublime-typescript-template.git",
		},
	}

	cmd := &CreateCommand{
		Templates: templates,
	}
	createCmd := NewCreateCmd(cmd)
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&cmd.Name, "name", "", "Lib or package name [REQUIRED]")
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().StringVar(&cmd.Type, "type", "", "Type of package (lib or pkg) [REQUIRED]")
	createCmd.MarkFlagRequired("type")

	createCmd.Flags().StringVar(&cmd.Template, "template", "lit", "Kind of template: (lit, solid, vue, typescript)")
	createCmd.Flags().StringVar(&cmd.Description, "description", "", "Description about your package")
}

func NewCreateCmd(cmdCreate *CreateCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create libs or packages from lit, solid, vue or typescript",
		PreRun: func(cmd *cobra.Command, _ []string) {
			sublime := core.GetSublime()

			supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, sublime.Author.Token, "production")

			isUserOrganization, err := supabase.ValidateUserOrganization(sublime.Organization)
			if err != nil {
				cmdCreate.ErrorOut(err, "Invalid user")
			}

			if !isUserOrganization {
				cmdCreate.ErrorOut(errors.New("Invalid organization user"), "User is not valid to this organization")
			}

			response, err := supabase.FindUserWorkspace(sublime.Name)
			if err != nil {
				cmdCreate.ErrorOut(err, "Workspace not found in this organization")
			}

			workspaces := []core.Workspace{}

			err = json.Unmarshal([]byte(response), &workspaces)
			if err != nil {
				cmdCreate.ErrorOut(err, "Unable to parse workspace")
			}

			var hasWorkspace bool = false
			for i := range workspaces {
				if workspaces[i].Name == sublime.Name {
					hasWorkspace = true
					cmdCreate.WorkspaceID = workspaces[i].ID
					break
				}
			}

			if !hasWorkspace {
				cmdCreate.ErrorOut(errors.New("Invalid workspace"), "Workspace is not valid to this organization")
			}
		},
		Run: func(cmd *cobra.Command, _ []string) {
			cmdCreate.Run(cmd)
			cmdCreate.CreateCloudPackage()
			cmdCreate.UpdateRepoFiles()
			cmdCreate.YarnLink()
		},
	}
}

func (ctx *CreateCommand) Run(cmd *cobra.Command) {
	color.Info.Println("🛢 Init package creation")

	sublime := core.GetSublime()
	var libType = "libs"

	if ctx.Type == "pkg" {
		libType = "packages"
	}

	ctx.PackageDir = filepath.Join(sublime.Root, libType, slug.Make(ctx.Name))

	scope := fmt.Sprintf("@%s", sublime.Organization)
	libNamespace := strings.Join([]string{scope, slug.Make(ctx.Name)}, "/")
	viteRel, _ := filepath.Rel(ctx.PackageDir, filepath.Join(sublime.Root, "libs/vite"))

	var template = ""
	var link = ""
	for i := range ctx.Templates {
		if ctx.Templates[i].Type == ctx.Template {
			template = ctx.Templates[i].Type
			link = ctx.Templates[i].Link
			break
		}
	}

	if link == "" {
		ctx.ErrorOut(errors.New("Template error"), "Unable to determine template. Valid types are: lit, solid, vue or typescript")
	}

	gitCmd := exec.Command("git", "clone", link, ctx.PackageDir)
	_, err := gitCmd.Output()
	if err != nil {
		ctx.ErrorOut(err, fmt.Sprintf("Unable to clone: %s template type ", template))
	}

	color.Success.Println("🛢 Template: ", template, "cloned. Initializing config files")

	var libPackageJson = "templates/lib-package.json"
	var libTsconfigJson = "templates/tsconfig-lib.json"
	var libViteConfigJson = "templates/vite-config-lit.json"

	if ctx.Template == "solid" {
		libPackageJson = "templates/lib-package-solid.json"
		libTsconfigJson = "templates/tsconfig-lib-solid.json"
		libViteConfigJson = "templates/vite-config-solid.json"
	}

	if ctx.Template == "vue" {
		libPackageJson = "templates/lib-package-vue.json"
		libTsconfigJson = "templates/tsconfig-lib-vue.json"
		libViteConfigJson = "templates/vite-config-vue.json"
	}

	if ctx.Template == "typescript" {
		libViteConfigJson = "templates/vite-config-typescript.json"
	}

	packageJson, _ := FileTemplates.ReadFile(libPackageJson)
	apiExtractorJson, _ := FileTemplates.ReadFile("templates/api-extractor-lib.json")
	tsConfigJson, _ := FileTemplates.ReadFile(libTsconfigJson)
	viteConfigJson, _ := FileTemplates.ReadFile(libViteConfigJson)

	pkgJsonFile, _ := os.Create(filepath.Join(ctx.PackageDir, "package.json"))
	pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &utils.PackageJsonVars{
		Namespace: libNamespace,
		Repo:      sublime.Repo,
		Name:      slug.Make(ctx.Name),
		Scope:     scope,
		Type:      libType,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Package json created and configured!")

	apiExtractorFile, _ := os.Create(filepath.Join(ctx.PackageDir, "api-extractor.json"))
	apiExtractorFile.WriteString(utils.ProcessString(string(apiExtractorJson), &utils.ApiExtractorJsonVars{
		Name: slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Api extractor created and configured!")

	tsConfigFile, _ := os.Create(filepath.Join(ctx.PackageDir, "tsconfig.json"))
	tsConfigFile.WriteString(utils.ProcessString(string(tsConfigJson), &utils.TsConfigJsonVars{
		Namespace: libNamespace,
		Vite:      viteRel,
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Tsconfig created and configured!")

	viteConfigFile, _ := os.Create(filepath.Join(ctx.PackageDir, "vite.config.js"))
	viteConfigFile.WriteString(utils.ProcessString(string(viteConfigJson), &utils.ViteJsonVars{
		Scope: scope,
		Name:  slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Success.Println("❤️‍🔥 Vite config created and configured!")

	os.RemoveAll(filepath.Join(ctx.PackageDir, ".git"))
}

func (ctx *CreateCommand) UpdateRepoFiles() {
	color.Info.Println("❤️‍🔥 Updating workspace files")

	sublime := core.GetSublime()
	scope := fmt.Sprintf("@%s", sublime.Organization)
	var libType = "libs"

	if ctx.Type == "pkg" {
		libType = "packages"
	}

	update := map[string]interface{}{
		"namespace":    sublime.Namespace,
		"name":         sublime.Name,
		"repo":         sublime.Repo,
		"root":         "./",
		"organization": sublime.Organization,
		"id":           sublime.ID,
		"description":  sublime.Description,
		"packages": append(sublime.Packages, core.Packages{
			Name:        slug.Make(ctx.Name),
			Scope:       scope,
			Type:        core.PackageType(ctx.Type),
			Description: ctx.Description,
		}),
	}

	data, _ := json.MarshalIndent(update, "", " ")

	os.WriteFile(filepath.Join(sublime.Root, ".sublime.json"), data, 0644)

	color.Success.Println("❤️‍🔥 Sublime json updated!")

	tsConfigBase := sublime.GetTsconfig()

	tsConfigBase.References = append(tsConfigBase.References, core.TsConfigReferences{
		Path: filepath.Join("./", libType, slug.Make(ctx.Name)),
		Name: filepath.Join(scope, slug.Make(ctx.Name)),
	})

	tsconfig, _ := json.MarshalIndent(tsConfigBase, "", " ")
	os.WriteFile(filepath.Join(sublime.Root, "tsconfig.base.json"), tsconfig, 0644)

	color.Success.Println("❤️‍🔥 Tsconfig base updated!")
}

func (ctx *CreateCommand) YarnLink() {
	color.Info.Println("❤️‍🔥 Init yarn link on workspace")

	workspaceDir := core.GetSublime().Root

	_, err := utils.YarnInstall(workspaceDir)
	if err != nil {
		ctx.ErrorOut(err, fmt.Sprintf("Yarn wasn't installed on: %s", workspaceDir))
	}

	color.Success.Println("✅ Package created and ready to be used")
}

func (ctx *CreateCommand) CreateCloudPackage() {
	color.Info.Println("❤️‍🔥 Creating package on cloud platform")
	sublime := core.GetSublime()

	supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, sublime.Author.Token, "production")
	_, err := supabase.CreateWorkspacePackage(ctx.Name, ctx.Type, ctx.Description, ctx.WorkspaceID)
	if err != nil {
		ctx.ErrorOut(err, "Unable to create package on cloud. Try again.")
	}

	color.Success.Println("❤️‍🔥 Package created on cloud")
}

func (ctx *CreateCommand) ErrorOut(err error, msg string) {
	if ctx.PackageDir != "" {
		os.RemoveAll(ctx.PackageDir)
	}

	color.Error.Println(msg, err)
	cobra.CheckErr(err)
}
