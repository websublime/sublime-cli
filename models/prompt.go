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

import (
	"errors"

	"github.com/manifoldco/promptui"
)

type PromptContent struct {
	Error   string `json:"error,omitempty"`
	Label   string `json:"label,omitempty"`
	Default string `json:"default,omitempty"`
	Hide    bool   `json:"hide,omitempty"`
	Mask    rune   `json:"mask,omitempty"`
}

type PromptSelectContent struct {
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func PromptGetInput(content PromptContent, length int) (string, error) {
	validate := func(input string) error {
		if len(input) <= length {
			return errors.New(content.Error)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:       content.Label,
		Templates:   templates,
		Validate:    validate,
		Default:     content.Default,
		HideEntered: content.Hide,
		Mask:        content.Mask,
	}

	return prompt.Run()
}

func PromptGetSelect(content PromptSelectContent) (int, string, error) {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label: content.Label,
			Items: content.Items,
		}

		index, result, err = prompt.Run()
	}

	return index, result, err
}
