package parsers

import (
	"os"
	"strings"
)

func ChangedPaths(definedPaths []string) map[string]bool {
	changedPaths := strings.Split(os.Getenv("CHANGED_PATHS"), ",")
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
