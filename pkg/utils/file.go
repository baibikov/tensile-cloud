package utils

import "strings"

func FileType(name string) string {
	nodes := strings.Split(name, ".")
	if len(nodes) == 0 {
		return ""
	}

	return nodes[len(nodes)-1]
}
