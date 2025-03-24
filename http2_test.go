package http

import (
	"fmt"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/request"
)

func TestHtt2(t *testing.T) {
	serve := func() *HTTP {
		server := Server("127.0.0.1", 0)

		server.Route().Get("/", func(req *Request, res *Response) *Response {
			return res.SetStatus(HTTP_RESPONSE_OK).Json(map[string]string{
				"message": "Successfully",
			})
		})

		go server.Listen()

		return server
	}

	t.Run("HTTP2NegotiationRequest", func(t *testing.T) {
		server := serve()

		req := request.CreateRequest().
			SetProtocal("HTTP/1.1").
			SetHeaders(types.Headers{
				"Connection":     "Upgrade, HTTP2-Settings",
				"Upgrade":        "h2c",
				"HTTP2-Settings": "AAAACCCCAAAAZZZZZ"})

		http, err := req.Get(fmt.Sprintf("http://%s:%d", server.Address(), server.Port()))

		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		res, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Failed to parse http response: %v", err)
		}

		if Status(res.StatusCode) != HTTP_RESPONSE_SWITCHING_PROTOCOLS {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_SWITCHING_PROTOCOLS, res.StatusCode)
		}

		if res.GetHeader("Connection") != "Upgrade" {
			t.Fatalf("Expected connection header to be (%s) but got (%s)", "Upgrade", res.GetHeader("Connection"))
		}

		if res.GetHeader("Upgrade") != "HTTP/2.0" {
			t.Fatalf("Expected upgrade header to be (%s) but got (%s)", "HTTP/2.0", res.GetHeader("Upgrade"))
		}

		server.Close()
	})

}

// > GET / HTTP/1.1
// > Host: 127.0.0.1:6666
// > Connection: Upgrade, HTTP2-Settings
// > Upgrade: h2c
// > HTTP2-Settings: ZZZZZZZZZZZZZZZZZ

// Connection: Upgrade
// Upgrade: HTTP/2.0
