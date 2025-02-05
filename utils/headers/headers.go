package headers

import (
	"net/http"

	"github.com/lucas11776-golang/http/types"
)

// Comment
func ToHeader(headers types.Headers) http.Header {
	header := make(http.Header)

	for k, v := range headers {
		header[k] = []string{v}
	}

	return header
}
