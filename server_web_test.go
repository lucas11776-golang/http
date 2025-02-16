package http

import (
	"encoding/json"
	"io/fs"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/reader"
	req "github.com/lucas11776-golang/http/utils/request"
	str "github.com/lucas11776-golang/http/utils/strings"
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
		route.Middleware(IsGuest).Group("login", func(route *Router) {
			route.Get("/", func(req *Request, res *Response) *Response {
				return res.Html("<h1>Login page</h1>")
			})
			route.Post("/", func(req *Request, res *Response) *Response {
				user := users[0]

				res.Session.Set("user_id", strconv.Itoa(int(user.ID)))
				res.Session.Set("role", strconv.Itoa(int(user.Role)))

				return res.Redirect("dashboard")
			})
		})
	})

	server.Route().Middleware(IsUser, IsAdmin).Group("dashboard", func(route *Router) {
		route.Get("/", func(req *Request, res *Response) *Response {
			html := strings.Join([]string{"<h1>Welcome to dashboard user ", req.Session.Get("user_id"), "</h1>"}, "")

			return res.Html(html)
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

		// res := NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, make(types.Headers), []byte{}).Json(users)

		rt, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Something went wrong went trying to convert http to response: %s", err.Error())
		}

		if rt.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/json", rt.GetHeader("content-type"))
		}

		if rt.GetHeader("content-type") != "application/json" {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/json", rt.GetHeader("content-type"))
		}

		tBody, _ := json.Marshal(users)

		if rt.GetHeader("content-length") != strconv.Itoa(len(tBody)) {
			t.Fatalf("Expected content length to be (%d) but got (%s)", len(tBody), rt.GetHeader("content-length"))
		}

		body := make([]byte, len(tBody))

		n, err := rt.Body.Read(body)

		if err != nil {
			t.Fatalf("Something went wrong when trying to read body: %s", err.Error())
		}

		if n != len(tBody) {
			t.Fatalf("Expected body size to be (%d) but got (%d)", len(tBody), n)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected body size to be (%s) but got (%s)", string(tBody), string(body))
		}
	})

	t.Run("TestMiddlewareUserPost", func(t *testing.T) {
		r := req.CreateRequest().Header("content-type", "application/json")

		http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Something went wrong went trying to send request: %s", err.Error())
		}

		rt, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Something went wrong went trying to convert http to response: %s", err.Error())
		}

		if rt.StatusCode != int(HTTP_RESPONSE_UNAUTHORIZED) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_UNAUTHORIZED, rt.StatusCode)
		}

		if rt.GetHeader("content-type") != "application/json" {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/json", rt.GetHeader("content-type"))
		}

		tBody, _ := json.Marshal(unauthorizedAccessMessage)

		if rt.GetHeader("content-length") != strconv.Itoa(len(tBody)) {
			t.Fatalf("Expected content length to be (%d) but got (%s)", len(tBody), rt.GetHeader("content-length"))
		}

		body := make([]byte, len(tBody))

		n, err := rt.Body.Read(body)

		if err != nil {
			t.Fatalf("Something went wrong when trying to read body: %s", err.Error())
		}

		if n != len(tBody) {
			t.Fatalf("Expected body size to be (%d) but got (%d)", len(tBody), n)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected body size to be (%s) but got (%s)", string(tBody), string(body))
		}

		/**
		 *
		 * WITH KEY
		 *
		 */
		r = req.CreateRequest().Header("content-type", "application/json").Header("authorization", AuthKey)

		http, err = r.Post(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Something went wrong went trying to send request: %s", err.Error())
		}

		rt, err = HttpToResponse(http)

		if err != nil {
			t.Fatalf("Something went wrong went trying to convert http to response: %s", err.Error())
		}

		if rt.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_OK, rt.StatusCode)
		}

		if rt.GetHeader("content-type") != "application/json" {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/json", rt.GetHeader("content-type"))
		}

		tBody, _ = json.Marshal(userCreatedMessage)

		if rt.GetHeader("content-length") != strconv.Itoa(len(tBody)) {
			t.Fatalf("Expected content length to be (%d) but got (%s)", len(tBody), rt.GetHeader("content-length"))
		}

		body = make([]byte, len(tBody))

		n, err = rt.Body.Read(body)

		if err != nil {
			t.Fatalf("Something went wrong when trying to read body: %s", err.Error())
		}

		if n != len(tBody) {
			t.Fatalf("Expected body size to be (%d) but got (%d)", len(tBody), n)
		}

		if string(tBody) != string(body) {
			t.Fatalf("Expected body size to be (%s) but got (%s)", string(tBody), string(body))
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
		server.Session([]byte(str.Random(10)))

		r := req.CreateRequest()

		http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/authentication/login"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Something went wrong when trying to login: %s", err.Error())
		}

		res, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Something went wrong went trying convert http to response: %s", err.Error())
		}

		if res.StatusCode != int(HTTP_RESPONSE_TEMPORARY_REDIRECT) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", 307, res.StatusCode)
		}

		cookie, err := url.ParseQuery(strings.ReplaceAll(res.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatalf("Something went wrong when trying to convert set-cooke to query: %s", err.Error())
		}

		r = req.CreateRequest().Header("Cookie", strings.Join([]string{"session", cookie.Get("session")}, "="))

		http, err = r.Get(strings.Join([]string{"http://", server.Host(), "/dashboard"}, ""))

		if err != nil {
			t.Fatalf("Something went wrong when trying to get dashboard view: %s", err.Error())
		}

		res, err = HttpToResponse(http)

		if err != nil {
			t.Fatalf("Something went wrong went trying convert http to response: %s", err.Error())
		}

		if res.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_OK, res.StatusCode)
		}
	})

	err := server.Close()

	if err != nil {
		t.Fatalf("Something went wrong when trying to close server: %s", err.Error())
	}
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
	Message: "User was created successfully",
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

// Comment
func IsGuest(req *Request, res *Response, next Next) *Response {
	if req.Session.Get("user_id") != "" {
		return res.Redirect("/")
	}

	return next()
}

// Comment
func IsUser(req *Request, res *Response, next Next) *Response {
	if req.Session.Get("user_id") == "" {
		return res.Redirect("authentication/login")
	}

	return next()
}

// Comment
func IsAdmin(req *Request, res *Response, next Next) *Response {
	if req.Session.Get("role") != "1" {
		return res.Redirect("/")
	}

	return next()
}
