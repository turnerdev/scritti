package server

import (
	"bufio"
	"html/template"
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
	// asset, err := p.store.Get(core.AssetKey{core.ComponentType, "main"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Render output
	// buffer := new(bytes.Buffer)
	base := getTemplate()
	// err = core.RenderComponent(buffer, asset.(core.Component), p.store.Get)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	data := map[string]interface{}{
		// "Body": template.HTML(buffer.String()),
	}
	tmpl := template.Must(template.New("main").Parse(base))
	tmpl.Execute(w, data)
}
