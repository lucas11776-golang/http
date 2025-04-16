package main

// import (
// 	"fmt"
// 	"net"
// 	"net/http"

// 	"golang.org/x/net/http2"
// )

// type Handler struct {
// }

// type contextKey struct {
// 	key string
// }

// var ConnContextKey = &contextKey{"http-conn"}

// var conns = make(map[string]*net.Conn)

// func ConnStateEvent(conn net.Conn, event http.ConnState) {
// 	if event == http.StateActive {
// 		conns[conn.RemoteAddr().String()] = &conn
// 	} else if event == http.StateHijacked || event == http.StateClosed {
// 		delete(conns, conn.RemoteAddr().String())
// 	}
// }

// func GetConn(r *http.Request) *net.Conn {
// 	return conns[r.RemoteAddr]
// }

// // Comment
// func (ctx *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	conn := GetConn(r)

// 	fmt.Println("Yes: ", conn, r.Proto)

// 	// (*conn).Write([]byte(
// 	// 	strings.Join([]string{
// 	// 		"HTTP/1.1 200 Ok",
// 	// 		"Content-Type: application/json",
// 	// 		"Content-Lenght: 15",
// 	// 		"\r\n",
// 	// 		`{"name": "jeo"}`,
// 	// 	}, "\r\n"),
// 	// ))

// 	// (*conn).Close()

// 	w.Write([]byte("<h1>Hello World</h1>"))

// 	// w.Write([]byte("<h1>Hello World"))

// }

// func main() {

// 	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

// 	server := http.Server{
// 		Addr:      "127.0.0.1:2222",
// 		ConnState: ConnStateEvent,
// 		Handler:   &Handler{},
// 	}

// 	err := http2.ConfigureServer(&server, &http2.Server{})

// 	if err != nil {
// 		panic(err)
// 	}

// 	panic(server.ListenAndServeTLS("main/host.cert", "main/host.key"))
// }

// // func main() {
// // 	server := http.ServerTLS("127.0.0.1", 2222, "main/host.cert", "main/host.key").
// // 		SetStatic("assets").
// // 		SetView("main/views", "htmp")

// // 	server.Route().Group("/", func(route *http.Router) {
// // 		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
// // 			return res.View("home", http.ViewData{})

// // 			// return res.SetStatus(http.HTTP_RESPONSE_OK).Json(map[string]string{
// // 			// 	"message": "Hello World!!!, How are you today",
// // 			// })

// // 		})
// // 	})

// // 	fmt.Println("Running Server On 127.0.0.1:6666")

// // 	server.Listen()
// // }
