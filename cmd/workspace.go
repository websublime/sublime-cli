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
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
)

type WorkSpaceCmd struct {
	Name  string
	Scope string
	Repo  string
}

//go:embed layout/*
var layoutFiles embed.FS

func init() {
	cmd := &WorkSpaceCmd{}
	workspaceCmd := NewWorkspaceCmd(cmd)
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.Flags().StringVar(&cmd.Name, "name", "", "Workspace folder name")
	workspaceCmd.MarkFlagRequired("name")

	workspaceCmd.Flags().StringVar(&cmd.Scope, "scope", "", "Workspace scope name")
	workspaceCmd.MarkFlagRequired("scope")

	workspaceCmd.Flags().StringVar(&cmd.Repo, "repo", "", "Github repo shortcut (you/repo)")
	workspaceCmd.MarkFlagRequired("repo")
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
	color.Info.Println("Creating new workspace: ", ctx.Name)

	workspaceDir := filepath.Join(core.GetEnvironment().WorkspaceRoot, slug.Make(ctx.Name))
	if err := os.Mkdir(workspaceDir, 0755); err != nil {
		color.Error.Printf("Error creating workspace: %s", ctx.Name)
		os.Exit(1)
	}

	fs.WalkDir(layoutFiles, "layout", func(path string, info fs.DirEntry, err error) error {
		relativePath := strings.ReplaceAll(path, "layout/", "")
		workspacePath := filepath.Join(workspaceDir, relativePath)
		hasLayout := strings.Contains(relativePath, "layout")

		if info.IsDir() && !hasLayout {
			os.Mkdir(workspacePath, 0755)
			color.Success.Println("Copy: ", relativePath, "To: ", workspaceDir)
		} else if !hasLayout {
			os.Create(workspacePath)
			color.Success.Println("Copy: ", relativePath, "To: ", workspaceDir)
		}

		return nil
	})

	//editorConfigFile, _ := os.Create(filepath.Join(workspaceDir, ".editorconfig"))
	//eslintIgnoreFile, _ := os.Create(filepath.Join(workspaceDir, ".eslintignore"))
	//eslintFile, _ := os.Create(filepath.Join(workspaceDir, ".eslintrc.json"))
	//gitIgnoreFile, _ := os.Create(filepath.Join(workspaceDir, ".gitignore"))
	//prettierFile, _ := os.Create(filepath.Join(workspaceDir, ".prettierrc.json"))

	//layouts.CreateEditorFile(editorConfigFile)
	//layouts.CreateEslintIgnoreFile(eslintIgnoreFile)
	//layouts.CreateEslintFile(eslintFile)
	//layouts.CreateGitIgnoreFile(gitIgnoreFile)
	//layouts.CreatePrettierFile(prettierFile)
}
