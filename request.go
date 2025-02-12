package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
	h "github.com/lucas11776-golang/http/utils/headers"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Request struct {
	*http.Request
	Conn     *connection.Connection
	Server   *server.Server
	Response *Response
}

type HttpRequestHeader struct {
	method   string
	path     string
	protocol string
}

type HttpRequestContent struct {
	host    string
	headers types.Headers
	body    []byte
}

// Comment
func NewRequest(method string, path string, protocol string, headers types.Headers, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(method, path, body)

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
func httpHeader(http []string) (*HttpRequestHeader, error) {
	header := strings.Split(http[0], " ")

	if len(header) != 3 {
		return nil, fmt.Errorf("Invalid http header: %s", http[0])
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
			content.body = []byte(strings.Trim(strings.Join(http[i:], "\r\n"), "\r\n"))
			break
		}

		header := strings.Split(line, ":")

		if len(header) < 2 {
			return nil, fmt.Errorf("Invalid header %s", header[0])
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
		header.method,
		header.path,
		header.protocol,
		content.headers,
		bytes.NewReader(content.body),
	)

	if err != nil {
		return nil, err
	}

	req.Host = content.host

	err = req.ParseForm()

	if err != nil {
		return nil, err
	}

	return req, nil
}
