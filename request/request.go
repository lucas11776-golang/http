package request

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/types"
)

type Request struct {
	Conn   *net.Conn
	Server *server.Server
	*http.Request
}

// Comment
func toHeader(headers types.Headers) http.Header {
	h := make(http.Header)

	for k, v := range headers {
		h[k] = []string{v}
	}

	return h
}

// Comment
func Create(method string, path string, protocol string, headers types.Headers, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(method, path, strings.NewReader(""))

	if err != nil {
		return nil, err
	}

	req := &Request{Request: r}

	// url, _ := url.Parse(path)

	// req.Method = method
	// req.URL = url
	// req.RequestURI = path
	// req.Proto = protocol
	req.Header = toHeader(headers)

	return req, nil
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
	header, ok := ctx.Header[key]

	if !ok {
		return ""
	}

	return strings.Join(header, ",")
}
