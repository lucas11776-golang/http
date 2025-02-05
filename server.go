package http

import (
	"log"
	"net/http"
	"reflect"

	"github.com/lucas11776-golang/http/server"
	serve "github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
)

const MAX_REQUEST_SIZE = 1024 * 1000

type HTTP struct {
	*serve.Server
}

// Comment
func newConnection(htp *HTTP, conn *connection.Connection) {
	conn.Message(func(r *http.Request, data []byte) {

		req, err := ParseHttpRequest(string(data))

		if err != nil {
			// Invalid request page
			return
		}

		req.Request = r
		req.Server = htp.Server

		route := htp.Router().MatchWebRoute(req)

		res := NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, make(types.Headers), []byte{})

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
func (ctx *HTTP) Router() *RouterGroup {
	return ctx.Get("router").(*RouterGroup)
}

// comment
func (ctx *HTTP) Route() *Router {
	return ctx.Router().Router()
}

// Comment
func (ctx *HTTP) SetView(views string, extension string) *HTTP {
	ctx.Set("view", InitView(ViewReader(views), extension))

	return ctx
}

// Comment
func Server(address string, port int32) *HTTP {
	server, err := server.Serve(address, port)

	if err != nil {
		log.Fatal(err)
	}

	http := &HTTP{
		Server: server,
	}

	http.Set("router", InitRouter())

	http.Connection(func(server *serve.Server, conn *connection.Connection) {
		newConnection(http, conn)
	})

	return http
}
