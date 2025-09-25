package routes

import (
	"goserve/middlewares"
	"goserve/responses"
	"net/http"
)

type RouteInfo interface {
	WithTags(tags ...string) RouteInfo
	WithMeta(key string, value interface{}) RouteInfo
	WithMiddleware(middlewares ...middlewares.MiddlewareFunc) RouteInfo

	GetPath() string
	GetMethod() Http_Method
	GetHandler() *ActionFunc
	GetMiddlewares() []middlewares.MiddlewareFunc
}

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
	tags        []string
	action      ActionFunc
	middlewares []middlewares.MiddlewareFunc
	meta        map[string]interface{}
	// queryParams []string
}

type ActionFunc func(r *http.Request) (responses.ApiResponse, error)

func CreateRoute(method Http_Method, path string, action ActionFunc) RouteInfo {
	return &Route{
		method:      method,
		path:        path,
		tags:        []string{},
		action:      action,
		middlewares: []middlewares.MiddlewareFunc{},
		meta:        make(map[string]interface{}),
	}
}

func CreateGET(path string, handler ActionFunc) RouteInfo {
	return CreateRoute(GET, path, handler)
}

func CreatePOST(path string, handler ActionFunc) RouteInfo {
	return CreateRoute(POST, path, handler)
}

func CreatePUT(path string, handler ActionFunc) RouteInfo {
	return CreateRoute(PUT, path, handler)
}

func CreateDELETE(path string, handler ActionFunc) RouteInfo {
	return CreateRoute(DELETE, path, handler)
}

func (r *Route) GetMethod() Http_Method {
	return r.method
}

func (r *Route) GetPath() string {
	return r.path
}

func (r *Route) GetHandler() *ActionFunc {
	return &r.action
}

func (r *Route) GetMiddlewares() []middlewares.MiddlewareFunc {
	return r.middlewares
}

func (r *Route) WithMeta(key string, value interface{}) RouteInfo {
	if r.meta == nil {
		r.meta = make(map[string]interface{})
	}
	r.meta[key] = value
	return r
}

func (r *Route) WithTags(tags ...string) RouteInfo {
	r.tags = append(r.tags, tags...)
	return r
}

func (r *Route) WithMiddleware(middlewares ...middlewares.MiddlewareFunc) RouteInfo {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}
