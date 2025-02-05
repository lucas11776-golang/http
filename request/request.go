package request

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/types"
	h "github.com/lucas11776-golang/http/utils/headers"
)

type Request struct {
	Conn   *net.Conn
	Server *server.Server
	*http.Request
}

// Comment
func Create(method string, path string, protocol string, headers types.Headers, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, path, strings.NewReader(""))

	if err != nil {
		return nil, err
	}

	req.Header = h.ToHeader(headers)

	return &Request{Request: req}, nil
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
