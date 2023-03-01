package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version   string = "v0.1.9"
	BuildTime string = fmt.Sprintf("%d", time.Now().Unix())
)

func init() {
	rootCommand.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sublime",
	Long:  `All software has versions. This is Sublime's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sublime CLI", Version)
	},
}
