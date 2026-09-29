package typst

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	nethtml "golang.org/x/net/html"
)

// Body parses the Typst-compiled file at htmlPath and returns the children
// of its <body> element re-serialized back to HTML, along with a table of
// contents built from its headings.
func Body(htmlPath string) (content template.HTML, toc template.HTML, err error) {
	f, err := os.Open(htmlPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	doc, err := nethtml.Parse(f)
	if err != nil {
		return "", "", err
	}

	body := findNode(doc, "body")
	if body == nil {
		return "", "", fmt.Errorf("typst: no <body> element found in %s", htmlPath)
	}

	AssignHeadingIDs(body)

	var flat []TocEntry
	collectHeadings(body, &flat)
	toc, err = RenderTOC(nestHeadings(flat))
	if err != nil {
		return "", "", err
	}

	wrapHeadingSections(body)

	var buf bytes.Buffer
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if err := nethtml.Render(&buf, c); err != nil {
			return "", "", err
		}
	}
	return template.HTML(buf.String()), toc, nil
}

// wrapHeadingSections rewrites parent's children in place, moving each h2
// element and the siblings that follow it, up to (but not including) the
// next h2, into a new <div class="sticky"> wrapper. Siblings before the
// first h2 are left untouched.
func wrapHeadingSections(parent *nethtml.Node) {
	for c := parent.FirstChild; c != nil; {
		if c.Type != nethtml.ElementNode || c.Data != "h2" {
			c = c.NextSibling
			continue
		}

		section := &nethtml.Node{
			Type: nethtml.ElementNode,
			Data: "div",
			Attr: []nethtml.Attribute{{Key: "class", Val: "sticky"}},
		}
		parent.InsertBefore(section, c)

		next := c.NextSibling
		parent.RemoveChild(c)
		section.AppendChild(c)
		for next != nil && !(next.Type == nethtml.ElementNode && next.Data == "h2") {
			sibling := next.NextSibling
			parent.RemoveChild(next)
			section.AppendChild(next)
			next = sibling
		}

		c = next
	}
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
