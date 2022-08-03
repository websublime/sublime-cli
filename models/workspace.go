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

import "github.com/websublime/sublime-cli/utils"

type PackageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Typings         string            `json:"typings"`
	Module          string            `json:"module"`
	Scripts         map[string]string `json:"scripts"`
	Keywords        []string          `json:"keywords"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

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

type SublimePackages struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Scope       string            `json:"scope"`
	Type        utils.PackageType `json:"type"`
	Description string            `json:"description"`
}

type SublimeJsonFileProps struct {
	Name         string            `json:"name"`
	Repo         string            `json:"repo"`
	Namespace    string            `json:"namespace"`
	Root         string            `json:"root"`
	Organization string            `json:"organization"`
	ID           string            `json:"id"`
	Description  string            `json:"description"`
	Packages     []SublimePackages `json:"packages"`
}

type SublimeViperProps struct {
	Name         string            `mapstructure:"name"`
	Repo         string            `mapstructure:"repo"`
	Namespace    string            `mapstructure:"namespace"`
	Root         string            `mapstructure:"root"`
	Organization string            `mapstructure:"organization"`
	ID           string            `mapstructure:"id"`
	Description  string            `mapstructure:"description"`
	Packages     []SublimePackages `mapstructure:"packages"`
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

type ApiExtractorFileProps struct {
	Name string
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

type Workspace struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Repo           string `json:"repo,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Private        bool   `json:"private,omitempty"`
}

type WorkspacesByOrganizationResponse struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Repo           string `json:"repo,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Private        bool   `json:"private,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
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
