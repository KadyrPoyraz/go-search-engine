package utils

import "os"

var PathToTestFiles = GetCurrentDirPath() + "/files_testdata/GrabTextFromFileTest/"

func GetCurrentDirPath() string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path
}

