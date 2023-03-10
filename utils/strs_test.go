package utils

import (
	"testing"
)

func TestClearSpaces(t *testing.T) {
	tests := []struct{
		name   string
		input  string
		result string
	}{
		{
			name: "Test of clearing spaces from strings",
			input: "   Some   text   for       testing ",
			result: "Some text for testing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearedText := ClearSpaces(test.input)

			if test.result != clearedText {
				t.Errorf("FAILED: expected %s, got %s\n", test.result, clearedText)
			}
		})
	}
}