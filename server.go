package http

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

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
func responseWrite(req *Request, http []byte) {
	err := req.Conn.Write(http)

	if err != nil {
		return
	}
}

const (
	ESTABLISH_CONNECTION_PAYLOAD_SIZE = 2048
	SEC_WEB_SOCKET_ACCEPT_STATIC      = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

var ERROR_INVALID_WEBSOCKET_REQUEST = errors.New("Invalid http request")

// Comment
func handShakeReplay(req *Request) ([]byte, error) {
	res := NewResponse(req.Protocol(), HTTP_RESPONSE_SWITCHING_PROTOCOLS, make(types.Headers), []byte{})

	secWebsocketKey := req.GetHeader("sec-websocket-key")

	if secWebsocketKey == "" {
		return nil, ERROR_INVALID_WEBSOCKET_REQUEST
	}

	alg := sha1.New()

	alg.Write([]byte(strings.Join([]string{secWebsocketKey, SEC_WEB_SOCKET_ACCEPT_STATIC}, "")))

	hashed := base64.StdEncoding.EncodeToString(alg.Sum(nil))

	// res.SetStatus(101)
	res.SetHeader("Upgrade", "websocket")
	res.SetHeader("Connection", "Upgrade")
	res.SetHeader("Sec-WebSocket-Accept", hashed)

	return []byte(ParseHttpResponse(res)), nil
}

// Comment
func handleHTTP1_1(htp *HTTP, req *Request) {
	if strings.ToLower(req.GetHeader("Upgrade")) == "websocket" {
		route := htp.Router().MatchWsRoute(req)

		if route == nil {
			// TODO Error response
			return
		}

		reply, err := handShakeReplay(req)

		if err != nil {
			// TODO Error response
			return
		}

		err = req.Conn.Write(reply)

		if err != nil {
			// TODO Error response
			return
		}

		ws := InitWs(req.Conn)

		ws.Request = req

		route.Call(reflect.ValueOf(req), reflect.ValueOf(ws))

		ws.Emit(EVENT_READY, []byte{})

		ws.Listen()

		return
	}

	route := htp.Router().MatchWebRoute(req)

	if route == nil {
		if htp.Server.Get("static") != nil {
			err := handleStatic(req.Conn, htp.Get("static").(*Static), req)

			if err == nil {
				req.Conn.Close()

				return
			}
		}

		// TODO Error response
		return
	}

	for _, middleware := range route.middleware {
		req.Response.Next = false

		res := middleware(req, req.Response, func() *Response {
			req.Response.Next = true

			return req.Response
		})

		if !res.Next {
			responseWrite(req, []byte(ParseHttpResponse(res)))

			return
		}
	}

	responseWrite(req, route.Call(reflect.ValueOf(req), reflect.ValueOf(req.Response)))

	req.Conn.Close()
}

// Comment
func newConnection(htp *HTTP, conn *connection.Connection) {
	r, err := http.ReadRequest(bufio.NewReader(bufio.NewReaderSize(conn.Conn(), int(htp.MaxRequestSize))))

	if err != nil {
		return
	}

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
