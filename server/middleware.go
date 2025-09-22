package server

type DecoratorFunc func(HandlerFunc) RouteHandler

type MiddlewareDecoratorFunc func(...interface{}) DecoratorFunc

type MiddlewareInfo struct {
	Name       string
	Middleware MiddlewareFunc
}
