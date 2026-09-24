// Package site walks a directory of Typst posts (nested under year
// directories, e.g. blogDir/2025/eduroam-cert.typ), compiles them into a
// separate output directory that mirrors that structure, and builds the
// root and per-keyword index pages linking to them. It also compiles
// arbitrary, non-indexed pages (e.g. an "about" page) via RenderPages.
//
// Every URL this package produces is root-relative (e.g. "/2025/foo.html"),
// on the assumption that outDir is served as the site's root.
package site

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	pagetemplate "render/template"
	"render/typst"
)

// Entry pairs a compiled post's metadata with the URL it's served at.
type Entry struct {
	typst.Metadata
	URL string
}

// Build finds every .typ file under blogDir, compiles each into outDir
// (mirroring blogDir's directory structure), wraps it in the page shell
// loaded from templateDir, and returns one Entry per post.
func Build(blogDir, outDir, templateDir string) ([]Entry, error) {
	paths, err := findTypstFiles(blogDir)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(paths))
	for _, typstPath := range paths {
		rel, err := filepath.Rel(blogDir, typstPath)
		if err != nil {
			return nil, err
		}
		htmlRel := strings.TrimSuffix(rel, filepath.Ext(rel)) + ".html"

		meta, err := compile(typstPath, filepath.Join(outDir, htmlRel), templateDir, "post.html")
		if err != nil {
			return nil, err
		}
		entries = append(entries, Entry{Metadata: meta, URL: "/" + filepath.ToSlash(htmlRel)})
	}
	return entries, nil
}

// RenderPages compiles every .typ file under dir (e.g. an "about" page)
// into outDir (mirroring dir's directory structure) and wraps it in the page
// shell loaded from templateDir. Unlike Build, these pages aren't collected
// into any index. It returns the number of pages compiled.
func RenderPages(dir, outDir, templateDir string) (int, error) {
	paths, err := findTypstFiles(dir)
	if err != nil {
		return 0, err
	}
	for _, typstPath := range paths {
		rel, err := filepath.Rel(dir, typstPath)
		if err != nil {
			return 0, err
		}
		htmlRel := strings.TrimSuffix(rel, filepath.Ext(rel)) + ".html"

		if _, err := compile(typstPath, filepath.Join(outDir, htmlRel), templateDir, "page.html"); err != nil {
			return 0, err
		}
	}
	return len(paths), nil
}

// compile compiles typstPath to htmlPath, creating htmlPath's directory as
// needed, and wraps the result in the page shell loaded from templateName
// (e.g. "post.html" or "page.html") in templateDir.
func compile(typstPath, htmlPath, templateDir, templateName string) (typst.Metadata, error) {
	if err := os.MkdirAll(filepath.Dir(htmlPath), 0o755); err != nil {
		return typst.Metadata{}, err
	}
	if err := typst.CompileHTML(typstPath, htmlPath); err != nil {
		return typst.Metadata{}, err
	}
	meta, err := typst.GetMetadata(typstPath, htmlPath)
	if err != nil {
		return typst.Metadata{}, err
	}
	if err := typst.ApplyLayout(htmlPath, meta, templateDir, templateName); err != nil {
		return typst.Metadata{}, err
	}
	return meta, nil
}

// findTypstFiles returns the path of every .typ file under dir.
func findTypstFiles(dir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".typ" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

// Keywords returns the distinct set of keywords used across entries, sorted
// alphabetically.
func Keywords(entries []Entry) []string {
	seen := make(map[string]bool)
	var keywords []string
	for _, e := range entries {
		for _, k := range e.Keywords {
			if !seen[k] {
				seen[k] = true
				keywords = append(keywords, k)
			}
		}
	}
	sort.Strings(keywords)
	return keywords
}

// recentPostsLimit caps how many posts the root index lists.
const recentPostsLimit = 5

// WriteIndexes renders the root index — the recentPostsLimit most recent
// posts — and one index per keyword listing every matching post, writing
// them under outDir as index.html and keywords/<keyword>.html.
func WriteIndexes(outDir string, entries []Entry, templateDir string) error {
	root, err := renderIndex("index.html", "Home", entries, recentPostsLimit, templateDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "index.html"), root, 0o644); err != nil {
		return err
	}

	keywords := Keywords(entries)
	if len(keywords) == 0 {
		return nil
	}

	keywordDir := filepath.Join(outDir, "keywords")
	if err := os.MkdirAll(keywordDir, 0o755); err != nil {
		return err
	}
	for _, keyword := range keywords {
		matched := withKeyword(entries, keyword)
		page, err := renderIndex("keyword.html", keyword, matched, 0, templateDir)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(keywordDir, keyword+".html"), page, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// withKeyword returns the entries tagged with keyword.
func withKeyword(entries []Entry, keyword string) []Entry {
	var matched []Entry
	for _, e := range entries {
		for _, k := range e.Keywords {
			if k == keyword {
				matched = append(matched, e)
				break
			}
		}
	}
	return matched
}

// renderIndex sorts entries most-recent-first, keeps at most limit of them
// (or all of them if limit is 0), and renders them via templateName (e.g.
// "index.html" or "keyword.html") in templateDir.
func renderIndex(templateName, title string, entries []Entry, limit int, templateDir string) ([]byte, error) {
	sorted := append([]Entry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})
	if limit > 0 && len(sorted) > limit {
		sorted = sorted[:limit]
	}

	items := make([]template.HTML, len(sorted))
	for i, e := range sorted {
		item := pagetemplate.IndexItem{Title: e.Title, URL: e.URL}
		if !e.Date.IsZero() {
			item.Date = e.Date.Format("2006-01-02")
		}
		html, err := pagetemplate.RenderEntry(item)
		if err != nil {
			return nil, err
		}
		items[i] = html
	}

	return pagetemplate.RenderIndex(templateDir, templateName, pagetemplate.IndexData{
		Title: title,
		Items: items,
	})
}
