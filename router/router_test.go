package router

import (
	"fmt"
	"http/request"
	"http/response"
	"http/ws"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	t.Run("TestRouterAddWebRouteUsingRoute", func(t *testing.T) {
		router := &RouterGroup{}

		route := router.Router().Route("GET", "/products/1", func(req *request.Request, res *response.Response) *response.Response {
			return res.Body([]byte("<h1>Hello World</h1>")).Header("content-type", "text/html")
		})

		if len(router.web) != 1 {
			t.Errorf("Route is not add to web routes")
		}

		if route != router.web[0] {
			t.Errorf("Route in web routes does not match web route: %p return route: %p", route, router.web[0])
		}

		fmt.Println("ROUTER NAME:", route.path)

		if strings.Trim(route.path, "/") != "products/1" {
			t.Errorf("Route path is not %s is %s", "products/1", route.path)
		}

		if route.method != "GET" {
			t.Errorf("Route is not a %s method is %s method", "GET", route.method)
		}

		route.Middleware(func(req *request.Request, res *response.Response, next Next) *response.Response {
			return next()
		})

		if len(route.middleware) != 1 && len(router.web[0].middleware) != 1 {
			t.Errorf("Middleware is not add to route")
		}
	})

	t.Run("TestRouterAddWsRouteUsingRoute", func(t *testing.T) {
		router := &RouterGroup{}

		route := router.Router().Ws("/position", func(req *request.Request, ws *ws.Ws) {

		})

		if len(router.ws) != 1 {
			t.Errorf("Route is not add to web routes")
		}

		if route != router.ws[0] {
			t.Errorf("Route in web routes does not match web route: %p return route: %p", route, router.ws[0])
		}

		if strings.Trim(route.path, "/") != "position" {
			t.Errorf("Route path is not %s is %s", "position", route.path)
		}

		if route.method != "GET" {
			t.Errorf("Route is not a %s method is %s method", "GET", route.method)
		}

		route.Middleware(func(req *request.Request, res *response.Response, next Next) *response.Response {
			return next()
		})

		if len(route.middleware) != 1 && len(router.ws[0].middleware) != 1 {
			t.Errorf("Middleware is not add to route")
		}
	})

	t.Run("TestRouterGroup", func(t *testing.T) {
		router := &RouterGroup{}

		guard := func(req *request.Request, res *response.Response, next Next) *response.Response {
			return next()
		}

		router.Router().Middleware(guard).Group("/api", func(router *Router) {
			router.Post("/authentication/register", func(req *request.Request, res *response.Response) *response.Response {
				return res.Body([]byte("<h1>Hello World</h1>")).Header("content-type", "text/html")
			})
		})

		if len(router.web) != 1 {
			t.Errorf("Route is not added to web router")
		}

		if router.web[0].method != "POST" {
			t.Errorf("Route is not %s method", "POST")
		}

		if len(router.web[0].middleware) != 1 {
			t.Errorf("Group route middleware is not added")
		}
	})

	t.Run("TestRouterMatch", func(t *testing.T) {
		router := &RouterGroup{}

		router.Router().Group("/api", func(route *Router) {
			route.Post("/products/{id}", func(req *request.Request, res *response.Response) *response.Response {
				return res.Body([]byte("<h1>Hello World</h1>")).Header("content-type", "text/html")
			})
		})

		router.Router().Group("/position", func(router *Router) {
			router.Ws("/move", func(req *request.Request, ws *ws.Ws) {
			})
		})

		// Match routes
	})
}
