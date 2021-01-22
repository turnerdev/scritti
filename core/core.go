package core

import (
	"bufio"
	"strings"
)

// Component tree
type Component struct {
	name     string
	parent   *Component
	children []*Component
}

// Append a child to a component
func (component *Component) Append(name string) *Component {
	child := Component{name, component, nil}
	component.children = append(component.children, &child)
	return &child
}

// readLine accepts a string and returns the count of leading
// white space characters and the trimmed input
func readLine(line string) (int, string) {
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
		indent, value := readLine(lines[i])

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

// func RenderComponent
