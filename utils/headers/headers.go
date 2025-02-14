package headers

import (
	"net/http"

	"github.com/lucas11776-golang/http/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Comment
func ToHeader(headers types.Headers) http.Header {
	header := make(http.Header)

	for k, v := range headers {
		header[cases.Title(language.English).String(k)] = []string{v}
	}

	return header
}
