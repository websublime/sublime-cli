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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/core/clients"
	"github.com/websublime/sublime-cli/utils"
)

type RootCommand struct {
	ConfigFile string
	Root       string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCmd()

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sublime",
		Short: "CLI tool to manage FE packages",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCommand := &RootCommand{}

	cobra.OnInitialize(func() {
		executionFlagsValidation(rootCmd)
		executionExpirationValidation(rootCmd)
		initConfig(rootCommand)
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&rootCommand.ConfigFile, "config", "", "config file (default is .sublime.json)")
	rootCmd.PersistentFlags().StringVar(&rootCommand.Root, "root", "", "Project working dir, default to current dir")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executionFlagsValidation(_ *cobra.Command) {
	flags := os.Args[1:]

	if !(utils.Contains(flags, "action") || utils.Contains(flags, "login") || utils.Contains(flags, "register")) {
		sublime := core.GetSublime()

		rcFile := filepath.Join(sublime.HomeDir, ".sublime/rc.json")
		rcJson, err := os.ReadFile(rcFile)
		if err != nil {
			color.Error.Println("Authentication file not found. Please register first then login to cloud service.")
			cobra.CheckErr(err)
		}

		authorMetadata := &utils.RcJsonVars{}

		err = json.Unmarshal(rcJson, &authorMetadata)
		if err != nil {
			color.Error.Println("Unable to parse author rc file", err)
			cobra.CheckErr(err)
		}

		if authorMetadata.Token == "" {
			color.Error.Println("You are not loggedin. Please login using the cli login command")
			cobra.CheckErr(errors.New("empty token provided"))
		}

		sublime.SetAuthor(authorMetadata)
	}
}

func executionExpirationValidation(_ *cobra.Command) {
	flags := os.Args[1:]

	if !(utils.Contains(flags, "action") || utils.Contains(flags, "login") || utils.Contains(flags, "register")) {
		sublime := core.GetSublime()

		now := time.Now()
		expiration := time.Unix(sublime.Author.Expire, 0)

		if now.After(expiration) {
			color.Warn.Println("⛑ Your token is expired. Starting refresh token")

			supabase := clients.NewSupabase(utils.ApiUrl, utils.ApiKey, sublime.Author.Token, "production")
			response, err := supabase.RefreshToken(sublime.Author.Refresh)
			if err != nil {
				color.Error.Println("Unable to refresh token. Please make login command", err)
				cobra.CheckErr(err)
			}

			login := &core.Auth{}

			json.Unmarshal([]byte(response), &login)

			err = sublime.UpdateAuthorMetadata(login)
			if err != nil {
				color.Error.Println(err)
				cobra.CheckErr(err)
			}

			color.Success.Println("👣 Author rc file updated.")
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(rootCommand *RootCommand) {
	sublime := core.GetSublime()

	if rootCommand.Root != "" {
		sublime.SetRoot(rootCommand.Root)
	}

	if rootCommand.ConfigFile != "" {
		viper.SetConfigFile(rootCommand.ConfigFile)
	} else {
		// Search config in home directory with name ".sublime" (without extension).
		viper.AddConfigPath(sublime.Root)
		viper.SetConfigType("json")
		viper.SetConfigName(".sublime")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		configFile := viper.ConfigFileUsed()

		if !filepath.IsAbs(configFile) {
			configFile = filepath.Join(sublime.Root, configFile)
		}

		data, _ := os.ReadFile(configFile)
		json.Unmarshal(data, &sublime)

		sublime.Root = filepath.Dir(configFile)

		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
