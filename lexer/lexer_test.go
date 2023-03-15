package lexer

import (
	"strings"
	"testing"
)

type Test struct {
	name string
	lexer Lexer
	result []string
}

func TestLexer_GetNextToken(t *testing.T) {
	tests := []Test{
		{
			name: "Test if lexer.GetNextToken returns correct data",
			lexer: Lexer{Content: strings.Split("test(someArgument) {} 1234 5", "")},
			result: []string{"TEST", "(", "SOMEARGU", ")", "{", "}", "1234", "5"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := -1
			for test.lexer.GetNextToken() {
				i += 1
				expectedValue := test.result[i]
				actualValue := test.lexer.Value

				if expectedValue != actualValue {
					t.Errorf("FAILED: expected %s, got %s\n", expectedValue, actualValue)
				}
			}
		})
	}
}

func TestLexer_TrimWhiteSpaces(t *testing.T) {
	tests := []Test{
		{
			name: "Lexer TrimWhiteSpaces works correctly",
			lexer: Lexer{Content: []string{" ", " ", "test", " ", " ", "test2"}},
			result: []string{"TEST", "TEST2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := -1
			for test.lexer.GetNextToken() {
				i += 1
				expectedValue := test.result[i]
				actualValue := test.lexer.Value

				if expectedValue != actualValue {
					t.Errorf("FAILED: expected %s, got %s\n", expectedValue, actualValue)
				}
			}
		})
	}
}