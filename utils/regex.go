package utils

import "regexp"

func NewMatcher(s string) *regexp.Regexp {
	return regexp.MustCompile(s)
}
