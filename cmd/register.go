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
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/api"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type RegisterFlags struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	HomeDir  string `json:"-"`
}

func init() {
	//config := core.GetConfig()

	//

	registerFlags := &RegisterFlags{}
	registerCmd := NewRegisterCommand(registerFlags)

	rootCommand.AddCommand(registerCmd)

}

func NewRegisterCommand(cmdReg *RegisterFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandRegister,
		Short: utils.MessageCommandRegisterShort,
		Long:  utils.MessageCommandRegisterLong,
		Run: func(cmd *cobra.Command, _ []string) {
			cmdReg.Run(cmd)
		},
	}
}

func (ctx *RegisterFlags) Run(cmd *cobra.Command) {
	config := core.GetConfig()

	nameContent := models.PromptContent{
		Error: utils.MessageErrorCommandRegisterNamePrompt,
		Label: utils.MessageCommandRegisterNamePrompt,
		Hide:  false,
	}

	userContent := models.PromptContent{
		Error: utils.MessageErrorCommandRegisterUsernamePrompt,
		Label: utils.MessageCommandRegisterUsernamePrompt,
		Hide:  false,
	}

	emailContent := models.PromptContent{
		Error: utils.MessageErrorCommandRegisterEmailPrompt,
		Label: utils.MessageCommandRegisterEmailPrompt,
		Hide:  false,
	}

	passwordContent := models.PromptContent{
		Error: utils.MessageErrorCommandRegisterPasswordPrompt,
		Label: utils.MessageCommandRegisterPasswordPrompt,
		Mask:  '*',
	}

	name, err := models.PromptGetInput(nameContent, 5)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}
	username, err := models.PromptGetInput(userContent, 5)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}
	email, err := models.PromptGetInput(emailContent, 5)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}
	password, err := models.PromptGetInput(passwordContent, 8)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	go config.Progress.Render()

	time.Sleep(time.Millisecond)
	config.Tracker.UpdateMessage("Start registration process")
	config.Tracker.Increment(10)

	supabase := api.NewSupabase("", "", "", "")
	supabase.RegisterAuthor(name, username, email, password)

	time.Sleep(time.Millisecond)
	config.Tracker.UpdateMessage("Author registered in platform. Init local config")
	config.Tracker.Increment(10)

	ctx.HomeDir = filepath.Join(config.HomeDir, ".sublime")
	config.Tracker.Increment(10)

	time.Sleep(time.Millisecond)
	config.Tracker.MarkAsDone()
}
