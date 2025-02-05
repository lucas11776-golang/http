package http

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/router"
	"github.com/lucas11776-golang/http/types"
	req "github.com/lucas11776-golang/http/utils/request"
)

func TestServerWeb(t *testing.T) {
	server := Server("127.0.0.1", 0)

	users := []User{
		(User{ID: 1, Email: "jane@doe.com"}),
		(User{ID: 2, Email: "jeo@doe.com"}),
	}

	server.Route().Group("api", func(route *router.Router) {
		route.Group("users", func(route *router.Router) {
			route.Get("/", func(req *request.Request, res *response.Response) *response.Response {
				return res.Json(users)
			})
			route.Post("/", func(req *request.Request, res *response.Response) *response.Response {
				return res.Status(response.HTTP_RESPONSE_OK).Json(userCreatedMessage)
			}).Middleware(AuthorizationGuard)
			route.Group("{id}", func(route *router.Router) {
				route.Get("/", func(req *request.Request, res *response.Response) *response.Response {
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

		res := response.Create("HTTP/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{}).Json(users)
		expectedHttp := response.ParseHttp(res)

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
func AuthorizationGuard(req *request.Request, res *response.Response, next router.Next) *response.Response {
	if req.GetHeader("authorization") != AuthKey {
		return res.Status(response.HTTP_RESPONSE_UNAUTHORIZED).Json(userCreatedMessage)
	}

	return next()
}
