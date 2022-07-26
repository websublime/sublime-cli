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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

var app = NewApp()

type App struct {
	Author models.AuthorFileProps `json:"author"`
}

func NewApp() *App {
	return &App{}
}

func GetApp() *App {
	return app
}

func (ctx *App) UpdateAuthorMetadata(author *models.AuthorFileProps) error {
	config := GetConfig()
	rcFile := filepath.Join(config.HomeDir, ".sublime/rc.json")
	rcJson, err := os.ReadFile(rcFile)
	if err != nil {
		return errors.New(utils.MessageErrorAuthorFileMissing)
	}

	authorMetadata := &models.AuthorFileProps{}

	err = json.Unmarshal(rcJson, &authorMetadata)
	if err != nil {
		return errors.New(utils.MessageErrorParseFile)
	}

	authorMetadata.Token = author.Token
	authorMetadata.Expire = author.Expire
	authorMetadata.Refresh = author.Refresh

	data, err := json.MarshalIndent(authorMetadata, "", " ")
	if err != nil {
		return errors.New(utils.MessageErrorIndentFile)
	}

	err = os.WriteFile(filepath.Join(config.HomeDir, ".sublime/rc.json"), data, 0644)
	if err != nil {
		return errors.New(utils.MessageErrorWriteFile)
	}

	return nil
}
