package server

import (
	"fmt"
	"http/request"
	"http/response"
	"http/router"
	"http/types"
	"net"
	"reflect"
	"strconv"
	"strings"
)

const MAX_REQUEST_SIZE = (1024 * 1000)

type Server struct {
	address  string
	port     int32
	listener net.Listener
	router   *router.RouterGroup
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
		address:  arr[0],
		port:     int32(prt),
		listener: listener,
		router:   &router.RouterGroup{},
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

// comment
func (ctx *Server) Router() *router.Router {
	return ctx.router.Router()
}

// Comment
func (ctx *Server) Listen() {
	for {
		conn, err := ctx.listener.Accept()

		if err != nil {
			continue
		}

		go func() {
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

			res := response.Create("HTTP/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{})

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
