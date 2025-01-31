# HTTP


## Getting Started

### Prerequisites

HTTP requests [Go]() version [1.23]()

### Getting HTTP

### Running HTTP

Create a basic example create a `go` file called `main.go` and paste the below code.

```go
package main

import (
	"fmt"
	"http/request"
	"http/response"
	"http/server"
	"log"
)

func main() {
	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Get("/", func(req *request.Request, res *response.Response) *response.Response {
		return res.Body([]byte("<h1>Hello World</h1>"))
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
```

To run the code, use the `go run` command, like:

```sh
$ go run example
```

Then open your favorite browser and visit [`127.0.0.1:8080`](http://12.0.0.1:8080) you should see `Hello World`


### HTTP Websocket

HTTP support websocket without third party packages

```go
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
```
