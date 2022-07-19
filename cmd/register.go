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
	Name         string
	Organization string
	Username     string
	Email        string
	Password     string
}

type UserMetadata struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Register struct {
	Id           string       `json:"id"`
	Aud          string       `json:"aud"`
	Role         string       `json:"role"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Confirmation string       `json:"confirmation_sent_at"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	Metadata     UserMetadata `json:"user_metadata"`
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

	registerCmd.Flags().StringVar(&cmd.Organization, "organization", "", "Github organization name [REQUIRED]")
	registerCmd.MarkFlagRequired("organization")
}

func NewRegisterCmd(cmdReg *RegisterCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register author on cloud platform",
		Run: func(cmd *cobra.Command, args []string) {
			cmdReg.Run(cmd)
		},
	}
}

func (ctx *RegisterCommand) Run(cmd *cobra.Command) {
	color.Info.Println("👣 Start registration process.")

	supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, "production")
	register := &Register{}

	response, err := supabase.Register(ctx.Email, ctx.Password, ctx.Name, ctx.Username)
	if err != nil {
		color.Error.Println("Error registering author:", err)
		cobra.CheckErr(err)
	}

	json.Unmarshal([]byte(response), &register)

	color.Success.Println("👣 Author registered with success.")
	color.Info.Println("👣 Creating user rc file.")

	userDir := filepath.Join(core.GetSublime().HomeDir, ".sublime")
	if err := os.Mkdir(userDir, 0755); err != nil {
		color.Error.Println("Error creating rc user directory", err)
		cobra.CheckErr(err)
	}

	rcJson, err := FileTemplates.ReadFile("templates/rc-template.json")
	if err != nil {
		color.Error.Println("Unable to read rc template file:", err)
		cobra.CheckErr(err)
	}

	rcFile, err := os.Create(filepath.Join(userDir, "rc.json"))
	if err != nil {
		color.Error.Println("Unable to create rc author file:", err)
		cobra.CheckErr(err)
	}

	rcFile.WriteString(utils.ProcessString(string(rcJson), &utils.RcJsonVars{
		Name:         register.Metadata.Name,
		Username:     register.Metadata.Username,
		Email:        register.Email,
		Organization: ctx.Organization,
		Token:        "",
		ID:           register.Id,
		Refresh:      "",
	}, "{{", "}}"))

	color.Success.Println("👣 Author data persisted with success.")
	fmt.Println()

	tabular := table.NewWriter()
	tabular.SetStyle(table.StyleBold)
	tabular.AppendHeader(table.Row{"Sublime cloud registration"})
	tabular.AppendRow(table.Row{"Author", ctx.Name})
	tabular.AppendRow(table.Row{"Author email", ctx.Email})
	tabular.AppendRow(table.Row{"Author organization", ctx.Organization})
	tabular.AppendFooter(table.Row{"Please check your email and finalize the registration process."})

	fmt.Println(tabular.Render())
}
