package core

import (
	"bufio"
	"fmt"
	"strings"
)

// TreeNode TODO
type TreeNode interface {
	Append() *TreeNode
}

// Component tree
type Component struct {
	name     string
	parent   *Component
	children []*Component
}

// GetName of component
func (component *Component) GetName() string {
	return component.name
}

// Append a child to a component
func (component *Component) Append(name string) *Component {
	child := Component{name, component, nil}
	component.children = append(component.children, &child)
	return &child
}

// parseLine accepts a string and returns the count of leading
// white space characters and the trimmed input
func parseLine(line string) (int, string) {
	trimmed := strings.TrimSpace(line)
	return strings.Index(line, trimmed), trimmed
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
	builder.WriteString(fmt.Sprintf("<div name=%q>", component.name))
	for i := range component.children {
		builder.WriteString(RenderComponent(component.children[i]))
	}
	builder.WriteString("</div>")
	return builder.String()
}
