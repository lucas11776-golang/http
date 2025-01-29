package router

import (
	"http/request"
	"http/response"
	"http/types"
	"reflect"
	"testing"
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

		route := router.Router().Get("/api/user", func(req *request.Request, res *response.Response) *response.Response {
			return res.Json(data)
		})

		req := request.Create("GET", "/", make(types.Query), "HTTP/1.1", make(types.Headers), []byte(""))
		res := response.Create("HTTP/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte(""))

		httpExpected := response.ParseHttp(res.Json(data))
		httpRoute := string(route.Call(reflect.ValueOf(req), reflect.ValueOf(res)))

		if httpExpected != httpRoute {
			t.Errorf("Excepted http json (%s) but got (%s)", httpExpected, httpRoute)
		}
	})
}
