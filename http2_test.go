package http

import "testing"

func TestHtt2(t *testing.T) {
	serve := func() *HTTP {
		server := Server("127.0.0.1", 0)

		server.Route().Get("/", func(req *Request, res *Response) *Response {
			return res.SetStatus(HTTP_RESPONSE_OK).Json(map[string]string{
				"message": "Successfully",
			})
		})

		return server
	}

	t.Run("HTTP2NegotiationRequest", func(t *testing.T) {
		server := serve()

		server.Close()
	})

}

// > GET / HTTP/1.1
// > Host: 127.0.0.1:6666
// > Connection: Upgrade, HTTP2-Settings
// > Upgrade: h2c
// > HTTP2-Settings: ZZZZZZZZZZZZZZZZZ
