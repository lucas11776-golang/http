package path

import (
	"strings"
)

// Comment
func Path(path ...string) string {
	pth := []string{}

	for _, p := range path {
		pth = append(pth, strings.Trim(strings.ReplaceAll(p, "/", "\\"), "\\"))
	}

	return strings.Join(pth, "\\")
}
