package udp

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type Http3RequestHandler struct {
}

// Comment
func (ctx *Http3RequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("<h1>Hello World</h1>"))
}

type Server struct {
	server   *http3.Server
	callback func(conn *connection.Connection, w http.ResponseWriter, r *http.Request)
}

// Comment
func (ctx *Server) Host() string {
	return ctx.server.Addr
}

// Comment
func (ctx *Server) Port() int {
	return ctx.server.Port
}

// Comment
func (ctx *Server) OnRequest(callback func(conn *connection.Connection, w http.ResponseWriter, r *http.Request)) {
	ctx.callback = callback
}

// Comment
func (ctx *Server) Listen() error {
	return nil
}

// Comment
func (ctx *Server) Close() error {
	return nil
}

// Comment
func Serve(host string, port int) *Server {
	return &Server{
		server: &http3.Server{
			Handler:    &Http3RequestHandler{},
			Addr:       fmt.Sprintf("%s:%d", host, port),
			TLSConfig:  http3.ConfigureTLSConfig(&tls.Config{}), // use your tls.Config here
			QUICConfig: &quic.Config{},
		},
	}
}

// ------------------------------------------------------------------------------------------------------ //
// var (
// 	ErrInvalidHost = errors.New("invalid host")
// )

// type Handler struct {
// }

// // Comment
// func (ctx *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
// 	fmt.Println("Request......")

// 	res.Write([]byte("<h1>Hello World</h1>"))

// }

// // Comment
// func Connect(host string, port int) {
// 	server := http3.Server{
// 		Handler:    &Handler{},
// 		Addr:       fmt.Sprintf("%s:%d", host, port),
// 		TLSConfig:  http3.ConfigureTLSConfig(&tls.Config{}), // use your tls.Config here
// 		QUICConfig: &quic.Config{},
// 	}

// 	// err := server.ListenAndServe()

// 	err := server.ListenAndServeTLS("host.cert", "host.key")

// 	if err != nil {
// 		panic(err)
// 	}

// }
