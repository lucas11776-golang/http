package main

import (
	"fmt"
	"http/request"
	"http/response"
	"http/router"
	"http/server"
	"log"
)

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func http() {
	users := []User{
		(User{ID: 1, Email: "jane@doe.com"}),
		(User{ID: 2, Email: "jeo@doe.com"}),
	}

	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Group("api", func(route *router.Router) {
		route.Get("/", func(req *request.Request, res *response.Response) *response.Response {
			fmt.Println("REQUEST MADE GET: ", req.Header("cookie"))
			return res.Json(users)
		})
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
