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
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
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
	BaseURL     string
	ApiKey      string
	Environment string
	HTTPClient  *http.Client
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func NewSupabase(baseURL string, supabaseKey string, env string) *Supabase {
	return &Supabase{
		BaseURL:     baseURL,
		ApiKey:      supabaseKey,
		Environment: env,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (ctx *Supabase) Upload(bucket string, filePath string, destination string) string {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err != nil {
		panic(fmt.Errorf("os.Open: %v", err))
	}
	defer file.Close()

	fileExtension := filepath.Ext(filePath)
	mime := utils.GetMimeType(strings.TrimPrefix(fileExtension, "."))

	stat, err := file.Stat()
	if err != nil {
		panic(fmt.Errorf("Couldn't get stat from file: %v", err))
	}

	payload := new(bytes.Buffer)
	writer := multipart.NewWriter(payload)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename="%s"`,
			escapeQuotes(filepath.Base(file.Name()))))
	h.Set("Content-Type", mime)

	formfile, err := writer.CreatePart(h)
	if err != nil {
		panic(fmt.Errorf("Couldn't create form for file: %v", err))
	}

	_, err = io.Copy(formfile, file)
	if err != nil {
		panic(fmt.Errorf("Couldn't copy file to form: %v", err))
	}

	writer.Close()

	uri := fmt.Sprintf("%s/storage/v1/object/%s/%s/%s", ctx.BaseURL, bucket, destination, filepath.Base(file.Name()))

	req, _ := http.NewRequest("POST", uri, payload)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiKey))
	req.Header.Add("apikey", ctx.ApiKey)
	req.Header.Add("x-upsert", "true")
	req.Header.Add("Content-Type", writer.FormDataContentType())
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

	return string(body)
}
