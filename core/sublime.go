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
package core

import (
	"os"
	"path/filepath"
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
	Name      string     `json:"name"`
	Scope     string     `json:"scope"`
	Repo      string     `json:"repo"`
	Namespace string     `json:"namespace"`
	Root      string     `json:"root"`
	Packages  []Packages `json:"packages"`
}

var sublime = NewSublime()

func NewSublime() *Sublime {
	dir, _ := os.Getwd()

	return &Sublime{
		Root:     dir,
		Packages: []Packages{},
	}
}

func GetSublime() *Sublime {
	return sublime
}

func (ctx *Sublime) SetRoot(path string) {
	dir, _ := os.Getwd()

	if filepath.IsAbs(path) {
		ctx.Root = path
	} else {
		ctx.Root = filepath.Join(dir, path)
	}
}
