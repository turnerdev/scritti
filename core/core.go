package core

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type Asset interface{}

// Style asset
type Style struct {
	name    string
	classes []string
}

// Component asset
type Component struct {
	name     string
	parent   *Component
	children []*Component
	style    *Style
}

// GetName of component
func (component *Component) GetName() string {
	return component.name
}

// Append a child to a component
func (component *Component) Append(name string) *Component {
	child := Component{name: name, parent: component}
	component.children = append(component.children, &child)
	return &child
}

// parseLine accepts a string and returns the count of leading
// white space characters and the trimmed input
func parseLine(line string) (int, string) {
	trimmed := strings.TrimSpace(line)
	return strings.Index(line, trimmed), trimmed
}

// CompileComponent walks a component tree and loads dependencies
func CompileComponent(component *Component, store AssetStore) {
	var wg sync.WaitGroup
	wg.Add(len(component.children))

	for _, child := range component.children {
		go func(component *Component) {
			defer wg.Done()
			CompileComponent(component, store)
		}(child)
	}

	fmt.Printf("Compiling %q\n", component.name)
	style, err := store.Get(StyleType, component.name)
	fmt.Printf("Style %q\n", style)
	if err != nil {
		fmt.Printf("Error %q", err)
	}
	if style != nil {
		component.style = style.(*Style)
	}

	wg.Wait()
}

// ParseStyle parses style source to a Style Asset
func ParseStyle(body string) *Style {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	fmt.Printf("Parsing styles %q", body)
	return &Style{
		classes: lines,
	}
}

// ParseComponent parses component source code to a component tree
func ParseComponent(body string) *Component {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	var components []*Component
	depth := map[*Component]int{}

	for i := range lines {
		var component *Component
		indent, value := parseLine(lines[i])

		if len(components) > 0 {
			previous := components[len(components)-1]
			for {
				if depth[previous] >= indent {
					previous = previous.parent
				} else {
					break
				}
			}
			component = previous.Append(value)
		} else {
			component = &Component{name: value}
		}

		depth[component] = indent
		components = append(components, component)
	}

	return components[0]
}

// RenderComponent generate HTML from a Component
func RenderComponent(component *Component) string {
	builder := strings.Builder{}
	var classes string

	if component.style != nil {
		classes = strings.Join(component.style.classes, " ")
	}

	builder.WriteString(fmt.Sprintf("<div name=%q class=%q>", component.name, classes))
	for i := range component.children {
		builder.WriteString(RenderComponent(component.children[i]))
	}
	builder.WriteString("</div>")
	return builder.String()
}
