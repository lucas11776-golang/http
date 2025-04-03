package path

import (
	pt "path"
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
