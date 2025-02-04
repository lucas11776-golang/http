package response

import (
	"embed"
	"encoding/json"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/view"
)

//go:embed views/*
var views embed.FS

type ViewReader struct{}

func (ctx *ViewReader) Open(name string) (fs.File, error) {
	return views.Open(strings.Join([]string{"views", name}, "/"))
}

func TestHttpResponse(t *testing.T) {
	t.Run("TestHttpResponseOk", func(t *testing.T) {
		res := Init().Status(HTTP_RESPONSE_OK)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Length: 0\r\n\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseContentTypeHeader", func(t *testing.T) {
		res := Init().Status(HTTP_RESPONSE_OK).Header("content-type", "application/json")

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 0\r\n\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseBody", func(t *testing.T) {
		body := []byte(`{"id": 1, "email": "jeo@doe.com"}`)
		res := Init().Status(HTTP_RESPONSE_OK).Header("content-type", "application/json").Body(body)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseHtml", func(t *testing.T) {
		html := []byte(
			strings.Join([]string{
				`<!DOCTYPE html>`,
				`<head>`,
				`  <style>`,
				`    h1 { font-size: 5em; color: green; }`,
				`  </style>`,
				`</head>`,
				`<body>`,
				`  <h1>Hello World!!!</h1>`,
				`</body>`,
				`</html>`,
			}, "\r\n"),
		)

		res := Init().Status(HTTP_RESPONSE_OK).Html(string(html))

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(html))}, ": ") + "\r\n",
			string(html) + "\r\n",
		}, "\r\n")

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseJson", func(t *testing.T) {
		j := struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
		}{
			Id:    1,
			Title: "Response With Json Body",
		}

		body, _ := json.Marshal(j)
		res := Init().Status(HTTP_RESPONSE_OK).Json(j)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseRedirect", func(t *testing.T) {
		html := []byte(
			strings.Join([]string{
				`<!DOCTYPE html>`,
				`<head>`,
				`  <meta http-equiv="Refresh" content="0, url='authentication/login'">`,
				`</head>`,
				`<body>`,
				`  <p>You will be redirected to authentication/login</p>`,
				`</body>`,
				`</html>`,
			}, "\r\n"),
		)

		res := Init().Status(HTTP_RESPONSE_OK).Redirect("authentication/login")

		httpExpected := strings.Join([]string{
			"HTTP/1.1 307 Temporary Redirect",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(html))}, ": ") + "\r\n",
			string(html) + "\r\n",
		}, "\r\n")

		http := ParseHttp(res)

		// TODO Test will fail sometimes because map does not go by order on loop
		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseDownload", func(t *testing.T) {
		file := []byte("Hello World: " + string(strconv.Itoa(int(rand.Float64()*1000))))

		res := Init().Download("text/plain; charset: utf-8", "hello.txt", file)

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Disposition: attachment; filename=\"hello.txt\"",
			"Content-Type: text/plain; charset: utf-8",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(file))}, ": ") + "\r\n",
			string(file) + "\r\n",
		}, "\r\n")

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseView", func(t *testing.T) {
		name := strings.Join([]string{"lucas", strconv.Itoa(int(rand.Float64() * 1000))}, "")

		res := Init()

		res.Request = request.Create("GET", "/", make(types.Query), "HTTP/1.1", make(types.Headers), []byte{})

		res.Request.Server = server.Init("127.0.0.1", 8080, nil).Set("view", view.Init(&ViewReader{}, "html"))

		res.View("index", view.Data{
			"name": name,
		})

		body := string(
			strings.Join([]string{
				`<!DOCTYPE html>`,
				`<html lang="en">`,
				`<head>`,
				`  <meta charset="UTF-8">`,
				`  <meta name="viewport" content="width=device-width, initial-scale=1.0">`,
				`  <title>View Render Page</title>`,
				`</head>`,
				`<body>`,
				strings.Join([]string{"  <h1>Hello user ", name, "</h1>"}, ""),
				`</body>`,
				`</html>`,
			}, "\r\n"),
		)

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len([]byte(body)))}, ": ") + "\r\n",
			body + "\r\n",
		}, "\r\n")

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})
}
