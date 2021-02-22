package core

import (
	"bufio"
	"log"
	"strings"
)

// Asset represents the common interface implemented by all Asset types
type Asset interface{}

type node struct {
	id       int
	parent   *node
	children []*node
}

// Style asset
type Style struct {
	classes []string
}

// Element of a component
type Element struct {
	style    string
	children []Element
}

// Component asset
type Component struct {
	Element
}

// ComponentSourceLine represents a single line of a Component source
type ComponentSourceLine struct {
	indent int
	style  string
}

func parseElement(line string) ComponentSourceLine {
	trimmed := strings.TrimSpace(line)
	return ComponentSourceLine{
		strings.Index(line, trimmed),
		trimmed,
	}
}

func getDependencyKeys(asset Asset) []AssetKey {
	keys := []AssetKey{}
	switch v := asset.(type) {
	case Component:
		keys = append(keys, AssetKey{StyleType, v.style})
		for _, child := range v.children {
			keys = append(keys, getDependencyKeys(child)...)
		}
	case Element:
		keys = append(keys, AssetKey{StyleType, v.style})
		for _, child := range v.children {
			keys = append(keys, getDependencyKeys(child)...)
		}
	}
	return keys
}

// NewAssetFactory generates a new Asset instance of the provided type
func NewAssetFactory(assetType AssetType, source string) (Asset, error) {
	switch assetType {
	case ComponentType:
		return NewComponent(source)
	case StyleType:
		return NewStyle(source)
	}
	panic("Not implemented")
}

// NewStyle constructs a new Style instance from provided source
func NewStyle(source string) (Style, error) {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(source))
	for sc.Scan() {
		lines = append(lines, strings.TrimSpace(sc.Text()))
	}
	return Style{
		lines,
	}, nil
}

// NewComponent constructs a new Component instance from provided source
func NewComponent(source string) (Component, error) {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(source))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	// Parsed source lines
	sourceLines := []ComponentSourceLine{}

	// Tree structure for tracking source hierarchy, id = sourceLines index
	ns := []*node{}
	n := &node{id: 0}

	log.Println("Source: ", source)

	for _, line := range lines {
		sourceLine := parseElement(line)

		if len(sourceLines) > 0 {
			previous := len(ns) - 1
			for {
				if sourceLines[previous].indent >= sourceLine.indent {
					previous--
				} else {
					break
				}
			}
			n = &node{id: len(ns), children: []*node{}, parent: ns[previous]}
			ns[previous].children = append(ns[previous].children, n)
		}

		sourceLines = append(sourceLines, sourceLine)
		ns = append(ns, n)
	}

	var build func(*node) Element
	build = func(n *node) Element {
		element := Element{
			style:    sourceLines[n.id].style,
			children: []Element{},
		}
		for _, child := range n.children {
			element.children = append(element.children, build(child))
		}
		return element
	}

	return Component{build(ns[0])}, nil
}
