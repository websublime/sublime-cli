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
package clients

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Github struct {
	BaseURL     string
	ApiKey      string
	Environment string
	HTTPClient  *http.Client
}

type Payload struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Sha     string `json:"sha,omitempty"`
}

type Content struct {
	Sha string `json:"sha"`
}

func NewGithub(baseURL string, supabaseKey string, env string) *Github {
	return &Github{
		BaseURL:     baseURL,
		ApiKey:      supabaseKey,
		Environment: env,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (ctx *Github) Upload(bucket string, filePath string, destination string) string {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err != nil {
		panic(fmt.Errorf("os.Open: %v", err))
	}
	file.Close()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("os.Open: %v", err))
	}

	base64data := base64.StdEncoding.EncodeToString(data)

	currentTime := time.Now()
	uri := fmt.Sprintf("%s/%s/%s/%s", ctx.BaseURL, bucket, destination, filepath.Base(file.Name()))
	sha := ctx.GetContent(uri)

	payload := Payload{
		Message: fmt.Sprintf("Upload artifact %s at %s", filepath.Base(file.Name()), currentTime.Format("2006.01.02-15:04:05")),
		Content: base64data,
	}

	if sha != "" {
		payload.Sha = sha
	}

	json, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(json))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiKey))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := ctx.HTTPClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("Couldn't upload file: %v", err))
	}

	defer response.Body.Close()

	status := fmt.Sprintf("Upload artifact %s with status: %s", filepath.Base(file.Name()), response.Status)
	fmt.Println(status)

	return status
}

func (ctx *Github) GetContent(contentUri string) string {
	content := Content{
		Sha: "",
	}

	req, _ := http.NewRequest("GET", contentUri, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiKey))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := ctx.HTTPClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("Couldn't upload file: %v", err))
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal(body, &content)
	}

	return content.Sha
}
