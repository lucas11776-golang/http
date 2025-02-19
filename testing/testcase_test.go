package testing

import (
	"testing"

	"github.com/lucas11776-golang/http"
)

func TestTestCaseRequest(t *testing.T) {
	testcase := NewTestCase(t, http.Server("127.0.0.1", 0), false)

	body := "<h1>Home Page</h1>"

	testcase.http.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.Html(body)
	})

	res := testcase.Request().Get("/")

	res.AssertOk()
	res.AssertHeader("content-type", "text/html")
	res.AssertBody([]byte(body))

	testcase.Cleanup()
}

func TestTestCaseWs(t *testing.T) {
	testcase := NewTestCase(t, http.Server("127.0.0.1", 0), false)

	testcase.Cleanup()
}
