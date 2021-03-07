package server

import (
	"fmt"
	"log"
	"net/http"
	core "scritti/core"
	"scritti/filesystem"
)

// Server initiates a web server on the given port
func Server(port int) {
	mux := http.NewServeMux()

	fs := filesystem.NewOSFileSystem()
	store := core.NewFileStore(fs, "sampledata")
	defer store.Close()

	server := NewComponentServer(store)

	mux.Handle("/wasm/", http.StripPrefix("/wasm/", http.FileServer(http.Dir("./www"))))
	mux.Handle("/js/", http.StripPrefix("", http.FileServer(http.Dir("./www"))))
	mux.HandleFunc("/ws", server.HandleHotReload)
	mux.HandleFunc("/", server.ServeHTTP)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Starting sort on port %d", port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
