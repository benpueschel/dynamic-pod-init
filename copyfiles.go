package main

import (
	"os"
	"path"
)

func CopyDir(dir string, targetDir string, recursive bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range files {
		filePath := path.Join(dir, entry.Name())
		targetPath := path.Join(targetDir, entry.Name())

		fileInfo, err := entry.Info()
		if err != nil {
			return err
		}

		if recursive && entry.IsDir() {
			os.MkdirAll(targetPath, fileInfo.Mode())
			CopyDir(filePath, targetPath, recursive)
		}

		if !entry.IsDir() {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			err = os.WriteFile(targetPath, data, fileInfo.Mode())
		}
	}

	return nil
}

