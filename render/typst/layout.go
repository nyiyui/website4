package typst

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	nethtml "golang.org/x/net/html"

	pagetemplate "render/template"
)

// ApplyLayout wraps the Typst-compiled HTML at htmlPath in the page shell
// loaded from templateName (e.g. "post.html" or "page.html") in templateDir
// (pagetemplate.DefaultDir if templateDir is empty), rewriting htmlPath in
// place. meta supplies the <head> fields; the original document's <body>
// content is carried over unchanged as the template's Content.
func ApplyLayout(htmlPath string, meta Metadata, templateDir, templateName string) error {
	content, err := serializeBody(htmlPath)
	if err != nil {
		return err
	}

	out, err := pagetemplate.Render(templateDir, templateName, pagetemplate.PageData{
		Title:       meta.Title,
		Description: meta.Description,
		Authors:     strings.Join(meta.Authors, ", "),
		Keywords:    strings.Join(meta.Keywords, ", "),
		Content:     template.HTML(content),
	})
	if err != nil {
		return err
	}

	return os.WriteFile(htmlPath, out, 0o644)
}

// serializeBody parses the Typst-compiled file at htmlPath and re-serializes
// the children of its <body> element back to an HTML string.
func serializeBody(htmlPath string) (string, error) {
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
	return buf.String(), nil
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
