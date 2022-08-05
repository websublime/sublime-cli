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
package core

import (
	"os"
	"path/filepath"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/websublime/sublime-cli/utils"
)

var config = NewConfig()

type Config struct {
	RootDir  string            `json:"root,omitempty"`
	HomeDir  string            `json:"home,omitempty"`
	Progress progress.Writer   `json:"-"`
	Tracker  *progress.Tracker `json:"-"`
}

func NewConfig() *Config {
	dir, err := os.Getwd()
	if err != nil {
		utils.ErrorOut(utils.MessageErrorCurrentDirectory, utils.ErrorMissingDirectory)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		utils.ErrorOut(utils.MessageErrorHomeDirectory, utils.ErrorMissingDirectory)
	}

	pg := progress.NewWriter()
	pg.SetAutoStop(false)
	pg.SetTrackerLength(25)
	pg.SetMessageWidth(50)
	pg.SetNumTrackersExpected(1)
	pg.SetTrackerPosition(progress.PositionRight)
	pg.SetStyle(progress.StyleBlocks)
	pg.SetUpdateFrequency(time.Nanosecond)
	pg.Style().Colors = progress.StyleColorsExample
	pg.Style().Visibility.Time = false

	units := &progress.UnitsDefault
	tracker := progress.Tracker{
		Total:   100,
		Message: "",
		Units:   *units,
	}

	pg.AppendTracker(&tracker)

	return &Config{
		RootDir:  dir,
		HomeDir:  homeDir,
		Progress: pg,
		Tracker:  &tracker,
	}
}

func GetConfig() *Config {
	return config
}

func (ctx *Config) SetRootDir(path string) {
	dir, err := os.Getwd()
	if err != nil {
		utils.ErrorOut(utils.MessageErrorCurrentDirectory, utils.ErrorMissingDirectory)
	}

	if filepath.IsAbs(path) {
		ctx.RootDir = path
	} else {
		ctx.RootDir = filepath.Join(dir, path)
	}
}

func (ctx *Config) UpdateProgress(message string, increment int64) {
	ctx.Tracker.UpdateMessage(message)

	for idx := int64(1); idx <= increment; idx++ {
		ctx.Tracker.Increment(idx * 10)
	}

	time.Sleep(time.Microsecond * 100)
}

func (ctx *Config) TerminateProgress() {
	ctx.Tracker.MarkAsDone()
	ctx.Progress.Stop()
}

func (ctx *Config) DoneProgress() {
	ctx.Tracker.MarkAsDone()
}

func (ctx *Config) AddTracker() {
	units := &progress.UnitsDefault
	tracker := progress.Tracker{
		Total:   100,
		Message: "",
		Units:   *units,
	}

	ctx.Progress.AppendTracker(&tracker)

	ctx.Tracker = &tracker
}

func (ctx *Config) TerminateErrorProgress(message string) {
	ctx.Tracker.UpdateMessage(message)
	ctx.Tracker.MarkAsErrored()
	time.Sleep(time.Microsecond * 100)
}
