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

type SignData struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type Sign struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Data     SignData `json:"data"`
}

type SignResponse struct {
	ID             string `json:"id"`
	Aud            string `json:"aud"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ConfirmationAt string `json:"confirmation_sent_at"`
	AppMetadata    struct {
		Provider  string   `json:"provider"`
		Providers []string `json:"providers"`
	} `json:"app_metadata"`
	UserMetadata struct {
		Author string `json:"author"`
		Name   string `json:"name"`
	} `json:"user_metadata"`
	Identities []struct {
		ID           string `json:"id"`
		UserID       string `json:"user_id"`
		IdentityData struct {
			Sub string `json:"sub"`
		} `json:"identity_data"`
		Provider   string `json:"provider"`
		LastSignAt string `json:"last_sign_in_at"`
		CreatedAt  string `json:"created_at"`
	} `json:"identities"`
}

func NewSignUp(email string, password string, name string, username string) *Sign {
	return &Sign{
		Email:    email,
		Password: password,
		Data: SignData{
			Name:   name,
			Author: username,
		},
	}
}
