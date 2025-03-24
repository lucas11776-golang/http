package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 6666)

	server.Route().Group("/", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			return res.SetStatus(http.HTTP_RESPONSE_OK).Json(map[string]string{
				"message": "Hello World!!!, How are you today",
			})
		})
	})

	fmt.Println("Running Server On 127.0.0.1:6666")

	server.Listen()
}
