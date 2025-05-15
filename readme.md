# HTTP


## Getting Started


### Prerequisites

HTTP requests [Go](https://go.dev) version [1.23](https://go.dev/doc/devel/release#go1.22.0) or above

**Http key features:**

- Router         - `Group`, `Subdomain`
- Response Types - `body`, `html`, `json`, `redirect`, `download` and `view`
- Static Assets
- WebSocket
- Middleware
- Session


## Getting with HTTP


### Running HTTP server

Create a basic example create a `go` file called `main.go` and paste the below code.

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.Html("<h1>Hello World</h1>")
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```

To run the code, use the `go run` command, like:

```sh
$ go run main.go
```

Then open your favorite browser and visit [`127.0.0.1:8080`](http://12.0.0.1:8080) you should see `Hello World`

#### Route Grouping

HTTP allows simple grouping of routes using `Group` method

```go
package main

import (
	"fmt"
	
	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Group("api", func(route *http.Router) {
		route.Group("products", func(route *http.Router) {
			// Some routes
		})
		route.Group("invoices", func(route *http.Router) {
			// Some routes
		})
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```

#### Route Parameters

HTTP supports route params

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Group("products", func(route *http.Router) {
		route.Group("{product}", func(route *http.Router) {
			route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
				return res.Body([]byte("<h1>Product: " + "</h1>")).Header("content-type", "text/html")
			})
		})
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```


### Response Types

HTTP has several response type which are `body`, `html`, `json`, `redirect`, `download` and `view` which will be explained in the next section.

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Group("/", func(route *http.Router) {
		route.Get("body", func(req *http.Request, res *http.Response) *http.Response {
			return res.Body([]byte("Hello World!!!")).Header("content-type", "text/plain; charset: utf-8")
		})
		route.Get("html", func(req *http.Request, res *http.Response) *http.Response {
			return res.Html("<h1 style='color: green; font-size: 3em;'>Hello World!!!</h1>")
		})
		route.Get("view", func(req *http.Request, res *http.Response) *http.Response {
			return res.View("hello_world")
		})
		route.Get("json", func(req *http.Request, res *http.Response) *http.Response {
			return res.Json(struct {
				Message string `json:"message"`
			}{Message: "Hello World!!!"})
		})
		route.Get("redirect", func(req *http.Request, res *http.Response) *http.Response {
			return res.Redirect("http://www.google.com/")
		})
		route.Get("download", func(req *http.Request, res *http.Response) *http.Response {
			return res.Download("text/plain; charset=utf-8", "hello.txt", []byte("Hello World!!!"))
		})
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```


#### Response View

HTTP `view` response uses [`Scriggo`](https://scriggo.com/templates) in order to use `view` in HTTP we have to tell application where to look for `views`

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.SetView("views", "html")

	server.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.View("index", http.ViewData{
			"name": "lucas11776",
		})
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```

Then create a folder in current working directory called `views` and create a file called `views/index.html` and put below html

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>View Reader</title>
</head>
<body>
  <h1>Hello {{ name }} this is home page.</h1>
</body>
</html>
```


### Static Assets

HTTP static allow allows us to specify a folder containing all webpage assets like `CSS`, `JavaScript`, `Images` etc.

To get start lets create `static` folder and in `static` folder create a file called `main.css` - `static/main.css`.

```css
body {
  margin: 0 !important;
  padding: 0 !important;
  background-color: limegreen;
}

h1 {
  font-size: 5em;
  color: #fff;
  text-align: center;
  text-decoration: underline;
  font-family: fantasy;
  margin: 5px 0px !important;
}
```

Now let create view called `home.html` in the `views` - `views/home.html`.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Home Page</title>
  <link rel="stylesheet" href="main.css">
</head>
<body>
  <h1>Hello World!!!</h1>
</body>
</html>
```

Lets add `static` to our `server` but specifying `static` path.

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.SetView("views", "html").SetStatic("static")

	server.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.View("home", http.ViewData{})
	})

	fmt.Println("Server running ", server.Host())

	server.Listen()
}
```


#### Route Subdomain

HTTP allows subdomain to break up services e.g api and web.
If you are running you code you need to change you `hosts` file to 

If you are running you code on local machine you need to change you host file e.g

- Linux/MacOs - hosts file path `/etc/hosts` or `/private/etc/hosts`
127.0.0.1 api.example.com

- Windows - hosts file path `C:\Windows\System32\drivers\etc\hosts`
127.0.0.1:80 api.example.com

```go
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
		route.Get("users", func(req *http.Request, res *http.Response) *http.Response {
			return res.Json([]User{
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
			})
		})
	})

	fmt.Printf("Running server on %s", server.Host())

	server.Listen()
}
```

Subdomain also support dynamic parameters here is a example.

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func main() {
	server := http.Server("127.0.0.1", 80)

	server.Route().Subdomain("{company}", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			return res.Json(map[string]string{
				"company": req.Parameters.Get("company"),
			})
		})
	})

	fmt.Printf("Running server on %s", server.Host())

	server.Listen()
}
```


### Websocket

HTTP support websocket without third party packages

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/ws"
)

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Ws("/", func(req *http.Request, ws *ws.Ws) {
		ws.OnReady(func(ws *ws.Ws) {
			ws.OnMessage(func(data []byte) {
				fmt.Println("On Message:", string(data))
			})

			ws.OnPing(func(data []byte) {
				fmt.Println("On Ping:", string(data))
			})

			ws.OnPong(func(data []byte) {
				fmt.Println("On Pong:", string(data))
			})

			ws.OnClose(func(data []byte) {
				fmt.Println("On Close:", string(data))
			})

			ws.OnError(func(data []byte) {
				fmt.Println("On Error:", string(data))
			})
		})
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```

Lets example websocket callback which are:

- `OnReady`   - This event is called when then websocket handshake is complete.
- `OnMessage` - This event is called when the is when data send through websocket connection.
- `OnPing`    - This event is called when the is `ping` event through websocket connection.
- `OnPong`    - This event is called when the is `pong` event through websocket connection.
- `OnClose`   - This event is called when websocket connection is closed.
- `OnError`   - This event is called the is error through websocket connection.

The above is all about receiving data/events now lets create a simple websocket that send random coordinates every two seconds.

```go
package main
import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lucas11776-golang/http"
)

type Coordinate struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Altitude  int     `json:"altitude"`
}

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Ws("coordinate", func(req *http.Request, ws *http.Ws) {
		ws.OnReady(func(ws *http.Ws) {
			for {
				if !ws.Alive {
					break
				}

				time.Sleep(time.Second * 2)

				ws.WriteJson(Coordinate{
					Longitude: rand.Float32() * 360,
					Latitude:  rand.Float32() * 180,
					Altitude:  int(rand.Float32() * 100),
				})
			}
		})
	})

	fmt.Println("Server running ", server.Host())

	server.Listen()
}
```

Here is simple `JavaScript` code to listen to coordinates from the server.

```javascript
const ws = new WebSocket("ws://127.0.0.1:8080/coordinate");

ws.onopen = () => {
	ws.onmessage = e => console.log(JSON.parse(e.data));
}

ws.onclose = e => console.log(e)
```

The are three type of writes which are:

- `Write`       - This will write/send text payload to websocket connection.
- `WriteBinary` - This will write/send binary payload to websocket connection. 
- `WriteJson`   - This will convert `struct`/`map` to json string and write/send to connection.


### Middleware

What`s an application without middleware/guard to protected routes from unauthorized request or unwanted request below is simple route with middleware.

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

// Comment
func IsAuth(req *http.Request, res *http.Response, next http.Next) *http.Response {
	// Auth logic
	return next()
}

// Comment
func IsGuest(req *http.Request, res *http.Response, next http.Next) *http.Response {
	// Guest Logic
	return next()
}

func main() {
	server := http.Server("127.0.0.1", 8080)

	server.Route().Middleware(IsGuest).Group("authentication", func(route *http.Router) {
		route.Group("login", func(route *http.Router) {
			// Login routes...
		})
		route.Group("register", func(route *http.Router) {
			// Register routes...
		})
		route.Post("logout", func(req *http.Request, res *http.Response) *http.Response {
			return res.Redirect("authentication/login")
		}, IsGuest)
		// OR
		route.Post("logout", func(req *http.Request, res *http.Response) *http.Response {
			return res.Redirect("authentication/login")
		}, IsGuest).Middleware(IsAuth)
	})

	server.Route().Group("dashboard", func(route *http.Router) {
		// Dashboard routes...
	}, IsAuth)

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```

If you `visit` [127.0.0.1:8080](http://127.0.0.1:8080) with Postman or you favorite API testing tool without header `Auth-Key` with value of `test@123` you will get code status `401` with message `Unauthorized Access`.


Before we forget `subdomain` also supports middleware e.g

```go
package main

import (
	"fmt"

	"github.com/lucas11776-golang/http"
)

func IsAuth(req *http.Request, res *http.Response, next http.Next) *http.Response {
	// Auth logic
	return next()
}

func IsEmployee(req *http.Request, res *http.Response, next http.Next) *http.Response {
	// Auth logic
	return next()
}

func main() {
	server := http.Server("127.0.0.1", 80)

	server.Route().Middleware(IsAuth, IsEmployee).Subdomain("{company}", func(route *http.Router) {
		// Company routes
	}, IsAuth, IsEmployee)

	// OR

	server.Route().Subdomain("{company}", func(route *http.Router) {
		// Company routes
	}, IsAuth, IsEmployee)

	fmt.Printf("Running server on %s", server.Host())

	server.Listen()
}
```

### Session

HTTP support session to allow us to store user `data` like user ID, user role etc below is a simple code of how session works we will example everything below the sample code. 

```go
package main

import (
	"fmt"
	"os"

	"github.com/lucas11776-golang/http"
)

// Comment
func IsAuth(req *http.Request, res *http.Response, next http.Next) *http.Response {
	if req.Session.Get("user_id") != "" {
		return res.Redirect("/")
	}

	return next()
}

// Comment
func IsGuest(req *http.Request, res *http.Response, next http.Next) *http.Response {
	if req.Session.Get("user_id") != "" {
		return res.Redirect("/")
	}

	return next()
}

func main() {
	server := http.Server("127.0.0.1", 8080)

	// Initialize application session
	server.Session([]byte(os.Getenv("SESSION_KEY")))

	server.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.Html("<h1>Home Page</h1>")
	})

	server.Route().Group("authentication", func(route *http.Router) {
		route.Group("login", func(route *http.Router) {
			route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
				return res.Html("<h1>Add Post form to login</h1>")
			})
			route.Post("/", func(req *http.Request, res *http.Response) *http.Response {
				res.Session.Set("user_id", "1")

				return res.Redirect("dashboard")
			})
		}, IsGuest)

		route.Post("logout", func(req *http.Request, res *http.Response) *http.Response {
			res.Session.Remove("user_id")

			return res.Redirect("authentication/login")
		}).Middleware(IsAuth)
	})

	server.Route().Middleware(IsAuth).Group("dashboard", func(route *http.Router) {
		route.Get("/", func(req *http.Request, res *http.Response) *http.Response {
			return res.Html("<h1>Dashboard Page Can Be Viewed By Login Users</h1>")
		})
	})

	fmt.Printf("Server running %s", server.Host())

	server.Listen()
}
```

## Issues

Having issues with HTTP framework contact me on:

- Email    - [thembangubeni04@gmail.com](mailto:thembangubeni04@gmail.com)
- Linkedin - [lucas11776](https://linkedin.com/)