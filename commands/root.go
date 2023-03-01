// Package commands /*
package commands

import (
	"fmt"

	"github.com/gookit/color"

	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/utils"
)

type RootFlags struct {
	ConfigFile string `json:"config_file"`
	Root       string `json:"root"`
}

// rootCmd represents the base command when called without any subcommands
var rootCommand = NewRootCommand()

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   utils.CommandRoot,
		Short: utils.MessageCommandRootShort,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				banner()
				cmd.Help()
			}
		},
	}
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		// utils.ErrorOut(utils.MessageErrorCommandExecution, utils.ErrorCmdExecution)
	}
}

func init() {
	rootFlags := &RootFlags{}

	cobra.OnInitialize(func() {})

	rootCommand.PersistentFlags().StringVar(&rootFlags.ConfigFile, utils.CommandFlagConfig, "", utils.MessageCommandConfigUsage)
	rootCommand.PersistentFlags().StringVar(&rootFlags.Root, utils.CommandFlagRoot, "", utils.MessageCommandRootUsage)
}

func banner() {
	gray := color.Gray.Render
	blue := color.Blue.Render
	green := color.Green.Render

	fmt.Printf("%s %s %s%s%s\n", gray(utils.BAR_START), green("✦ Build tools for the future"), gray(utils.BAR_H), gray(utils.BAR_H), blue("●"))
	fmt.Printf("%s\n", gray(utils.BAR))
	fmt.Printf("%s %s\n", gray("◆"), "SUBLIME CLI")
	fmt.Printf("%s Version: %s-%s\n", gray(utils.BAR), Version, BuildTime)
	fmt.Printf("%s\n", gray(utils.BAR_END))
}
