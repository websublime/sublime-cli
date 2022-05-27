/*
Copyright © 2022 Miguel Ramos miguel.marques.ramos@gmail.com

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
)

type CreateCommand struct {
	Name string
	Type string
}

func init() {
	cmd := &CreateCommand{}
	createCmd := NewCreateCmd(cmd)
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&cmd.Name, "name", "", "Lib or package name [REQUIRED]")
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().StringVar(&cmd.Type, "type", "", "Type of package (lib or pkg) [REQUIRED]")
	createCmd.MarkFlagRequired("name")
}

func NewCreateCmd(cmdCreate *CreateCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create libs or packages",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}

func (ctx *CreateCommand) Library(cmd *cobra.Command) {}

func (ctx *CreateCommand) Package(cmd *cobra.Command) {}
