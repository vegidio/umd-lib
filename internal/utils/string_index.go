package utils

import "strings"

// LastRightOf returns the substring of `s` that is to the right of the last occurrence of `substring`.
// If `substring` is not found in `s`, it returns `s` unchanged.
func LastRightOf(s string, substring string) string {
	lastSlashIndex := strings.LastIndex(s, substring)
	if lastSlashIndex == -1 {
		return s
	}

	return s[lastSlashIndex+1:]
}
