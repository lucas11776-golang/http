package main

import (
	"encoding/json"
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.ServerTLS("127.0.0.1", 8080, "./main/host.cert", "./main/host.key").SetView("main/views", "html")
	// server := http.Server("127.0.0.1", 8080).SetView("main/views", "html")

	// server := udp.ServerTLS("127.0.0.1", 8080, "main/host.cert", "main/host.key")

	server.Route().Group("/", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			return res.View("home", http.ViewData{})
		})

		route.Group("products", func(route *http.Router) {
			route.Get("{slug}", func(req *http.Request, res *http.Response) *http.Response {
				return res.Html(fmt.Sprintf("<h1>%s</h1>", req.Parameters.Get("slug")))
			})
		})
	})

	server.Route().Group("/", func(route *http.Router) {
		route.Ws("", func(req *http.Request, ws *http.Ws) {
			ws.OnReady(func(ws *http.Ws) {
				ws.OnMessage(func(data []byte) {
					message := make(map[string]string)

					json.Unmarshal(data, &message)

					ws.WriteJson(message)
				})
			})
		})
	})

	fmt.Printf("Running server on %s", server.Host())

	// server := udp.ServerTLS("127.0.0.1", 8080, "./main/host.cert", "./main/host.key")

	// server.OnRequest(func(conn *connection.Connection, w hp.ResponseWriter, r *hp.Request) {
	// 	fmt.Println("Connection")

	// 	w.Write([]byte("Hello World..."))
	// })

	server.Listen()
}
