package path

import (
	"fmt"
	pt "path"
	"strings"
)

// Comment
func Path(path ...string) string {
	return pt.Join(path...)
}

// Comment
func FileRealPath(file string, extension string) string {
	return fmt.Sprintf("%s.%s", strings.ReplaceAll(strings.ReplaceAll(file, "\\", "/"), ".", "/"), extension)
}
