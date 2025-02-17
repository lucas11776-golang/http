package http

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/lucas11776-golang/http/types"
)

func TestRouterRouteResponse(t *testing.T) {
	t.Run("TestJsonResponse", func(t *testing.T) {
		data := struct {
			Id    int64  `id:"id"`
			Email string `email:"email"`
		}{
			Id:    1,
			Email: "jeo@doe.com",
		}

		router := &RouterGroup{}

		router.Router().Get("/api/user", func(req *Request, res *Response) *Response {
			return res.Json(data)
		})

		req, err := NewRequest("GET", "/api/user", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte("")))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		route := router.MatchWebRoute(req)

		if route == nil {
			t.Fatalf("The route %s does not exist", req.Path())
		}

		res := NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, types.Headers{
			"content-type": "application/json",
		}, []byte(""))

		httpExpected := ParseHttpResponse(res.Json(data))
		httpResponse := ParseHttpResponse(route.Call(reflect.ValueOf(req), reflect.ValueOf(res)))

		if httpExpected != httpResponse {
			t.Fatalf("Excepted http json (%s) but got (%s)", httpExpected, httpResponse)
		}
	})
}
