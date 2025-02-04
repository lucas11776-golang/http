package request

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/types"
)

type Request struct {
	*http.Request
	method   string
	path     string
	query    types.Query
	protocol string
	headers  types.Headers
	body     []byte
	Conn     *net.Conn
	Server   *server.Server
}

type headerInfo struct {
	method   string
	path     string
	query    types.Query
	protocol string
}

// Comment
func Create(method string, path string, query types.Query, protocol string, headers types.Headers, body []byte) *Request {
	return &Request{
		method:   strings.ToUpper(method),
		path:     strings.Trim(path, "/"),
		query:    query,
		protocol: protocol,
		headers:  headers,
		body:     body,
	}
}

// Comment
func getHeaderInfo(header string) (*headerInfo, error) {
	arr := strings.Split(header, " ")

	if len(arr) != 3 {
		return nil, fmt.Errorf("Invalid http header: %s", header)
	}

	u, err := url.Parse(strings.Trim(arr[1], "/"))

	if err != nil {
		return nil, err
	}

	query := make(types.Query)

	for k := range u.Query() {
		query[k] = u.Query().Get(k)
	}

	return &headerInfo{
		method:   strings.ToUpper(arr[0]),
		path:     strings.Trim(u.Path, "/"),
		protocol: strings.ToUpper(arr[2]),
		query:    query,
	}, nil
}

type content struct {
	headers types.Headers
	body    []byte
}

// Comment
func getContent(arr []string) (*content, error) {
	content := content{
		headers: make(types.Headers),
	}

	for i, line := range arr {
		if line == "" {
			content.body = []byte(strings.Trim(strings.Join(arr[i:], "\r\n"), "\r\n"))
			break
		}

		header := strings.Split(line, ":")

		if len(header) < 2 {
			return nil, fmt.Errorf("Invalid header %s", header[0])
		}

		content.headers[strings.ToLower(header[0])] = strings.Trim(strings.Join(header[1:], ":"), " ")
	}

	return &content, nil
}

// Comment
func ParseRequest(request *http.Request) *Request {

	return &Request{}
}

// Comment
func ParseHttp(http string) (*Request, error) {
	arr := strings.Split(http, "\r\n")
	header, err := getHeaderInfo(arr[0:1][0])

	if err != nil {
		return nil, err
	}

	content, err := getContent(arr[1:])

	if err != nil {
		return nil, err
	}

	return &Request{
		method:   header.method,
		path:     header.path,
		protocol: header.protocol,
		query:    header.query,
		headers:  content.headers,
		body:     content.body,
	}, nil
}

// Comment
func (ctx *Request) Method() string {
	return ctx.method
}

// Comment
func (ctx *Request) Path() string {
	return ctx.path
}

// Comment
func (ctx *Request) Protocol() string {
	return ctx.protocol
}

// Comment
func (ctx *Request) Query(key string) string {
	query, ok := ctx.query[key]

	if !ok {
		return ""
	}

	return query
}

// Comment
func (ctx *Request) Header(header string) string {
	h, ok := ctx.headers[header]

	if !ok {
		return ""
	}

	return h
}

// Comment
func (ctx *Request) Body() []byte {
	return ctx.body
}
