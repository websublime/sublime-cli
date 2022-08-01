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

type ActionFlags struct {
	Type        utils.GitType `json:"type"`
	Environment utils.EnvType `json:"environment"`
}

func init() {
	actionFlags := &ActionFlags{}
	actionCmd := NewActionCmd(actionFlags)

	rootCommand.AddCommand(actionCmd)

	//actionCmd.Flags().StringVar(&actionFlags.Type, "type", "branch", "Kind of action (branch or tag)")
	//actionCmd.Flags().StringVar(&actionFlags.Environment, "env", "develop", "Environment")
}

func NewActionCmd(cmdAction *ActionFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandAction,
		Short: utils.MessageCommandCreateShort,
		Long:  utils.MessageCommandCreateLong,
		Run: func(cmd *cobra.Command, _ []string) {
			//cmdAction.Run(cmd)
		},
	}
}
