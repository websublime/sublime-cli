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
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/gookit/color"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
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
			Link: "",
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

	createCmd.Flags().StringVar(&cmd.Template, "template", "lit", "Kind of template (lit, react, solid, vue)")
	createCmd.MarkFlagRequired("template")
}

func NewCreateCmd(cmdCreate *CreateCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create libs or packages",
		Run: func(cmd *cobra.Command, args []string) {
			if cmdCreate.Type == "lib" {
				cmdCreate.Library(cmd)
			}

			if cmdCreate.Type == "pkg" {
				cmdCreate.Package((cmd))
			}
		},
	}
}

func (ctx *CreateCommand) Library(cmd *cobra.Command) {
	idx := sort.Search(len(ctx.Templates), func(index int) bool {
		return string(ctx.Templates[index].Type) >= ctx.Template
	})

	template := ctx.Templates[idx].Type
	link := ctx.Templates[idx].Link

	gitCmd := exec.Command("git", "clone", link, filepath.Join("libs", slug.Make(ctx.Name)))
	_, err := gitCmd.Output()
	if err != nil {
		color.Error.Println("Unable to clone: ", template, " template type")
		cobra.CheckErr(err)
	}
}

func (ctx *CreateCommand) Package(cmd *cobra.Command) {}
