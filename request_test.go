package http

import (
	"testing"
)

// Comment
func TestRequest(t *testing.T) {
	const httpText = "POST /authentication/login?redirect=cart&ref=lucas11776 HTTP/1.1\r\n" +
		"HOST: example.com\r\n" +
		"Content-Type: application/x-www-form-urlencoded\r\n" +
		"Content-Length: 6\r\n\r\n" +
		"user=1\r\n"

	req, _ := ParseHttpRequest(httpText)

	t.Run("TestMethod", func(t *testing.T) {
		if req.Method != "POST" {
			t.Fatalf("Expected method to be (%s) but got (%s)", "POST", req.Method)
		}
	})

	t.Run("TestPath", func(t *testing.T) {
		if req.Path() != "authentication/login" {
			t.Fatalf("Expected path to be (%s) but got (%s)", "authentication/login", req.Path())
		}
	})

	t.Run("TestProtocol", func(t *testing.T) {
		if req.Protocol() != "HTTP/1.1" {
			t.Fatalf("Expected protocol to be (%s) but got (%s)", "HTTP/1.1", req.Protocol())
		}
	})

	t.Run("TestQuery", func(t *testing.T) {
		if req.GetQuery("redirect") != "cart" {
			t.Fatalf("Expected query to be (%s) but got (%s)", "cart", req.GetQuery("redirect"))
		}

		if req.GetQuery("ref") != "lucas11776" {
			t.Fatalf("Expected query to be (%s) but got (%s)", "lucas11776", req.GetQuery("lucas11776"))
		}
	})

	t.Run("TestHost", func(t *testing.T) {
		if req.Host != "example.com" {
			t.Fatalf("Expected host to be (%s) but got (%s)", "example.com", req.Host)
		}
	})

	t.Run("TestHeaders", func(t *testing.T) {
		if req.GetHeader("content-type") != "application/x-www-form-urlencoded" {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/x-www-form-urlencoded", req.GetHeader("content-type"))
		}

		if req.GetHeader("content-length") != "6" {
			t.Fatalf("Expected content length to be (%s) but got (%s)", "6", req.GetHeader("content-length"))
		}
	})

	// TODO Build body to http.Request
	// t.Run("TestBody", func(t *testing.T) {
	// 	if req.FormValue("user") != "1" {
	// 		t.Fatalf("Expected form value user to be (%s) but got (%s)", "1", req.FormValue("user"))
	// 	}
	// })
}
