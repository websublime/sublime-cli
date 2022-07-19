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
