package utils

import (
	"os"
	"path/filepath"
)

func GrabTextFromFile(filePath string) string {
	fileExtension := filepath.Ext(filePath)


	if fileExtension == ".xhtml" {
		bytesFile, err := os.ReadFile(filePath)

		if err != nil {
			panic(err)
		}

		fileContent := string(bytesFile)
		content := ParseXHTML(fileContent)

		return content
	}

	return ""
}
