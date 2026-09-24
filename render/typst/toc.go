package typst

import (
	"os"

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
