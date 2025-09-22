package server

import (
	"goserve/configuration"
	"net/http"
)

type ServerBuilder interface {
	// Add a configuration to the server
	WithConfiguration(config configuration.Configuration) ServerBuilder

	// Set the port for the server (default: 8080)
	SetPort(port int) ServerBuilder

	// Add a single route to the server
	AddRoute(method Http_Method, path string, handler HandlerFunc) ServerBuilder
	// Add multiple routes to the server
	AddRoutes(routes []Route) ServerBuilder
	// Create and add a GET route to the server
	GET(path string, handler HandlerFunc) ServerBuilder
	// Create and add a POST route to the server
	POST(path string, handler HandlerFunc) ServerBuilder
	// Create and add a PUT route to the server
	PUT(path string, handler HandlerFunc) ServerBuilder
	// Create and add a DELETE route to the server
	DELETE(path string, handler HandlerFunc) ServerBuilder

	// Add middleware for all routes
	AddGlobalMiddleware(name string, middleware MiddlewareFunc) ServerBuilder
	// Add middleware for specific paths
	AddRouteMiddleware(name string, middleware MiddlewareFunc, Path []string) ServerBuilder
	// Add middleware for specific methods
	AddMethodMiddleware(name string, middleware MiddlewareFunc, Method []string) ServerBuilder

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
