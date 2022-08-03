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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type ActionFlags struct {
	Type        string                   `json:"type"`
	Environment string                   `json:"environment"`
	Sublime     models.SublimeViperProps `json:"-"`
}

func init() {
	actionFlags := &ActionFlags{
		Sublime: models.SublimeViperProps{},
	}
	actionCmd := NewActionCmd(actionFlags)

	rootCommand.AddCommand(actionCmd)

	actionCmd.Flags().StringVar(&actionFlags.Type, utils.CommandFlagActionType, "branch", "Type of action (branch or tag)")
	actionCmd.Flags().StringVar(&actionFlags.Environment, utils.CommandFlagActionEnv, "develop", "Environment")
}

func NewActionCmd(cmdAction *ActionFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandAction,
		Short: utils.MessageCommandActionShort,
		Long:  utils.MessageCommandActionLong,
		PreRun: func(cmd *cobra.Command, _ []string) {
			isCiEnv := viper.GetBool("CI")

			err := viper.Unmarshal(&cmdAction.Sublime)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidEnvironment)
			}

			if !isCiEnv {
				utils.ErrorOut(utils.MessageErrorCommandActionEnv, utils.ErrorInvalidEnvironment)
			}
		},
		Run: func(cmd *cobra.Command, _ []string) {
			cmdAction.Run(cmd)
		},
	}
}

func getBranchDiffPackages(dir string, pkgs []models.SublimePackages) []models.SublimePackages {
	packages := []models.SublimePackages{}

	changedList, _ := utils.GetBranchList(dir)

	for _, pkg := range pkgs {
		founded := strings.Contains(changedList, pkg.Name)

		if founded {
			pkgs = append(pkgs, pkg)
		}
	}

	return packages
}

func (ctx *ActionFlags) Run(cmd *cobra.Command) {
	config := core.GetConfig()

	//types := utils.GitType(ctx.Type)
	//env := utils.EnvType(ctx.Environment)
	packages := getBranchDiffPackages(config.RootDir, ctx.Sublime.Packages)

	fmt.Println(packages)
}
