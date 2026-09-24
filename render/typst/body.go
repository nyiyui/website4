package typst

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	nethtml "golang.org/x/net/html"
)

// Body parses the Typst-compiled file at htmlPath and re-serializes the
// children of its <body> element back to HTML.
func Body(htmlPath string) (template.HTML, error) {
	f, err := os.Open(htmlPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	doc, err := nethtml.Parse(f)
	if err != nil {
		return "", err
	}

	body := findNode(doc, "body")
	if body == nil {
		return "", fmt.Errorf("typst: no <body> element found in %s", htmlPath)
	}

	var buf bytes.Buffer
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if err := nethtml.Render(&buf, c); err != nil {
			return "", err
		}
	}
	return template.HTML(buf.String()), nil
}

// findNode returns the first descendant of n (or n itself) with the given
// tag name, or nil if none is found.
func findNode(n *nethtml.Node, tag string) *nethtml.Node {
	if n.Type == nethtml.ElementNode && n.Data == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findNode(c, tag); found != nil {
			return found
		}
	}
	return nil
}
