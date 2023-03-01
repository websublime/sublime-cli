package utils

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/pkg"
)

// error out message
func ErrorOut(message string, code pkg.ErrorType) {
	red := color.Red.Render
	cobra.CheckErr(errors.New(fmt.Sprintf("%s [TYPE: %s] - %s", red("✘"), code, message)))
}

// info out message
func InfoOut(message string) {
	blue := color.Blue.Render
	fmt.Println(fmt.Sprintf("%s %s", blue("●"), message))
}

// success out message
func SuccessOut(message string) {
	green := color.Green.Render
	fmt.Println(fmt.Sprintf("%s %s", green("✔︎"), message))
}

// warning out message
func WarningOut(message string) {
	yellow := color.Yellow.Render
	fmt.Println(fmt.Sprintf("%s %s", yellow("▲"), message))
}
