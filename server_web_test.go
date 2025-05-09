package http

import (
	"encoding/json"
	"io"
	"math/rand"
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
	var AuthKey = "KEY-" + strconv.Itoa(int(rand.Float32()*10000))

	type Message struct {
		Message string `json:"message"`
	}

	type User struct {
		ID    int64  `json:"id"`
		Role  byte   `json:"role"`
		Email string `json:"email"`
	}

	users := []User{
		(User{ID: 1, Role: 1, Email: "jane@doe.com"}),
		(User{ID: 2, Role: 0, Email: "jeo@doe.com"}),
	}

	t.Run("TestApiGetUsers", func(t *testing.T) {
		server := Server("127.0.0.1", 0)

		server.Route().Group("api", func(route *Router) {
			route.Group("users", func(route *Router) {
				route.Get("/", func(req *Request, res *Response) *Response {
					return res.Json(users)
				})
			})
		})

		go server.Listen()

		r := req.CreateRequest().
			SetHeaders(types.Headers{
				"content-type": "application/json",
				"host":         server.Host(),
			})

		http, err := r.Get(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""))

		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		res, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Failed to parse http: %v", err)
		}

		if res.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_OK, res.StatusCode)
		}

		if res.GetHeader("Content-Type") != "application/json" {
			t.Fatalf("Expected header content-type to be (%s) but got (%s)", "application/json", res.GetHeader("Content-Type"))
		}

		tBody, _ := json.Marshal(users)
		body, _ := io.ReadAll(res.Body)

		if string(tBody) != string(body) {
			t.Fatalf("Expected body to be (%s) but got (%s)", string(tBody), string(body))
		}

		server.Close()
	})

	t.Run("TestMiddlewareUserPost", func(t *testing.T) {
		var unauthorizedMessage = Message{
			Message: "Authorization key is invalid",
		}

		var createdMessage = Message{
			Message: "Authorization key is invalid",
		}

		auth := func(req *Request, res *Response, next Next) *Response {
			if req.GetHeader("authorization") != AuthKey {
				return res.SetStatus(HTTP_RESPONSE_UNAUTHORIZED).Json(unauthorizedMessage)
			}
			return next()
		}

		server := Server("127.0.0.1", 0)

		server.Route().Group("api", func(route *Router) {
			route.Group("users", func(route *Router) {
				route.Post("/", func(req *Request, res *Response) *Response {
					return res.SetStatus(HTTP_RESPONSE_OK).Json(createdMessage)
				}).Middleware(auth)
			})
		})

		go server.Listen()

		// ------------------------------ WITHOUT KEY ------------------------------ //

		r := req.CreateRequest().
			SetHeaders(types.Headers{
				"content-type": "application/json",
				"host":         "127.0.0.1:4567",
			})

		http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Failed to send request: %s", err.Error())
		}

		res, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Failed to parse http: %v", err)
		}

		if res.StatusCode != int(HTTP_RESPONSE_UNAUTHORIZED) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_UNAUTHORIZED, res.StatusCode)
		}

		tBody, _ := json.Marshal(unauthorizedMessage)
		body, _ := io.ReadAll(res.Body)

		if string(tBody) != string(body) {
			t.Fatalf("Expected body to be (%s) but got (%s)", string(tBody), string(body))
		}

		// ----------------------------------- WITHOUT KEY ----------------------------------- //

		r = req.CreateRequest().
			SetHeaders(types.Headers{
				"content-type":  "application/json",
				"host":          "127.0.0.1:4567",
				"authorization": AuthKey,
			})

		http, err = r.Post(strings.Join([]string{"http://", server.Host(), "/api/users"}, ""), []byte{})

		if err != nil {
			t.Fatalf("Failed to send request: %s", err.Error())
		}

		res, err = HttpToResponse(http)

		if err != nil {
			t.Fatalf("Failed to parse http: %v", err)
		}

		if res.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_OK, res.StatusCode)
		}

		if res.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_OK, res.StatusCode)
		}

		tBody, _ = json.Marshal(createdMessage)
		body, _ = io.ReadAll(res.Body)

		if string(tBody) != string(body) {
			t.Fatalf("Expected body to be (%s) but got (%s)", string(tBody), string(body))
		}

		server.Close()
	})

	t.Run("TestStatic", func(t *testing.T) {
		var (
			fileName = "assets/css/main.css"
			tBody    = strings.Join([]string{
				"body { margin: 0px !important; padding: 0px !important; background-color: green; }",
			}, "\r\n")
		)

		server := Server("127.0.0.1", 0)

		server.Set("static", InitStatic(reader.NewTestingReader(scriggo.Files{
			fileName: []byte(tBody),
		})))

		go server.Listen()

		r := req.CreateRequest().
			SetHeaders(types.Headers{
				"content-type": "application/json",
				"host":         "127.0.0.1:4567",
			})

		http, err := r.Get(strings.Join([]string{"http://", server.Host(), "/", fileName}, ""))

		if err != nil {
			t.Fatalf("Failed to send request: %s", err.Error())
		}

		res, err := HttpToResponse(http)

		if err != nil {
			t.Fatalf("Failed to parse http: %v", err)
		}

		if res.StatusCode != int(HTTP_RESPONSE_OK) {
			t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_UNAUTHORIZED, res.StatusCode)
		}

		body, _ := io.ReadAll(res.Body)

		if string(tBody) != string(body) {
			t.Fatalf("Expected body to be (%s) but got (%s)", string(tBody), string(body))
		}

		server.Close()
	})

	t.Run("TestSession", func(t *testing.T) {
		server := Server("127.0.0.1", 0)

		server.Session([]byte(str.Random(10)))

		// Comment
		IsGuest := func(req *Request, res *Response, next Next) *Response {
			if req.Session.Get("user_id") != "" {
				return res.Redirect("/")
			}

			return next()
		}

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

		go server.Listen()

		r := req.CreateRequest().SetHeader("host", "127.0.0.1:4567")

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

		// cookie, err := url.ParseQuery(strings.ReplaceAll(res.GetHeader("Set-Cookie"), "; ", "&"))

		// if err != nil {
		// 	t.Fatalf("Something went wrong when trying to convert set-cooke to query: %s", err.Error())
		// }

		// r = req.CreateRequest().SetHeader("Cookie", strings.Join([]string{"session", cookie.Get("session")}, "="))

		// http, err = r.Get(strings.Join([]string{"http://", server.Host(), "/dashboard"}, ""))

		// if err != nil {
		// 	t.Fatalf("Something went wrong when trying to get dashboard view: %s", err.Error())
		// }

		// res, err = HttpToResponse(http)

		// if err != nil {
		// 	t.Fatalf("Something went wrong went trying convert http to response: %s", err.Error())
		// }

		// if res.StatusCode != int(HTTP_RESPONSE_OK) {
		// 	t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_OK, res.StatusCode)
		// }

		server.Close()
	})

	t.Run("TestHTTP3", func(t *testing.T) {
		server := Server("127.0.0.1", 0)

		server.Route().Get("/", func(req *Request, res *Response) *Response {
			return res.SetBody([]byte(req.Proto))
		})

		go server.Listen()

		// url := "https://cloudflare-quic.com" // Make sure the server supports HTTP/3

		// // Create HTTP/3 Transport
		// transport := &http3.Transport{}

		// // Create a custom client using the HTTP/3 transport
		// client := &http.Client{
		// 	Transport: transport,
		// }

		// // Create request (with context)
		// req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		// if err != nil {
		// 	log.Fatalf("Failed to create request: %v", err)
		// }

		// // Execute request
		// resp, err := client.Do(req)
		// if err != nil {
		// 	log.Fatalf("HTTP/3 request failed: %v", err)
		// }
		// defer resp.Body.Close()

		// // Read and print the response
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Fatalf("Failed to read response: %v", err)
		// }

		// fmt.Println("Response Status:", resp.Status)
		// fmt.Println("Response Body:")
		// fmt.Println(string(body))

		server.Close()
	})
}
