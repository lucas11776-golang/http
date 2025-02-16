package http

import (
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
)

// Comment
func TestRequest(t *testing.T) {
	body := []byte("username=test123&password=pass1234")

	httpText := strings.Join([]string{
		"POST /authentication/login?redirect=cart&ref=lucas11776 HTTP/1.1",
		"HOST: example.com",
		"Content-Type: application/x-www-form-urlencoded",
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n",
		string(body) + "\r\n",
	}, "\r\n")

	req, _ := ParseHttpRequest(httpText)

	t.Run("TestMethod", func(t *testing.T) {
		if req.Method != "POST" {
			t.Fatalf("Expected method to be (%s) but got (%s)", "POST", req.Method)
		}
	})

	t.Run("TestPath", func(t *testing.T) {
		if req.Path() != "authentication/login" {
			t.Fatalf("Expected path to be (%s) but got (%s)", "authentication/login", req.Path())
		}
	})

	t.Run("TestProtocol", func(t *testing.T) {
		if req.Protocol() != "HTTP/1.1" {
			t.Fatalf("Expected protocol to be (%s) but got (%s)", "HTTP/1.1", req.Protocol())
		}
	})

	t.Run("TestQuery", func(t *testing.T) {
		if req.GetQuery("redirect") != "cart" {
			t.Fatalf("Expected query to be (%s) but got (%s)", "cart", req.GetQuery("redirect"))
		}

		if req.GetQuery("ref") != "lucas11776" {
			t.Fatalf("Expected query to be (%s) but got (%s)", "lucas11776", req.GetQuery("lucas11776"))
		}
	})

	t.Run("TestHost", func(t *testing.T) {
		if req.Host != "example.com" {
			t.Fatalf("Expected host to be (%s) but got (%s)", "example.com", req.Host)
		}
	})

	t.Run("TestHeaders", func(t *testing.T) {
		if req.GetHeader("content-type") != "application/x-www-form-urlencoded" {
			t.Fatalf("Expected content type to be (%s) but got (%s)", "application/x-www-form-urlencoded", req.GetHeader("content-type"))
		}

		if req.GetHeader("content-length") != strconv.Itoa(len(body)) {
			t.Fatalf("Expected content length to be (%s) but got (%s)", "6", req.GetHeader("content-length"))
		}
	})

	t.Run("TestParseBodyX_WWW_FORM_URLENCODED", func(t *testing.T) {
		if req.Form.Get("username") != "test123" {
			t.Fatalf("Expected username to be (%s) but go (%s)", "test123", req.Form.Get("username"))
		}

		if req.Form.Get("password") != "pass1234" {
			t.Fatalf("Expected username to be (%s) but go (%s)", "test123", req.Form.Get("username"))
		}
	})

	t.Run("TestParseBodyFormData", func(t *testing.T) {
		filename := "testing.txt"
		fileContent := "This is my TODO list"

		body := strings.Join([]string{
			`----------------------------392745197678564846379504`,
			`Content-Disposition: form-data; name="name"` + "\r\n",
			`My list`,
			`----------------------------392745197678564846379504`,
			`Content-Disposition: form-data; name="file"; filename="` + filename + `"`,
			`Content-Type: text/plain` + "\r\n",
			fileContent,
			`----------------------------392745197678564846379504--`,
		}, "\r\n")

		headers := types.Headers{
			"content-type":   "multipart/form-data; boundary=--------------------------392745197678564846379504",
			"content-length": strconv.Itoa(len(body)),
		}

		req, err := NewRequest("POST", "login", "HTTP/1.1", headers, strings.NewReader(body))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		if req.FormValue("name") != "My list" {
			t.Fatalf("Expected name to be (%s) but go (%s)", "My list", req.FormValue("name"))
		}

		file, header, err := req.FormFile("file")

		if err != nil {
			t.Fatalf("Something went wrong when trying get file: %s", err.Error())
		}

		if header.Filename != filename {
			t.Fatalf("Expected filename to be (%s) but go (%s)", filename, header.Filename)
		}

		if header.Header["Content-Type"][0] != "text/plain" {
			t.Fatalf("Expected file content type to be (%s) but go (%s)", "text/plain", header.Header["Content-Type"])
		}

		if header.Size != int64(len(fileContent)) {
			t.Fatalf("Expected file size to be (%d) but go (%d)", int64(len(fileContent)), header.Size)
		}

		content := make([]byte, header.Size)

		_, err = file.Read(content)

		if err != nil {
			t.Fatalf("Something went wrong when trying to read file: %s", err.Error())
		}

		if string(content) != fileContent {
			t.Fatalf("Expected file content to be (%s) but go (%s)", fileContent, string(content))
		}
	})
}
