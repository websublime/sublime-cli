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
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/websublime/sublime-cli/models"
	"github.com/websublime/sublime-cli/utils"
)

func (ctx *Supabase) CreateOrganizationBucket(name string, public bool) (models.Bucket, error) {
	bucket := models.NewBucket(name, slug.Make(name), public)
	model := models.Bucket{}

	payload, err := json.Marshal(bucket)
	if err != nil {
		return model, err
	}

	uri := fmt.Sprintf("%s/%s/bucket", ctx.BaseURL, StorageEndpoint)

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

	return model, err
}

func (ctx *Supabase) Upload(bucket string, filePath string, destination string) (models.BucketUpload, error) {
	model := models.BucketUpload{}
	file, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err != nil {
		return model, err
	}
	defer file.Close()

	fileExtension := filepath.Ext(filePath)
	mime := utils.GetMimeType(strings.TrimPrefix(fileExtension, "."))

	stat, err := file.Stat()
	if err != nil {
		return model, err
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
		return model, err
	}

	_, err = io.Copy(formfile, file)
	if err != nil {
		return model, err
	}

	writer.Close()

	uri := fmt.Sprintf("%s/%s/object/%s/%s/%s", ctx.BaseURL, StorageEndpoint, bucket, destination, filepath.Base(file.Name()))

	req, _ := http.NewRequest("POST", uri, payload)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.ApiToken))
	req.Header.Add("apikey", ctx.ApiKey)
	req.Header.Add("x-upsert", "true")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Content-Length", fmt.Sprintf("%d", stat.Size()))

	// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
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
