package response

import (
	"testing"
)

func TestResponse(t *testing.T) {
	t.Run("TestResponseOk", func(t *testing.T) {
		res := Init().SetStatus(HTTP_RESPONSE_OK)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Length: 0\r\n\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestResponseContentTypeHeader", func(t *testing.T) {
		res := Init().SetStatus(HTTP_RESPONSE_OK).SetHeader("content-type", "application/json")

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 0\r\n\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Fatalf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})
}
