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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/websublime/sublime-cli/api"
	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
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
				cmd.Help()
			}
		},
	}
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		utils.ErrorOut(utils.MessageErrorCommandExecution, utils.ErrorCmdExecution)
	}
}

func init() {
	rootFlags := &RootFlags{}

	cobra.OnInitialize(func() {
		banner()
		executeAuthorValidation()
		executeTokenExpirationValidation()
		initializeCommand(rootFlags)
	})

	rootCommand.PersistentFlags().StringVar(&rootFlags.ConfigFile, utils.CommandFlagConfig, "", utils.MessageCommandConfigUsage)
	rootCommand.PersistentFlags().StringVar(&rootFlags.Root, utils.CommandFlagRoot, "", utils.MessageCommandRootUsage)
}

func banner() {
	banner := `
	
███████╗██╗   ██╗██████╗ ██╗     ██╗███╗   ███╗███████╗     ██████╗██╗     ██╗
██╔════╝██║   ██║██╔══██╗██║     ██║████╗ ████║██╔════╝    ██╔════╝██║     ██║
███████╗██║   ██║██████╔╝██║     ██║██╔████╔██║█████╗      ██║     ██║     ██║
╚════██║██║   ██║██╔══██╗██║     ██║██║╚██╔╝██║██╔══╝      ██║     ██║     ██║
███████║╚██████╔╝██████╔╝███████╗██║██║ ╚═╝ ██║███████╗    ╚██████╗███████╗██║
╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝╚══════╝╚═╝
                                                                              

|-----------------------------------------------------------------------------|	
`
	color.Green.Println(banner)
	color.Blue.Println(fmt.Sprintf("Version: %s", Version))
}

func initializeCommand(rootFlags *RootFlags) {
	config := core.GetConfig()

	if rootFlags.Root != "" {
		config.SetRootDir(rootFlags.Root)
	}

	if rootFlags.ConfigFile != "" {
		viper.SetConfigFile(rootFlags.ConfigFile)
	} else {
		viper.AddConfigPath(config.RootDir)
		viper.SetConfigType("json")
		viper.SetConfigName(".sublime")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// configFile := viper.ConfigFileUsed()
	}
}

func executeAuthorValidation() {
	flags := os.Args[1:]

	if !isCommandExclude(flags) {
		app := core.GetApp()

		err := app.InitAuthor()
		if err != nil {
			utils.ErrorOut(err.Error(), utils.ErrorInvalidAuthor)
		}

		if app.Author.Token == "" {
			utils.ErrorOut(utils.MessageErrorAuthorTokenMissing, utils.ErrorInvalidAuthor)
		}
	}
}

func executeTokenExpirationValidation() {
	flags := os.Args[1:]

	if !isCommandExclude(flags) {
		app := core.GetApp()

		now := time.Now()
		expiration := time.Unix(app.Author.Expire, 0)

		if now.After(expiration) {
			supabase := api.NewSupabase(utils.ApiUrl, utils.ApiKey, utils.ApiKey, "production")
			refresh, err := supabase.RefreshToken(app.Author.Refresh)
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidToken)
			}

			tokenStrings := strings.Split(refresh.Token, ".")
			claimString, _ := base64.StdEncoding.DecodeString(tokenStrings[1])
			regex := regexp.MustCompile(`(?m)(exp*.":)(\d*)`)

			parts := regex.FindStringSubmatch(string(claimString))
			expires, _ := strconv.ParseInt(parts[2], 10, 64)

			refresh.Expires = expires

			err = app.UpdateAuthorMetadata(&models.AuthorFileProps{
				Expire:  refresh.Expires,
				Token:   refresh.Token,
				Refresh: refresh.RefreshToken,
			})
			if err != nil {
				utils.ErrorOut(err.Error(), utils.ErrorInvalidAuthor)
			}
		}
	}
}

func isCommandExclude(flags []string) bool {
	if utils.Contains(flags, "action") || utils.Contains(flags, "login") || utils.Contains(flags, "register") {
		return true
	} else {
		return false
	}
}
