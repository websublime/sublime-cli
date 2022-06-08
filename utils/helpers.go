package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PackageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Typings         string            `json:"typings"`
	Module          string            `json:"module"`
	Scripts         map[string]string `json:"scripts"`
	Keywords        []string          `json:"keywords"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func PathWalk(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking path: %v", err)
	}

	return files, nil
}
