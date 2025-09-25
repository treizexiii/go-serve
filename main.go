package main

import (
	"fmt"
	"goserve/configuration"
	"goserve/responses"
	"goserve/routes"
	"goserve/server"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")

	configBuilder := configuration.New()
	configuration, err := configBuilder.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	hello := routes.
		CreateRoute(routes.GET, "/hello", func(r *http.Request) (responses.ApiResponse, error) {
			return responses.Ok(map[string]string{
				"message": "Hello, World!",
			}), nil
		})

	demo := routes.CreateRoute(routes.GET, "/demo", contextDemoHandler)

	builder := server.New().
		WithConfiguration(configuration).
		// WithJSONSerialization().
		WithLogging(true, true).
		AddRoutes([]routes.RouteInfo{
			hello,
			demo,
		}).
		POST("/echo", func(r *http.Request) (responses.ApiResponse, error) {
			body := make([]byte, r.ContentLength)
			r.Body.Read(body)
			return responses.Ok(body), nil
		})

	server := builder.Build()

	err = server.Start()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func contextDemoHandler(r *http.Request) (responses.ApiResponse, error) {
	// Créer des données via le contexte (approche alternative)
	responseData := map[string]interface{}{
		"message":     "Données passées via le contexte",
		"method":      r.Method,
		"path":        r.URL.Path,
		"user_agent":  r.UserAgent(),
		"remote_addr": r.RemoteAddr,
	}

	// Utiliser le helper pour structurer la réponse
	return responses.Ok(responseData), nil
}
