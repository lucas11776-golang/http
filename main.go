package main

import (
	"fmt"
	"http/request"
	"http/server"
	"http/ws"
	"log"
)

func main() {
	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Ws("/", func(req *request.Request, socket *ws.Ws) {
		socket.OnReady(func(socket *ws.Ws) {

			socket.OnMessage(func(data []byte) {
				fmt.Println("On Message:", string(data))
			})

			socket.OnPing(func(data []byte) {
				fmt.Println("On Ping:", string(data))
			})

			socket.OnPong(func(data []byte) {
				fmt.Println("On Pong:", string(data))
			})

			socket.OnClose(func(data []byte) {
				fmt.Println("On Close:", string(data))
			})

			socket.OnError(func(data []byte) {
				fmt.Println("On Error:", string(data))
			})

		})
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
