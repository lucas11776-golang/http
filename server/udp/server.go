package udp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type Http3RequestHandler struct {
	server *Server
}

// Comment
func (ctx *Http3RequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if ctx.server.callback == nil {
		return
	}

	ctx.server.callback(nil, res, req)
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
func (ctx *Server) Address() string {
	return strings.Split(ctx.server.Addr, ":")[0]
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
	return ctx.server.ListenAndServe()
}

// Comment
func (ctx *Server) Close() error {
	return ctx.server.Close()
}

// Comment
func Serve(host string, port int) *Server {
	server := &Server{
		server: &http3.Server{
			Addr:       fmt.Sprintf("%s:%d", host, port),
			TLSConfig:  http3.ConfigureTLSConfig(&tls.Config{}), // use your tls.Config here
			QUICConfig: &quic.Config{},
		},
	}

	server.server.Handler = &Http3RequestHandler{server: server}

	return server
}

// Comments
func ServerTLS(host string, port int, certFile string, keyFile string) *Server {
	var err error
	config := &tls.Config{}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	config.InsecureSkipVerify = true

	if err != nil {
		panic(err)
	}

	server := &Server{
		server: &http3.Server{
			Handler:    &Http3RequestHandler{},
			Addr:       fmt.Sprintf("%s:%d", host, port),
			QUICConfig: &quic.Config{},
			TLSConfig:  http3.ConfigureTLSConfig(config),
		},
	}

	server.server.Handler = &Http3RequestHandler{server: server}

	return server
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
