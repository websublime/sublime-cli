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

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/core/clients"
	"github.com/websublime/sublime-cli/utils"
)

type LoginCommand struct {
	Email    string
	Password string
}

func init() {
	cmd := &LoginCommand{}
	loginCmd := NewLoginCmd(cmd)
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&cmd.Email, "email", "", "Author email")

	loginCmd.Flags().StringVar(&cmd.Password, "password", "", "Author password [REQUIRED]")
	loginCmd.MarkFlagRequired("username")
}

func NewLoginCmd(cmdLogin *LoginCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login author on cloud platform",
		Run: func(cmd *cobra.Command, _ []string) {
			cmdLogin.Run(cmd)
		},
	}
}

func (ctx *LoginCommand) Run(cmd *cobra.Command) {
	color.Info.Println("👣 Author login process.")
	sublime := core.GetSublime()

	email := ctx.Email
	if email == "" {
		email = sublime.Author.Email
	}
	login := &core.Auth{}

	supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, "production")
	response, err := supabase.Login(email, ctx.Password)
	if err != nil {
		color.Error.Println("Error login author:", err)
		cobra.CheckErr(err)
	}

	color.Success.Println("👣 Login successful.")
	color.Info.Println("👣 Updating author rc file.")

	json.Unmarshal([]byte(response), &login)

	err = sublime.UpdateAuthorMetadata(login)
	if err != nil {
		color.Error.Println(err)
		cobra.CheckErr(err)
	}

	color.Success.Println("👣 Author rc file updated.")
}
