package response

import (
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Comment
func ParseHttp(res *Response) string {
	http := []string{}

	http = append(http, strings.Join([]string{res.Proto, res.Status}, " "))

	keys := make([]string, 0, len(res.Header))

	for k := range res.Header {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		k := cases.Title(language.English).String(key)
		v := strings.Join(res.Header[key], ";")

		http = append(http, strings.Join([]string{k, v}, ": "))
	}

	http = append(http, strings.Join([]string{"Content-Length", strconv.Itoa(len(res.body))}, ": "))

	if len(res.body) == 0 {
		return strings.Join(append(http, "\r\n"), "\r\n")
	}

	return strings.Join(append(http, strings.Join([]string{"\r\n", string(res.body), "\r\n"}, "")), "\r\n")
}
