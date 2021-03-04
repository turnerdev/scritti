package core

import (
	"io"
	"log"
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
			log.Println(err)
		} else {
			classes = strings.Join(style.(Style).classes, " ")
		}
	}

	tag := "div"
	if len(element.tag) > 0 {
		tag = element.tag
	}

	node := &html.Node{
		Type: html.ElementNode,
		Data: tag,
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: classes,
			},
		},
	}

	if len(element.text) > 0 {
		textNode := &html.Node{
			Type: html.TextNode,
			Data: element.text,
		}
		node.AppendChild(textNode)
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
