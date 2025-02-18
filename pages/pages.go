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
func NotFound(path string) string {
	return strings.Join([]string{
		`<!DOCTYPE html>`,
		`<html lang="en">`,
		`<head>`,
		`  <meta charset="UTF-8">`,
		`  <meta name="viewport" content="width=device-width, initial-scale=1.0">`,
		`  <title>Route Not Found</title>`,
		`</head>`,
		`<body>`,
		`  <h1>Route ` + path + ` not found</h1>`,
		`</body>`,
		`</html>`,
	}, "\r\n")
}

// Comment
func ServerError() string {
	return ""
}
