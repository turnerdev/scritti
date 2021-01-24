package server

import (
	"html/template"
	"log"
	"net/http"
	core "scritti/core"
)

// TestData structure
type TestData struct {
}

// ComponentServer TODO
type ComponentServer struct {
	store core.AssetStore
}

// ServeHTTP test
func (p ComponentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	asset, err := p.store.Get(core.ComponentType, "main")
	if err != nil {
		log.Fatal(err)
		return
	}

	component := asset.(*core.Component)
	core.CompileComponent(component, p.store)

	html := `<!doctype html><html><head><link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet"></head><body>`
	html = html + core.RenderComponent(component)
	html = html + `</body></html>`

	data := TestData{}
	tmpl, err := template.New("Test").Parse(html)
	if err != nil {
		log.Fatal("Error parsing template", err)
		return
	}
	tmpl.Execute(w, data)
}
