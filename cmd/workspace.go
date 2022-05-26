/*
Copyright © 2022 Miguel Ramos miguel.marques.ramos@gmail.com

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
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gookit/color"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
)

type WorkSpaceCmd struct {
	Name  string
	Scope string
}

func init() {
	cmd := &WorkSpaceCmd{}
	workspaceCmd := NewWorkspaceCmd(cmd)
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.Flags().StringVar(&cmd.Name, "name", "", "Workspace folder name")
	workspaceCmd.MarkFlagRequired("name")

	workspaceCmd.Flags().StringVar(&cmd.Scope, "scope", "", "Workspace scope name")
	workspaceCmd.MarkFlagRequired("scope")
}

func NewWorkspaceCmd(cmd *WorkSpaceCmd) *cobra.Command {

	return &cobra.Command{
		Use:   "workspace",
		Short: "Print the version number of sublime",
		Long:  `All software has versions. This is Sublime's`,
		Run:   cmd.Run,
	}
}

func (ctx *WorkSpaceCmd) Run(cmd *cobra.Command, args []string) {
	color.Info.Printf("Creating new workspace: %s", ctx.Name)

	workspaceDir := filepath.Join(core.GetEnvironment().WorkspaceRoot, slug.Make(ctx.Name))
	if err := os.Mkdir(workspaceDir, 0755); err != nil {
		color.Error.Printf("Error creating workspace: %s", ctx.Name)
		os.Exit(1)
	}

	vars := make(map[string]interface{})
	vars["Name"] = "Brienne"
	vars["House"] = "Tarth"
	vars["Traits"] = []string{"Brave", "Loyal"}

	file, _ := os.Create(filepath.Join(workspaceDir, "my.txt"))
	resultA := ProcessString("{{.Name}} of house {{.House}}", vars)
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

func ProcessString(str string, vars interface{}) string {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}
