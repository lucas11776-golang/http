package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var (
	ErrInvalidHost = errors.New("invalid host")
)

type Handler struct {
}

// Comment
func (ctx *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Request......")

	res.Write([]byte("<h1>Hello World</h1>"))

}

// Comment
func Connect(host string, port int) {
	server := http3.Server{
		Handler:    &Handler{},
		Addr:       fmt.Sprintf("%s:%d", host, port),
		TLSConfig:  http3.ConfigureTLSConfig(&tls.Config{}), // use your tls.Config here
		QUICConfig: &quic.Config{},
	}

	// err := server.ListenAndServe()

	err := server.ListenAndServeTLS("host.cert", "host.key")

	if err != nil {
		panic(err)
	}

}
