package server

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	core "scritti/core"
	"strings"
)

// TestData structure
type TestData struct {
}

// ComponentServer TODO
type ComponentServer struct {
	reload chan bool
	store  core.AssetStore
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
	// WebSocket

	// Compile components
	asset, err := p.store.Get(core.ComponentType, "main")

	go p.store.Watch(core.ComponentType, "main", p.reload)

	if err != nil {
		log.Fatal(err)
		return
	}
	component := asset.(*core.Component)
	core.CompileComponent(component, p.store)

	// Render output
	base := getTemplate()
	html := core.RenderComponent(component, 1)
	data := map[string]interface{}{
		"Body": template.HTML(html),
	}
	tmpl := template.Must(template.New("main").Parse(base))
	tmpl.Execute(w, data)
}
