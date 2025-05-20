package http

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/pages"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/response"
	"github.com/open2b/scriggo"
)

func TestResponse(t *testing.T) {
	t.Run("TestResponseOk", func(t *testing.T) {
		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK)
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		if res.StatusCode != 200 {
			t.Fatalf("Expected response status code to be (%d) but got (%d)", 200, res.StatusCode)
		}

		if res.GetHeader("Content-Type") == "application/json" {
			t.Fatalf("Expected content-type header to be (%s) but got (%s)", "0", res.GetHeader("Content-Length"))
		}
	})

	t.Run("TestResponseContentTypeHeader", func(t *testing.T) {
		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK).
			SetHeader("content-type", "application/json")
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		if res.GetHeader("Content-Type") != "application/json" {
			t.Fatalf("Expected content-type header to be (%s) but got (%s)", "application/json", res.GetHeader("Content-Type"))
		}
	})

	t.Run("TestResponseBody", func(t *testing.T) {
		tBody := []byte(`{"id": 1, "email": "jeo@doe.com"}`)
		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK).
			SetHeader("content-type", "application/json").
			SetBody(tBody)
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		body, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected response body to be (%s) but got (%s)", string(tBody), string(body))
		}
	})

	t.Run("TestResponseHtml", func(t *testing.T) {
		tBody := strings.Join([]string{
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
		}, "\r\n")

		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK).
			SetHeader("content-type", "application/json").
			Html(tBody)
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		body, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected response body to be (%s) but got (%s)", tBody, string(body))
		}

		if res.GetHeader("Content-Type") != "text/html" {
			t.Fatalf("Expected content-type header to be (%s) but got (%s)", "text/html", res.GetHeader("Content-Type"))
		}
	})

	t.Run("TestResponseJson", func(t *testing.T) {
		type Movie struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
		}

		j := &Movie{
			Id:    1,
			Title: "Response With Json Body",
		}

		tBody, _ := json.Marshal(j)
		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK).
			SetHeader("content-type", "application/json").
			Json(j)
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		body, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected response body to be (%s) but got (%s)", tBody, string(body))
		}

		if res.GetHeader("Content-Type") != "application/json" {
			t.Fatalf("Expected content-type header to be (%s) but got (%s)", "0", res.GetHeader("Content-Type"))
		}
	})

	t.Run("TestResponseRedirect", func(t *testing.T) {
		uri := "authentication/login"
		tBody := pages.Redirect(uri)
		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK).
			SetHeader("content-type", "application/json").
			Redirect(uri)
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		body, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		if res.StatusCode != int(HTTP_RESPONSE_TEMPORARY_REDIRECT) {
			t.Fatalf("Expected response status code to be (%d) but got (%d)", 200, res.StatusCode)
		}

		if tBody != string(body) {
			t.Fatalf("Expected response body to be (%s) but got (%s)", tBody, string(body))
		}
	})

	t.Run("TestResponseDownload", func(t *testing.T) {
		tBody := []byte("Hello World: " + string(strconv.Itoa(int(rand.Float64()*1000))))
		reply := InitResponse().SetStatus(HTTP_RESPONSE_OK).
			Download("text/plain; charset: utf-8", "hello.txt", tBody)
		res, _ := HttpToResponse(response.ResponseToHttp(reply.Response))

		body, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		if res.GetHeader("Content-Disposition") != `attachment; filename="hello.txt"` {
			t.Fatalf(
				"Expected content-disposition header to be (%s) but got (%s)",
				`attachment; filename=\"hello.txt\"`,
				res.GetHeader("Content-Disposition"),
			)
		}

		if res.GetHeader("Content-Type") != "text/plain; charset: utf-8" {
			t.Fatalf(
				"Expected content-type header to be (%s) but got (%s)",
				"text/plain; charset: utf-8",
				res.GetHeader("Content-Type"),
			)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected response body to be (%s) but got (%s)", tBody, string(body))
		}
	})

	t.Run("TestResponseView", func(t *testing.T) {
		tBody := []byte(strings.Join([]string{
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
		}, "\r\n"))

		var responseViewTest = &ViewReaderTest{
			Ext:   "html",
			Files: scriggo.Files{"home.html": tBody},
		}

		name := strings.Join([]string{"lucas", strconv.Itoa(int(rand.Float64() * 1000))}, "")

		res := InitResponse()

		req, _ := NewRequest("GET", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		res.Request = req

		vw := NewView(responseViewTest)

		res.Request.Server = Server("127.0.0.1", 0)

		res.Request.Server.Set("view", vw)

		res.View("home", ViewData{"name": name})

		tRes, _ := HttpToResponse(response.ResponseToHttp(res.Response))

		body, err := io.ReadAll(tRes.Body)

		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		if res.GetHeader("Content-Type") != "text/html" {
			t.Fatalf("Expected content-type header to be (%s) but got (%s)", "text/html", tRes.GetHeader("Content-Type"))
		}

		if strings.ReplaceAll(string(tBody), "{{ name }}", name) != string(body) {
			t.Fatalf("Expected response body to be (%s) but got (%s)", strings.ReplaceAll(string(tBody), "{{ name }}", name), string(body))
		}

		req.Server.Close()
	})
}
