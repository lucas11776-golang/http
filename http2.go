package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

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
			NewResponse("HTTP/1.1", HTTP_RESPONSE_SWITCHING_PROTOCOLS, types.Headers{
				"Connection": "Upgrade",
				"Upgrade":    "HTTP/2.0",
			}, []byte{}),
		),
	))
}

// // Comment
func (ctx *HTTP2) listen(conn *connection.Connection) {
	for {
		// Do some HTTP2 greement....
		data, err := conn.Read(make([]byte, 1024*4))

		fmt.Printf("Connection: \r\n%s\r\n", string(data))

		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}

		// TODO: lazy work must write better way
		if strings.Split(string(data), "\r\n")[0] != "PRI * HTTP/2.0" {
			break
		}

		// Do some HTTP2 greement....
		data, err = conn.Read(make([]byte, 1024*4))

		fmt.Printf("New: \r\n%s\r\n", string(data))

		if err != nil {
			continue
		}

		conn.Write([]byte(
			strings.Join([]string{
				":status: 200",
				":scheme: http",
				"\r\n",
			}, "\r\n"),
		))

	}

	conn.Close()
}

// Comment
func (ctx *HTTP2) Init(conn *connection.Connection, req *http.Request) {

	fmt.Println("Request", req.URL.Path, req.Method)

	if err := ctx.upgrade(conn); err == nil {
		ctx.listen(conn)
	}
}
