package tcp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/server/connection"
	"golang.org/x/net/http2"
)

type ConnContextHolder struct{}

type Server struct {
	server   *http.Server
	listener net.Listener
	callback func(conn *connection.Connection, w http.ResponseWriter, r *http.Request)
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

type Handler struct {
	Server *Server
}

// Comment
func (ctx *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ctx.Server.callback == nil {
		return
	}

	c, ok := r.Context().Value(ConnContextHolder{}).(net.Conn)

	if !ok {
		log.Fatalf("%s: Could not find the request connection.", r.RemoteAddr)

		return
	}

	ctx.Server.callback(connection.Init(&c), w, r)
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
	}

	httpServer := &http.Server{
		TLSConfig: tlsConfig,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, ConnContextHolder{}, c)
		},
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
	var err error
	config := &tls.Config{}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	config.InsecureSkipVerify = true

	if err != nil {
		panic(err)
	}

	return initialize(tls.NewListener(listener(host, port), config), config)
}
