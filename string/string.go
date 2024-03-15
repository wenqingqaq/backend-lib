package string

import "strings"

func TrimColon(s string) string {
	return strings.ReplaceAll(s, ":", "")
}
