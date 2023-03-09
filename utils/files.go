package utils

import (
	"os"
	"path/filepath"
)

func GrabTextFromFile(filePath string) string {
	fileExtension := filepath.Ext(filePath)

	switch fileExtension {
	case ".xhtml":
		bytesFile, err := os.ReadFile(filePath)

		if err != nil {
			panic(err)
		}

		fileContent := string(bytesFile)
		content := ParseXHTML(fileContent)

		return content
	case ".txt":
		bytesFile, err := os.ReadFile(filePath)

		if err != nil {
			panic(err)
		}

		fileContent := string(bytesFile)
		content := ClearSpaces(fileContent)

		return content
	default:
		return ""
	}
}
