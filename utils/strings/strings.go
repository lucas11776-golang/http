package strings

import (
	"math/rand"
	"strings"
)

// Comment
func JoinPath(path ...string) string {
	pth := []string{}

	for i := range path {
		if path[i] == "" || path[i] == "/" {
			continue
		}

		pth = append(pth, strings.Trim(path[i], "/"))
	}

	return strings.Join(pth, "/")
}

// Comment
func Random(size int) string {
	str := ""

	for i := 0; i < size; i++ {
		str += string(byte(65 + (rand.Float32() * 58)))
	}

	return str
}
