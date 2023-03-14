package utils

import (
	"testing"
)

type GrabTextFromFileTest struct {
	name 		  string
	filepath 	  string
	expected 	  string
	expectedError error
}

var pathToTestFiles = GetCurrentDirPath() + "/files_testdata/GrabTextFromFileTest/"

var grabTextFromFileTests = []GrabTextFromFileTest{
	{
		name: "Test grabbing from xhtml files",
		filepath: pathToTestFiles + "testFile.xhtml",
		expected: "Expected data is: Some test data expected after grabbing text from .xhtml files",
	},
	{
		name: "Test grabbing from txt files",
		filepath: pathToTestFiles + "testFile.txt",
		expected: "Expected data is: Some test data expected after grabbing text from .txt files",
	},
	{
		name: "Test grabbing from txt files",
		filepath: "ErrorPath",
		expected: "",
	},
}

func TestGrabTextFromFile(t *testing.T) {
	tests := grabTextFromFileTests

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := GrabTextFromFile(test.filepath)

			if test.expected != actual {
				t.Errorf("FAILED: expected %s, got %s\n", test.expected, actual)
			}
		})
	}
}