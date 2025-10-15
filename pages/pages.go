package pages

import (
	"strings"

	"github.com/lucas11776-golang/http/utils/helper"
)

// Comment
func RedirectPage(to string) string {
	return strings.Join([]string{
		`<!DOCTYPE html>`,
		`<head>`,
		`  <meta http-equiv="Refresh" content="0, url='` + helper.GetUrl(to) + `'">`,
		`</head>`,
		`<body>`,
		`</body>`,
		`</html>`,
	}, "\r\n")
}

// Comment
func ServerErrorPage() string {
	return ""
}

// Comment
func CsrfExpired() string {
	return strings.Join([]string{
		`<!DOCTYPE html>`,
		`<head>`,
		`</head>`,
		`<body>`,
		`	<h1>CSRF token has expired</h1>`,
		`</body>`,
		`</html>`,
	}, "\r\n")
}
