package http

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/lucas11776-golang/http/config"
	"github.com/lucas11776-golang/http/server"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/slices"
	str "github.com/lucas11776-golang/http/utils/strings"
)

type Dependency interface{}

type Dependencies map[string]Dependency

type HttpServer interface {
	Host() string
	// Address() string
	Port() int
	OnRequest(callback func(conn *connection.Connection, w http.ResponseWriter, r *http.Request))
	GetConnection(r *http.Request) *connection.Connection
	Listen()
	Close() error
}

type HTTP struct {
	TCP                     HttpServer
	UDP                     interface{}
	Debug                   bool
	MaxWebSocketPayloadSize int
	dependency              Dependencies
	MaxRequestSize          int64
}

type HttpHandler interface {
	Init(conn *connection.Connection, req *http.Request)
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
func (ctx *HTTP) Set(name string, dependency Dependency) *HTTP {
	ctx.dependency[name] = dependency

	return ctx
}

// Comment
func (ctx *HTTP) Get(name string) Dependency {
	dependency, ok := ctx.dependency[name]

	if !ok {
		return nil
	}

	return dependency
}

// Comment
func (ctx *HTTP) handleStatic(req *Request) *Response {
	res, err := ctx.Get("static").(*Static).HandleRequest(req)

	if err != nil {
		// TODO: must improve the checking is temp
		if len(strings.Split(slices.End(strings.Split(req.Path(), "/")), ".")) > 1 {
			return req.Response.SetStatus(HTTP_RESPONSE_NOT_FOUND)
		}

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

	res := htp.handleRouteMiddleware(route, req)

	if res != nil {
		return nil
	}

	route.Call(reflect.ValueOf(req), reflect.ValueOf(ws))

	ws.Emit(EVENT_READY, []byte{})

	ws.Listen()

	return nil
}

// Comment
func (ctx *HTTP) routeNotFound(req *Request) *Response {
	if ctx.Get("static") != nil {
		res := ctx.handleStatic(req)

		if res != nil {
			return res
		}
	}

	return ctx.Router().fallback(req, req.Response)
}

// Comment
func (ctx *HTTP) handleRouteMiddleware(route *Route, req *Request) *Response {
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
func (ctx *HTTP) initRequestSession(req *Request) *Request {
	req.Response.Request = req
	req.Session = ctx.Get("session").(SessionsManager).Session(req)
	req.Response.Session = req.Session

	return req
}

// Comment
func (ctx *HTTP) HandleRequest(req *Request) *Response {
	req = ctx.initRequestSession(req)

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

	res := ctx.handleRouteMiddleware(route, req)

	if res != nil {
		return res
	}

	res = route.Call(reflect.ValueOf(req), reflect.ValueOf(req.Response))

	req.Session.Save()

	return res
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
func (ctx *HTTP) negotiation(req *http.Request) HttpHandler {
	if strings.ToLower(req.Header.Get("upgrade")) == "h2c" {
		return InitHttp2(ctx)
	}

	return InitHttp1(ctx)
}

// comment
func (ctx *HTTP) readRequest(conn *connection.Connection) (*http.Request, error) {
	return http.ReadRequest(
		bufio.NewReader(
			bufio.NewReaderSize(conn.Conn(), int(ctx.MaxRequestSize)),
		),
	)
}

// Comment
func (ctx *HTTP) newConnection(conn *connection.Connection) {
	// TODO: must check if connection is direct APLN/H2C or just HTTP/1.1
	req, err := ctx.readRequest(conn)

	if err == nil {
		ctx.negotiation(req).Init(conn, req)
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

// Comment
func defaultRouteFallback(req *Request, res *Response) *Response {
	return res.SetStatus(HTTP_RESPONSE_NOT_FOUND)
}

func Init(tcp HttpServer) *HTTP {
	server := &HTTP{
		MaxWebSocketPayloadSize: MAX_WEBSOCKET_PAYLOAD,
		dependency: Dependencies{
			"config": config.Init(),
		},
	}

	server.TCP = tcp
	server.UDP = nil

	server.Set("router", InitRouter()).Get("router").(*RouterGroup).fallback = defaultRouteFallback
	server.Session([]byte(str.Random(10)))

	// http.TCP.Connection(func(conn *connection.Connection) { http.newConnection(conn) })

	server.TCP.OnRequest(func(conn *connection.Connection, w http.ResponseWriter, r *http.Request) {
		// fmt.Println("YES.....", r.Proto, conn, r.RemoteAddr)

		// r.Header.Set("Content-Type", "text/html")

		w.Write([]byte("<h1>Hello World</h1>"))

		// req := server.NewRequest(r, conn)

		// res := server.HandleRequest(req)

		// if

		// res.

		// w.WriteHeader(200)
	})

	return server
}

// Comment
func ServerTLS(host string, port int, cert string, key string) *HTTP {
	// TODO: must bind address to QUIC/UDP server here
	return Init(
		server.ServeTLS(host, port, cert, key),
	)
}

// Comment
func Server(address string, port int) *HTTP {
	// TODO: must bind address to QUIC/UDP server here
	return Init(
		server.Serve(address, port),
	)
}

// Comm
func (ctx *HTTP) Host() string {
	return ctx.TCP.Host()
}

// Comment
// func (ctx *HTTP) Address() string {
// 	return ctx.TCP.Address()
// }

// Comment
func (ctx *HTTP) Port() int {
	return ctx.TCP.Port()
}

// Comment
func (ctx *HTTP) Listen() {
	ctx.TCP.Listen()
}

// Comment
func (ctx *HTTP) Close() (tcp error, udp error) {
	return ctx.TCP.Close(), nil
}
