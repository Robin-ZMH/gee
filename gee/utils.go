package gee

import "strings"

type H map[string]any

func parsePattern(path string) (parts []string) {
	for _, part := range strings.Split(path, "/") {
		if part != "" {
			parts = append(parts, part)
			if strings.HasPrefix(part, "*") {
				return
			}
		}
	}
	return
}
