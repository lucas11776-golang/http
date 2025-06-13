package pages

import "strings"

// Comment
func Redirect(to string) string {
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
