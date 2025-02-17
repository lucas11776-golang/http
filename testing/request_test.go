package testing

import (
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
)

func TestRequest(t *testing.T) {

	t.Run("TestSetProtocolPathMethodHeadersBody", func(t *testing.T) {
		req := NewRequest(&TestCase{
			HTTP: http.Server("127.0.0.1", 0),
		})

		body := `{"id":1,"email":"jeo@doe.com"}`

		req.Protocol("HTTP/1.1")
		req.setPath("api/products")
		req.setMethod(http.METHOD_POST)
		req.SetHeader("content-type", "application/json")
		req.SetHeaders(types.Headers{"user-agent": "Mozilla/5.0"})
		req.setBody([]byte(body))

		r, err := req.make()

		if err != nil {
			t.Fatalf("Something want wrong when trying to create request: %v", err)
		}

		if r.Protocol() != "HTTP/1.1" {
			t.Fatalf("Expected request protocol to be (%s) but got (%s)", "HTTP/1.1", r.Protocol())
		}

		if r.Path() != "api/products" {
			t.Fatalf("Expected request path to be (%s) but got (%s)", "api/products", r.Path())
		}

		if r.Method != string(http.METHOD_POST) {
			t.Fatalf("Expected request method to be (%s) but got (%s)", http.METHOD_POST, r.Method)
		}

		if r.GetHeader("content-type") != "application/json" {
			t.Fatalf("Expected request content-type to be (%s) but got (%s)", "application/json", r.GetHeader("content-type"))
		}

		if r.GetHeader("user-agent") != "Mozilla/5.0" {
			t.Fatalf("Expected request user-agent to be (%s) but got (%s)", "Mozilla/5.0", r.GetHeader("user-agent"))
		}

		tBody := make([]byte, len(body))

		_, err = r.Body.Read(tBody)

		if err != nil {
			t.Fatalf("Something want wrong when trying to read request body: %v", err)
		}

		if string(tBody) != body {
			t.Fatalf("Expected request body to be (%s) but got (%s)", body, tBody)
		}

		req.TestCase.Cleanup()
	})

	t.Run("TestSendJson", func(t *testing.T) {
		server := http.Server("127.0.0.1", 0)

		req := NewRequest(&TestCase{
			HTTP: server,
		})

		user := struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}{
			ID:    1,
			Email: "jeo@doe.com",
		}

		server.Route().Get("users/{id}", func(req *http.Request, res *http.Response) *http.Response {
			return res.Json(user)
		})

		// req := NewRequest(&Testing{
		// 	HTTP: server,
		// })

		// req.Headers(types.Headers{
		// 	"authorization": strings.Join([]string{str.Random(50)}, "."),
		// })

		req.TestCase.Cleanup()
	})

}
