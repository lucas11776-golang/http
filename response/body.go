package response

import (
	"encoding/json"
)

// Comment
func BodyRedirect(res *Response, path string) *Response {
	// res.body = []byte(
	// 	strings.Join([]string{
	// 		`<!DOCTYPE html>`,
	// 		`<head>`,
	// 		`  <meta name="viewport" content="0, url='` + string(res.body) + `'">`,
	// 		`</head>`,
	// 		`<body>`,
	// 		`  <p>You will be redirected to ` + string(res.body) + `</p>`,
	// 		`</body>`,
	// 		`</html>`,
	// 	}, "\r\n"),
	// )

	return res
}

// Comment
func BodyJson(res *Response, value any) *Response {
	res.Header("content-type", "application/json")

	data, err := json.Marshal(value)

	if err != nil {
		res.body = []byte("{}")

		return res
	}

	res.body = data

	return res
}

// Comment
func BodyDownload(res *Response, contentType string, filename string, binary []byte) *Response {
	return res
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
