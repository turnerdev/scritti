package core

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// RenderComponent renders a Component type Asset to HTML
func RenderComponent(w io.Writer, component Component, fn func(AssetKey) (Asset, error)) error {
	node, err := renderElement(component.Element, fn)
	if err != nil {
		return err
	}
	html.Render(w, node)
	return nil
}

// RenderElement generates HTML for an Element.
func renderElement(element Element, fn func(AssetKey) (Asset, error)) (*html.Node, error) {
	var classes string
	if len(element.style) > 0 {
		style, err := fn(AssetKey{StyleType, element.style})
		if err != nil {
			return nil, err
		}
		classes = strings.Join(style.(Style).classes, " ")
	}

	node := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: classes,
			},
		},
	}

	for i := range element.children {
		childNode, err := renderElement(element.children[i], fn)
		if err != nil {
			return nil, err
		}
		node.AppendChild(childNode)
	}

	return node, nil
}
