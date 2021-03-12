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
	var svgsource string

	// If element specified style, fetch classes from Style asset
	if len(element.style) > 0 {
		style, err := fn(AssetKey{StyleType, element.style})
		if err != nil {
			log.Println(err)
		} else {
			classes = strings.Join(style.(Style).classes, " ")
		}
	}

	// If element tag is 'svg', fetch source from SVG asset
	if element.tag == "svg" && len(element.style) > 0 {
		svg, err := fn(AssetKey{SVGType, element.style})
		if err != nil {
			log.Println(err)
		} else {
			svgsource = svg.(SVG).Source
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

	if len(svgsource) > 0 {
		reader := strings.NewReader(svgsource)
		nodes, err := html.Parse(reader)
		if err != nil {
			log.Println(err)
		}
		body := nodes.FirstChild.FirstChild.NextSibling
		node = body.FirstChild
		body.RemoveChild(node)
		node.Attr = append(node.Attr, html.Attribute{
			Key: "class",
			Val: classes,
		})
	}

	if len(element.text) > 0 {
		textNode := &html.Node{
			Type: html.TextNode,
			Data: strings.ReplaceAll(element.text, "\\n", "\n"),
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
