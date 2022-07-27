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
package models

type PackageJsonFileProps struct {
	Name      string
	Namespace string
	Repo      string
	Username  string
	Email     string
	Version   string
	Scope     string
	Type      string
}

type ViteJsonFileProps struct {
	Namespace string
	Scope     string
	Name      string
}

type TsConfigJsonFileProps struct {
	Namespace string
	Vite      string
}

type SublimeJsonFileProps struct {
	Name         string `json:"name"`
	Repo         string `json:"repo"`
	Namespace    string `json:"namespace"`
	Root         string `json:"root"`
	Organization string `json:"organization"`
	ID           string `json:"id"`
	Description  string `json:"description"`
}

type ReadmeFileProps struct {
	Name         string
	Repo         string
	Organization string
}

type ReleaseYamlFileProps struct {
	Username string
	Email    string
	Scope    string
}

type ArtifactsYamlFileProps struct {
	Version string
}

type Workspace struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Repo           string `json:"repo,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Private        bool   `json:"private,omitempty"`
}

func NewWorkspace(name string, repo string, description string, orgID string) *Workspace {
	return &Workspace{
		Name:           name,
		Repo:           repo,
		Description:    description,
		OrganizationID: orgID,
		Private:        false,
	}
}
