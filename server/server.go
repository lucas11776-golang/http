package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/config"
	"github.com/lucas11776-golang/http/server/connection"
)

const MAX_REQUEST_SIZE int64 = 1024 * 1000

var (
	ErrInvalidCertificates = errors.New("invalid certificates")
	ErrInvalidPort         = errors.New("invalid port")
)

type ConnectionCallback func(server *Server, conn *connection.Connection)

type Dependency interface{}

type Dependencies map[string]Dependency

type Server struct {
	address        string
	port           int
	listener       net.Listener
	connection     ConnectionCallback
	MaxRequestSize int64
	dependency     Dependencies
}

// Comment
func Init(address string, port int, listener net.Listener) *Server {
	return &Server{
		address:        address,
		port:           port,
		MaxRequestSize: MAX_REQUEST_SIZE,
		listener:       listener,
		dependency: Dependencies{
			"config": config.Init(),
		},
	}
}

// Comment
func ServerTLS(host string, port int, certFile string, keyFile string) *Server {
	if certFile != "" || keyFile != "" {
		panic(ErrInvalidCertificates)
	}

	var err error

	config := &tls.Config{}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		panic(err)
	}

	return Init(host, port, listener)
}

// Comment
func Serve(host string, port int) *Server {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		panic(err)
	}

	prt, err := strconv.Atoi(strings.Split(listener.Addr().String(), ":")[1])

	if err != nil {
		panic(ErrInvalidPort)
	}

	return Init(host, prt, listener)
}

// Comment
func (ctx *Server) Address() string {
	return ctx.address
}

// Comment
func (ctx *Server) Port() int {
	return ctx.port
}

// Comment
func (ctx *Server) Host() string {
	return ctx.listener.Addr().String()
}

// Comment
func (ctx *Server) Connection(callback ConnectionCallback) *Server {
	ctx.connection = callback

	return ctx
}

// Comment
func (ctx *Server) SetMaxRequestSize(size int64) *Server {
	ctx.MaxRequestSize = size

	return ctx
}

// Comment
func (ctx *Server) Set(name string, dependency Dependency) *Server {
	ctx.dependency[name] = dependency

	return ctx
}

// Comment
func (ctx *Server) Get(name string) Dependency {
	dependency, ok := ctx.dependency[name]

	if !ok {
		return nil
	}

	return dependency
}

// Comment
func (ctx *Server) Listen() {
	for {
		conn, err := ctx.listener.Accept()

		if err != nil {
			continue
		}

		go func() {
			ctx.connection(ctx, connection.Init(&conn, ctx.MaxRequestSize))
		}()
	}
}

// Comment
func (ctx *Server) Close() error {
	return ctx.listener.Close()
}
