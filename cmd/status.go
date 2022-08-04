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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type StatusCommand struct{}

func init() {
	cmd := &StatusCommand{}
	statusCmd := NewStatusCmd(cmd)
	rootCommand.AddCommand(statusCmd)
}

func NewStatusCmd(cmdStatus *StatusCommand) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandStatus,
		Short: utils.MessageCommandStatusShort,
		Run: func(cmd *cobra.Command, _ []string) {
			cmdStatus.Run(cmd)
		},
	}
}

func (ctx *StatusCommand) Run(cmd *cobra.Command) {
	app := core.GetApp()
	sublime := models.SublimeViperProps{}

	err := viper.Unmarshal(&sublime)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorInvalidWorkspace)
	}

	tabular := table.NewWriter()
	tabular.SetStyle(table.StyleBold)
	tabular.AppendHeader(table.Row{fmt.Sprintf("Project: %s", sublime.Namespace)})
	tabular.AppendRow(table.Row{"ID", sublime.ID})
	tabular.AppendRow(table.Row{"Repo", fmt.Sprintf("https://github.com/%s", sublime.Repo)})
	tabular.AppendRow(table.Row{"Issues", fmt.Sprintf("https://github.com/%s/issues", sublime.Repo)})
	tabular.AppendRow(table.Row{"Docs", fmt.Sprintf("https://websublime.dev/organization/%s/%s", app.Organization, sublime.Name)})

	for _, pkg := range sublime.Packages {
		var libType = "libs"

		if pkg.Type == "pkg" {
			libType = "packages"
		}

		pkgDir := filepath.Join(sublime.Root, libType, pkg.Name)

		pkgJson, err := os.ReadFile(filepath.Join(pkgDir, "package.json"))
		if err != nil {
			utils.ErrorOut(err.Error(), utils.ErrorReadFile)
		}

		packageJson := models.PackageJson{}

		err = json.Unmarshal(pkgJson, &packageJson)
		if err != nil {
			utils.ErrorOut(err.Error(), utils.ErrorReadFile)
		}

		rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
		tabular.AppendSeparator()
		tabular.AppendRow(table.Row{"Package name", pkg.Name}, rowConfigAutoMerge)
		tabular.AppendRow(table.Row{"Package type", pkg.Type}, rowConfigAutoMerge)
		tabular.AppendRow(table.Row{"Package Version", packageJson.Version}, rowConfigAutoMerge)
		tabular.AppendRow(table.Row{"Package Description", packageJson.Description}, rowConfigAutoMerge)
	}

	tabular.AppendFooter(table.Row{fmt.Sprintf("Packages on organization: %s", sublime.Organization)})

	fmt.Println(tabular.Render())
}
