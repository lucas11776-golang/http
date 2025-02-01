package strings

import "strings"

// Comment
func JoinPath(path ...string) string {
	pth := []string{}

	for i := range path {
		if path[i] == "" || path[i] == "/" {
			continue
		}

		pth = append(pth, strings.Trim(path[i], "/"))
	}

	// if len(pth) == 0 {
	// 	pth = append(pth, "")
	// }

	return strings.Join(pth, "/")
}
