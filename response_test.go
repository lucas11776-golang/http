package http

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"github.com/open2b/scriggo"
)

func TestResponse(t *testing.T) {
	t.Run("TestResponseOk", func(t *testing.T) {
		res := InitResponse().SetStatus(HTTP_RESPONSE_OK)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Length: 0\r\n\r\n"

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseContentTypeHeader", func(t *testing.T) {
		res := InitResponse().SetStatus(HTTP_RESPONSE_OK).SetHeader("content-type", "application/json")

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 0\r\n\r\n"

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseBody", func(t *testing.T) {
		body := []byte(`{"id": 1, "email": "jeo@doe.com"}`)
		res := InitResponse().SetStatus(HTTP_RESPONSE_OK).SetHeader("content-type", "application/json").SetBody(body)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseHtml", func(t *testing.T) {
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

		res := InitResponse().SetStatus(HTTP_RESPONSE_OK).Html(string(html))

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(html))}, ": ") + "\r\n",
			string(html) + "\r\n",
		}, "\r\n")

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseJson", func(t *testing.T) {
		j := struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
		}{
			Id:    1,
			Title: "Response With Json Body",
		}

		body, _ := json.Marshal(j)
		res := InitResponse().SetStatus(HTTP_RESPONSE_OK).Json(j)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseRedirect", func(t *testing.T) {
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

		res := InitResponse().SetStatus(HTTP_RESPONSE_OK).Redirect("authentication/login")

		httpExpected := strings.Join([]string{
			"HTTP/1.1 307 Temporary Redirect",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(html))}, ": ") + "\r\n",
			string(html) + "\r\n",
		}, "\r\n")

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseDownload", func(t *testing.T) {
		file := []byte("Hello World: " + string(strconv.Itoa(int(rand.Float64()*1000))))

		res := InitResponse().Download("text/plain; charset: utf-8", "hello.txt", file)

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Disposition: attachment; filename=\"hello.txt\"",
			"Content-Type: text/plain; charset: utf-8",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(file))}, ": ") + "\r\n",
			string(file) + "\r\n",
		}, "\r\n")

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseView", func(t *testing.T) {
		name := strings.Join([]string{"lucas", strconv.Itoa(int(rand.Float64() * 1000))}, "")

		res := InitResponse()

		req, err := NewRequest("GET", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		res.Request = req

		vw := InitView(responseViewTest, "html")

		res.Request.Server = Server("127.0.0.1", 0)

		res.Request.Server.Set("view", vw)

		res.View("home", ViewData{
			"name": name,
		})

		body := strings.Join([]string{
			`<!DOCTYPE html>`,
			`<html lang="en">`,
			`<head>`,
			`  <meta charset="UTF-8">`,
			`  <meta name="viewport" content="width=device-width, initial-scale=1.0">`,
			`  <title>Home Page</title>`,
			`</head>`,
			`<body>`,
			`  <h1>Hello user ` + name + `</h1>`,
			`</body>`,
			`</html>`,
		}, "\r\n")

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n",
			body + "\r\n",
		}, "\r\n")

		http := ParseHttpResponse(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})
}

var responseViewTest = &ViewReaderTest{
	Files: scriggo.Files{
		"home.html": []byte(strings.Join([]string{
			`<!DOCTYPE html>`,
			`<html lang="en">`,
			`<head>`,
			`  <meta charset="UTF-8">`,
			`  <meta name="viewport" content="width=device-width, initial-scale=1.0">`,
			`  <title>Home Page</title>`,
			`</head>`,
			`<body>`,
			`  <h1>Hello user {{ name }}</h1>`,
			`</body>`,
			`</html>`,
		}, "\r\n")),
	},
	cache: make(scriggo.Files),
}
