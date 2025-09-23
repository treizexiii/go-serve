package main

import (
	"fmt"
	"goserve/configuration"
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

	hello := server.
		CreateRoute(server.GET, "/hello", func(w http.ResponseWriter, r *http.Request)x {
			w.Write([]byte("Hello, World!"))
		})

	demo := server.CreateRoute(server.GET, "/demo", contextDemoHandler())

	builder := server.New().
		WithConfiguration(configuration).
		WithJSONSerialization().
		WithLogging(true, true).
		AddRoutes([]server.RouteInfo{
			hello,
		}).
		POST("/echo", func(w http.ResponseWriter, r *http.Request) {
			body := make([]byte, r.ContentLength)
			r.Body.Read(body)

			result := fmt.Sprintf("Echo: %s", string(body))

			w.Write([]byte(result))
		})

	server := builder.Build()

	err = server.Start()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func contextDemoHandler(w http.ResponseWriter, r *http.Request) any {
	// Créer des données via le contexte (approche alternative)
	responseData := map[string]interface{}{
		"message":     "Données passées via le contexte",
		"method":      r.Method,
		"path":        r.URL.Path,
		"user_agent":  r.UserAgent(),
		"remote_addr": r.RemoteAddr,
	}

	// Utiliser le helper pour structurer la réponse
	return responseData
}
