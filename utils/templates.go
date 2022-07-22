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
package utils

import (
	"bytes"
	"html/template"
)

type EmptyVars struct{}

type ArtifactsVars struct {
	Version string
}

type PackageJsonVars struct {
	Name      string
	Namespace string
	Repo      string
	Username  string
	Email     string
	Version   string
	Scope     string
	Type      string
}

type RcJsonVars struct {
	Name     string
	Username string
	Email    string
	Token    string
	ID       string
	Expire   int64
	Refresh  string
}

type ViteJsonVars struct {
	Namespace string
	Scope     string
	Name      string
}

type SublimeJsonVars struct {
	Name         string
	Repo         string
	Namespace    string
	Root         string
	Organization string
	ID           string
	Description  string
}

type ReadmeVars struct {
	Name         string
	Repo         string
	Organization string
}

type ReleaseYamlVars struct {
	Username string
	Email    string
	Scope    string
}

type ApiExtractorJsonVars struct {
	Name string
}

type TsConfigJsonVars struct {
	Namespace string
	Vite      string
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

func ProcessString(str string, vars interface{}, delimLeft string, delimRight string) string {
	tmpl, err := template.New("tmpl").Delims(delimLeft, delimRight).Parse(str)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}
