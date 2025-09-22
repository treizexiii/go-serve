package server

import (
	"context"
	"fmt"
	"goserve/configuration"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server *http.Server
	mux    *http.ServeMux
	config configuration.Configuration
}

func (s *Server) GetHttpServer() *http.Server {
	return s.server
}

func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

func (s *Server) Start() error {
	if s.server == nil {
		return fmt.Errorf("server not built, call Build() before Start()")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Start listening on %v:%d", s.config.GetAddress(), s.config.GetPort())
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %v: %v", s.server.Addr, err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	return nil
}
