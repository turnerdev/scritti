package server

import (
	"html/template"
	"log"
	"net/http"
	core "scritti/core"
	"scritti/store"
)

// TestData structure
type TestData struct {
}

// ComponentServer TODO
type ComponentServer struct {
	store store.IComponentStore
}

// ServeHTTP test
func (p *ComponentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component, err := p.store.Get("main")
	if err != nil {
		log.Fatal(err)
		return
	}
	html := core.RenderComponent(component)
	data := TestData{}
	tmpl, err := template.New("Test").Parse(html)
	if err != nil {
		log.Fatal("Error parsing template", err)
		return
	}
	tmpl.Execute(w, data)
}
