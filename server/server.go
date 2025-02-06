package server

import (
	"net"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/config"
	"github.com/lucas11776-golang/http/server/connection"
)

const MAX_REQUEST_SIZE int64 = 1024 * 1000

type ConnectionCallback func(server *Server, conn *connection.Connection)

type Dependency interface{}

type Dependencies map[string]Dependency

type Server struct {
	address        string
	port           int32
	listener       net.Listener
	connection     []ConnectionCallback
	MaxRequestSize int64
	dependency     Dependencies
}

// Comment
func Init(address string, port int32, listener net.Listener) *Server {
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
func Serve(host string, port int32) (*Server, error) {
	listener, err := net.Listen("tcp", strings.Join([]string{host, strconv.Itoa(int(port))}, ":"))

	if err != nil {
		return nil, err
	}

	arr := strings.Split(listener.Addr().String(), ":")
	prt, _ := strconv.Atoi(arr[1])

	return Init(arr[0], int32(prt), listener), nil
}

// Comment
func (ctx *Server) Address() string {
	return ctx.address
}

// Comment
func (ctx *Server) Port() int32 {
	return ctx.port
}

// Comment
func (ctx *Server) Host() string {
	return ctx.listener.Addr().String()
}

// Comment
func (ctx *Server) Connection(callback ConnectionCallback) *Server {
	ctx.connection = append(ctx.connection, callback)

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
			for _, callback := range ctx.connection {
				go func() {
					callback(ctx, connection.Init(&conn, ctx.MaxRequestSize))
				}()
			}
		}()
	}
}

// Comment
func (ctx *Server) Close() error {
	return ctx.listener.Close()
}
