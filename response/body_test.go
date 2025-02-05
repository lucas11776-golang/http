package response

import (
	"bytes"
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
	"github.com/open2b/scriggo"
)

//go:embed views/*
var viewsBody embed.FS

type ViewReaderBody struct {
	cache scriggo.Files
}

func (ctx *ViewReaderBody) Open(name string) (fs.File, error) {
	return viewsBody.Open(strings.Join([]string{"views", name}, "/"))
}

// Comment
func (ctx *ViewReaderBody) Views(name string) (scriggo.Files, error) {
	return view.ReadViewCache(ctx, ctx.cache, name)
}

func TestBody(t *testing.T) {
	t.Run("TestResponseBody", func(t *testing.T) {
		body := []byte(`{"id": 1, "email": "jeo@doe.com"}`)
		res := Init().SetStatus(HTTP_RESPONSE_OK).SetHeader("content-type", "application/json").SetBody(body)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttp(res)

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

		res := Init().SetStatus(HTTP_RESPONSE_OK).Html(string(html))

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

	t.Run("TestResponseJson", func(t *testing.T) {
		j := struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
		}{
			Id:    1,
			Title: "Response With Json Body",
		}

		body, _ := json.Marshal(j)
		res := Init().SetStatus(HTTP_RESPONSE_OK).Json(j)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttp(res)

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

		res := Init().SetStatus(HTTP_RESPONSE_OK).Redirect("authentication/login")

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

	t.Run("TestResponseDownload", func(t *testing.T) {
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

	t.Run("TestResponseView", func(t *testing.T) {
		name := strings.Join([]string{"lucas", strconv.Itoa(int(rand.Float64() * 1000))}, "")

		res := Init()

		req, err := request.Create("GET", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		res.Request = req

		vw := view.Init(&ViewReaderBody{
			cache: make(scriggo.Files),
		}, "html")

		res.Request.Server = server.Init("127.0.0.1", 8080, nil).Set("view", vw)

		res.View("index", view.Data{
			"name": name,
		})

		body := strings.Join([]string{"<h1>Hello user ", name, "</h1>"}, "")

		httpExpected := strings.Join([]string{
			"HTTP/1.1 200 Ok",
			"Content-Type: text/html",
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n",
			body + "\r\n",
		}, "\r\n")

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})
}
