package utils

import "strings"

func SplitByRune(s string, r rune) []string {
	return strings.Split(s, string(r))
}

func ParseTags(tags string) []string {
	tags = strings.TrimSpace(tags)
	if tags == "" {
		return []string{}
	}

	return strings.Fields(tags)
}
