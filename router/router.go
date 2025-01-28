package router

import (
	"http/request"
	"http/response"
	"http/ws"
	"reflect"
	"strings"
)

type Next func() *response.Response

type Middleware func(req *request.Request, res *response.Response, next Next) *response.Response

type Route struct {
	path       string
	method     string
	middleware []Middleware
	router     *Router
	callback   reflect.Value
}

type Routes []*Route

type RouterGroup struct {
	web Routes
	ws  Routes
}

type Router struct {
	path       string
	middleware []Middleware
	routes     *RouterGroup
}

type Group func(router *Router)

type Web func(req *request.Request, res *response.Response) *response.Response

type Ws func(req *request.Request, ws *ws.Ws)

// Comment
func JoinPath(path ...string) string {
	arr := []string{}

	for _, p := range path {
		if p == "" {
			continue
		}

		arr = append(arr, strings.Trim(p, "/"))
	}

	return strings.Join(arr, "/")
}

// Comment
func (ctx *Route) Middleware(middleware ...Middleware) *Route {
	ctx.middleware = append(ctx.middleware, middleware...)

	return ctx
}

// Comment
func (ctx *Router) getRoute(router *Router, method string, uri string, callback reflect.Value) *Route {
	return &Route{
		path:     strings.Trim(uri, "/"),
		method:   strings.ToUpper(method),
		router:   router,
		callback: callback,
	}
}

// Comment
func (ctx *RouterGroup) Router() *Router {
	return &Router{routes: ctx}
}

// Comment
func (ctx *Router) Route(method string, uri string, callback Web) *Route {
	route := ctx.getRoute(ctx, method, JoinPath(ctx.path, uri), reflect.ValueOf(callback))

	ctx.routes.web = append(ctx.routes.web, route)

	return route
}

// Comment
func (ctx *Router) Group(prefix string, group Group) {
	group(&Router{
		path:       JoinPath(ctx.path, prefix),
		routes:     ctx.routes,
		middleware: ctx.middleware,
	})
}

// Comment
func (ctx *Router) Middleware(middlewares ...Middleware) *Router {
	ctx.middleware = append(ctx.middleware, middlewares...)

	return ctx
}

// Comment
func (ctx *Router) Get(uri string, callback Web) *Route {
	return ctx.Route("GET", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Post(uri string, callback Web) *Route {
	return ctx.Route("POST", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Put(uri string, callback Web) *Route {
	return ctx.Route("PUT", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Patch(uri string, callback Web) *Route {
	return ctx.Route("PATCH", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Delete(uri string, callback Web) *Route {
	return ctx.Route("DELETE", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Head(uri string, callback Web) *Route {
	return ctx.Route("HEAD", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Options(uri string, callback Web) *Route {
	return ctx.Route("OPTIONS", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Connect(uri string, callback Web) *Route {
	return ctx.Route("CONNECT", JoinPath(ctx.path, uri), callback).Middleware(ctx.middleware...)
}

// Comment
func (ctx *Router) Ws(uri string, callback Ws) *Route {
	route := ctx.getRoute(ctx, "GET", uri, reflect.ValueOf(callback))

	ctx.routes.ws = append(ctx.routes.ws, route)

	return route
}
