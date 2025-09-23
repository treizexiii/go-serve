package server

import (
	"net/http"
)

type Http_Method string

const (
	GET    Http_Method = "GET"
	POST   Http_Method = "POST"
	PUT    Http_Method = "PUT"
	DELETE Http_Method = "DELETE"
)

type Route struct {
	method      Http_Method
	path        string
	handler     RouteHandler
	queryParams []string
	tags        []string
}

type RouteHandler struct {
	handler     HandlerFunc
	middlewares []MiddlewareFunc
	meta        map[string]interface{}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)
type MiddlewareFunc func(http.Handler) http.Handler

func CreateRoute(method Http_Method, path string, handler HandlerFunc) RouteInfo {
	return &Route{
		method: method,
		path:   path,
		tags:   []string{},
		handler: RouteHandler{
			handler:     handler,
			middlewares: []MiddlewareFunc{},
			meta:        make(map[string]interface{}),
		},
	}
}

func CreateGET(path string, handler HandlerFunc) RouteInfo {
	return CreateRoute(GET, path, handler)
}

func CreatePOST(path string, handler HandlerFunc) RouteInfo {
	return CreateRoute(POST, path, handler)
}

func CreatePUT(path string, handler HandlerFunc) RouteInfo {
	return CreateRoute(PUT, path, handler)
}

func CreateDELETE(path string, handler HandlerFunc) RouteInfo {
	return CreateRoute(DELETE, path, handler)
}

func (r *Route) GetMethod() Http_Method {
	return r.method
}

func (r *Route) GetPath() string {
	return r.path
}

func (r *Route) GetHandler() *RouteHandler {
	return &r.handler
}

// WithMeta implements RouteInfo.
func (r *Route) WithMeta(key string, value interface{}) RouteInfo {
	if r.handler.meta == nil {
		r.handler.meta = make(map[string]interface{})
	}
	r.handler.meta[key] = value
	return r
}

// WithTags implements RouteInfo.
func (r *Route) WithTags(tags ...string) RouteInfo {
	r.tags = append(r.tags, tags...)
	return r
}

func (r *Route) WithMiddleware(middlewares ...MiddlewareFunc) RouteInfo {
	r.handler.middlewares = append(r.handler.middlewares, middlewares...)
	return r
}
