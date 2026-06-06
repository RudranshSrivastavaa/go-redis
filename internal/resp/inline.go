package resp

import "strings"

func ParseInline(line string) []string {
	return strings.Fields(line)
}