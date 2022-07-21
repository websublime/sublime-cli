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
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/core/clients"
	"github.com/websublime/sublime-cli/utils"
)

type RegisterCommand struct {
	Name     string
	Username string
	Email    string
	Password string
	HomeDir  string
}

func init() {
	cmd := &RegisterCommand{}
	registerCmd := NewRegisterCmd(cmd)
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVar(&cmd.Name, "name", "", "Workspace folder name [REQUIRED]")
	registerCmd.MarkFlagRequired("name")

	registerCmd.Flags().StringVar(&cmd.Username, "username", "", "Git username [REQUIRED]")
	registerCmd.MarkFlagRequired("username")

	registerCmd.Flags().StringVar(&cmd.Email, "email", "", "Git email [REQUIRED]")
	registerCmd.MarkFlagRequired("email")

	registerCmd.Flags().StringVar(&cmd.Password, "password", "", "User password [REQUIRED]")
	registerCmd.MarkFlagRequired("password")
}

func NewRegisterCmd(cmdReg *RegisterCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register author on cloud platform",
		Long:  "Register command will register you as author on cloud platform and also create you local rc file to use on future workspaces and commands",
		Run: func(cmd *cobra.Command, _ []string) {
			cmdReg.Run(cmd)
		},
	}
}

func (ctx *RegisterCommand) Run(cmd *cobra.Command) {
	color.Info.Println("👣 Start registration process.")
	sublime := core.GetSublime()

	supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, utils.ApiKey, "production")
	register := &core.Register{}

	response, err := supabase.Register(ctx.Email, ctx.Password, ctx.Name, ctx.Username)
	if err != nil {
		ctx.ErrorOut(err, "Error registering author")
	}

	json.Unmarshal([]byte(response), &register)

	color.Success.Println("👣 Author registered with success.")
	color.Info.Println("👣 Creating user rc file.")

	ctx.HomeDir = filepath.Join(sublime.HomeDir, ".sublime")
	if err := os.Mkdir(ctx.HomeDir, 0755); err != nil {
		ctx.HomeDir = ""
		ctx.ErrorOut(err, "Error creating rc user directory")
	}

	rcJson, err := FileTemplates.ReadFile("templates/rc-template.json")
	if err != nil {
		ctx.ErrorOut(err, "Unable to read rc template file")
	}

	rcFile, err := os.Create(filepath.Join(ctx.HomeDir, "rc.json"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to create rc author file")
	}

	_, err = rcFile.WriteString(utils.ProcessString(string(rcJson), &utils.RcJsonVars{
		Name:     register.Metadata.Name,
		Username: register.Metadata.Username,
		Email:    register.Email,
		Token:    "",
		ID:       register.Id,
	}, "{{", "}}"))
	if err != nil {
		ctx.ErrorOut(err, "Unable to write rc author file")
	}

	color.Success.Println("👣 Author data persisted with success.")
	fmt.Println()

	tabular := table.NewWriter()
	tabular.SetStyle(table.StyleBold)
	tabular.AppendHeader(table.Row{"Sublime cloud registration"})
	tabular.AppendRow(table.Row{"Author", ctx.Name})
	tabular.AppendRow(table.Row{"Author email", ctx.Email})
	tabular.AppendFooter(table.Row{"Please check your email and finalize the registration process."})
	tabular.AppendFooter(table.Row{"If you are creating a new workspace for the first time please go to the platform and create the organization."})

	fmt.Println(tabular.Render())
}

func (ctx *RegisterCommand) ErrorOut(err error, msg string) {
	if ctx.HomeDir != "" {
		os.RemoveAll(ctx.HomeDir)
	}

	color.Error.Println(msg, err)
	cobra.CheckErr(err)
}
