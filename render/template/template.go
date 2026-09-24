// Package template renders documents into a page shell, independent of
// whatever produced the document content (e.g. Typst). Each page type has
// its own template file: post.html for blog posts, page.html for arbitrary
// pages, index.html for the root index, and keyword.html for per-keyword
// indexes. entry.html renders a single post's metadata into a link fragment
// used by index.html and keyword.html.
//
// entry.html is embedded into the binary at build time; the rest are read
// from disk at render time so they can be edited without a rebuild.
package template

import (
	"bytes"
	"embed"
	"html/template"
	"path/filepath"
)

// DefaultDir is the template directory used when dir is empty.
const DefaultDir = "./render/template/"

//go:embed entry.html
var entryFS embed.FS

var entryTmpl = template.Must(template.ParseFS(entryFS, "entry.html"))

// PageData is the data rendered into post.html and page.html.
type PageData struct {
	Title       string
	Description string
	Authors     string
	Keywords    string
	Content     template.HTML
}

// IndexItem is one post's metadata, rendered into a link fragment by
// RenderEntry.
type IndexItem struct {
	Title string
	URL   string
	Date  string
}

// IndexData is the data rendered into index.html and keyword.html.
type IndexData struct {
	Title string
	Items []template.HTML
}

// Render parses templateName (e.g. "post.html" or "page.html") out of dir
// (DefaultDir if dir is empty) and executes it with data, returning the
// resulting page.
func Render(dir, templateName string, data PageData) ([]byte, error) {
	return execute(dir, templateName, data)
}

// RenderIndex parses templateName (e.g. "index.html" or "keyword.html") out
// of dir (DefaultDir if dir is empty) and executes it with data, returning
// the resulting page.
func RenderIndex(dir, templateName string, data IndexData) ([]byte, error) {
	return execute(dir, templateName, data)
}

// RenderEntry executes the embedded entry.html template with item, returning
// the HTML fragment that links to that post.
func RenderEntry(item IndexItem) (template.HTML, error) {
	var buf bytes.Buffer
	if err := entryTmpl.Execute(&buf, item); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func execute(dir, templateName string, data any) ([]byte, error) {
	if dir == "" {
		dir = DefaultDir
	}

	tmpl, err := template.ParseFiles(filepath.Join(dir, templateName))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
