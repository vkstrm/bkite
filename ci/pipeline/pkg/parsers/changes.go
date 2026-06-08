package parsers

import (
	"strings"
)

func ChangedPaths(changes string, definedPaths []string) map[string]bool {
	changedPaths := strings.Split(changes, ",")
	changeMap := map[string]bool{}
	for _, changedPath := range changedPaths {
		for _, definedPath := range definedPaths {
			if strings.HasPrefix(changedPath, definedPath) {
				changeMap[definedPath] = true
			}
		}
	}
	return changeMap
}
