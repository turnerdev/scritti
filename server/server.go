package server

import (
	"fmt"
	"log"
	"net/http"
	"scritti/store"
)

// Server initiates a web server on the given port
func Server(port int) {
	mux := http.NewServeMux()

	server := ComponentServer{
		store.NewFileStore(),
	}

	mux.HandleFunc("/", server.ServeHTTP)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
