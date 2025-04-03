package main

import (
	"encoding/json"
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	// server := http.ServerTLS("127.0.0.1", 2222, "main/host.cert", "main/host.key").SetView("main/views", "html")
	server := http.Server("127.0.0.1", 2222).SetView("main/views", "html")

	server.Route().Group("/", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			// return res.Json(map[string]string{
			// 	"message": "Hello World",
			// })

			return res.View("home", http.ViewData{})
		})
	})

	server.Route().Group("/", func(route *http.Router) {
		route.Ws("", func(req *http.Request, ws *http.Ws) {

			ws.OnReady(func(ws *http.Ws) {

				ws.OnMessage(func(data []byte) {

					message := make(map[string]string)

					json.Unmarshal(data, &message)

					// fmt.Printf("\r\n\r\n\r\n %s \r\n\r\n\r\n ", message)

					ws.WriteJson(message)
				})

				go func() {
					// time.Sleep(time.Second * 2)

					// ws.Write([]byte("Hello World"))
				}()
			})
		})
	})

	fmt.Printf("Running server on %s", server.Host())

	server.Listen()
}
