package core

import (
	"bufio"
	"log"
	"regexp"
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
	Source  string
	classes []string
}

// SVG asset
type SVG struct {
	Source string
}

// Element of a component
type Element struct {
	text     string
	tag      string
	style    string
	children []Element
}

// Component asset
type Component struct {
	Source string
	Element
}

// ComponentSourceLine represents a single line of a Component source
type ComponentSourceLine struct {
	indent int
	tag    string
	text   string
	style  string
}

func parseElement(line string) ComponentSourceLine {
	pattern := `^(?P<indent>\s*)((?P<tag>\w+)\.)?(?P<style>\w+)?(\s*"(?P<text>[^"\\]*(\\.[^"\\]*)*)")?\s*$`
	pathMetadata := regexp.MustCompile(pattern)

	matches := pathMetadata.FindStringSubmatch(line)
	names := pathMetadata.SubexpNames()
	groups := make(map[string]string)

	for i, match := range matches {
		if len(names[i]) > 0 {
			groups[names[i]] = match
		}
	}
	return ComponentSourceLine{
		len(groups["indent"]),
		groups["tag"],
		strings.ReplaceAll(groups["text"], `\"`, `"`),
		groups["style"],
	}
}

func getDependencyKeys(asset Asset) []AssetKey {
	distinct := make(map[AssetKey]bool)
	keys := []AssetKey{}
	switch v := asset.(type) {
	case Component:
		distinct[AssetKey{StyleType, v.style}] = true
		// keys = append(keys, AssetKey{StyleType, v.style})
		for _, child := range v.children {
			for _, childKey := range getDependencyKeys(child) {
				distinct[childKey] = true
			}
			// keys = append(keys, getDependencyKeys(child)...)
		}
	case Element:
		distinct[AssetKey{StyleType, v.style}] = true
		// keys = append(keys, AssetKey{StyleType, v.style})
		for _, child := range v.children {
			for _, childKey := range getDependencyKeys(child) {
				distinct[childKey] = true
			}
			// keys = append(keys, getDependencyKeys(child)...)
		}
	}
	for k := range distinct {
		keys = append(keys, k)
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
	case SVGType:
		return NewSVG(source)
	}
	panic("Not implemented")
}

// NewSVG construts a new SVG instance from provided source
func NewSVG(source string) (SVG, error) {
	return SVG{source}, nil
}

// NewStyle constructs a new Style instance from provided source
func NewStyle(source string) (Style, error) {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(source))
	for sc.Scan() {
		lines = append(lines, strings.TrimSpace(sc.Text()))
	}
	return Style{
		Source:  source,
		classes: lines,
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

	log.Println("Source: ", sourceLines)

	for _, line := range lines {
		sourceLine := parseElement(line)

		if len(sourceLines) > 0 {
			previous := len(ns) - 1
			for {
				if sourceLines[previous].indent >= sourceLine.indent && previous > 0 {
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
		source := sourceLines[n.id]
		element := Element{
			text:     source.text,
			tag:      source.tag,
			style:    source.style,
			children: []Element{},
		}

		for _, child := range n.children {
			element.children = append(element.children, build(child))
		}
		return element
	}

	if len(ns) == 0 {
		return Component{}, nil
	}

	return Component{
		Source:  source,
		Element: build(ns[0]),
	}, nil
}
