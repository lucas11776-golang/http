package response

import (
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Comment
func ParseHttpResponse(res *http.Response) string {
	http := []string{}
	http = append(http, strings.Join([]string{res.Proto, res.Status}, " "))
	keys := make([]string, 0, len(res.Header))
	body, err := io.ReadAll(res.Body)

	if err != nil {
		body = []byte{}
	}

	res.Header.Set("Content-Length", strconv.Itoa(len(body)))

	for k := range res.Header {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		http = append(http, strings.Join([]string{
			cases.Title(language.English).String(key),
			strings.Join(res.Header[key], ";"),
		}, ": "))
	}

	if len(body) == 0 {
		return strings.Join(append(http, "\r\n"), "\r\n")
	}

	return strings.Join(append(http, strings.Join([]string{"\r\n", string(body), "\r\n"}, "")), "\r\n")
}
