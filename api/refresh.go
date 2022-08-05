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

func (ctx *Supabase) RefreshToken(token string) (models.RefreshResponse, error) {
	refresh := models.NewRefresh(token)
	model := models.RefreshResponse{}

	payload, err := json.Marshal(refresh)
	if err != nil {
		return model, err
	}

	uri := fmt.Sprintf("%s/%s/token?grant_type=refresh_token", ctx.BaseURL, AuthEndpoint)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		return model, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiToken))
	req.Header.Add("apikey", ctx.ApiKey)

	response, err := ctx.HTTPClient.Do(req)
	if err != nil {
		return model, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model, err
	}

	if response.StatusCode >= 400 {
		return model, errors.New(string(body))
	}

	err = json.Unmarshal(body, &model)
	if err != nil {
		return model, err
	}

	return model, nil
}
