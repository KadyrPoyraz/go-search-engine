package utils

import "os"

func GetCurrentDirPath() string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path
}

