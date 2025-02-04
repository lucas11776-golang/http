package request

import (
	"testing"
)

// Comment
func TestParseHttpToRequest(t *testing.T) {
	const httpText = "POST /authentication/login?redirect=cart&ref=lucas11776 HTTP/1.1\r\n" +
		"HOST: example.com\r\n" +
		"Content-Type: application/json\r\n" +
		"Content-Length: 10\r\n\r\n" +
		"{\"user\":1}\r\n"

	req, _ := ParseHttp(httpText)

	t.Run("TestHttpMethod", func(t *testing.T) {
		if req.Method() != "POST" {
			t.Fatalf("Expected method to be (%s) but got (%s)", "POST", req.Method())
		}
	})

	t.Run("TestHttpPath", func(t *testing.T) {
		if req.Path() != "authentication/login" {
			t.Fatalf("Expected path to be (%s) but got (%s)", "authentication/login", req.Path())
		}
	})

	t.Run("TestHttpProtocol", func(t *testing.T) {
		if req.Protocol() != "HTTP/1.1" {
			t.Fatalf("Expected protocol to be (%s) but got (%s)", "HTTP/1.1", req.Protocol())
		}
	})

	t.Run("TestHttpQuery", func(t *testing.T) {
		if req.Query("redirect") != "cart" {
			t.Fatalf("Expected query to be (%s) but got (%s)", "cart", req.Query("redirect"))
		}

		if req.Query("ref") != "lucas11776" {
			t.Fatalf("Expected query to be (%s) but got (%s)", "lucas11776", req.Query("lucas11776"))
		}
	})

	t.Run("TestHttpHeaders", func(t *testing.T) {
		if req.Header("host") != "example.com" {
			t.Fatalf("Expected host to be (%s) but got (%s)", "example.com", req.Header("host"))
		}

		if req.Header("content-type") != "application/json" {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/json", req.Header("content-type"))
		}
		if req.Header("content-length") != "10" {
			t.Fatalf("Expected content length to be (%s) but got (%s)", "10", req.Header("content-length"))
		}
	})

	t.Run("TestHttpBody", func(t *testing.T) {
		if string(req.Body()) != "{\"user\":1}" {
			t.Fatalf("Expected body to be (%s) but got (%s)", "{\"user\":1}", req.Body())
		}
	})
}
