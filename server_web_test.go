package http

import (
	"errors"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/reader"
	req "github.com/lucas11776-golang/http/utils/request"
	"github.com/open2b/scriggo"
)

func TestServerWeb(t *testing.T) {
	server := Server("127.0.0.1", 0)

	users := []User{
		(User{ID: 1, Role: 1, Email: "jane@doe.com"}),
		(User{ID: 2, Role: 0, Email: "jeo@doe.com"}),
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

	server.Route().Group("authentication", func(route *Router) {
		route.Group("login", func(route *Router) {
			route.Post("/", func(req *Request, res *Response) *Response {
				// res.Session.Set("user_id", strconv.Itoa(int(users[0].ID)))

				return res.Redirect("dashboard")
			})
		})
	})

	server.Route().Group("dashboard", func(route *Router) {
		route.Get("/", func(req *Request, res *Response) *Response {
			html := strings.Join([]string{"<h1>Welcome to dashboard user ", req.Session.Get("user_id"), "</h1>"}, "")

			return res.Html(html)
		})
	}).Middleware()

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

	t.Run("TestMiddlewareUserPost", func(t *testing.T) {
		r := req.CreateRequest().Header("content-type", "application/json")

		http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Something went wrong went trying to send request: %s", err.Error())
		}

		res := NewResponse("HTTP/1.1", HTTP_RESPONSE_UNAUTHORIZED, make(types.Headers), []byte{}).Json(unauthorizedAccessMessage)
		expectedHttp := ParseHttpResponse(res)

		if expectedHttp != http {
			t.Fatalf("Expected response to be (%s) but got (%s), (%d,%d)", expectedHttp, http, len(expectedHttp), len(http))
		}

		// With key
		r = req.CreateRequest().Header("content-type", "application/json").Header("authorization", AuthKey)

		http, err = r.Post(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Something went wrong went trying to send request: %s", err.Error())
		}

		res = NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, make(types.Headers), []byte{}).Json(userCreatedMessage)
		expectedHttp = ParseHttpResponse(res)

		if expectedHttp != http {
			t.Fatalf("Expected response to be (%s) but got (%s), (%d,%d)", expectedHttp, http, len(expectedHttp), len(http))
		}
	})

	t.Run("TestStatic", func(t *testing.T) {
		server.Set("static", InitStatic(&webServerReaderTest{
			cache: make(scriggo.Files),
		}))

		r := req.CreateRequest().Header("Accept", "text/css")

		http, err := r.Get(strings.Join([]string{"http://", server.Host(), "/", cssNameWebServer}, ""))

		if err != nil {
			t.Fatalf("Something went wrong when trying get static asset: %s", err.Error())
		}

		headers := types.Headers{"content-type": "text/css"}
		body := []byte(cssContentWebServer)

		res := NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, headers, []byte{}).SetBody(body)
		expectedHttp := ParseHttpResponse(res)

		if expectedHttp != http {
			t.Fatalf("Expected response to be (%s) but got (%s)", expectedHttp, http)
		}
	})

	t.Run("TestSession", func(t *testing.T) {
		// server.Session([]byte(str.Random(10)))

		// r := req.CreateRequest()

		// http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/authentication/login"}, ""), []byte{})

		// if err != nil {
		// 	t.Fatalf("Something went wrong when trying to login: %s", err.Error())
		// }

		// fmt.Println(http)
	})

	err := server.Close()

	if err != nil {
		t.Fatalf("Something went wrong when trying to close server: %s", err.Error())
	}
}

var (
	INVALID_HTTP_RESPONSE = errors.New("Invalid http response")
)

// Comment
func HttpToResponse(http string) (*Request, error) {
	hp := strings.Split(http, "\r\n")

	if len(hp) < 2 {
		return nil, INVALID_HTTP_RESPONSE
	}

	header := strings.Split(hp[0], " ")

	if len(header) < 3 {
		return nil, INVALID_HTTP_RESPONSE
	}

	// headers := make(types.Headers)
	// body := make([]byte{}, 0)

	for _, v := range hp[1:] {

		if v == "" {
			// y :=
		}

	}

	return nil, nil
}

var AuthKey = "KEY-" + strconv.Itoa(int(rand.Float32()*10000))

type Message struct {
	Message string `json:"message"`
}

type User struct {
	ID    int64  `json:"id"`
	Role  byte   `json:"role"`
	Email string `json:"email"`
}

var unauthorizedAccessMessage = Message{
	Message: "Authorization key is invalid",
}

var userCreatedMessage = Message{
	Message: "Authorization key is invalid",
}

// Comment
func AuthorizationGuard(req *Request, res *Response, next Next) *Response {
	if req.GetHeader("authorization") != AuthKey {
		return res.SetStatus(HTTP_RESPONSE_UNAUTHORIZED).Json(unauthorizedAccessMessage)
	}

	return next()
}

var cssNameWebServer = "assets/css/main.css"

var cssContentWebServer = strings.Join([]string{
	"body { margin: 0px !important; padding: 0px !important; background-color: green; }",
}, "\r\n")

var webServerReaderTestFS = scriggo.Files{
	cssName: []byte(cssContentWebServer),
}

type webServerReaderTest struct {
	cache scriggo.Files
}

// Comment
func (ctx *webServerReaderTest) Open(name string) (fs.File, error) {
	return webServerReaderTestFS.Open(name)
}

// Comment
func (ctx *webServerReaderTest) Cache(name string) (scriggo.Files, error) {
	return reader.ReadCache(ctx, ctx.cache, name)
}
