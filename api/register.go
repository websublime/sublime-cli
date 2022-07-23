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
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/websublime/sublime-cli/models"
)

func (ctx *Supabase) RegisterAuthor(name string, username string, email string, password string) (string, error) {
	signup := models.NewSignUp(email, password, name, username)

	payload, err := json.Marshal(signup)
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("%s/%s/signup", ctx.BaseURL, AuthEndpoint)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiToken))
	req.Header.Add("apikey", ctx.ApiKey)

	response, err := ctx.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 400 {
		return "", errors.New(string(body))
	}

	return string(body), err
}
