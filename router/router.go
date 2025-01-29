package router

import (
	"http/request"
	"http/response"
	"http/ws"
	"reflect"
	"regexp"
	"strings"
)

type Next func() *response.Response

type Middleware func(req *request.Request, res *response.Response, next Next) *response.Response

type Parameters map[string]string

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
		if p == "" || p == "/" {
			continue
		}

		arr = append(arr, strings.Trim(p, "/"))
	}

	if len(arr) == 0 {
		arr[0] = ""
	}

	return strings.Join(arr, "/")
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
func (ctx *Route) setParameters(parameters Parameters) *Route {
	ctx.parameters = parameters

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
func (ctx *Route) Call(value ...reflect.Value) []byte {
	rt := ctx.callback.Call(value)

	if len(rt) == 0 {
		return []byte("")
	}

	switch rt[0].Type().String() {
	case "*response.Response":
		return []byte(response.ParseHttp((rt[0].Interface().(*response.Response))))
	default:
		return []byte("")
	}
}

// comment
func routeMatch(routes Routes, method string, uri string) (*Route, Parameters) {
	parameters := make(Parameters)
	path := strings.Split(strings.Trim(uri, "/"), "/")

	for _, route := range routes {
		if strings.ToUpper(method) != route.Method() {
			continue
		}

		if route.Path() == "*" {
			return route, parameters
		}

		regexGlobal, _ := regexp.Compile("[\\*]")

		if len(path) != len(route.path) && (path[0] == route.path[0] && regexGlobal.Match([]byte(route.Path())) == false) {
			continue
		}

		if strings.Trim(uri, "/") == route.Path() {
			return route, parameters
		}

		regexParameter, _ := regexp.Compile("\\{[a-zA-Z_]+\\}")

		for i, segment := range route.path {
			if segment == "*" {
				return route, parameters
			}

			if segment == path[i] {
				continue
			}

			if regexParameter.Match([]byte(segment)) {
				parameters[strings.Trim(strings.Trim(segment, "{"), "}")] = path[i]

				continue
			}

			return nil, nil
		}

		return route, parameters
	}

	return nil, nil
}

// Comment
func (ctx *RouterGroup) MatchWebRoute(method string, uri string) *Route {
	route, parameters := routeMatch(ctx.web, method, uri)

	if route == nil {
		return nil
	}

	return route.setParameters(parameters)
}

// Comment
func (ctx *RouterGroup) MatchWsRoute(uri string) *Route {
	route, parameters := routeMatch(ctx.ws, "GET", uri)

	if route == nil {
		return nil
	}

	return route.setParameters(parameters)

}

// Comment
func (ctx *Router) getRoute(router *Router, method string, uri string, callback reflect.Value) *Route {
	return &Route{
		method:     strings.ToUpper(method),
		path:       strings.Split(JoinPath(ctx.path, uri), "/"),
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
func (ctx *Router) Route(method string, uri string, callback Web) *Route {
	route := ctx.getRoute(ctx, method, uri, reflect.ValueOf(callback))

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
	return ctx.Route("GET", uri, callback)
}

// Comment
func (ctx *Router) Post(uri string, callback Web) *Route {
	return ctx.Route("POST", uri, callback)
}

// Comment
func (ctx *Router) Put(uri string, callback Web) *Route {
	return ctx.Route("PUT", uri, callback)
}

// Comment
func (ctx *Router) Patch(uri string, callback Web) *Route {
	return ctx.Route("PATCH", uri, callback)
}

// Comment
func (ctx *Router) Delete(uri string, callback Web) *Route {
	return ctx.Route("DELETE", uri, callback)
}

// Comment
func (ctx *Router) Head(uri string, callback Web) *Route {
	return ctx.Route("HEAD", uri, callback)
}

// Comment
func (ctx *Router) Options(uri string, callback Web) *Route {
	return ctx.Route("OPTIONS", uri, callback)
}

// Comment
func (ctx *Router) Connect(uri string, callback Web) *Route {
	return ctx.Route("CONNECT", uri, callback)
}

// Comment
func (ctx *Router) Ws(uri string, callback Ws) *Route {
	route := ctx.getRoute(ctx, "GET", uri, reflect.ValueOf(callback))

	ctx.routes.ws = append(ctx.routes.ws, route)

	return route
}
