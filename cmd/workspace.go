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
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/utils"
)

type CreateWorkspace struct {
	Name         string `json:"name"`
	Repo         string `json:"repo"`
	Organization string `json:"organization"`
	Description  string `json:"description"`
}

func init() {
	createWorkspace := &CreateWorkspace{}
	workspaceCmd := NewWorkspaceCmd(createWorkspace)

	workspaceCmd.Flags().StringVar(&createWorkspace.Organization, "organization", "", "Github organization name [REQUIRED]")
	workspaceCmd.MarkFlagRequired("organization")

	rootCommand.AddCommand(workspaceCmd)
}

func NewWorkspaceCmd(cmdWorkspace *CreateWorkspace) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandWorkspace,
		Short: utils.MessageCommandWorkspaceShort,
		Long:  utils.MessageCommandWorkspaceLong,
		Run: func(cmd *cobra.Command, _ []string) {
			cmdWorkspace.Run(cmd)
		},
	}
}

func (ctx *CreateWorkspace) Run(cmd *cobra.Command) {}
