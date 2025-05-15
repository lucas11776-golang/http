package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	server := http.Server("127.0.0.1", 80)

	server.Route().Subdomain("api", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			return res.Json(map[string]string{
				"company": req.Parameters.Get("company"),
			})
		})
	})

	fmt.Printf("Running server on %s", server.Host())

	server.Listen()
}
