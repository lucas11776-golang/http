package response

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Comment
func ParseHttpToResponse(text string) (protocol string, statusCode int, headers types.Headers, body []byte, err error) {
	array := strings.Split(text, "\r\n")

	if len(array) < 2 {
		return "", 0, nil, nil, errors.New("")
	}

	h := strings.Split(array[0], " ")

	if len(h) < 3 {
		return "", 0, nil, nil, errors.New("")
	}

	status, err := strconv.Atoi(h[1])

	if err != nil {
		return "", 0, nil, nil, errors.New("")
	}

	hdrs := make(types.Headers)
	data := []byte{}

	for i, line := range array[1:] {
		if line == "" {
			data = append(data, []byte(strings.TrimRight(strings.Join(array[i+2:], "\r\n"), "\r\n"))...)

			break
		}

		header := strings.Split(line, ":")
		hdrs[cases.Title(language.English).String(header[0])] = strings.Trim(strings.Join(header[1:], ":"), " ")
	}

	return h[0], status, hdrs, data, nil
}

// Comment
func ResponseToHttp(res *http.Response) string {
	text := []string{}
	text = append(text, strings.Join([]string{res.Proto, res.Status}, " "))
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
		text = append(text, fmt.Sprintf(
			"%s: %s",
			cases.Title(language.English).String(key),
			strings.Join(res.Header[key], ";"),
		))
	}

	if len(body) == 0 {
		return strings.Join(append(text, "\r\n"), "\r\n")
	}

	return strings.Join(
		append(text, strings.Join([]string{"\r\n", string(body), "\r\n"}, "")),
		"\r\n",
	)
}
