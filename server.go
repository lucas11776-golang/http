package http

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/lucas11776-golang/http/pages"
	"github.com/lucas11776-golang/http/server"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
	str "github.com/lucas11776-golang/http/utils/strings"
)

type HTTP struct {
	*server.Server
	Debug                   bool
	MaxWebSocketPayloadSize int
}

const (
	SEC_WEB_SOCKET_ACCEPT_STATIC = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	SESSION_NAME                 = "session"
)

var (
	ErrWebsocketRequest = errors.New("invalid websocket request")
	ErrHttpRequest      = errors.New("invalid websocket request")
)

// Comment
func (ctx *HTTP) handleStatic(req *Request) *Response {
	res, err := ctx.Get("static").(*Static).HandleRequest(req)

	if err != nil {
		// TODO Check request is asset e.g (.js,.css and etc.) and return 404 not found request with empty body
		return nil
	}

	return res
}

// Comment
func websocketHandshake(req *Request) error {
	secWebsocketKey := req.GetHeader("sec-websocket-key")

	if secWebsocketKey == "" {
		return ErrWebsocketRequest
	}

	alg := sha1.New()

	_, err := alg.Write([]byte(strings.Join([]string{secWebsocketKey, SEC_WEB_SOCKET_ACCEPT_STATIC}, "")))

	if err != nil {
		return err
	}

	headers := types.Headers{
		"upgrade":              "websocket",
		"connection":           "Upgrade",
		"sec-webSocket-accept": base64.StdEncoding.EncodeToString(alg.Sum(nil)),
	}

	res := NewResponse(req.Protocol(), HTTP_RESPONSE_SWITCHING_PROTOCOLS, headers, []byte{})

	return req.Conn.Write([]byte(ParseHttpResponse(res)))
}

// Comment
func webSocketRequestHandler(htp *HTTP, req *Request) *Response {
	route := htp.Router().MatchWsRoute(req)

	if route == nil {
		req.Conn.Close()

		return nil
	}

	if websocketHandshake(req) != nil {
		req.Conn.Close()

		return nil
	}

	ws := InitWs(req.Conn)

	req.Ws = ws
	req.Response.Ws = ws
	ws.Request = req

	route.Call(reflect.ValueOf(req), reflect.ValueOf(ws))

	ws.Emit(EVENT_READY, []byte{})

	ws.Listen()

	return nil
}

// Comment
func (ctx *HTTP) routeNotFound(req *Request) *Response {
	if ctx.Server.Get("static") != nil {
		res := ctx.handleStatic(req)

		if res != nil {
			return res
		}
	}

	return ctx.Router().fallback(req, req.Response)
}

// Comment
func (ctx *HTTP) handleWebRouteMiddleware(route *Route, req *Request) *Response {
	for _, middleware := range route.middleware {
		next := false

		res := middleware(req, req.Response, func() *Response {
			next = true

			return req.Response
		})

		if !next {
			req.Session.Save()

			return res
		}
	}

	return nil
}

// Comment
func (ctx *HTTP) handleHTTP1_1(req *Request) *Response {
	if req.FormValue("__METHOD__") != "" {
		req.Method = strings.ToUpper(req.FormValue("__METHOD__"))
	}

	if strings.ToLower(req.GetHeader("upgrade")) == "websocket" {
		return webSocketRequestHandler(ctx, req)
	}

	route := ctx.Router().MatchWebRoute(req)

	if route == nil {
		return ctx.routeNotFound(req)
	}

	res := ctx.handleWebRouteMiddleware(route, req)

	if res != nil {
		return res
	}

	return route.Call(reflect.ValueOf(req), reflect.ValueOf(req.Response))
}

// Comment
func (ctx *HTTP) initRequestSession(req *Request) *HTTP {
	req.Response.Request = req
	req.Session = ctx.Get("session").(SessionsManager).Session(req)
	req.Response.Session = req.Session

	return ctx
}

// Comment
func (ctx *HTTP) HandleRequest(req *Request) *Response {
	switch req.Protocol() {
	case "HTTP/1.1":
		res := ctx.initRequestSession(req).handleHTTP1_1(req)

		req.Session.Save()

		return res
	default:
		return nil
	}
}

// Comment
func (ctx *HTTP) NewRequest(rq *http.Request, conn *connection.Connection) *Request {
	return &Request{
		Request:  rq,
		Server:   ctx,
		Response: NewResponse(rq.Proto, HTTP_RESPONSE_OK, make(types.Headers), []byte{}),
		Conn:     conn,
	}
}

// Comment
func (ctx *HTTP) negotiations(conn *connection.Connection, req *http.Request) (*Request, error) {
	return ctx.NewRequest(req, conn), nil
}

// Comment
func (ctx *HTTP) newConnection(conn *connection.Connection) {
	r, err := http.ReadRequest(
		bufio.NewReader(
			bufio.NewReaderSize(conn.Conn(), int(ctx.MaxRequestSize)),
		),
	)

	if err != nil {
		conn.Close()

		return
	}

	req, err := ctx.negotiations(conn, r)

	if err != nil {
		conn.Close()

		return
	}

	res := ctx.HandleRequest(req)

	if res != nil {
		conn.Write([]byte(ParseHttpResponse(res)))
	}

	conn.Close()
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
func (ctx *HTTP) SetMaxWebsocketPayload(size int) *HTTP {
	ctx.MaxWebSocketPayloadSize = size

	return ctx
}

// Comment
func (ctx *HTTP) Session(key []byte) SessionsManager {
	ctx.Set("session", InitSession(SESSION_NAME, key))

	return ctx.Get("session").(SessionsManager)
}

var (
	ErrInvalidCertificates = errors.New("invalid certificates")
)

func Init(tcp *server.Server) *HTTP {
	http := &HTTP{
		Server:                  tcp,
		MaxWebSocketPayloadSize: MAX_WEBSOCKET_PAYLOAD,
	}

	http.Set("router", InitRouter()).Get("router").(*RouterGroup).fallback = defaultRouteFallback

	http.Connection(func(server *server.Server, conn *connection.Connection) { http.newConnection(conn) })
	http.Session([]byte(str.Random(10)))

	return http
}

// Comment
func ServerTLS(host string, port int, certFile string, keyFile string) *HTTP {
	// TODO: must bind address to QUIC/UDP server here
	return Init(
		server.ServerTLS(host, port, certFile, keyFile),
	)
}

// Comment
func Server(address string, port int) *HTTP {
	// TODO: must bind address to QUIC/UDP server here
	return Init(
		server.Serve(address, port),
	)
}

type message struct {
	Message string `json:"message"`
}

// Comment
func defaultRouteFallback(req *Request, res *Response) *Response {
	res.SetStatus(HTTP_RESPONSE_NOT_FOUND)

	if req.contentType() == "application/json" {
		return res.Json(message{Message: "Route " + req.Path() + " is not found"})
	}

	return res.Html(pages.NotFound(req.Path()))
}
