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

	"github.com/websublime/sublime-cli/core"
	"github.com/websublime/sublime-cli/models"
)

func (ctx *Supabase) CreateOrganizationWorkspace(name string, repo string, description string, orgID string) ([]models.Workspace, error) {
	workspace := models.NewWorkspace(name, repo, description, orgID)
	model := []models.Workspace{}

	payload, err := json.Marshal(workspace)
	if err != nil {
		return model, err
	}

	uri := fmt.Sprintf("%s/%s/workspaces", ctx.BaseURL, RestEndpoint)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		return model, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiToken))
	req.Header.Add("apikey", ctx.ApiKey)
	req.Header.Add("Prefer", "return=representation")

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

func (ctx *Supabase) GetOrganizationByUser(userID string) ([]models.OrganizationByUserResponse, error) {
	uri := fmt.Sprintf("%s/%s/organization_users?user_id=eq.%s&select=organization(id,name,created_at)", ctx.BaseURL, RestEndpoint, userID)
	model := []models.OrganizationByUserResponse{}

	req, err := http.NewRequest("GET", uri, nil)
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

func (ctx *Supabase) ValidateUserOrganization(organization string) (bool, error) {
	var isUserOrganization bool = false

	app := core.GetApp()

	organizations, err := ctx.GetOrganizationByUser(app.Author.ID)
	if err != nil {
		return isUserOrganization, err
	}

	for _, org := range organizations {
		if org.Organization.Name == organization {
			isUserOrganization = true
			app.OrganizationID = org.Organization.ID
			app.Organization = org.Organization.Name
			break
		}
	}

	return isUserOrganization, nil
}
