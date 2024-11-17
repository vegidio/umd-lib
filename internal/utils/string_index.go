package utils

import "strings"

func LastRightOf(s string, substring string) string {
	lastSlashIndex := strings.LastIndex(s, substring)
	if lastSlashIndex == -1 {
		return s
	}

	return s[lastSlashIndex+1:]
}
