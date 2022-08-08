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
	"strings"

	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

var app = NewApp()

type App struct {
	Author         *models.AuthorFileProps `json:"author"`
	Organization   string                  `json:"organization"`
	OrganizationID string                  `json:"Organization_id"`
}

func NewApp() *App {
	return &App{}
}

func GetApp() *App {
	return app
}

func (ctx *App) InitAuthor() error {
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

	ctx.Author = authorMetadata

	return nil
}

func (ctx *App) UpdateAuthorMetadata(author *models.AuthorFileProps) error {
	config := GetConfig()

	sublimeDir := filepath.Join(config.HomeDir, ".sublime")
	rcFile := filepath.Join(sublimeDir, "rc.json")
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		utils.WarningOut(utils.MessageErrorAuthorFileMissing)
	}

	data, err := json.MarshalIndent(author, "", " ")
	if err != nil {
		return errors.New(utils.MessageErrorIndentFile)
	}

	if _, err := os.Stat(sublimeDir); os.IsNotExist(err) {
		if err := os.Mkdir(sublimeDir, 0755); err != nil {
			return errors.New(utils.MessageErrorCommandRegisterHomeDir)
		}
	}

	err = os.WriteFile(rcFile, data, 0644)
	if err != nil {
		return errors.New(utils.MessageErrorWriteFile)
	}

	return nil
}

func (ctx *App) UpdateWorkspace(workspace *models.Workspace) error {
	config := GetConfig()
	sublimeFile := filepath.Join(config.RootDir, workspace.Name, ".sublime.json")
	sublimeJson, err := os.ReadFile(sublimeFile)
	if err != nil {
		return errors.New(utils.MessageErrorReadFile)
	}

	sublimeMetadata := models.SublimeJsonFileProps{}

	err = json.Unmarshal(sublimeJson, &sublimeMetadata)
	if err != nil {
		return errors.New(utils.MessageErrorParseFile)
	}

	sublimeMetadata.ID = workspace.ID

	data, err := json.MarshalIndent(sublimeMetadata, "", " ")
	if err != nil {
		return errors.New(utils.MessageErrorIndentFile)
	}

	err = os.WriteFile(sublimeFile, data, 0644)
	if err != nil {
		return errors.New(utils.MessageErrorWriteFile)
	}

	return nil
}

func (ctx *App) UpdatePackage(packages *models.Package) error {
	config := GetConfig()
	sublimeFile := filepath.Join(config.RootDir, ".sublime.json")
	sublimeJson, err := os.ReadFile(sublimeFile)
	if err != nil {
		return errors.New(utils.MessageErrorAuthorFileMissing)
	}

	sublimeMetadata := models.SublimeJsonFileProps{}

	err = json.Unmarshal(sublimeJson, &sublimeMetadata)
	if err != nil {
		return errors.New(utils.MessageErrorParseFile)
	}

	for i := range sublimeMetadata.Packages {
		if sublimeMetadata.Packages[i].Name == packages.Name {
			sublimeMetadata.Packages[i].ID = packages.ID
			break
		}
	}

	data, err := json.MarshalIndent(sublimeMetadata, "", " ")
	if err != nil {
		return errors.New(utils.MessageErrorIndentFile)
	}

	err = os.WriteFile(sublimeFile, data, 0644)
	if err != nil {
		return errors.New(utils.MessageErrorWriteFile)
	}

	return nil
}

func (ctx *App) RemoveConfigurationsOnPackageError(packageName string, packageType string) error {
	config := GetConfig()

	tsFile := filepath.Join(config.RootDir, "tsconfig.base.json")
	sublimeFile := filepath.Join(config.RootDir, ".sublime.json")
	sublimeJson, err := os.ReadFile(sublimeFile)
	if err != nil {
		return errors.New(utils.MessageErrorAuthorFileMissing)
	}

	sublimeMetadata := models.SublimeJsonFileProps{}

	err = json.Unmarshal(sublimeJson, &sublimeMetadata)
	if err != nil {
		return errors.New(utils.MessageErrorParseFile)
	}

	tsconfig, err := ctx.GetTsconfig()
	if err != nil {
		return err
	}

	var sublimePackages []models.SublimePackages = []models.SublimePackages{}
	var configReferences []models.TsConfigReferences = []models.TsConfigReferences{}

	for _, pkg := range sublimeMetadata.Packages {
		if pkg.Name != packageName {
			sublimePackages = append(sublimePackages, pkg)
		}
	}

	sublimeMetadata.Packages = sublimePackages

	for _, cfg := range tsconfig.References {
		if cfg.Path != strings.Join([]string{packageType, packageName}, "/") {
			configReferences = append(configReferences, cfg)
		}
	}

	tsconfig.References = configReferences

	sublimeData, err := json.MarshalIndent(sublimeMetadata, "", " ")
	if err != nil {
		return errors.New(utils.MessageErrorIndentFile)
	}

	tsData, err := json.MarshalIndent(tsconfig, "", " ")
	if err != nil {
		return errors.New(utils.MessageErrorIndentFile)
	}

	err = os.WriteFile(sublimeFile, sublimeData, 0644)
	if err != nil {
		return errors.New(utils.MessageErrorWriteFile)
	}

	err = os.WriteFile(tsFile, tsData, 0644)
	if err != nil {
		return errors.New(utils.MessageErrorWriteFile)
	}

	return nil
}

func (ctx *App) GetTsconfig() (*models.TsconfigBase, error) {
	config := GetConfig()
	tsconfig := &models.TsconfigBase{}
	data, err := os.ReadFile(filepath.Join(config.RootDir, "tsconfig.base.json"))

	if err != nil {
		return nil, errors.New(string(utils.ErrorReadFile))
	}

	json.Unmarshal(data, &tsconfig)

	return tsconfig, nil
}
