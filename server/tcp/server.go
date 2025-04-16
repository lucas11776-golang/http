package tcp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/lucas11776-golang/http/server/connection"
	"golang.org/x/net/http2"
)

type Server struct {
	server      *http.Server
	listener    net.Listener
	connections *connections
	callback    func(conn *connection.Connection, w http.ResponseWriter, r *http.Request)
}

type connections struct {
	connections map[string]*connection.Connection
	mutex       sync.Mutex
}

// Comment
func (ctx *connections) add(conn net.Conn) {
	if _, ok := ctx.connections[conn.RemoteAddr().String()]; !ok {
		ctx.mutex.Lock()
		ctx.connections[conn.RemoteAddr().String()] = connection.Init(&conn)
		ctx.mutex.Unlock()
	}
}

// Comment
func (ctx *connections) remove(conn net.Conn) {
	ctx.mutex.Lock()
	delete(ctx.connections, conn.RemoteAddr().String())
	ctx.mutex.Unlock()
}

// Comment
func (ctx *connections) ConnStateEvent(conn net.Conn, event http.ConnState) {
	switch event {

	case http.StateActive:
		ctx.add(conn)

	case http.StateHijacked, http.StateClosed:
		ctx.remove(conn)

	}
}

// Comment
func (ctx *connections) GetConnection(req *http.Request) *connection.Connection {
	if conn, ok := ctx.connections[req.RemoteAddr]; ok {
		return conn
	}
	return nil
}

// Comment
func (ctx *Server) Host() string {
	return ctx.listener.Addr().String()
}

// Comment
func (ctx *Server) Address() string {
	return strings.Split(ctx.listener.Addr().String(), ":")[0]
}

// Comment
func (ctx *Server) Port() int {
	port, _ := strconv.Atoi(strings.Split(ctx.listener.Addr().String(), ":")[1])

	return port
}

// Comment
func (ctx *Server) OnRequest(callback func(conn *connection.Connection, w http.ResponseWriter, r *http.Request)) {
	ctx.callback = callback
}

// Comment
func (ctx *Server) GetConnection(r *http.Request) *connection.Connection {
	return ctx.connections.GetConnection(r)
}

type Handler struct {
	Server *Server
}

// Comment
func (ctx *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ctx.Server.callback != nil {
		ctx.Server.callback(ctx.Server.GetConnection(r), w, r)
	}
}

// Comment
func (ctx *Server) Listen() error {
	return ctx.server.Serve(ctx.listener)
}

// Comment
func (ctx *Server) Close() error {
	return ctx.server.Close()
}

// Comment
func listener(host string, port int) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		panic(err)
	}

	return listener
}

// Comment
func initialize(listener net.Listener, tlsConfig *tls.Config) *Server {
	server := &Server{
		listener: listener,
		connections: &connections{
			connections: make(map[string]*connection.Connection),
		},
	}

	httpServer := &http.Server{
		TLSConfig: tlsConfig,
		ConnState: server.connections.ConnStateEvent,
	}

	httpServer.Handler = &Handler{Server: server}

	if err := http2.ConfigureServer(httpServer, nil); err != nil {
		panic(err)
	}

	server.server = httpServer

	return server
}

// Comment
func Serve(host string, port int) *Server {
	return initialize(listener(host, port), nil)
}

// Comment
func ServeTLS(host string, port int, certFile string, keyFile string) *Server {
	var config *tls.Config = nil

	var err error

	config = &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		certs, err := tls.X509KeyPair([]byte(certFile), []byte(keyFile))

		if err != nil {
			panic(err)
		}

		config.Certificates[0] = certs
	}

	return initialize(tls.NewListener(listener(host, port), config), config)
}
