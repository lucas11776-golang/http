package response

import (
	"encoding/json"
	"strings"

	"github.com/lucas11776-golang/http/view"
)

// Comment
func BodyDefault(res *Response, data []byte) *Response {
	res.body = data

	return res
}

// Comment
func BodyHtml(res *Response, html string) *Response {
	return res.SetHeader("content-type", "text/html").SetBody([]byte(html))
}

// Comment
func BodyJson(res *Response, value any) *Response {
	res.SetHeader("content-type", "application/json")

	data, err := json.Marshal(value)

	if err != nil {
		return res.SetBody([]byte("{}"))
	}

	return res.SetBody(data)
}

// Comment
func BodyRedirect(res *Response, path string) *Response {
	return res.Html(strings.Join([]string{
		`<!DOCTYPE html>`,
		`<head>`,
		`  <meta http-equiv="Refresh" content="0, url='` + path + `'">`,
		`</head>`,
		`<body>`,
		`  <p>You will be redirected to ` + path + `</p>`,
		`</body>`,
		`</html>`,
	}, "\r\n")).SetStatus(HTTP_RESPONSE_TEMPORARY_REDIRECT)
}

// Comment
func BodyDownload(res *Response, contentType string, filename string, binary []byte) *Response {
	res.SetHeader("content-type", contentType).SetHeader("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return res.SetBody(binary)
}

// Comment
func BodyView(res *Response, v string, data view.Data) *Response {
	html, err := res.Request.Server.Get("view").(*view.View).Read(v, data)

	if err != nil {
		// TODO Error page 500
		return res
	}

	return res.Html(string(html))
}
