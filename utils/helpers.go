package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func PathWalk(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking path: %v", err)
	}

	return files, nil
}
