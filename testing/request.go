package testing

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	h "net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/headers"
)

type File struct {
	Name     string
	Filename string
	Type     string
	Data     []byte
}

type Values map[string]string

type Files []*File

type RequestReadCloser struct {
	io.Reader
}

type Request struct {
	TestCase *TestCase
	Testing  *Testing
	Request  *http.Request
	session  Values
	protocol string
	path     string
	method   http.Method
	headers  types.Headers
	body     []byte
}

type MultiPartForm struct {
	request  *Request
	values   Values
	files    Files
	boundary string
}

// Comment
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
	ctx.SetHeader("content-length", strconv.Itoa(len(ctx.body)))

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

// Comment
func (ctx *Request) MultipartForm() *MultiPartForm {
	return &MultiPartForm{
		request:  ctx,
		values:   make(Values),
		boundary: fmt.Sprintf("--------------------------%d", 100*(1000+rand.Int()*9999)),
	}
}

// comment
func (ctx *MultiPartForm) File(name string, filename string, contentType string, data []byte) *MultiPartForm {
	ctx.files = append(ctx.files, &File{
		Name:     name,
		Filename: filename,
		Type:     contentType,
		Data:     data,
	})

	return ctx
}

// Comment
func (ctx *MultiPartForm) Value(name string, value string) *MultiPartForm {
	ctx.values[name] = value

	return ctx
}

// Comment
func (ctx *MultiPartForm) Session(key string, value string) *MultiPartForm {
	ctx.request.Session(key, value)

	return ctx
}

// Comment
func (ctx *MultiPartForm) Sessions(sessions Values) *MultiPartForm {
	ctx.request.Sessions(sessions)

	return ctx
}

// Comment
func (ctx *MultiPartForm) value(name string, value string) string {
	return fmt.Sprintf(
		strings.Join([]string{
			"Content-Disposition: form-data; name=\"%s\"\r\n",
			value,
		}, "\r\n"), name,
	)
}

// Comment
func (ctx *MultiPartForm) file(file *File) string {
	return fmt.Sprintf(
		strings.Join([]string{
			"Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"",
			"Content-Type: %s\r\n",
			string(file.Data),
		}, "\r\n"), file.Name, file.Filename, file.Type,
	)
}

// Comment
func (ctx *MultiPartForm) body() string {
	if len(ctx.values) == 0 && len(ctx.files) == 0 {
		return ""
	}

	body := []string{}

	for name, value := range ctx.values {
		body = append(body, fmt.Sprintf("--%s", ctx.boundary), ctx.value(name, value))
	}

	for _, file := range ctx.files {
		body = append(body, fmt.Sprintf("--%s", ctx.boundary), ctx.file(file))
	}

	return strings.Join(append(body, fmt.Sprintf("--%s--", ctx.boundary)), "\r\n")
}

// Comment
func (ctx *MultiPartForm) Send(method http.Method, uri string) *Response {
	switch method {
	case http.METHOD_POST, http.METHOD_PUT, http.METHOD_PATCH:
		ctx.request.SetHeader("content-type", fmt.Sprintf("multipart/form-data; boundary=%s", ctx.boundary))

		return ctx.Value("__METHOD__", string(method)).request.Call(http.METHOD_POST, uri, []byte(ctx.body()))
	default:
		ctx.request.Testing.Fatalf(
			"Multipart form does not support (%s) it only support (%v)",
			method,
			[]http.Method{http.METHOD_POST, http.METHOD_PUT, http.METHOD_PATCH},
		)
		return nil
	}
}

// Comment
type FormUrlencoded struct {
	request *Request
	values  Values
}

// Comment
func (ctx *Request) FormUrlencoded() *FormUrlencoded {
	return &FormUrlencoded{
		request: ctx,
		values:  make(Values),
	}
}

// Comment
func (ctx *FormUrlencoded) Value(key string, value string) *FormUrlencoded {
	ctx.values[key] = value

	return ctx
}

// Comment
func (ctx *FormUrlencoded) Values(values Values) *FormUrlencoded {
	for k, v := range values {
		ctx.Value(k, v)
	}

	return ctx
}

// Comment
func (ctx *FormUrlencoded) body() string {
	query := []string{}

	for v, k := range ctx.values {
		query = append(query, strings.Join([]string{v, k}, "="))
	}

	return strings.Join(query, "&")
}

// Comment
func (ctx *FormUrlencoded) Send(method http.Method, uri string) *Response {
	ctx.request.SetHeader("content-type", "application/x-www-form-urlencoded")

	return ctx.request.Call(method, uri, []byte(ctx.body()))
}
