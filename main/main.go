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

	var Users = []User{
		{
			ID:    1,
			Name:  "Jeo Doe",
			Email: "jeo@deo.com",
		},
		{
			ID:    2,
			Name:  "Jane Doe",
			Email: "jane@deo.com",
		},
	}

	server := http.Server("127.0.0.1", 9090).SetView("main/views", "html").SetStatic("main/static")

	server.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.View("home", http.ViewData{
			"users": &Users,
		})
	})

	fmt.Printf("Running server on %s", server.Host())

	server.Listen()
}
