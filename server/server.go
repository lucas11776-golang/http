package server

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/router"
	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
)

const MAX_REQUEST_SIZE int64 = (1024 * 1000)

type ConnectionCallback func(server *Server, conn *connection.Connection)

type Server struct {
	address        string
	port           int32
	listener       net.Listener
	router         *router.RouterGroup
	connection     []ConnectionCallback
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
		router:         &router.RouterGroup{},
		MaxRequestSize: MAX_REQUEST_SIZE,
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
func (ctx *Server) Router() *router.RouterGroup {
	return ctx.router
}

// comment
func (ctx *Server) Route() *router.Router {
	return ctx.router.Router()
}

// Comment
func (ctx *Server) Connection(callback ConnectionCallback) *Server {
	ctx.connection = append(ctx.connection, callback)

	return ctx
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

		go func() {

			return

			// Must be read the hole request
			http := make([]byte, MAX_REQUEST_SIZE)

			n, err := conn.Read(http)

			if err != nil {
				return
			}

			req, err := request.ParseHttp(string(http[:n]))

			if err != nil {
				return
			}

			route := ctx.router.MatchWebRoute(req.Method(), req.Path())

			if route == nil {
				// Not found page
				return
			}

			res := response.Create("github.com/lucas11776-golang/http/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{})

			r := route.Call(reflect.ValueOf(req), reflect.ValueOf(res))

			_, err = conn.Write(r)

			if err != nil {
				fmt.Println("Failed To Send Response", string(err.Error()))
				return
			}

			conn.Close()
		}()
	}
}

// Comment
func (ctx *Server) Close() error {
	return ctx.listener.Close()
}
