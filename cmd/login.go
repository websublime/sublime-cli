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
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/api"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

type LoginFlags struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func init() {
	loginFlags := &LoginFlags{}
	loginCmd := NewLoginCmd(loginFlags)

	//loginCmd.Flags().StringVar(&loginFlags.Email, "email", "", "Email")
	//loginCmd.MarkFlagRequired("email")

	//loginCmd.Flags().StringVar(&loginFlags.Password, "password", "", "Password")
	//loginCmd.MarkFlagRequired("password")

	rootCommand.AddCommand(loginCmd)
}

func NewLoginCmd(cmdLogin *LoginFlags) *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandLogin,
		Short: utils.MessageCommandLoginShort,
		Long:  utils.MessageCommandLoginLong,
		Run: func(cmd *cobra.Command, _ []string) {
			cmdLogin.Run(cmd)
			cmdLogin.LoginAuthor()
		},
	}
}

func (ctx *LoginFlags) Run(cmd *cobra.Command) {
	emailContent := models.PromptContent{
		Error: utils.MessageErrorCommandLoginEmailPrompt,
		Label: utils.MessageCommandLoginEmailPrompt,
		Hide:  false,
	}

	passwordContent := models.PromptContent{
		Error: utils.MessageErrorCommandLoginPasswordPrompt,
		Label: utils.MessageCommandLoginPasswordPrompt,
		Mask:  '*',
	}

	email, err := models.PromptGetInput(emailContent, 3)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}
	password, err := models.PromptGetInput(passwordContent, 8)
	if err != nil {
		utils.ErrorOut(err.Error(), utils.ErrorPromptInvalid)
	}

	ctx.Email = email
	ctx.Password = password
}

func (ctx *LoginFlags) LoginAuthor() {
	config := core.GetConfig()
	app := core.GetApp()

	go config.Progress.Render()

	config.UpdateProgress(utils.MessageCommandLoginProgressInit, 2)
	supabase := api.NewSupabase(utils.ApiUrl, utils.ApiKey, utils.ApiKey, "production")
	author, err := supabase.Login(ctx.Email, ctx.Password)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidAuthor)
	}

	config.UpdateProgress(utils.MessageCommandLoginAuthor, 2)

	tokenStrings := strings.Split(author.Token, ".")
	claimString, _ := base64.StdEncoding.DecodeString(tokenStrings[1])
	regex := regexp.MustCompile(`(?m)(exp*.":)(\d*)`)

	parts := regex.FindStringSubmatch(string(claimString))
	expires, _ := strconv.ParseInt(parts[2], 10, 64)

	author.Expires = expires

	user, err := supabase.GetUser(author.Token)
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidAuthor)
	}

	config.UpdateProgress(utils.MessageCommandLoginAuthor, 2)
	err = app.UpdateAuthorMetadata(&models.AuthorFileProps{
		Expire:   author.Expires,
		Token:    author.Token,
		Refresh:  author.RefreshToken,
		Name:     user.UserMetadata.Name,
		Username: user.UserMetadata.Author,
		Email:    user.Email,
		ID:       user.ID,
	})
	if err != nil {
		ctx.CommandError(err.Error(), utils.ErrorInvalidAuthor)
	}

	config.UpdateProgress(utils.MessageCommandLoginSuccess, 2)
	config.TerminateProgress()
}

func (ctx *LoginFlags) CommandError(message string, errorType utils.ErrorType) {
	config := core.GetConfig()

	config.TerminateErrorProgress(fmt.Sprintf("Error: %s", errorType))
	utils.ErrorOut(message, errorType)
}
