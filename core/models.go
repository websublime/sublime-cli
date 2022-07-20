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

type UserMetadata struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Register struct {
	Id           string       `json:"id"`
	Aud          string       `json:"aud"`
	Role         string       `json:"role"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Confirmation string       `json:"confirmation_sent_at"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	Metadata     UserMetadata `json:"user_metadata"`
}

type Auth struct {
	Token        string `json:"access_token"`
	Type         string `json:"token_type"`
	Expires      int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpData struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type SignUp struct {
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Data     SignUpData `json:"data,omitempty"`
}

type Refresh struct {
	Token string `json:"refresh_token"`
}

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
