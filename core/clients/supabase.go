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
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/websublime/sublime-cli/utils"
)

const (
	AuthEndpoint = "auth/v1"
	RestEndpoint = "rest/v1"
)

type Supabase struct {
	BaseURL string
	// apiKey can be a client API key or a service key
	apiKey     string
	HTTPClient *http.Client
}

func NewSupabase(baseURL string, supabaseKey string) *Supabase {
	return &Supabase{
		BaseURL: baseURL,
		apiKey:  supabaseKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (ctx *Supabase) Upload(bucket string, directory string, destination string) {
	fileList, err := utils.PathWalk(directory)
	if err != nil {
		panic(err)
	}

	for idx := range fileList {
		file, err := os.OpenFile(fileList[idx], os.O_RDWR, 0755)
		if err != nil {
			panic(fmt.Errorf("os.Open: %v", err))
		}
		defer file.Close()

		fileContents, err := ioutil.ReadFile(fileList[idx])
		if err != nil {
			panic(fmt.Errorf("Couldn't read file: %v", err))
		}

		fileExtension := filepath.Ext(fileList[idx])
		mime := utils.GetMimeType(strings.TrimPrefix(fileExtension, "."))

		stat, err := file.Stat()
		if err != nil {
			panic(fmt.Errorf("Couldn't get stat from file: %v", err))
		}

		payload := new(bytes.Buffer)
		writer := multipart.NewWriter(payload)
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			panic(fmt.Errorf("Couldn't create form for file: %v", err))
		}

		part.Write(fileContents)

		writer.Close()

		uri := fmt.Sprintf("%s/storage/v1/object/%s/%s/%s", ctx.BaseURL, bucket, destination, filepath.Base(file.Name()))

		req, _ := http.NewRequest("POST", uri, payload)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.apiKey))
		req.Header.Add("apikey", ctx.apiKey)
		req.Header.Add("x-upsert", "true")
		req.Header.Add("x-cache-control", "3600")
		req.Header.Add("Content-Type", mime)
		req.Header.Add("Content-Length", fmt.Sprintf("%d", stat.Size()))

		// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
		response, err := ctx.HTTPClient.Do(req)
		if err != nil {
			panic(fmt.Errorf("Couldn't upload file: %v", err))
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(fmt.Errorf("Couldn't read body response: %v", err))
		}

		fmt.Println(string(body))
	}
}
