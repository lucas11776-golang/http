package router

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
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

		route := router.Router().Get("/api/user", func(req *request.Request, res *response.Response) *response.Response {
			return res.Json(data)
		})

		req, err := request.Create("GET", "/", "github.com/lucas11776-golang/http/1.1", make(types.Headers), bytes.NewReader([]byte("")))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		res := response.Create("github.com/lucas11776-golang/http/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte(""))

		httpExpected := response.ParseHttp(res.Json(data))
		httpRoute := string(route.Call(reflect.ValueOf(req), reflect.ValueOf(res)))

		if httpExpected != httpRoute {
			t.Fatalf("Excepted http json (%s) but got (%s)", httpExpected, httpRoute)
		}
	})
}
