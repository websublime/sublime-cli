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

	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/core"
)

type RegisterFlags struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	HomeDir  string `json:"-"`
}

func init() {
	registerFlags := &RegisterFlags{}
	registerCmd := NewRegisterCommand(registerFlags)

	rootCommand.AddCommand(registerCmd)

	registerCmd.Flags().StringVar(&registerFlags.Name, "name", "", "Workspace folder name [REQUIRED]")
	registerCmd.MarkFlagRequired("name")

	registerCmd.Flags().StringVar(&registerFlags.Username, "username", "", "Git username [REQUIRED]")
	registerCmd.MarkFlagRequired("username")

	registerCmd.Flags().StringVar(&registerFlags.Email, "email", "", "Git email [REQUIRED]")
	registerCmd.MarkFlagRequired("email")

	registerCmd.Flags().StringVar(&registerFlags.Password, "password", "", "User password [REQUIRED]")
	registerCmd.MarkFlagRequired("password")
}

func NewRegisterCommand(cmdReg *RegisterFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register author on cloud platform",
		Long:  "Register command will register you as author on cloud platform and also create you local rc file to use on future workspaces and commands",
		Run: func(cmd *cobra.Command, _ []string) {
			cmdReg.Run(cmd)
		},
	}
}

func (ctx *RegisterFlags) Run(cmd *cobra.Command) {
	config := core.GetConfig()
	config.Tracker.UpdateMessage("Start registration process")
	config.Tracker.Increment(10)

	ctx.HomeDir = filepath.Join(config.HomeDir, ".sublime")
	config.Tracker.MarkAsDone()
}
