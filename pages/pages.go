package pages

import (
	"net/url"
	"strings"

	"github.com/lucas11776-golang/http/utils/helper"
)

// Comment
func Redirect(to string) string {
	if _, err := url.Parse(to); err != nil {
		to = helper.Url(to)
	}

	return strings.Join([]string{
		`<!DOCTYPE html>`,
		`<head>`,
		`  <meta http-equiv="Refresh" content="0, url='` + to + `'">`,
		`</head>`,
		`<body>`,
		`  <p>You will be redirected to ` + to + `</p>`,
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
