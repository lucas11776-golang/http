package http

import (
	"net/http"

	"github.com/lucas11776-golang/http/server/connection"
)

type HTTP1 struct {
	*HTTP
}

// Comment
func InitHttp1(http *HTTP) *HTTP1 {
	return &HTTP1{HTTP: http}
}

// Comment
func (ctx *HTTP1) Init(conn *connection.Connection, req *http.Request) {
	res := ctx.HandleRequest(ctx.NewRequest(req, conn))

	if res != nil {
		conn.Write([]byte(ParseHttpResponse(res)))
	}

	conn.Close()
}
