package http

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	req "github.com/lucas11776-golang/http/utils/request"
)

func TestServerWeb(t *testing.T) {
	server := Server("127.0.0.1", 0)

	users := []User{
		(User{ID: 1, Email: "jane@doe.com"}),
		(User{ID: 2, Email: "jeo@doe.com"}),
	}

	server.Route().Group("api", func(route *Router) {
		route.Group("users", func(route *Router) {
			route.Get("/", func(req *Request, res *Response) *Response {
				return res.Json(users)
			})
			route.Post("/", func(req *Request, res *Response) *Response {
				return res.SetStatus(HTTP_RESPONSE_OK).Json(userCreatedMessage)
			}).Middleware(AuthorizationGuard)
			route.Group("{id}", func(route *Router) {
				route.Get("/", func(req *Request, res *Response) *Response {
					return res
				})
			})
		})
	})

	go func() {
		server.Listen()
	}()

	t.Run("TestApiGetUsers", func(t *testing.T) {
		r := req.CreateRequest().Header("content-type", "application/json")

		http, err := r.Get(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""))

		if err != nil {
			t.Fatalf("Something went wrong went trying to send request: %s", err.Error())
		}

		res := NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, make(types.Headers), []byte{}).Json(users)
		expectedHttp := ParseHttpResponse(res)

		if expectedHttp != http {
			t.Fatalf("Expected response to be (%s) but got (%s), (%d,%d)", expectedHttp, http, len(expectedHttp), len(http))
		}
	})

	server.Close()
}

var AuthKey = "KEY-" + strconv.Itoa(int(rand.Float32()*10000))

type Message struct {
	Message string `json:"message"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

var userCreatedMessage = Message{
	Message: "User has been created successfully",
}

// Comment
func AuthorizationGuard(req *Request, res *Response, next Next) *Response {
	if req.GetHeader("authorization") != AuthKey {
		return res.SetStatus(HTTP_RESPONSE_UNAUTHORIZED).Json(userCreatedMessage)
	}

	return next()
}
