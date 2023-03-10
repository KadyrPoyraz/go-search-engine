package utils

import (
	"testing"
)

type Test struct {
	name 	string
	input 	string
	result	string
}

func TestStripTags(t *testing.T) {
	tests := []Test{
		{
			name: "Test of stripping tags from string",
			input: "<someTag> Some test text <someTag/>",
			result: " Some test text ",
		},
		{
			name: "Test of stripping inner tags from string",
			input: "<someTag> Some <anotherTag> aboba <anotherTag/> test text <someTag/>",
			result: " Some  aboba  test text ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stripedString := StripTags(test.input)

			if test.result != stripedString {
				t.Errorf("FAILED: expected %s, got %s\n", test.result, stripedString)
			}

		})
	}
}

func TestParseXHTML(t *testing.T) {
	tests := []Test{
		{
			name: "Test parsing of xhtml",
			input: "<someTag> Some test text <someTag/>",
			result: "Some test text",
		},
		{
			name: "Test parsing of xhtml with inner tags",
			input: "   <someTag> Some <anotherTag> aboba <anotherTag/> test text <someTag/>  ",
			result: "Some aboba test text",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stripedString := ParseXHTML(test.input)

			if test.result != stripedString {
				t.Errorf("FAILED: expected %s, got %s\n", test.result, stripedString)
			}

		})
	}
}