package testing

import (
	"fmt"
	"io"
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TestingT struct {
	T      *testing.T
	catch  bool
	errors []string
}

type Response struct {
	TestCase *TestCase
	Request  *Request
	Response *http.Response
	Testing  *TestingT
}

// Comment
func InitTestingResponse(t *testing.T, catchErrors bool) *TestingT {
	return &TestingT{T: t, catch: catchErrors}
}

// Comment
func InitResponse(req *Request, res *http.Response, catchErrors bool) *Response {
	return &Response{
		TestCase: req.TestCase,
		Request:  req,
		Response: res,
		Testing:  InitTestingResponse(req.TestCase.Testing, catchErrors),
	}
}

// Comment
func (ctx *Response) AssertProtocol(protocol string) *Response {
	if ctx.Response.Protocol() != protocol {
		ctx.Testing.Fatalf("Expected response protocol to be (%s) but got (%s)", ctx.Response.Protocol(), protocol)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertStatusCode(status http.Status) *Response {
	if ctx.Response.StatusCode != int(status) {
		ctx.Testing.Fatalf("Expected response status code to be (%d) but got (%d)", status, ctx.Response.StatusCode)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertOk() *Response {
	if !(ctx.Response.StatusCode >= int(http.HTTP_RESPONSE_OK) && ctx.Response.StatusCode <= int(http.HTTP_RESPONSE_ACCEPTED)) {
		ctx.Testing.Fatalf("Expected response status code to be (200, 201, 202) but got (%d)", ctx.Response.StatusCode)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertNotFound() *Response {
	return ctx.AssertStatusCode(http.HTTP_RESPONSE_NOT_FOUND)
}

// Comment
func (ctx *Response) AssertUnauthorized() *Response {
	return ctx.AssertStatusCode(http.HTTP_RESPONSE_UNAUTHORIZED)
}

// Comment
func (ctx *Response) AssertHasHeader(header string) *Response {
	_, ok := ctx.Response.Header[cases.Title(language.English).String(header)]

	if !ok {
		ctx.Testing.Fatalf("Expected response to contain header (%s)", header)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertHeader(header string, value string) *Response {
	if ctx.Response.GetHeader(header) != value {
		ctx.Testing.Fatalf("Expected response header (%s) to be (%s) but got (%s)", header, value, ctx.Response.GetHeader(header))
	}

	return ctx
}

// Comment
func (ctx *Response) AssertHeaders(headers types.Headers) *Response {
	for k, v := range headers {
		ctx.AssertHeader(k, v)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertBody(body []byte) *Response {
	tBody, err := io.ReadAll(ctx.Response.Body)

	if err != nil {
		ctx.Testing.Fatalf("Something went wrong when trying to read body: %v", err)
	}

	if string(tBody) != string(body) {
		ctx.Testing.Fatalf("Expected response body to be (%s) but got (%s)", string(body), string(body))
	}

	return ctx
}

// Comment
func (ctx *Response) AssertIsRedirect() *Response {
	if ctx.Response.RedirectBag == nil {
		ctx.Testing.Fatalf("Expected response to be redirect")
	}

	return ctx
}

// Comment
func (ctx *TestingT) Fatalf(format string, args ...any) {
	if ctx.catch {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.T.Fatalf(format, args...)
}

// Comment
func (ctx *TestingT) hasError() bool {
	return len(ctx.errors) > 0
}

// Comment
func (ctx *TestingT) popError() string {
	if !ctx.hasError() {
		return ""
	}

	err := ctx.errors[len(ctx.errors)-1]

	ctx.errors = ctx.errors[:len(ctx.errors)-1]

	return err
}
