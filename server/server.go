package server

import (
	"net"
	"strconv"
	"strings"
)

type Server struct {
	address  string
	port     int32
	listener net.Listener
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
func (ctx *Server) Listen() {

}

// Comment
func (ctx *Server) Close() error {
	return ctx.listener.Close()
}
