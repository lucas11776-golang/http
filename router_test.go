package http

import (
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
)

// Comment
func makeRequest(method Method, path string) *Request {
	req, _ := NewRequest(method, path, "HTTP/1.1", make(types.Headers), strings.NewReader(""))

	return req
}

func TestRouter(t *testing.T) {
	t.Run("TestRouterAddWebRouteUsingRoute", func(t *testing.T) {
		router := &RouterGroup{}

		route := router.Router().Route("GET", "/products/{id}", func(req *Request, res *Response) *Response {
			return res.SetBody([]byte("<h1>Hello World</h1>")).SetHeader("content-type", "text/html")
		})

		route.Middleware(func(req *Request, res *Response, next Next) *Response {
			return next()
		})

		TestingRoute(t, router.web, route, 0, "GET", "products/{id}", 1)
	})

	t.Run("TestRouterAddWsRouteUsingRoute", func(t *testing.T) {
		router := &RouterGroup{}

		route := router.Router().Ws("/position/moving", func(req *Request, ws *Ws) {
			// Ws staff
		})

		route.Middleware(func(req *Request, res *Response, next Next) *Response {
			return next()
		})

		TestingRoute(t, router.ws, route, 0, "GET", "position/moving", 1)
	})

	t.Run("TestRouterWebGroup", func(t *testing.T) {
		router := &RouterGroup{}

		middleware := func(req *Request, res *Response, next Next) *Response {
			return next()
		}

		router.Router().Middleware(middleware).Group("/api", func(router *Router) {
			router.Options("/*", func(req *Request, res *Response) *Response {
				return res
			})
			router.Group("/products", func(router *Router) {
				router.Get("/", func(req *Request, res *Response) *Response {
					return res
				})
				router.Post("/", func(req *Request, res *Response) *Response {
					return res
				})
				router.Put("/", func(req *Request, res *Response) *Response {
					return res
				})
				router.Patch("/{id}", func(req *Request, res *Response) *Response {
					return res
				})
				router.Delete("/{id}", func(req *Request, res *Response) *Response {
					return res
				})
				router.Head("/{id}", func(req *Request, res *Response) *Response {
					return res
				})
			})
		})

		router.Router().Connect("/*", func(req *Request, res *Response) *Response {
			return res
		})

		TestingRoute(t, router.web, router.web[0], 0, "OPTIONS", "api/*", 1)
		TestingRoute(t, router.web, router.web[1], 1, "GET", "api/products", 1)
		TestingRoute(t, router.web, router.web[2], 2, "POST", "api/products", 1)
		TestingRoute(t, router.web, router.web[3], 3, "PUT", "api/products", 1)
		TestingRoute(t, router.web, router.web[4], 4, "PATCH", "api/products/{id}", 1)
		TestingRoute(t, router.web, router.web[5], 5, "DELETE", "api/products/{id}", 1)
		TestingRoute(t, router.web, router.web[6], 6, "HEAD", "api/products/{id}", 1)
		TestingRoute(t, router.web, router.web[7], 7, "CONNECT", "*", 0)
	})

	t.Run("TestRouterWsGroup", func(t *testing.T) {
		router := &RouterGroup{}

		middleware := func(req *Request, res *Response, next Next) *Response {
			return next()
		}

		router.Router().Middleware(middleware).Group("/position", func(router *Router) {
			router.Ws("/change", func(req *Request, ws *Ws) {
				// Ws staff
			})
		})

		TestingRoute(t, router.ws, router.ws[0], 0, "GET", "position/change", 1)
	})

	t.Run("TestRouterWebMatch", func(t *testing.T) {
		router := &RouterGroup{}

		router.Router().Connect("/*", func(req *Request, res *Response) *Response {
			return res
		})
		router.Router().Group("/api", func(router *Router) {
			router.Options("/*", func(req *Request, res *Response) *Response {
				return res
			})
			router.Group("/products", func(router *Router) {
				router.Get("/", func(req *Request, res *Response) *Response {
					return res
				})
				router.Group("{id}", func(router *Router) {
					router.Get("/", func(req *Request, res *Response) *Response {
						return res
					})
					router.Post("/", func(req *Request, res *Response) *Response {
						return res
					})
				})
			})
		})

		TestingRoute(t, router.web, router.MatchWebRoute(makeRequest("CONNECT", "api/products")), 0, "CONNECT", "*", 0)
		TestingRoute(t, router.web, router.MatchWebRoute(makeRequest("OPTIONS", "api/products/50")), 1, "OPTIONS", "api/*", 0)
		TestingRoute(t, router.web, router.MatchWebRoute(makeRequest("GET", "api/products")), 2, "GET", "api/products", 0)
		TestingRoute(t, router.web, router.MatchWebRoute(makeRequest("GET", "api/products/20")), 3, "GET", "api/products/{id}", 0)
		TestingRoute(t, router.web, router.MatchWebRoute(makeRequest("POST", "api/products/20")), 4, "POST", "api/products/{id}", 0)

		route := router.MatchWebRoute(makeRequest("POST", "api/products/203"))

		if route.Parameter("id") != "203" {
			t.Fatalf("The route parameter is id is not %s but got %s", "203", route.Parameter("id"))
		}
	})

	t.Run("TestRouterWsMatch", func(t *testing.T) {
		router := &RouterGroup{}

		router.Router().Group("devices", func(router *Router) {
			router.Group("/{device}", func(router *Router) {
				router.Ws("/", func(req *Request, ws *Ws) {
					// Ws staff
				})
				router.Ws("/position", func(req *Request, ws *Ws) {
					// Ws staff
				})
			})
		})

		TestingRoute(t, router.ws, router.MatchWsRoute(makeRequest("GET", "devices/R833WC0GL3CF")), 0, "GET", "devices/{device}", 0)
		TestingRoute(t, router.ws, router.MatchWsRoute(makeRequest("GET", "devices/R833WC0GL3CF/position")), 1, "GET", "devices/{device}/position", 0)

		route := router.MatchWsRoute(makeRequest("GET", "devices/R833WC0GL3CF/position"))

		if route.Parameter("device") != "R833WC0GL3CF" {
			t.Fatalf("The route parameter is id is not %s but got %s", "1", route.Parameter("id"))
		}
	})
}

// Comment
func TestingRoute(t *testing.T, routes Routes, route *Route, index int, method string, path string, middlewares int) {
	if len(routes) == 0 {
		t.Fatalf("Route is not add to web routes")
	}

	if route != routes[index] {
		t.Fatalf("Route in web routes does not match web route: %p return route: %p", route, routes[index])
	}

	if route.Method() != strings.ToUpper(method) {
		t.Fatalf("Route is not a %s method is %s method", strings.ToUpper(method), route.Method())
	}

	if route.Path() != path {
		t.Fatalf("Route path is not %s is %s", path, route.Path())
	}

	if len(route.Middlewares()) < middlewares {
		t.Fatalf("Middleware is not add to route")
	}
}
