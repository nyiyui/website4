package typst

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// TocEntry is a single header and its children in a table of contents.
type TocEntry struct {
	Level    int
	ID       string
	Node     *html.Node
	Children []TocEntry
}

var headingLevels = map[string]int{
	"h1": 1,
	"h2": 2,
	"h3": 3,
	"h4": 4,
	"h5": 5,
	"h6": 6,
}

// GetTOC extracts a table of contents from Typst's compiled HTML output by
// walking the document for heading elements and nesting them by level.
func GetTOC(htmlPath string) ([]TocEntry, error) {
	f, err := os.Open(htmlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	doc, err := html.Parse(f)
	if err != nil {
		return nil, err
	}

	var flat []TocEntry
	collectHeadings(doc, &flat)
	return nestHeadings(flat), nil
}

// collectHeadings traverses the HTML node tree and appends a flat, in-order
// list of heading entries found in the document.
func collectHeadings(n *html.Node, headings *[]TocEntry) {
	if n.Type == html.ElementNode {
		if level, ok := headingLevels[n.Data]; ok {
			*headings = append(*headings, TocEntry{
				Level: level,
				ID:    attr(n, "id"),
				Node:  n,
			})
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectHeadings(c, headings)
	}
}

// nestHeadings turns a flat, in-order list of headings into a tree based on
// heading level.
func nestHeadings(flat []TocEntry) []TocEntry {
	var root []TocEntry
	var stack []*TocEntry

	for _, entry := range flat {
		for len(stack) > 0 && stack[len(stack)-1].Level >= entry.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			root = append(root, entry)
			stack = append(stack, &root[len(root)-1])
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, entry)
			stack = append(stack, &parent.Children[len(parent.Children)-1])
		}
	}

	return root
}

// attr returns the value of the named attribute on n, or "" if not present.
func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// setAttr sets key to val on n, overwriting any existing attribute with that
// key.
func setAttr(n *html.Node, key, val string) {
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: val})
}

// AssignHeadingIDs walks n for heading elements that have no id attribute
// (Typst only emits one for headings with a referenced label) and gives them
// one derived from their text content, deduplicating repeats by appending
// -2, -3, and so on.
func AssignHeadingIDs(n *html.Node) {
	assignHeadingIDs(n, make(map[string]int))
}

func assignHeadingIDs(n *html.Node, seen map[string]int) {
	if n.Type == html.ElementNode {
		if _, ok := headingLevels[n.Data]; ok && attr(n, "id") == "" {
			setAttr(n, "id", uniqueSlug(textContent(n), seen))
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		assignHeadingIDs(c, seen)
	}
}

// uniqueSlug slugifies s and, if that collides with a slug already in seen,
// appends -2, -3, etc. until it doesn't.
func uniqueSlug(s string, seen map[string]int) string {
	base := slugify(s)
	if base == "" {
		base = "section"
	}
	seen[base]++
	if n := seen[base]; n > 1 {
		return fmt.Sprintf("%s-%d", base, n)
	}
	return base
}

// slugify lowercases s and collapses runs of non-alphanumeric characters
// into a single hyphen, trimming any leading/trailing hyphen.
func slugify(s string) string {
	var b strings.Builder
	lastHyphen := true
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
			lastHyphen = false
		} else if !lastHyphen {
			b.WriteRune('-')
			lastHyphen = true
		}
	}
	return strings.TrimRight(b.String(), "-")
}

// textContent returns the concatenated text of n and its descendants.
func textContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(textContent(c))
	}
	return b.String()
}

// RenderTOC renders entries as a nested <ul> of links to each heading's id,
// preserving each heading's inline markup as the link text. It returns ""
// if entries is empty.
func RenderTOC(entries []TocEntry) (template.HTML, error) {
	if len(entries) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, tocList(entries)); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// tocList builds a <ul> element with a <li> per entry, each linking to the
// entry's heading (with its inline markup cloned in as the link text) and
// nesting its children as a further <ul>.
func tocList(entries []TocEntry) *html.Node {
	ul := &html.Node{Type: html.ElementNode, Data: "ul"}
	for _, e := range entries {
		a := &html.Node{
			Type: html.ElementNode,
			Data: "a",
			Attr: []html.Attribute{{Key: "href", Val: "#" + e.ID}},
		}
		for c := e.Node.FirstChild; c != nil; c = c.NextSibling {
			a.AppendChild(cloneNode(c))
		}

		li := &html.Node{Type: html.ElementNode, Data: "li"}
		li.AppendChild(a)
		if len(e.Children) > 0 {
			li.AppendChild(tocList(e.Children))
		}
		ul.AppendChild(li)
	}
	return ul
}

// cloneNode returns a parentless deep copy of n, suitable for attaching
// elsewhere in the tree.
func cloneNode(n *html.Node) *html.Node {
	clone := &html.Node{
		Type: n.Type,
		Data: n.Data,
		Attr: append([]html.Attribute(nil), n.Attr...),
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		clone.AppendChild(cloneNode(c))
	}
	return clone
}
