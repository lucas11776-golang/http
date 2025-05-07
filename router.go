package http

import (
	"reflect"
	"regexp"
	"strings"

	str "github.com/lucas11776-golang/http/utils/strings"
)

const (
	ParameterRegex string = "\\{[a-zA-Z_]+\\}"
)

type Next func() *Response

type Middleware func(req *Request, res *Response, next Next) *Response

type Parameters map[string]string

// Comment
func (ctx Parameters) Get(key string) string {
	value, ok := ctx[key]

	if !ok {
		return ""
	}

	return value
}

type Route struct {
	method     string
	path       []string
	parameters Parameters
	middleware []Middleware
	router     *Router
	callback   reflect.Value
}

type Routes []*Route

type RouterGroup struct {
	web      Routes
	ws       Routes
	fallback WebCallback
}

type Router struct {
	path       string
	middleware []Middleware
	routes     *RouterGroup
}

type GroupCallback func(route *Router)

type WebCallback func(req *Request, res *Response) *Response

type WsCallback func(req *Request, ws *Ws)

// Comment
func InitRouter() *RouterGroup {
	return &RouterGroup{}
}

// Comment
func (ctx *Route) Path() string {
	return strings.Join(ctx.path, "/")
}

// Comment
func (ctx *Route) Method() string {
	return ctx.method
}

// Comment
func (ctx *Route) Middlewares() []Middleware {
	return ctx.middleware
}

// Comment
func (ctx *Route) setParameters(req *Request, parameters Parameters) *Route {
	ctx.parameters = parameters
	req.Parameters = parameters

	return ctx
}

// Comment
func (ctx *Route) Parameters() Parameters {
	return ctx.parameters
}

// Comment
func (ctx *Route) Parameter(parameter string) string {
	param, ok := ctx.parameters[parameter]

	if !ok {
		return ""
	}

	return param
}

// Comment
func (ctx *Route) Middleware(middleware ...Middleware) *Route {
	ctx.middleware = append(ctx.middleware, middleware...)

	return ctx
}

// Comment
func (ctx *Route) Call(value ...reflect.Value) *Response {
	rt := ctx.callback.Call(value)

	if len(rt) == 0 {
		return nil
	}

	switch rt[0].Interface().(type) {
	case *Response:
		return rt[0].Interface().(*Response)
	default:
		return nil
	}
}

// Comment
func parametersRouteMatch(route *Route, path []string) (Parameters, bool) {
	regex, _ := regexp.Compile(ParameterRegex)
	parameters := make(Parameters)

	for i, seg := range path {
		if i >= len(route.path) {
			return nil, false
		}

		if route.path[i] == "*" {
			return parameters, true
		}

		if seg == route.path[i] {
			continue
		}

		if regex.Match([]byte(route.path[i])) {
			name := strings.Trim(strings.Trim(route.path[i], "{"), "}")
			value := path[i]

			// TODO Will add route parameter match regex in future
			// if !regex_match {
			// 	return nil, false
			// }

			parameters[name] = value

			continue
		}

		return nil, false
	}

	return parameters, true
}

// comment
func routeMatch(routes Routes, method string, uri string) (*Route, Parameters) {
	path := strings.Split(strings.Trim(uri, "/"), "/")

	for _, route := range routes {
		if strings.ToUpper(method) != route.Method() {
			continue
		}

		if route.Path() == "*" {
			return route, make(Parameters)
		}

		if strings.Trim(uri, "/") == route.Path() {
			return route, make(Parameters)
		}

		parameters, ok := parametersRouteMatch(route, path)

		if !ok {
			continue
		}

		return route, parameters
	}

	return nil, nil
}

// Comment
func (ctx *RouterGroup) MatchWebRoute(req *Request) *Route {
	route, parameters := routeMatch(ctx.web, req.Method, req.Path())

	if route == nil {
		return nil
	}

	return route.setParameters(req, parameters)
}

// Comment
func (ctx *RouterGroup) MatchWsRoute(req *Request) *Route {
	route, parameters := routeMatch(ctx.ws, req.Method, req.Path())

	if route == nil {
		return nil
	}

	return route.setParameters(req, parameters)

}

// Comment
func (ctx *Router) getRoute(router *Router, method string, uri string, callback reflect.Value) *Route {
	return &Route{
		method:     strings.ToUpper(method),
		path:       strings.Split(str.JoinPath(ctx.path, uri), "/"),
		parameters: make(Parameters),
		middleware: ctx.middleware,
		router:     router,
		callback:   callback,
	}
}

// Comment
func (ctx *RouterGroup) Router() *Router {
	return &Router{routes: ctx}
}

// Comment
func (ctx *Router) Route(method string, uri string, callback WebCallback) *Route {
	route := ctx.getRoute(ctx, method, uri, reflect.ValueOf(callback))

	ctx.routes.web = append(ctx.routes.web, route)

	return route
}

// Comment
func (ctx *Router) Group(prefix string, group GroupCallback) *Router {
	router := &Router{
		path:       str.JoinPath(ctx.path, prefix),
		routes:     ctx.routes,
		middleware: ctx.middleware,
	}

	group(router)

	return router
}

// Comment
func (ctx *Router) Middleware(middlewares ...Middleware) *Router {
	ctx.middleware = append(ctx.middleware, middlewares...)

	return ctx
}

// Comment
func (ctx *Router) Get(uri string, callback WebCallback) *Route {
	return ctx.Route("GET", uri, callback)
}

// Comment
func (ctx *Router) Post(uri string, callback WebCallback) *Route {
	return ctx.Route("POST", uri, callback)
}

// Comment
func (ctx *Router) Put(uri string, callback WebCallback) *Route {
	return ctx.Route("PUT", uri, callback)
}

// Comment
func (ctx *Router) Patch(uri string, callback WebCallback) *Route {
	return ctx.Route("PATCH", uri, callback)
}

// Comment
func (ctx *Router) Delete(uri string, callback WebCallback) *Route {
	return ctx.Route("DELETE", uri, callback)
}

// Comment
func (ctx *Router) Head(uri string, callback WebCallback) *Route {
	return ctx.Route("HEAD", uri, callback)
}

// Comment
func (ctx *Router) Options(uri string, callback WebCallback) *Route {
	return ctx.Route("OPTIONS", uri, callback)
}

// Comment
func (ctx *Router) Connect(uri string, callback WebCallback) *Route {
	return ctx.Route("CONNECT", uri, callback)
}

// Comment
func (ctx *Router) Ws(uri string, callback WsCallback) *Route {
	route := ctx.getRoute(ctx, "GET", uri, reflect.ValueOf(callback))

	ctx.routes.ws = append(ctx.routes.ws, route)

	return route
}

// Comment
func (ctx *Router) Fallback(fallback WebCallback) {
	ctx.routes.fallback = fallback
}
