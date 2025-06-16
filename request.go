package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
	h "github.com/lucas11776-golang/http/utils/headers"
	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Method string

const (
	METHOD_GET     Method = "GET"
	METHOD_POST    Method = "POST"
	METHOD_PUT     Method = "PUT"
	METHOD_PATCH   Method = "PATCH"
	METHOD_DELETE  Method = "DELETE"
	METHOD_HEAD    Method = "HEAD"
	METHOD_OPTIONS Method = "OPTIONS"
	METHOD_CONNECT Method = "CONNECT"
)

type Request struct {
	*http.Request
	Conn       *connection.Connection
	Server     *HTTP
	Response   *Response
	Session    SessionManager
	Ws         *Ws
	Parameters Parameters
}

type HttpRequestHeader struct {
	method   string
	path     string
	protocol string
}

type HttpRequestContent struct {
	host    string
	headers types.Headers
	body    io.Reader
}

// Comment
func NewRequest(method Method, path string, protocol string, headers types.Headers, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(string(method), path, body)

	if err != nil {
		return nil, err
	}

	r.Header = h.ToHeader(headers)

	res := InitResponse()

	res.Request = &Request{
		Request:  r,
		Response: InitResponse(),
	}

	return res.Request, nil
}

// Comment
func (ctx *Request) Path() string {
	return strings.Trim(ctx.URL.Path, "/")
}

// Comment
func (ctx *Request) Protocol() string {
	return ctx.Proto
}

// Comment
func (ctx *Request) GetQuery(key string) string {
	return ctx.URL.Query().Get(key)
}

// Comment
func (ctx *Request) GetHeader(key string) string {
	header, ok := ctx.Header[cases.Title(language.English).String(key)]

	if !ok {
		return ""
	}

	return strings.Join(header, ";")
}

// Comment
func (ctx *Request) IP() string {
	return ctx.Conn.IP()
}

// Comment
func (ctx *Request) contentType() string {
	header := strings.Split(ctx.GetHeader("content-type"), ";")

	if len(header) == 0 {
		return ""
	}

	return header[0]
}

// Comment
func (ctx *Request) parseBodyX_WWW_FORM_URLENCODED() {
	buf := make([]byte, ctx.ContentLength)

	n, err := ctx.Body.Read(buf)

	if err != nil {
		return
	}

	ctx.Form, _ = url.ParseQuery(string(buf[:n]))
}

// Comment
func (ctx *Request) addJsonField(value interface{}, names []string) {
	for i, name := range names {
		if i > 0 {
			names[i] = fmt.Sprintf("[%s]", name)
		}
	}

	ctx.Form.Set(strings.Join(names, ""), cast.ToString(value))
}

// Comment
func (ctx *Request) addJsonFields(structure map[string]interface{}, names []string) {
	for k, v := range structure {
		switch v.(type) {
		case map[string]interface{}:
			ctx.addJsonFields(v.(map[string]interface{}), append(names, k))

		default:
			ctx.addJsonField(v, append(names, k))
		}
	}
}

// Comment
func (ctx *Request) parseBodyJson() {
	body, err := io.ReadAll(ctx.Body)

	if err != nil {
		return
	}

	jsonMap := map[string]interface{}{}

	err = json.Unmarshal(body, &jsonMap)

	if err != nil {
		return
	}

	if ctx.Form == nil {
		ctx.Form = url.Values{}
	}

	ctx.addJsonFields(jsonMap, []string{})
}

// Comment
func (ctx *Request) parseBody() {
	switch strings.ToLower(ctx.contentType()) {
	case "application/x-www-form-urlencoded":
		ctx.parseBodyX_WWW_FORM_URLENCODED()
	case "multipart/form-data":
		ctx.ParseMultipartForm(ctx.ContentLength)
	case "application/json":
		ctx.parseBodyJson()
	}
}

// Comment
func httpHeader(http []string) (*HttpRequestHeader, error) {
	header := strings.Split(http[0], " ")

	if len(header) != 3 {
		return nil, fmt.Errorf("invalid http header: %s", http[0])
	}

	return &HttpRequestHeader{
		method:   strings.ToUpper(header[0]),
		path:     header[1],
		protocol: strings.ToUpper(header[2]),
	}, nil
}

// Comment
func httpContent(http []string) (*HttpRequestContent, error) {
	content := HttpRequestContent{
		headers: make(types.Headers),
	}

	for i, line := range http[1:] {
		if line == "" {
			content.body = strings.NewReader(strings.TrimRight(strings.Join(http[i+2:], "\r\n"), "\r\n"))

			break
		}

		header := strings.Split(line, ":")

		if len(header) < 2 {
			return nil, fmt.Errorf("invalid header %s", header[0])
		}

		key := cases.Title(language.English).String(header[0])
		value := strings.Trim(strings.Join(header[1:], ","), " ")

		if key == "Host" {
			content.host = value

			continue
		}

		content.headers[key] = value
	}

	return &content, nil
}

// Comment
func ParseHttpRequest(http string) (*Request, error) {
	hp := strings.Split(http, "\r\n")

	header, err := httpHeader(hp)

	if err != nil {
		return nil, err
	}

	content, err := httpContent(hp)

	if err != nil {
		return nil, err
	}

	req, err := NewRequest(
		Method(header.method),
		header.path,
		header.protocol,
		content.headers,
		content.body,
	)

	if err != nil {
		return nil, err
	}

	req.Host = content.host

	req.parseBody()

	return req, nil
}
