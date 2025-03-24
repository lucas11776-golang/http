package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
	"golang.org/x/net/http2"
)

type HTTP2 struct {
	*HTTP
}

// Comment
func InitHttp2(http *HTTP) *HTTP2 {
	return &HTTP2{HTTP: http}
}

// Comment
func (ctx *HTTP2) handlen(header http2.FrameHeader, req *Request) {
	fmt.Println("------- HANDLE -------")
}

// Comment
func (ctx *HTTP2) upgrade(conn *connection.Connection) error {
	return conn.Write([]byte(
		ParseHttpResponse(
			NewResponse("HTTP/1.1", HTTP_RESPONSE_SWITCHING_PROTOCOLS, types.Headers{}, []byte{}),
		),
	))
}

// // Comment
func (ctx *HTTP2) listen(conn *connection.Connection) {
	for {
		// Do some HTTP2 greement....

		_, err := conn.Read(make([]byte, 1024*4))

		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}
	}
}

// Comment
func (ctx *HTTP2) Init(conn *connection.Connection, req *http.Request) {
	if err := ctx.upgrade(conn); err == nil {
		ctx.listen(conn)
	}
}
