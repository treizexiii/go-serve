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
		CreateRoute(server.GET, "/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})

	builder := server.New().
		WithConfiguration(configuration).
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
