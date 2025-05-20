package path

import (
	"fmt"
	pt "path"
	"strings"
)

// Comment
func Path(path ...string) string {
	// pth := []string{}

	// for _, p := range path {
	// 	pth = append(pth, strings.Trim(strings.ReplaceAll(p, "/", "\\"), "\\"))
	// }

	// file, err := os.Open(pt.Join(path...))

	// fmt.Println("PATH ----> ", "FILE", file, "ERROR", err)

	// return "/" + pt.Join(strings.Join(pth, "\\"))

	return pt.Join(path...)
}

// Comment
func FileRealPath(file string, extension string) string {
	return fmt.Sprintf("%s.%s", strings.ReplaceAll(strings.ReplaceAll(file, "\\", "/"), ".", "/"), extension)
}
