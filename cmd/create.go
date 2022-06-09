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

type CreateCommand struct {
	Name      string
	Type      string
	Template  string
	Templates []LibTemplate
}

type LibTemplate struct {
	Type string
	Link string
}

func init() {
	templates := []LibTemplate{
		{
			Type: "vue",
			Link: "",
		},
		{
			Type: "lit",
			Link: "git@github.com:websublime/sublime-lit-template.git",
		},
		{
			Type: "solid",
			Link: "",
		},
		{
			Type: "react",
			Link: "",
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
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().StringVar(&cmd.Template, "template", "lit", "Kind of template (lit) incoming support to: (react, solid, vue)")
	// createCmd.MarkFlagRequired("template")
}

func NewCreateCmd(cmdCreate *CreateCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create libs or packages from lit, solid, vue or react",
		Run: func(cmd *cobra.Command, args []string) {
			if cmdCreate.Type == "lib" {
				cmdCreate.Library(cmd)
			}

			if cmdCreate.Type == "pkg" {
				cmdCreate.Package((cmd))
			}

			cmdCreate.YarnLink()
		},
	}
}

func (ctx *CreateCommand) Library(cmd *cobra.Command) {
	sublime := core.GetSublime()

	libNamespace := strings.Join([]string{sublime.Scope, slug.Make(ctx.Name)}, "/")
	libDirectory := filepath.Join(sublime.Root, "libs", slug.Make(ctx.Name))

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
		color.Error.Println("Unable to determine template. Valid types are: lit, solid, vue or react")
		cobra.CheckErr("Template error")
	}

	gitCmd := exec.Command("git", "clone", link, libDirectory)
	_, err := gitCmd.Output()
	if err != nil {
		color.Error.Println("Unable to clone: ", template, " template type")
		cobra.CheckErr(err)
	}

	color.Info.Println("🛢 Template: ", template, "cloned. Initializing config files")

	packageJson, _ := FileTemplates.ReadFile("templates/lib-package.json")
	apiExtractorJson, _ := FileTemplates.ReadFile("templates/api-extractor-lib.json")
	tsConfigJson, _ := FileTemplates.ReadFile("templates/tsconfig-lib.json")
	viteConfigJson, _ := FileTemplates.ReadFile("templates/vite-config-lit.json")

	pkgJsonFile, _ := os.Create(filepath.Join(libDirectory, "package.json"))
	pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &utils.PackageJsonVars{
		Namespace: libNamespace,
		Repo:      sublime.Repo,
		Name:      slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Package json created and configured!")

	apiExtractorFile, _ := os.Create(filepath.Join(libDirectory, "api-extractor.json"))
	apiExtractorFile.WriteString(utils.ProcessString(string(apiExtractorJson), &utils.ApiExtractorJsonVars{
		Name: slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Api extractor created and configured!")

	tsConfigFile, _ := os.Create(filepath.Join(libDirectory, "tsconfig.json"))
	tsConfigFile.WriteString(utils.ProcessString(string(tsConfigJson), &utils.TsConfigJsonVars{
		Namespace: libNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Tsconfig created and configured!")

	viteConfigFile, _ := os.Create(filepath.Join(libDirectory, "vite.config.js"))
	viteConfigFile.WriteString(utils.ProcessString(string(viteConfigJson), &utils.ViteJsonVars{
		Scope: sublime.Scope,
		Name:  slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Vite config created and configured!")

	sublime.Packages = append(sublime.Packages, core.Packages{
		Name:  slug.Make(ctx.Name),
		Scope: sublime.Scope,
		Type:  "lib",
	})

	data, _ := json.MarshalIndent(sublime, "", " ")

	os.WriteFile(filepath.Join(sublime.Root, ".sublime.json"), data, 0644)

	color.Info.Println("❤️‍🔥 Sublime json updated!")

	tsConfigBase := sublime.GetTsconfig()

	tsConfigBase.References = append(tsConfigBase.References, core.TsConfigReferences{
		Path: filepath.Join("./libs", slug.Make(ctx.Name)),
		Name: filepath.Join(sublime.Scope, slug.Make(ctx.Name)),
	})

	tsconfig, _ := json.MarshalIndent(tsConfigBase, "", " ")
	os.WriteFile(filepath.Join(sublime.Root, "tsconfig.base.json"), tsconfig, 0644)

	color.Info.Println("❤️‍🔥 Tsconfig base updated!")

	os.RemoveAll(filepath.Join(libDirectory, ".git"))
}

func (ctx *CreateCommand) Package(cmd *cobra.Command) {
	sublime := core.GetSublime()

	libNamespace := strings.Join([]string{sublime.Scope, slug.Make(ctx.Name)}, "/")
	libDirectory := filepath.Join(sublime.Root, "packages", slug.Make(ctx.Name))

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
		color.Error.Println("Unable to determine template. Valid types are: lit, solid, vue or react")
		cobra.CheckErr("Template error")
	}

	gitCmd := exec.Command("git", "clone", link, libDirectory)
	_, err := gitCmd.Output()
	if err != nil {
		color.Error.Println("Unable to clone: ", template, " template type")
		cobra.CheckErr(err)
	}

	color.Info.Println("🛢 Template: ", template, "cloned. Initializing config files")

	packageJson, _ := FileTemplates.ReadFile("templates/lib-package.json")
	apiExtractorJson, _ := FileTemplates.ReadFile("templates/api-extractor-lib.json")
	tsConfigJson, _ := FileTemplates.ReadFile("templates/tsconfig-lib.json")
	viteConfigJson, _ := FileTemplates.ReadFile("templates/vite-config-lit.json")

	pkgJsonFile, _ := os.Create(filepath.Join(libDirectory, "package.json"))
	pkgJsonFile.WriteString(utils.ProcessString(string(packageJson), &utils.PackageJsonVars{
		Namespace: libNamespace,
		Repo:      sublime.Repo,
		Name:      slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Package json created and configured!")

	apiExtractorFile, _ := os.Create(filepath.Join(libDirectory, "api-extractor.json"))
	apiExtractorFile.WriteString(utils.ProcessString(string(apiExtractorJson), &utils.ApiExtractorJsonVars{
		Name: slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Api extractor created and configured!")

	tsConfigFile, _ := os.Create(filepath.Join(libDirectory, "tsconfig.json"))
	tsConfigFile.WriteString(utils.ProcessString(string(tsConfigJson), &utils.TsConfigJsonVars{
		Namespace: libNamespace,
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Tsconfig created and configured!")

	viteConfigFile, _ := os.Create(filepath.Join(libDirectory, "vite.config.js"))
	viteConfigFile.WriteString(utils.ProcessString(string(viteConfigJson), &utils.ViteJsonVars{
		Scope: sublime.Scope,
		Name:  slug.Make(ctx.Name),
	}, "{{", "}}"))

	color.Info.Println("❤️‍🔥 Vite config created and configured!")

	sublime.Packages = append(sublime.Packages, core.Packages{
		Name:  slug.Make(ctx.Name),
		Scope: sublime.Scope,
		Type:  "pkg",
	})

	data, _ := json.MarshalIndent(sublime, "", " ")

	os.WriteFile(filepath.Join(sublime.Root, ".sublime.json"), data, 0644)

	color.Info.Println("❤️‍🔥 Sublime json updated!")

	tsConfigBase := sublime.GetTsconfig()

	tsConfigBase.References = append(tsConfigBase.References, core.TsConfigReferences{
		Path: filepath.Join("./packages", slug.Make(ctx.Name)),
		Name: filepath.Join(sublime.Scope, slug.Make(ctx.Name)),
	})

	tsconfig, _ := json.MarshalIndent(tsConfigBase, "", " ")
	os.WriteFile(filepath.Join(sublime.Root, "tsconfig.base.json"), tsconfig, 0644)

	color.Info.Println("❤️‍🔥 Tsconfig base updated!")

	os.RemoveAll(filepath.Join(libDirectory, ".git"))
}

func (ctx *CreateCommand) YarnLink() {
	workspaceDir := core.GetSublime().Root

	_, err := utils.YarnInstall(workspaceDir)

	if err != nil {
		color.Error.Println("Yarn wasn't installed on", workspaceDir, ". Please do it manually")
		color.Error.Println("Yarn error:", err.Error())
	}

	color.Success.Println("✅ Your app is updated. Yarn performed link on packages.")
}
