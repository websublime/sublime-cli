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
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/websublime/sublime-cli/utils"
)

type PackageType string

const (
	Library PackageType = "lib"
	Package PackageType = "pkg"
)

type Packages struct {
	Name  string      `json:"name"`
	Scope string      `json:"scope"`
	Type  PackageType `json:"type"`
}

type Sublime struct {
	Name      string           `json:"name"`
	Scope     string           `json:"scope"`
	Repo      string           `json:"repo"`
	Namespace string           `json:"namespace"`
	Root      string           `json:"root"`
	HomeDir   string           `json:"homeDir"`
	Packages  []Packages       `json:"packages"`
	Author    utils.RcJsonVars `json:"author"`
}

type TsconfigBase struct {
	CompilerOptions TsconfigCompilerOptions `json:"compilerOptions"`
	References      []TsConfigReferences    `json:"references"`
	Exclude         []string                `json:"exclude"`
}

type TsconfigCompilerOptions struct {
	Target                  string `json:"target"`
	UseDefineForClassFields bool   `json:"useDefineForClassFields"`
	Module                  string `json:"module"`
	ModuleResolution        string `json:"moduleResolution"`
	Strict                  bool   `json:"strict"`
	SourceMap               bool   `json:"sourceMap"`
	ResolveJsonModule       bool   `json:"resolveJsonModule"`
	EsModuleInterop         bool   `json:"esModuleInterop"`
	Declaration             bool   `json:"declaration"`
	SkipLibCheck            bool   `json:"skipLibCheck"`
	Composite               bool   `json:"composite"`
	Incremental             bool   `json:"incremental"`
}

type TsConfigReferences struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

var sublime = NewSublime()

func NewSublime() *Sublime {
	dir, err := os.Getwd()
	if err != nil {
		color.Error.Println("Unable to get current directory")
		cobra.CheckErr(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		color.Error.Println("Unable to get user home directory")
		cobra.CheckErr(err)
	}

	return &Sublime{
		Root:     dir,
		HomeDir:  homeDir,
		Packages: []Packages{},
	}
}

func GetSublime() *Sublime {
	return sublime
}

func (ctx *Sublime) SetRoot(path string) {
	dir, err := os.Getwd()
	if err != nil {
		color.Error.Println("Unable to get current directory")
		cobra.CheckErr(err)
	}

	if filepath.IsAbs(path) {
		ctx.Root = path
	} else {
		ctx.Root = filepath.Join(dir, path)
	}
}

func (ctx *Sublime) GetTsconfig() *TsconfigBase {
	tsconfig := &TsconfigBase{}
	data, err := os.ReadFile(filepath.Join(ctx.Root, "tsconfig.base.json"))

	if err != nil {
		color.Error.Println("Unable to get tsconfig.base.json file")
		cobra.CheckErr(err)
	}

	json.Unmarshal(data, &tsconfig)

	return tsconfig
}

func (ctx *Sublime) SetAuthor(metadata *utils.RcJsonVars) {
	ctx.Author = *metadata
}
