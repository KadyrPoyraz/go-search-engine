package utils

import (
	"strings"
	"unicode/utf8"
)

const (
	tagStart = 60 // Unicode `<`
	tagEnd   = 62 // Unicode `>`
)

func ParseXHTML(xhtml string) string {
	xhtml = StripTags(xhtml)
	xhtml = ClearSpaces(xhtml)

	return xhtml
}

func StripTags(from string) string {
	targetString := from
	var builder strings.Builder
	builder.Grow(len(targetString) + utf8.UTFMax)

	in := false
	start := 0
	end := 0

	for i, c := range targetString {
		if (i+1) == len(targetString) && end >= start {
			builder.WriteString(targetString[end:])
		}

		if c != tagStart && c != tagEnd {
			continue
		}

		if c == tagStart {
			if !in {
				start = i
			}
			in = true

			builder.WriteString(targetString[end:start])
			continue
		}
		in = false
		end = i + 1
	}
	targetString = builder.String()

	return targetString
}

