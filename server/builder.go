package server

import (
	"fmt"
	"goserve/configuration"
	"goserve/middlewares"
	"goserve/routes"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type builder struct {
	mux          *http.ServeMux
	port         int
	address      string
	routes       []routes.RouteInfo
	middlewares  []middlewares.MiddlewareInfo
	config       configuration.Configuration
	readTimeout  int
	writeTimeout int
	idleTimeout  int
}

// Constructor

func New() ServerBuilder {
	return &builder{
		mux:          http.NewServeMux(),
		routes:       make([]routes.RouteInfo, 0),
		address:      "",
		port:         8080, // Default port
		middlewares:  make([]middlewares.MiddlewareInfo, 0),
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

func (s *builder) AddRoute(method routes.Http_Method, path string, handler routes.ActionFunc) ServerBuilder {
	s.routes = append(s.routes, routes.CreateRoute(method, path, handler))

	return s
}

func (s *builder) AddRoutes(routes []routes.RouteInfo) ServerBuilder {
	s.routes = append(s.routes, routes...)
	return s
}

func (s *builder) GET(path string, handler routes.ActionFunc) ServerBuilder {
	return s.addRoute(routes.CreateGET(path, handler))
}

func (s *builder) POST(path string, handler routes.ActionFunc) ServerBuilder {
	return s.addRoute(routes.CreatePOST(path, handler))
}

func (s *builder) PUT(path string, handler routes.ActionFunc) ServerBuilder {
	return s.addRoute(routes.CreatePUT(path, handler))
}

func (s *builder) DELETE(path string, handler routes.ActionFunc) ServerBuilder {
	return s.addRoute(routes.CreateDELETE(path, handler))
}

func (s *builder) addRoute(route routes.RouteInfo) ServerBuilder {
	s.routes = append(s.routes, route)
	return s
}

func (s *builder) AddGlobalMiddleware(name string, middleware middlewares.MiddlewareFunc) ServerBuilder {
	s.middlewares = append(s.middlewares, middlewares.MiddlewareInfo{
		Name:       name,
		Middleware: middleware,
	})
	return s
}

func (s *builder) WithJSONSerialization() ServerBuilder {
	return s.AddGlobalMiddleware("JSONSerialization", JSONSerializationMiddleware())
}

func (s *builder) WithLogging(logRequests, logResponses bool) ServerBuilder {
	if logRequests || logResponses {
		requestId := uuid.New().String()
		s.AddGlobalMiddleware("Logging", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if logRequests {
					log.Printf("Request: %s %s %s", requestId, r.Method, r.URL.Path)
				}

				if logResponses {
					lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
					next.ServeHTTP(lrw, r)
					log.Printf("Response: %s %d", requestId, lrw.statusCode)
				} else {
					next.ServeHTTP(w, r)
				}
			})
		})
	}
	return s
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
		config: s.config,
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
			log.Printf("  - %s (%s)", mw.Name, scope)
		}
	}

	log.Printf("Registered Routes:")
	for _, route := range s.routes {
		log.Printf("  %s %s", route.GetMethod(), route.GetPath())
	}
}

func (s *builder) registerRoute(route routes.RouteInfo) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != string(route.GetMethod()) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		action := route.GetHandler()
		if action == nil {
			http.Error(w, "Not Implemented", http.StatusNotImplemented)
			return
		}
		result, err := (*action)(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", result.ContentType())
		w.WriteHeader(result.Code())
		if !result.IsDataEmpty() {
			json, err := result.JsonString()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to serialize response"))
				return
			}

			w.Write(json)
		}
	})

	finalHandler := s.applyMiddlewares(handler, route)

	s.mux.Handle(route.GetPath(), finalHandler)
}

func (s *builder) applyMiddlewares(handler http.Handler, route routes.RouteInfo) http.Handler {
	result := handler

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		middleware := s.middlewares[i]
		result = middleware.Middleware(result)
	}

	middlewares := route.GetMiddlewares()
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		result = middleware(result)
	}

	return result
}
