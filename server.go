package http

import (
	"log"
	"reflect"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/router"
	"github.com/lucas11776-golang/http/server"
	serve "github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
)

const MAX_REQUEST_SIZE = (1024 * 1000)

type HTTP struct {
	router *router.RouterGroup
	*serve.Server
}

// Comment
func newConnection(http *HTTP, conn *connection.Connection) {
	conn.Message(func(data []byte) {
		req, err := request.ParseHttp(string(data))

		if err != nil {
			// Invalid request page
			return
		}

		req.Server = http.Server

		route := http.Router().MatchWebRoute(req.Method(), req.Path())

		res := response.Create("HTTP/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{})

		res.Request = req

		if route == nil {
			// Not found page
			return
		}

		http := route.Call(reflect.ValueOf(req), reflect.ValueOf(res))

		err = conn.Write(http)

		if err != nil {
			// Log error
			return
		}

		conn.Close()
	})

	conn.Listen()
}

// Comment
func (ctx *HTTP) Router() *router.RouterGroup {
	return ctx.router
}

// comment
func (ctx *HTTP) Route() *router.Router {
	return ctx.router.Router()
}

// Comment
func Server(address string, port int32) *HTTP {
	server, err := server.Serve(address, port)

	if err != nil {
		log.Fatal(err)
	}

	http := &HTTP{
		Server: server,
		router: &router.RouterGroup{},
	}

	http.Connection(func(server *serve.Server, conn *connection.Connection) {
		newConnection(http, conn)
	})

	return http
}
