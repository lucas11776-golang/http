package server

import (
	"net"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/server/connection"
)

const MAX_REQUEST_SIZE int64 = (1024 * 1000)

type Configuration map[string]string

type ConnectionCallback func(server *Server, conn *connection.Connection)

type Server struct {
	address        string
	port           int32
	listener       net.Listener
	connection     []ConnectionCallback
	configuration  Configuration
	MaxRequestSize int64
}

// Comment
func Serve(host string, port int32) (*Server, error) {
	listener, err := net.Listen("tcp", strings.Join([]string{host, strconv.Itoa(int(port))}, ":"))

	if err != nil {
		return nil, err
	}

	arr := strings.Split(listener.Addr().String(), ":")
	prt, _ := strconv.Atoi(arr[1])

	return &Server{
		address:        arr[0],
		port:           int32(prt),
		listener:       listener,
		MaxRequestSize: MAX_REQUEST_SIZE,
		configuration:  make(Configuration),
	}, nil
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
func (ctx *Server) SetConfig(key string, value string) *Server {
	ctx.configuration[key] = value

	return ctx
}

// Comment
func (ctx *Server) GetConfig(key string) string {
	config, ok := ctx.configuration[key]

	if !ok {
		return ""
	}

	return config
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
