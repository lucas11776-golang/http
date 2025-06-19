package pages

import (
	"net/url"
	"strings"

	"github.com/lucas11776-golang/http/utils/helper"
)

// Comment
func getUrl(to string) string {
	parsed, err := url.Parse(to)

	if err != nil || parsed.Host == "" {
		return helper.Url(to)
	}

	return to
}

// Comment
func Redirect(to string) string {
	return strings.Join([]string{
		`<!DOCTYPE html>`,
		`<head>`,
		`  <meta http-equiv="Refresh" content="0, url='` + getUrl(to) + `'">`,
		`</head>`,
		`<body>`,
		`</body>`,
		`</html>`,
	}, "\r\n")
}

// Comment
func ServerError() string {
	return ""
}

// Comment
func CsrfExpired() string {
	return "<h1>Request token has expired refresh page</h1>"
}
