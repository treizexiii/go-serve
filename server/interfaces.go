package server

import (
	"goserve/configuration"
	"net/http"
)

type RouteInfo interface {
	WithTags(tags ...string) RouteInfo
	WithMeta(key string, value interface{}) RouteInfo
	WithMiddleware(middlewares ...MiddlewareFunc) RouteInfo

	GetPath() string
	GetMethod() Http_Method
	GetHandler() *RouteHandler
}

type ServerBuilder interface {
	// Add a configuration to the server
	WithConfiguration(config configuration.Configuration) ServerBuilder

	// Set the port for the server (default: 8080)
	SetPort(port int) ServerBuilder

	AddGlobalMiddleware(name string, middleware MiddlewareFunc) ServerBuilder
	WithJSONSerialization() ServerBuilder
	WithLogging(logRequests, logResponses bool) ServerBuilder

	// Add a single route to the server
	AddRoute(method Http_Method, path string, handler HandlerFunc) ServerBuilder
	// Add multiple routes to the server
	AddRoutes(routes []RouteInfo) ServerBuilder
	// Create and add a GET route to the server
	GET(path string, handler HandlerFunc) ServerBuilder
	// Create and add a POST route to the server
	POST(path string, handler HandlerFunc) ServerBuilder
	// Create and add a PUT route to the server
	PUT(path string, handler HandlerFunc) ServerBuilder
	// Create and add a DELETE route to the server
	DELETE(path string, handler HandlerFunc) ServerBuilder

	// Build and return the configured http server
	Build() HttpServer
}

type HttpServer interface {
	// Get the underlying http.Server instance
	GetHttpServer() *http.Server
	// Get the underlying http.ServeMux instance
	GetMux() *http.ServeMux

	// Start the HTTP server
	Start() error
}
