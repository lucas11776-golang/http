package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/reader"
	req "github.com/lucas11776-golang/http/utils/request"
	"github.com/lucas11776-golang/http/utils/response"
	"github.com/lucas11776-golang/http/utils/rsa"
	str "github.com/lucas11776-golang/http/utils/strings"
	"github.com/lucas11776-golang/http/validation"
	"github.com/open2b/scriggo"
	"github.com/quic-go/quic-go/http3"
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

	t.Run("TestFormValidation", func(t *testing.T) {
		server := Server("127.0.0.1", 0).ParseJson(true)
		sessions := server.Session([]byte(str.Random(10)))
		emailErrMsg := "The email is required"

		LoginRequest := func() Middleware {
			return FormRequest(validation.RulesBag{
				"email": validation.Rules{"required"},
			})
		}

		server.Route().Group("authentication", func(route *Router) {
			route.Group("login", func(route *Router) {
				route.Post("/", func(req *Request, res *Response) *Response {
					return res.Redirect("dashboard")
				}, LoginRequest())
				route.Get("/", func(req *Request, res *Response) *Response {
					return res.Html("Login Veiw")
				})
			})
		})

		go server.Listen()

		t.Run("TestWebRequest", func(t *testing.T) {
			r := req.CreateRequest().
				SetHeaders(types.Headers{
					"content-type": "application/x-www-form-urlencoded",
					"host":         "127.0.0.1:4567",
				})

			http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/authentication/login"}, ""), []byte{})

			if err != nil {
				t.Fatal(err)
			}

			_, _, h, _, err := response.ParseHttpToResponse(http)

			if err != nil {
				t.Fatal(err)
			}

			cookie, err := url.ParseQuery(strings.ReplaceAll(h.Get("Set-Cookie"), "; ", "&"))

			if err != nil {
				t.Fatal(err)
			}

			// Second Request
			headers := types.Headers{
				"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
			}

			req, err := NewRequest("GET", "/authentication/login", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

			if err != nil {
				t.Fatal(err)
			}

			session := sessions.Session(req)

			if msg := session.Error("email"); msg != emailErrMsg {
				t.Fatalf("Expected email to be (%s) but got (%s)", emailErrMsg, msg)
			}
		})

		t.Run("TestWebRequest", func(t *testing.T) {
			r := req.CreateRequest().
				SetHeaders(types.Headers{
					"content-type": "application/json",
					"host":         "127.0.0.1:4567",
				})

			http, err := r.Post(strings.Join([]string{"http://", server.Host(), "/authentication/login"}, ""), []byte{})

			if err != nil {
				t.Fatal(err)
			}

			_, statusCode, _, body, err := response.ParseHttpToResponse(http)

			if err != nil {
				t.Fatal(err)
			}

			if statusCode != int(HTTP_RESPONSE_UNPROCESSABLE_CONTENT) {
				t.Fatalf("Expected response statuts code to be (%d) but got (%d)", HTTP_RESPONSE_UNPROCESSABLE_CONTENT, statusCode)
			}

			var jsonErrorResponse JsonErrorResponse

			if err := json.Unmarshal(body, &jsonErrorResponse); err != nil {
				t.Fatal(err)
			}

			if jsonErrorResponse.Message != FormValidationErrorMessage {
				t.Fatalf("Expected json error message to be (%s) but got (%s)", FormValidationErrorMessage, jsonErrorResponse.Message)
			}

			if msg := jsonErrorResponse.Errors["email"]; msg != emailErrMsg {
				t.Fatalf("Expected email error to be (%s) but got (%s)", emailErrMsg, msg)
			}
		})

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

		// Comment
		IsAuth := func(req *Request, res *Response, next Next) *Response {
			if req.Session.Get("user_id") == "" {
				return res.Redirect("/dashboar")
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

		server.Route().Group("dashboard", func(route *Router) {
			route.Get("/", func(req *Request, res *Response) *Response {
				return res.Html("Dashboard")
			})
		}, IsAuth)

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

		cookie, err := url.ParseQuery(strings.ReplaceAll(res.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatalf("Something went wrong when trying to convert set-cooke to query: %s", err.Error())
		}

		r = req.CreateRequest().SetHeader("Cookie", strings.Join([]string{"session", cookie.Get("session")}, "="))

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

		server.Close()
	})

	t.Run("TestHTTP3", func(t *testing.T) {
		certFile := "./temp/host.cert"
		keyFile := "./temp/host.key"

		cert, err := os.Create(certFile)

		if err != nil {
			t.Fatal(err)
		}

		key, err := os.Create(keyFile)

		if err != nil {
			t.Fatal(err)
		}

		certData, keyData, err := rsa.GenerateCertificate("127.0.0.1")

		if err != nil {
			t.Fatal(err)
		}

		if _, err := cert.Write([]byte(certData)); err != nil {
			t.Fatal(err)
		}

		if _, err := key.Write([]byte(keyData)); err != nil {
			t.Fatal(err)
		}

		server := ServerTLS("127.0.0.1", 0, certFile, keyFile)

		responseBody := fmt.Sprintf("Hello HTTP3: %f", rand.Float32()*10000000)

		server.Route().Get("/", func(req *Request, res *Response) *Response {
			return res.SetBody([]byte(responseBody))
		})

		go server.Listen()
		client := &http.Client{Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

		req, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			fmt.Sprintf("https://%s", server.Host()),
			nil,
		)

		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		// Execute request
		res, err := client.Do(req)

		if err != nil {
			log.Fatalf("HTTP/3 request failed: %v", err)
		}

		defer res.Body.Close()

		// Read and print the response
		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}

		if string(body) != responseBody {
			t.Fatalf("Expected response body to be (%s) but got (%s)", responseBody, string(body))
		}

		server.Close()
	})
}
