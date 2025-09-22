package server

import (
	"fmt"
	"goserve/configuration"
	"log"
	"net/http"
	"time"
)

type builder struct {
	mux          *http.ServeMux
	port         int
	address      string
	routes       []Route
	middlewares  []Middleware
	config       configuration.Configuration
	readTimeout  int
	writeTimeout int
	idleTimeout  int
}

// Constructor

func New() ServerBuilder {
	return &builder{
		mux:          http.NewServeMux(),
		routes:       make([]Route, 0),
		address:      "",
		port:         8080, // Default port
		middlewares:  make([]Middleware, 0),
		config:       nil,
		readTimeout:  15,
		writeTimeout: 15,
		idleTimeout:  60,
	}
}

func (s *builder) WithConfiguration(config configuration.Configuration) ServerBuilder {
	s.config = config
	return s
}

func (s *builder) SetPort(port int) ServerBuilder {
	s.port = port
	return s
}

func (s *builder) AddRoute(method Http_Method, path string, handler HandlerFunc) ServerBuilder {
	s.routes = append(s.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})

	return s
}

func (s *builder) AddRoutes(routes []Route) ServerBuilder {
	s.routes = append(s.routes, routes...)
	return s
}

func (s *builder) GET(path string, handler HandlerFunc) ServerBuilder {
	return s.addRoute(CreateGET(path, handler))
}

func (s *builder) POST(path string, handler HandlerFunc) ServerBuilder {
	return s.addRoute(CreatePOST(path, handler))
}

func (s *builder) PUT(path string, handler HandlerFunc) ServerBuilder {
	return s.addRoute(CreatePUT(path, handler))
}

func (s *builder) DELETE(path string, handler HandlerFunc) ServerBuilder {
	return s.addRoute(CreateDELETE(path, handler))
}

func (s *builder) addRoute(route Route) ServerBuilder {
	s.routes = append(s.routes, route)
	return s
}

func (s *builder) AddGlobalMiddleware(name string, middleware MiddlewareFunc) ServerBuilder {
	s.middlewares = append(s.middlewares, Middleware{
		Name:       name,
		Middleware: middleware,
		Path:       nil,
		Method:     nil,
	})
	return s
}

func (s *builder) AddRouteMiddleware(name string, middleware MiddlewareFunc, Path []string) ServerBuilder {
	s.middlewares = append(s.middlewares, Middleware{
		Name:       name,
		Middleware: middleware,
		Path:       Path,
		Method:     nil,
	})
	return s
}

func (s *builder) AddMethodMiddleware(name string, middleware MiddlewareFunc, Method []string) ServerBuilder {
	s.middlewares = append(s.middlewares, Middleware{
		Name:       name,
		Middleware: middleware,
		Path:       nil,
		Method:     Method,
	})
	return s
}

func (s *builder) registerRoute(route Route) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != string(route.Method) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		route.Handler(w, r)
	})

	finalHandler := s.applyMiddlewares(handler, route.Path, string(route.Method))

	s.mux.Handle(route.Path, finalHandler)
}

func (s *builder) applyMiddlewares(handler http.Handler, path, method string) http.Handler {
	result := handler

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		middleware := s.middlewares[i]
		if middleware.Apply(path, method) {
			result = middleware.Middleware(result)
		}
	}

	return result
}

func (s *builder) Build() HttpServer {
	if s.config != nil {
		s.applyConfiguration()
	}

	mux := http.NewServeMux()

	for _, route := range s.routes {
		s.registerRoute(route)
	}

	s.logServerConfig()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.address, s.port),
		Handler:      s.mux,
		ReadTimeout:  time.Duration(s.readTimeout) * time.Second,
		WriteTimeout: time.Duration(s.writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.idleTimeout) * time.Second,
	}

	return &Server{
		server: httpServer,
		mux:    mux,
	}
}

func (s *builder) applyConfiguration() {
	if s.config.GetPort() != 0 {
		s.port = s.config.GetPort()
	}
	if s.config.GetAddress() != "" {
		s.address = s.config.GetAddress()
	}
	if s.config.GetReadTimeout() != 0 {
		s.readTimeout = s.config.GetReadTimeout()
	}
	if s.config.GetWriteTimeout() != 0 {
		s.writeTimeout = s.config.GetWriteTimeout()
	}
	if s.config.GetIdleTimeout() != 0 {
		s.idleTimeout = s.config.GetIdleTimeout()
	}
}

func (s *builder) logServerConfig() {
	if len(s.middlewares) > 0 {
		log.Printf("Registered Middlewares:")
		for _, mw := range s.middlewares {
			scope := "global"
			if len(mw.Path) > 0 || len(mw.Method) > 0 {
				scope = fmt.Sprintf("paths: %v, methods: %v", mw.Path, mw.Method)
			}
			log.Printf("  - %s (%s)", mw.Name, scope)
		}
	}

	log.Printf("Registered Routes:")
	for _, route := range s.routes {
		log.Printf("  %s %s", route.Method, route.Path)
	}
}
