// Package template renders data (e.g., list of documents) and documents to
// HTML.
//
// base.html is the shared page shell (doctype, head, nav); each page type
// (post.html, page.html, index.html, keyword.html) fills in base.html's
// "title", "head", and "content" blocks.
//
// Every template (including entries.html) is read from dir at render time,
// not embedded, so a whole templateDir (e.g. per-locale, for i18n) can be
// swapped in without a rebuild.
package template

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	"render/typst"
)

// PageData is the data rendered into post.html and page.html.
type PageData struct {
	Title       string
	Description string
	Authors     string
	Keywords    string
	Date        time.Time
	Content     template.HTML
	TOC         template.HTML
}

// Metadata pairs a post's Typst metadata with the URL it's served at.
type Metadata struct {
	typst.Metadata
	URL string
}

// IndexData is the data rendered into index.html and keyword.html.
type IndexData struct {
	Title   string
	Entries template.HTML
}

// YearGroup is one year's worth of posts, pre-rendered into a <ul>.
type YearGroup struct {
	Year    string
	Entries template.HTML
}

// AllPostsData is the data rendered into all-posts.html.
type AllPostsData struct {
	Title string
	Years []YearGroup
}

// Render well renders a single page :D
func Render(dir, templateName string, data PageData) ([]byte, error) {
	return execute(dir, templateName, data)
}

// RenderIndex is like Render is exclusively for the index page.
func RenderIndex(dir, templateName string, data IndexData) ([]byte, error) {
	return execute(dir, templateName, data)
}

// RenderAllPosts is like Render, exclusively for the all-posts page.
func RenderAllPosts(dir, templateName string, data AllPostsData) ([]byte, error) {
	return execute(dir, templateName, data)
}

// RenderEntries returns a <ul> listing every post in metas, most-recent order
// preserved as given.
func RenderEntries(dir string, metas []Metadata) (template.HTML, error) {
	tmpl, err := template.ParseFiles(filepath.Join(dir, "entries.html"))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, metas); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// execute renders templateName within the shared base.html layout: base.html
// provides the page shell, and templateName fills in the "title", "head",
// and "content" blocks it declares.
func execute(dir, templateName string, data any) ([]byte, error) {
	tmpl, err := template.ParseFiles(filepath.Join(dir, "base.html"), filepath.Join(dir, templateName))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
