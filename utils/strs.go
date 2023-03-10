package utils

import (
	"regexp"
	"strings"
)

func ClearSpaces(from string) string {
	space := regexp.MustCompile(`\s+`)
	from = space.ReplaceAllString(from , " ")

	return strings.Trim(from, " ")
}
