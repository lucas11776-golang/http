package response

import (
	"strconv"
	"strings"
	"testing"
)

func TestHttpResponse(t *testing.T) {
	t.Run("TestHttpResponseOk", func(t *testing.T) {
		res := Init().Status(HTTP_RESPONSE_OK)

		httpExpected := "HTTP/1.1 200 Ok\r\n\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Errorf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpRespondContentTypeHeader", func(t *testing.T) {
		res := Init().Status(HTTP_RESPONSE_OK).Header("content-type", "application/json")

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Errorf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})

	t.Run("TestHttpResponseSetBody", func(t *testing.T) {
		body := []byte(`{"id": 1, "email": "jeo@doe.com"}`)
		res := Init().Status(HTTP_RESPONSE_OK).Header("content-type", "application/json").Body(body)

		httpExpected := "HTTP/1.1 200 Ok\r\n" +
			"Content-Type: application/json\r\n" +
			strings.Join([]string{"Content-Length", strconv.Itoa(len(body))}, ": ") + "\r\n\r\n" +
			string(body) + "\r\n"

		http := ParseHttp(res)

		if httpExpected != http {
			t.Errorf("Expected response to be (%s) but go (%s)", httpExpected, http)
		}
	})
}
