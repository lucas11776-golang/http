package http

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/lucas11776-golang/http/config"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/server/tcp"
	"github.com/lucas11776-golang/http/server/udp"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/response"
	"github.com/lucas11776-golang/http/utils/slices"
	str "github.com/lucas11776-golang/http/utils/strings"
)

const (
	SEC_WEB_SOCKET_ACCEPT_STATIC = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	SESSION_NAME                 = "session"
)

var (
	ErrWebsocketRequest    = errors.New("invalid websocket request")
	ErrHttpRequest         = errors.New("invalid websocket request")
	ErrInvalidCertificates = errors.New("invalid certificates")
)

type Dependency interface{}

type Dependencies map[string]Dependency

type HttpServer interface {
	Host() string
	Port() int
	OnRequest(callback func(conn *connection.Connection, w http.ResponseWriter, r *http.Request))
	Listen() error
	Close() error
}

type HTTP struct {
	TCP                     HttpServer
	UDP                     HttpServer
	Debug                   bool
	MaxWebSocketPayloadSize int
	dependency              Dependencies
	MaxRequestSize          int64
}

type HttpHandler interface {
	Init(conn *connection.Connection, req *http.Request)
}

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
			return res
		}
	}

	return nil
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
func (ctx *HTTP) Router() *RouterGroup {
	return ctx.Get("router").(*RouterGroup)
}

// comment
func (ctx *HTTP) Route() *Router {
	return ctx.Router().Router()
}

// Comment
func (ctx *HTTP) SetView(path string, extension string) *HTTP {
	return ctx.Set("view", NewView(NewDefaultViewReader(path), extension))
}

// Comment
func (ctx *HTTP) SetStatic(path string) *HTTP {
	return ctx.Set("static", InitStatic(NewDefaultStaticReader(path)))
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

// Comment
func defaultRouteFallback(req *Request, res *Response) *Response {
	return res.SetStatus(HTTP_RESPONSE_NOT_FOUND)
}

type RequestHandler interface {
	Handle(connection *connection.Connection, w http.ResponseWriter, r *Request)
}

// Comment
func (ctx *HTTP) setupRequest(req *Request) *Request {
	req.Server = ctx

	req.Session = ctx.Get("session").(SessionsManager).Session(req)
	req.Response.Session = req.Session

	if method := req.FormValue("__METHOD__"); method != "" {
		req.Method = strings.ToUpper(method)
	}

	req.Response.Request = req

	return req
}

// Comment
func (ctx *HTTP) requestHandler(req *Request) *Response {
	route := ctx.Router().MatchWebRoute(req)

	if route == nil {
		return ctx.routeNotFound(req)
	}

	if res := ctx.handleRouteMiddleware(route, req); res != nil {
		return res
	}

	return route.Call(reflect.ValueOf(req), reflect.ValueOf(req.Response))
}

// Comment
func (ctx *HTTP) websocketHandshake(req *Request) error {
	secWebsocketKey := req.GetHeader("sec-websocket-key")

	if secWebsocketKey == "" {
		return ErrWebsocketRequest
	}

	alg := sha1.New()

	_, err := alg.Write([]byte(strings.Join([]string{secWebsocketKey, SEC_WEB_SOCKET_ACCEPT_STATIC}, "")))

	if err != nil {
		return err
	}

	res := NewResponse(req.Protocol(), HTTP_RESPONSE_SWITCHING_PROTOCOLS, types.Headers{
		"upgrade":              "websocket",
		"connection":           "Upgrade",
		"sec-webSocket-accept": base64.StdEncoding.EncodeToString(alg.Sum(nil)),
	}, []byte{})

	return req.Conn.Write([]byte(response.ResponseToHttp(res.Response)))
}

// Comment
func (ctx *HTTP) websocketHandler(req *Request) {
	route := ctx.Router().MatchWsRoute(req)

	if route == nil {
		return
	}

	if err := ctx.websocketHandshake(req); err != nil {
		return
	}

	ws := InitWs(req.Conn, req)

	req.Ws = ws
	req.Response.Ws = ws
	ws.Request = req

	if res := ctx.handleRouteMiddleware(route, req); res != nil {
		req.Conn.Close()

		return
	}

	route.Call(reflect.ValueOf(req), reflect.ValueOf(ws))

	ws.isReady()

	ws.Listen()
}

// Comment
func (ctx *HTTP) HandleRequest(req *Request) *Response {
	switch strings.ToLower(ctx.setupRequest(req).GetHeader("upgrade")) {

	case "websocket":
		ctx.websocketHandler(req)

		return nil

	default:
		return ctx.requestHandler(req)

	}
}

// Comment
func (ctx *HTTP) Handler(conn *connection.Connection, req *Request) {
	res := ctx.HandleRequest(req)

	if res == nil {
		return
	}

	for key, value := range res.Header {
		req.Response.Writer.Header().Set(key, value[0])
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		req.Response.Writer.WriteHeader(int(HTTP_RESPONSE_INTERNAL_SERVER_ERROR))
		req.Response.Writer.Write([]byte{})

		return
	}

	req.Response.Writer.WriteHeader(res.StatusCode)
	req.Response.Writer.Write(body)
}

// Comment
func Init(tcp HttpServer, udp HttpServer) *HTTP {
	server := &HTTP{
		MaxWebSocketPayloadSize: MAX_WEBSOCKET_PAYLOAD,
		dependency: Dependencies{
			"config": config.Init(),
		},
	}

	server.TCP = tcp
	server.UDP = udp

	server.Set("router", InitRouter()).Get("router").(*RouterGroup).fallback = defaultRouteFallback
	server.Session([]byte(str.Random(10)))

	server.TCP.OnRequest(server.onRequest) // HTTP/1.1 and HTTP/2.0 requests
	server.UDP.OnRequest(server.onRequest) // HTTP/3.0 requests

	return server
}

// Comment
func (ctx *HTTP) onRequest(conn *connection.Connection, w http.ResponseWriter, r *http.Request) {
	// TODO: Need to check http3 support websocket not sure.
	if strings.ToLower(r.Header.Get("upgrade")) == "websocket" && strings.ToUpper(r.Proto) == "HTTP/3.0" {
		return
	}

	req := ctx.NewRequest(r, conn)
	req.Response.Writer = w

	ctx.Handler(conn, req)
	req.Session.Save()
}

// Comment
func ServerTLS(host string, port int, certFile string, keyFile string) *HTTP {
	tcp := tcp.ServeTLS(host, port, certFile, keyFile)
	udp := udp.ServerTLS(host, tcp.Port(), certFile, keyFile)

	return Init(tcp, udp)
}

// Comment
func Server(address string, port int) *HTTP {
	tcp := tcp.Serve(address, port)
	udp := udp.Serve(address, tcp.Port())

	return Init(tcp, udp)
}

// Comm
func (ctx *HTTP) Host() string {
	return ctx.TCP.Host()
}

// Comment
func (ctx *HTTP) Port() int {
	return ctx.TCP.Port()
}

// Comment
func (ctx *HTTP) Listen() {
	go ctx.TCP.Listen()
	go ctx.UDP.Listen()
	select {}
}

// Comment
func (ctx *HTTP) Close() (tcp error, udp error) {
	return ctx.TCP.Close(), nil
}
