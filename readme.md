# HTTP


## Getting Started


### Prerequisites

HTTP requests [Go](https://go.dev) version [1.23](https://go.dev/doc/devel/release#go1.22.0) or above

**Http key features:**

- Router
- Router grouping
- Router parameters
- Response types `body`, `html`, `json`, `redirect`, `download` and `view`
- Websocket support
- Middleware


## Getting with HTTP


### Running HTTP server

Create a basic example create a `go` file called `main.go` and paste the below code.

```go
package main

import (
	"fmt"
	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/server"
	"log"
)

func main() {
	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Get("/", func(req *request.Request, res *response.Response) *response.Response {
		return res.Html("<h1>Hello World</h1>")
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
```

To run the code, use the `go run` command, like:

```sh
$ go run main.go
```

Then open your favorite browser and visit [`127.0.0.1:8080`](http://12.0.0.1:8080) you should see `Hello World`


### Route Grouping

HTTP allows simple grouping of routes using `Group` method

```go
package main

import (
	"fmt"
	"github.com/lucas11776-golang/http/router"
	"github.com/lucas11776-golang/http/server"
	"log"
)

func main() {
	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Group("api", func(route *router.Router) {
		route.Group("products", func(route *router.Router) {
			// Some routes
		})
		route.Group("invoices", func(route *router.Router) {
			// Some routes
		})
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
```


### Route Parameters

HTTP supports route params

```go
package main

import (
	"fmt"
	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/router"
	"github.com/lucas11776-golang/http/server"
	"log"
)

func main() {
	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Group("products", func(route *router.Router) {
		route.Group("{product}", func(route *router.Router) {
			route.Get("/", func(req *request.Request, res *response.Response) *response.Response {
				return res.Body([]byte("<h1>Product: " + "</h1>")).Header("content-type", "text/html")
			})
		})
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
```


### Response Types

HTTP has several response type which are `body`, `html`, `json`, `redirect`, `download` and `view` which will be explained in the next section.

```go
package main

import (
	"fmt"
	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/response"
	"github.com/lucas11776-golang/http/router"
	"github.com/lucas11776-golang/http/server"
	"log"
)

func main() {
	machine, err := server.Serve("127.0.0.1", 8080)

	if err != nil {
		log.Fatal(err)
	}

	machine.Router().Group("/", func(route *router.Router) {
		route.Get("body", func(req *request.Request, res *response.Response) *response.Response {
			return res.Body([]byte("Hello World!!!")).Header("content-type", "text/plain; charset: utf-8")
		})
		route.Get("html", func(req *request.Request, res *response.Response) *response.Response {
			return res.Html("<h1 style='color: green; font-size: 3em;'>Hello World!!!</h1>")
		})
		route.Get("json", func(req *request.Request, res *response.Response) *response.Response {
			return res.Html("<h1 style='color: green; font-size: 3em;'>Hello World!!!</h1>")
		})
		route.Get("redirect", func(req *request.Request, res *response.Response) *response.Response {
			return res.Redirect("http://www.google.com/")
		})
		route.Get("download", func(req *request.Request, res *response.Response) *response.Response {
			return res.Download("text/plain; charset=utf-8", "hello.txt", []byte("Hello World!!!"))
		})
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
```


### Response View

HTTP `view` response uses `Scriggo`(https://scrigoo.com/templates) in order to use `view` in HTTP we have to tell application where to look for `views`

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
		return res
	})

	fmt.Printf("Server running %s:%d", machine.Address(), machine.Port())

	machine.Listen()
}
```

Then create a folder in current working directory called `views` and create a file called `views/index.html` and put below html

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Golang View</title>
  <style>
    h1 { font-size: 2em; color: limegreen; }
    p { font-size: 18px; font-weight: bolder; color: gray }
  </style>
</head>
<body>
  <h1>Hello World!!!</h1>
  <p>This is my first view.</p>
</body>
</html>
```


### Websocket

HTTP support websocket without third party packages

```go
package main

import (
	"fmt"
	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/server"
	"github.com/lucas11776-golang/http/ws"
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