package server

import (
	"http/request"
	"http/response"
	"http/router"
	"http/types"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

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

var authKey = "KEY-" + strconv.Itoa(int(rand.Float32()*10000))

// Comment
func AuthKey(req *request.Request, res *response.Response, next router.Next) *response.Response {
	if req.Header("auth-key") != authKey {
		return res.Status(response.HTTP_RESPONSE_UNAUTHORIZED).Json(userCreatedMessage)
	}

	return next()
}

func TestServerWeb(t *testing.T) {
	users := []User{
		(User{ID: 1, Email: "jane@doe.com"}),
		(User{ID: 2, Email: "jeo@doe.com"}),
	}

	machine, err := Serve("127.0.0.1", 0)

	if err != nil {
		t.Fatalf("Something went wrong when trying to create to server: %s", err.Error())
	}

	machine.Router().Group("/", func(route *router.Router) {
		route.Get("/", func(req *request.Request, res *response.Response) *response.Response {
			return res
		})
	})

	machine.Router().Group("api", func(route *router.Router) {
		route.Group("users", func(route *router.Router) {
			route.Get("/", func(req *request.Request, res *response.Response) *response.Response {
				return res.Json(users)
			})
			route.Post("/", func(req *request.Request, res *response.Response) *response.Response {
				return res.Status(response.HTTP_RESPONSE_OK).Json(userCreatedMessage)
			}).Middleware(AuthKey)
		})
	})

	go func() {
		machine.Listen()
	}()

	t.Run("TestApiGetUsers", func(t *testing.T) {

		req := CreateRequest().Header("Content-Type", "application/json")

		http, err := req.Get("http://" + machine.listener.Addr().String() + "/api/users")

		if err != nil {
			t.Fatalf("Something went wrong went trying to send request: %s", err.Error())
		}

		res := response.Create("HTTP/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
		expectedHttp := response.ParseHttp(res.Json(users))

		if expectedHttp != http {
			t.Fatalf("Expected response to be (%s) but got (%s), (%d,%d)", expectedHttp, http, len(expectedHttp), len(http))
		}
	})

	machine.Close()
}

type Request struct {
	method  string
	headers types.Headers
	data    []byte
}

func CreateRequest() *Request {
	return &Request{
		headers: make(types.Headers),
	}
}

// Comment
func (ctx *Request) Header(key string, value string) *Request {
	ctx.headers[key] = value

	return ctx
}

// Comment
func (ctx *Request) Post(url string, data []byte) (string, error) {
	return ctx.Request("POST", url, data)
}

// Comment
func (ctx *Request) Get(url string) (string, error) {
	return ctx.Request("GET", url, []byte{})
}

func (ctx *Request) parse(method string, path string, data []byte) string {
	arr := []string{
		strings.Join([]string{strings.ToUpper(method), path, "HTTP/1.1"}, " "),
	}

	for key, value := range ctx.headers {
		arr = append(arr, strings.Join([]string{key, value}, ": "))
	}

	if len(data) == 0 {
		return strings.Join(append(arr, "\r\n"), "\r\n")
	}

	arr = append(arr, strings.Join([]string{"Content-Length", strconv.Itoa(len(data))}, ": "))

	return strings.Join(append(arr, strings.Join([]string{"\r\n", string(data), "\r\n"}, "")), "\r\n")
}

// Comment
func (ctx *Request) Request(method string, address string, data []byte) (string, error) {
	url, err := url.Parse(address)

	if err != nil {
		return "", err
	}

	listener, err := net.Dial("tcp", url.Host)

	if err != nil {
		return "", err
	}

	_, err = listener.Write([]byte(ctx.parse(method, url.Path, data)))

	if err != nil {
		return "", err
	}

	err = listener.SetDeadline(time.Now().Add(time.Second * 3))

	if err != nil {
		return "", err
	}

	http := make([]byte, MAX_REQUEST_SIZE)

	n, err := listener.Read(http)

	if err != nil {
		return "", err
	}

	listener.Close()

	return string(http[:n]), nil

	// return ""
}
