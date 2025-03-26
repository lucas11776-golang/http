package main

import "github.com/lucas11776-golang/http"

// import (
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/lucas11776-golang/http/server"
// )

// func main() {
// 	server := server.Serve("127.0.0.1", 2222)

// 	server.OnRequest(func(w http.ResponseWriter, r *http.Request) {

// 		fmt.Println("Requesting Connection")

// 		w.Write([]byte("<h1>Hello World</h1>"))

// 		time.Sleep(time.Second * 3)

// 		w.Write([]byte("<h1>Hello World 2</h1>"))

// 	})

// 	server.Listen()
// }

func main() {
	server := http.ServerTLS("127.0.0.1", 2222, "main/host.cert", "main/host.key").SetView("main/views", "html")
	// server := http.Server("127.0.0.1", 2222)

	server.Route().Group("/", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			// return res.Json(map[string]string{
			// 	"message": "Hello World",
			// })

			return res //.View("home", http.ViewData{})
		})
	})

	server.Listen()
}

// panic(server.ListenAndServeTLS("main/host.cert", "main/host.key"))
