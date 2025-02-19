package testing

import (
	"encoding/json"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
)

func TestRequest(t *testing.T) {

	t.Run("TestSetProtocolPathMethodHeadersBody", func(t *testing.T) {
		req := NewRequest(&TestCase{
			http: http.Server("127.0.0.1", 0),
		})

		body := `{"id":1,"email":"jeo@doe.com"}`

		req.Protocol("HTTP/1.1")
		req.setPath("api/products")
		req.setMethod(http.METHOD_POST)
		req.SetHeader("content-type", "application/json")
		req.SetHeaders(types.Headers{"user-agent": "Mozilla/5.0"})
		req.setBody([]byte(body))

		r, err := req.make()

		if err != nil {
			t.Fatalf("Something want wrong when trying to create request: %v", err)
		}

		if r.Protocol() != "HTTP/1.1" {
			t.Fatalf("Expected request protocol to be (%s) but got (%s)", "HTTP/1.1", r.Protocol())
		}

		if r.Path() != "api/products" {
			t.Fatalf("Expected request path to be (%s) but got (%s)", "api/products", r.Path())
		}

		if r.Method != string(http.METHOD_POST) {
			t.Fatalf("Expected request method to be (%s) but got (%s)", http.METHOD_POST, r.Method)
		}

		if r.GetHeader("content-type") != "application/json" {
			t.Fatalf("Expected request content-type to be (%s) but got (%s)", "application/json", r.GetHeader("content-type"))
		}

		if r.GetHeader("user-agent") != "Mozilla/5.0" {
			t.Fatalf("Expected request user-agent to be (%s) but got (%s)", "Mozilla/5.0", r.GetHeader("user-agent"))
		}

		tBody := make([]byte, len(body))

		_, err = r.Body.Read(tBody)

		if err != nil {
			t.Fatalf("Something want wrong when trying to read request body: %v", err)
		}

		if string(tBody) != body {
			t.Fatalf("Expected request body to be (%s) but got (%s)", body, tBody)
		}

		req.testCase.Cleanup()
	})
}

func TestRoute(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), false))

	user := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    1,
		Email: "jeo@doe.com",
	}

	req.testCase.http.Route().Get("users/{id}", func(req *http.Request, res *http.Response) *http.Response {
		return res.Json(user)
	})

	tBody, _ := json.Marshal(user)

	res := req.Json(http.METHOD_GET, "users/1", []byte{})

	res.AssertHeader("content-type", "application/json")
	res.AssertBody(tBody)

	req.testCase.Cleanup()
}

func TestSession(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), false))

	body := "<h1>Welcome to dashboard</h1>"

	isUser := func(req *http.Request, res *http.Response, next http.Next) *http.Response {
		if req.Session.Get("role") == "" {
			return res.Redirect("authentication/login")
		}

		return next()
	}

	isAdmin := func(req *http.Request, res *http.Response, next http.Next) *http.Response {
		if req.Session.Get("role") != "1" {
			return res.Redirect("/")
		}

		return next()
	}

	req.testCase.http.Route().Get("dashboard", func(req *http.Request, res *http.Response) *http.Response {
		return res.Html(body)
	}).Middleware(isUser, isAdmin)

	req.Sessions(map[string]string{
		"user_id": "1",
		"role":    "1",
	})

	res := req.Get("dashboard")

	res.AssertHeader("content-type", "text/html")
	res.AssertBody([]byte(body))

	req.testCase.Cleanup()
}

// Comment
func TestMultipartForm(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), false))

	value := &struct {
		name string
		data string
	}{
		name: "name",
		data: "Test Image",
	}

	file := &File{
		Name:     "picture",
		Filename: "image.jpeg",
		Type:     "image/jpeg",
		Data:     []byte{54, 34, 67, 46, 120, 255},
	}

	type response struct {
		Value string `json:"value"`
		File  string `json:"file"`
	}

	type message struct {
		Message string `json:"message"`
	}

	isUser := func(req *http.Request, res *http.Response, next http.Next) *http.Response {
		if req.Session.Get("user_id") == "" {
			return res.Redirect("/")
		}

		return next()
	}

	req.testCase.http.Route().Put("api/gallery", func(req *http.Request, res *http.Response) *http.Response {
		file, _, err := req.FormFile(file.Name)

		if err != nil {
			return res.SetStatus(http.HTTP_RESPONSE_UNPROCESSABLE_CONTENT).Json(message{
				Message: "The file is request",
			})
		}

		fileData, _ := io.ReadAll(file)

		return res.SetStatus(http.HTTP_RESPONSE_CREATED).Json(&response{
			Value: req.FormValue(value.name),
			File:  string(fileData),
		})
	}).Middleware(isUser)

	r := req.MultipartForm().Value(value.name, value.data).File(file.Name, file.Filename, file.Type, file.Data)

	r.Session("user_id", "1")

	res := r.Send(http.METHOD_PUT, "api/gallery")

	tBody, _ := json.Marshal(&response{
		Value: value.data,
		File:  string(file.Data),
	})

	res.AssertOk()
	res.AssertHeader("content-type", "application/json")
	res.AssertBody(tBody)

	req.testCase.Cleanup()
}

func TestFormUrlencoded(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), false))

	type response struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := &response{
		Email:    "jane@deo.com",
		Password: strings.Join([]string{"password", strconv.Itoa(int(rand.Float32() * 10000))}, "@"),
	}

	req.testCase.http.Route().Post("api/authentication/login", func(req *http.Request, res *http.Response) *http.Response {
		return res.Json(&response{
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		})
	})

	res := req.FormUrlencoded().Values(Values{
		"email":    body.Email,
		"password": body.Password,
	}).Send(http.METHOD_POST, "api/authentication/login")

	tBody, _ := json.Marshal(body)

	res.AssertOk()
	res.AssertHeader("content-type", "application/json")
	res.AssertBody(tBody)

	req.testCase.Cleanup()
}
