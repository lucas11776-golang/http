package response

import (
	"encoding/json"
	"strings"
)

// Comment
func BodyRedirect(res *Response, path string) *Response {
	res.Body([]byte(
		strings.Join([]string{
			`<!DOCTYPE html>`,
			`<head>`,
			`  <meta http-equiv="Refresh" content="0, url='` + path + `'">`,
			`</head>`,
			`<body>`,
			`  <p>You will be redirected to ` + path + `</p>`,
			`</body>`,
			`</html>`,
		}, "\r\n"),
	))
	return res.Header("content-type", "text/html; charset=utf-8").Status(HTTP_RESPONSE_TEMPORARY_REDIRECT)
}

// Comment
func BodyJson(res *Response, value any) *Response {
	res.Header("content-type", "application/json")

	data, err := json.Marshal(value)

	if err != nil {
		return res.Body([]byte("{}"))
	}

	return res.Body(data)
}

// Comment
func BodyDownload(res *Response, contentType string, filename string, binary []byte) *Response {
	res.Header("content-type", contentType)
	res.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return res.Body(binary)
}

func BodyHtml(res *Response, html string) *Response {
	res.Header("content-type", "text/html; charset=utf-8")
	return res.Body([]byte(html))
}

// Comment
func BodyView(res *Response, view string, data any) *Response {
	return res
}

// Comment
func BodyDefault(res *Response, data []byte) *Response {
	res.body = data

	return res
}
