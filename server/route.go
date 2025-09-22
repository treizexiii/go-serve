package server

import "net/http"

type Http_Method string

const (
	GET    Http_Method = "GET"
	POST   Http_Method = "POST"
	PUT    Http_Method = "PUT"
	DELETE Http_Method = "DELETE"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Route struct {
	Method  Http_Method
	Path    string
	Handler HandlerFunc
}

func CreateGET(path string, handler HandlerFunc) Route {
	return Route{
		Method:  GET,
		Path:    path,
		Handler: handler,
	}
}

func CreatePOST(path string, handler HandlerFunc) Route {
	return Route{
		Method:  POST,
		Path:    path,
		Handler: handler,
	}
}

func CreatePUT(path string, handler HandlerFunc) Route {
	return Route{
		Method:  PUT,
		Path:    path,
		Handler: handler,
	}
}

func CreateDELETE(path string, handler HandlerFunc) Route {
	return Route{
		Method:  DELETE,
		Path:    path,
		Handler: handler,
	}
}
