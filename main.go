package main

import (
	"fmt"
	"goserve/configuration"
	"goserve/server"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")

	configBuilder := configuration.New()
	configuration, err := configBuilder.LoadConfig()

	builder := server.New().WithConfiguration(configuration)

	builder.
		GET("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		}).
		POST("/echo", func(w http.ResponseWriter, r *http.Request) {
			body := make([]byte, r.ContentLength)
			r.Body.Read(body)

			result := fmt.Sprintf("Echo: %s", string(body))

			w.Write([]byte(result))
		}).
		AddGlobalMiddleware("RequestLogging", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
			})
		})

	server := builder.Build()

	err = server.Start()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
