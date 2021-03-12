package server

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	core "scritti/core"
	"scritti/filesystem"
	"strings"
)

type ComponentServer struct {
	store core.AssetStore
}

// NewComponentServer initializes a new server with a specified Asset Store
func NewComponentServer(store core.AssetStore) *ComponentServer {
	return &ComponentServer{
		store,
	}
}

func getTemplate() string {
	file, err := os.Open("www/index.html")
	if err != nil {
		panic("Unable to open template")
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

// ServeHTTP test
func (p ComponentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := getTemplate()
	data := map[string]interface{}{}
	tmpl := template.Must(template.New("main").Parse(base))
	tmpl.Execute(w, data)
}

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
