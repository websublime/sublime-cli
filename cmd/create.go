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

type CreateFlags struct {
	Name        string             `json:"name"`
	Type        utils.PackageType  `json:"type"`
	Template    utils.TemplateType `json:"template"`
	Description string             `json:"description"`
}

func init() {
	createFlags := &CreateFlags{}
	createCmd := NewCreateCmd(createFlags)

	rootCommand.AddCommand(createCmd)
}

func NewCreateCmd(cmdCreate *CreateFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandCreate,
		Short: utils.MessageCommandLoginShort,
		Long:  utils.MessageCommandLoginLong,
		Run: func(cmd *cobra.Command, _ []string) {
			cmdCreate.Run(cmd)
		},
	}
}

func (ctx *CreateFlags) Run(cmd *cobra.Command) {}
