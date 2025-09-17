package http

import (
	"net"
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
	subdomain   string
	path        string
	middlewares []Middleware
	routes      *RouterGroup
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
func parametersRouteMatch(routePath []string, requestPath []string) (found bool, parameters Parameters) {
	regex, _ := regexp.Compile(ParameterRegex)
	params := make(Parameters)

	for i, seg := range requestPath {
		if i >= len(routePath) {
			return false, nil
		}

		if routePath[i] == "*" {
			return true, params
		}

		if seg == routePath[i] {
			continue
		}

		if regex.Match([]byte(routePath[i])) {
			// TODO: requestPath will be empty == ""
			params[strings.Trim(strings.Trim(routePath[i], "{"), "}")] = requestPath[i]

			continue
		}

		return false, nil
	}

	return true, params
}

// Comment
func getSubdomain(host string) string {
	host = strings.Split(host, ":")[0]

	if net.ParseIP(host) != nil {
		return ""
	}

	parts := strings.Split(host, ".")

	if len(parts) < 3 {
		return ""
	}

	return strings.Join(parts[:len(parts)-2], ".")
}

// TODO: move to utils...
// Comment
func mergeMap(m map[string]string, maps ...map[string]string) map[string]string {
	for _, mp := range maps {
		for k, v := range mp {
			m[k] = v
		}
	}

	return m
}

// comment
func routeMatch(routes Routes, req *Request) (*Route, Parameters) {
	method := req.Method
	uri := req.Path()
	path := strings.Split(strings.Trim(uri, "/"), "/")
	requestSubdomain := strings.Split(getSubdomain(req.Host), ".")

	for _, route := range routes {
		if strings.ToUpper(method) != route.Method() {
			continue
		}

		parameters := make(Parameters)
		routeSubdomain := strings.Split(route.router.subdomain, ".")

		ok, params := parametersRouteMatch(routeSubdomain, requestSubdomain)

		if !ok {
			continue
		}

		parameters = mergeMap(parameters, params)

		if route.Path() == "*" {
			return route, parameters
		}

		if strings.Trim(uri, "/") == route.Path() {
			return route, parameters
		}

		ok, params = parametersRouteMatch(route.path, path)

		if !ok {
			continue
		}

		parameters = mergeMap(parameters, params)

		return route, parameters
	}

	return nil, nil
}

// Comment
func (ctx *RouterGroup) MatchWebRoute(req *Request) *Route {
	route, parameters := routeMatch(ctx.web, req)

	if route != nil {
		req.Parameters = parameters
	}
	return route
}

// Comment
func (ctx *RouterGroup) MatchWsRoute(req *Request) *Route {
	route, parameters := routeMatch(ctx.ws, req)

	if route != nil {
		req.Parameters = parameters
	}

	return route

}

// Comment
func (ctx *Router) getRoute(router *Router, method string, uri string, callback reflect.Value, middleware ...Middleware) *Route {
	return &Route{
		method:     strings.ToUpper(method),
		path:       strings.Split(str.JoinPath(ctx.path, uri), "/"),
		middleware: append(ctx.middlewares, middleware...),
		router:     router,
		callback:   callback,
	}
}

// Comment
func (ctx *RouterGroup) Router() *Router {
	return &Router{routes: ctx}
}

// Comment
func (ctx *Router) Route(method string, uri string, callback WebCallback, middleware ...Middleware) *Route {
	route := ctx.getRoute(ctx, method, uri, reflect.ValueOf(callback), middleware...)

	ctx.routes.web = append(ctx.routes.web, route)

	return route
}

func (ctx *Router) Callback(group GroupCallback, middleware ...Middleware) {
	group(&Router{
		path:        ctx.path,
		routes:      ctx.routes,
		middlewares: append(ctx.middlewares, middleware...),
	})
}

// Comment
func (ctx *Router) Subdomain(subdomain string, group GroupCallback, middleware ...Middleware) {
	group(&Router{
		subdomain:   subdomain,
		path:        ctx.path,
		routes:      ctx.routes,
		middlewares: append(ctx.middlewares, middleware...),
	})
}

// Comment
func (ctx *Router) Group(prefix string, group GroupCallback, middleware ...Middleware) {
	group(&Router{
		subdomain:   ctx.subdomain,
		path:        str.JoinPath(ctx.path, prefix),
		routes:      ctx.routes,
		middlewares: append(ctx.middlewares, middleware...),
	})
}

// Comment
func (ctx *Router) Middleware(middlewares ...Middleware) *Router {
	ctx.middlewares = append(ctx.middlewares, middlewares...)

	return ctx
}

// Comment
func (ctx *Router) Get(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("GET", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Post(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("POST", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Put(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("PUT", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Patch(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("PATCH", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Delete(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("DELETE", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Head(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("HEAD", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Options(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("OPTIONS", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Connect(uri string, callback WebCallback, middleware ...Middleware) *Route {
	return ctx.Route("CONNECT", uri, callback, middleware...)
}

// Comment
func (ctx *Router) Ws(uri string, callback WsCallback, middleware ...Middleware) *Route {
	route := ctx.getRoute(ctx, "GET", uri, reflect.ValueOf(callback), middleware...)

	ctx.routes.ws = append(ctx.routes.ws, route)

	return route
}

// Comment
func (ctx *Router) Fallback(fallback WebCallback) {
	ctx.routes.fallback = fallback
}
