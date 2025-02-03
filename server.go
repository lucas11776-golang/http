package http

import (
	"log"
	"reflect"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
)

const MAX_REQUEST_SIZE = (1024 * 1000)

// Comment
func newConnection(server *server.Server, conn *connection.Connection) {
	conn.Message(func(data []byte) {
		req, err := request.ParseHttp(string(data))

		if err != nil {
			// Invalid request page
			return
		}

		route := server.Router().MatchWebRoute(req.Method(), req.Path())

		res := response.Create("HTTP/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{})

		if route == nil {
			// Not found page
			return
		}

		http := route.Call(reflect.ValueOf(req), reflect.ValueOf(res))

		err = conn.Write(http)

		if err != nil {
			// Log error
			return
		}

		conn.Close()
	})

	conn.Listen()

	// Must be read the hole request
	// http := make([]byte, MAX_REQUEST_SIZE)

	// n, err := conn.Read(http)

	// if err != nil {
	// 	return
	// }

	// req, err := request.ParseHttp(string(http[:n]))

	// if err != nil {
	// 	return
	// }

	// route := ctx.router.MatchWebRoute(req.Method(), req.Path())

	// if route == nil {
	// 	// Not found page
	// 	return
	// }

	// res := response.Create("github.com/lucas11776-golang/http/1.1", response.HTTP_RESPONSE_OK, make(types.Headers), []byte{})

	// r := route.Call(reflect.ValueOf(req), reflect.ValueOf(res))

	// _, err = conn.Write(r)

	// if err != nil {
	// 	fmt.Println("Failed To Send Response", string(err.Error()))
	// 	return
	// }

	// conn.Close()
}

// Comment
func Server(address string, port int32) *server.Server {
	server, err := server.Serve(address, port)

	if err != nil {
		log.Fatal(err)
	}

	return server.Connection(newConnection)
}
