package testing

import (
	"bytes"
	"io"
	h "net/http"
	"net/url"
	"strings"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/headers"
)

type File struct {
	Name        string
	Filename    string
	ContentType string
	Content     []byte
}

type Values map[string]string

type Files []Files

type RequestReadCloser struct {
	io.Reader
}

type Request struct {
	TestCase *TestCase
	Testing  *Testing
	Request  *http.Request
	values   Values
	files    Files
	session  Values
	protocol string
	path     string
	method   http.Method
	headers  types.Headers
	body     []byte
}

func (ctx *RequestReadCloser) Close() error {
	return nil
}

// Comment
func NewRequest(testcase *TestCase) *Request {
	req := &Request{
		TestCase: testcase,
		Testing:  testcase.Testing,
		protocol: "HTTP/1.1",
		method:   "GET",
		headers:  make(types.Headers),
		values:   make(Values),
		session:  make(Values),
	}

	req.Request, _ = req.make()

	return req
}

// Comment
func (ctx *Request) Protocol(protocol string) *Request {
	return ctx
}

// Comment
func (ctx *Request) setPath(path string) *Request {
	ctx.path = path

	return ctx
}

// Comment
func (ctx *Request) setMethod(method http.Method) *Request {
	ctx.method = method

	return ctx
}

// Comment
func (ctx *Request) SetHeader(key string, value string) *Request {
	ctx.headers[key] = value

	return ctx
}

// Comment
func (ctx *Request) SetHeaders(headers types.Headers) *Request {
	for k, v := range headers {
		ctx.SetHeader(k, v)
	}

	return ctx
}

// Comment
func (ctx *Request) setBody(body []byte) *Request {
	ctx.body = body

	return ctx
}

// Comment
func (ctx *Request) make() (*http.Request, error) {
	r, err := h.NewRequest(
		string(ctx.method),
		ctx.path,
		bytes.NewReader(ctx.body),
	)

	if err != nil {
		return nil, err
	}

	req := ctx.TestCase.HTTP.NewRequest(r, nil)

	req.Proto = ctx.protocol
	req.Header = headers.ToHeader(ctx.headers)

	return req, nil
}

// Comment
func (ctx *Request) addSessionHeader(req *http.Request) *http.Request {
	if len(ctx.session) == 0 {
		return req
	}

	r, err := h.NewRequest("GET", "/", bytes.NewReader([]byte{}))

	if err != nil {
		ctx.Testing.Fatalf("Something went wrong when trying to create request for session: %v", err)
	}

	rq := ctx.TestCase.HTTP.NewRequest(r, nil)
	session := ctx.TestCase.HTTP.Get("session").(http.SessionsManager).Session(rq)

	for k, v := range ctx.session {
		session.Set(k, v)
	}

	session.Save()

	cookie, err := url.ParseQuery(strings.ReplaceAll(rq.Response.GetHeader("set-cookie"), ";", "&"))

	if err != nil {
		ctx.Testing.Fatalf("Something went wrong when trying to convert session to query: %v", err)
	}

	req.Header["Cookie"] = []string{strings.Join([]string{http.SESSION_NAME, cookie.Get(http.SESSION_NAME)}, "=")}

	return req
}

// Comment
func (ctx *Request) makeRequest(req *http.Request) *Response {
	ctx.Request = req

	res := ctx.TestCase.HTTP.HandleRequest(ctx.addSessionHeader(req))

	if res == nil {
		ctx.TestCase.Testing.Fatalf("Request does not support WebSocket request use Ws testing")
	}

	req.Response = res

	return NewResponse(ctx, res)
}

// Comment
func (ctx *Request) Call(method http.Method, uri string, body []byte) *Response {
	ctx.setMethod(method)
	ctx.setPath(uri)
	ctx.setBody(body)

	req, err := ctx.make()

	if err != nil {
		ctx.TestCase.Testing.Fatalf("Something went wrong when create request: %v", err)
	}

	return ctx.makeRequest(req)
}

// Comment
func (ctx *Request) Session(key string, value string) *Request {
	ctx.session[key] = value

	return ctx
}

// Comment
func (ctx *Request) Sessions(sessions Values) *Request {
	for k, v := range sessions {
		ctx.Session(k, v)
	}

	return ctx
}

// Comment
func (ctx *Request) Get(uri string) *Response {
	return ctx.Call(http.METHOD_GET, uri, []byte{})
}

// Comment
func (ctx *Request) Post(uri string, body []byte) *Response {
	return ctx.Call(http.METHOD_POST, uri, body)
}

// Comment
func (ctx *Request) Put(uri string, body []byte) *Response {
	return ctx.Call(http.METHOD_PUT, uri, body)
}

// Comment
func (ctx *Request) Patch(uri string, body []byte) *Response {
	return ctx.Call(http.METHOD_PATCH, uri, body)
}

// Comment
func (ctx *Request) Delete(uri string) *Response {
	return ctx.Call(http.METHOD_DELETE, uri, []byte{})

}

// Comment
func (ctx *Request) Head(uri string) *Response {
	return ctx.Call(http.METHOD_HEAD, uri, []byte{})
}

// Comment
func (ctx *Request) Options(uri string) *Response {
	return ctx.Call(http.METHOD_OPTIONS, uri, []byte{})
}

// Comment
func (ctx *Request) Connect(uri string, body []byte) *Response {
	return ctx.Call(http.METHOD_CONNECT, uri, body)
}

// Comment
func (ctx *Request) Json(method http.Method, uri string, body []byte) *Response {
	return ctx.Call(method, uri, body)
}

type MultiPartForm struct {
	request *Request
	values  Values
	files   Files
}

// Comment
func (ctx *Request) MultiPartForm() *MultiPartForm {
	return &MultiPartForm{request: ctx}
}

// comment
func (ctx *MultiPartForm) File(name string, filename string, contentType string, data []byte) *MultiPartForm {
	return ctx
}

// Comment
func (ctx *MultiPartForm) Value(name string, value string) *MultiPartForm {
	return ctx
}
