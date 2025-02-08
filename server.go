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

type HTTP struct {
	*serve.Server
}

// Comment
func handleStatic(conn *connection.Connection, static *Static, req *Request) error {
	res, err := static.HandleRequest(req)

	if err != nil {
		return err
	}

	return conn.Write([]byte(ParseHttpResponse(res)))
}

// Comment
func handleHTTP1_1(htp *HTTP, req *Request) {
	route := htp.Router().MatchWebRoute(req)

	if route == nil {
		if htp.Server.Get("static") != nil {
			err := handleStatic(req.Conn, htp.Get("static").(*Static), req)

			if err == nil {
				req.Conn.Close()

				return
			}
		}

		// TODO Resource not found
		return
	}

	http := route.Call(reflect.ValueOf(req), reflect.ValueOf(req.Response))

	err := req.Conn.Write(http)

	if err != nil {
		return
	}

	req.Conn.Close()
}

// Comment
func newConnection(htp *HTTP, conn *connection.Connection) {
	conn.Message(func(r *http.Request) {
		req := &Request{
			Request:  r,
			Server:   htp.Server,
			Response: NewResponse(r.Proto, HTTP_RESPONSE_OK, make(types.Headers), []byte{}),
			Conn:     conn,
		}

		req.Response.Request = req

		switch req.Protocol() {
		case "HTTP/1.1":
			handleHTTP1_1(htp, req)
			break
		default:
			conn.Close()
		}
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
	ctx.Set("view", InitView(DefaultViewReader(views), extension))

	return ctx
}

// Comment
func (ctx *HTTP) SetStatic(statics string) *HTTP {
	ctx.Set("static", InitStatic(DefaultStaticReader(statics)))

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
