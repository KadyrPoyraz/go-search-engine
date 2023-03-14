package lexer

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/dchest/stemmer/porter2"
)

type Lexer struct {
	Content []string
	Value    string
}

func (l *Lexer) TrimWhiteSpaces() {
	for ; len(l.Content) > 0; {
		isWhitespacePresent := regexp.MustCompile(`\s`).MatchString(l.Content[0])
		if isWhitespacePresent {
			l.Content = l.Content[1:]
		} else {
			break
		}
	}
}

func (l *Lexer) validateToken(token []string) string {
	eng := porter2.Stemmer

	validatedToken := strings.Join(token, "")
	validatedToken = eng.Stem(validatedToken)
	validatedToken = strings.ToUpper(validatedToken)

	return validatedToken
}

func (l *Lexer) chop(n int) {
	token := l.Content[0:n]
	l.Content = l.Content[n:]
	l.Value = l.validateToken(token)
}

func (l *Lexer) GetNextToken() bool {
	l.TrimWhiteSpaces()

	if len(l.Content) == 0 {
		return false
	}

	for _, r := range l.Content[0] {
		if unicode.IsNumber(r) {
			n := 0
			for ; n < len(l.Content) ; {
				_, err := strconv.ParseFloat(l.Content[n], 64)
				if err == nil {
					n += 1
				} else {
					break
				}
			}
			l.chop(n)
			return true
		}

		if unicode.IsLetter(r) {
			n := 0
			for ; n < len(l.Content) ; {
				//is_alphanumeric := regexp.MustCompile(``).MatchString(l.Content[n])
				isWhitespacePresent := regexp.MustCompile(`\s`).MatchString(l.Content[n])
				isAlphanumeric := regexp.MustCompile(`[$&+,:;=?@#|'<>.^*()%!-]`).MatchString(l.Content[n])
				if !isAlphanumeric && !isWhitespacePresent {
					n += 1
				} else {
					break
				}
			}
			l.chop(n)
			return true
		}
	}

	l.chop(1)

	return true
}
