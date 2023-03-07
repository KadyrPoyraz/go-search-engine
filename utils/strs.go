package utils

import "regexp"

func ClearSpaces(from string) string {
	space := regexp.MustCompile(`\s+`)
	from = space.ReplaceAllString(from , " ")

	return from
}
