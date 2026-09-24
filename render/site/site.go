// Package site walks a directory of blog posts (possibly nested) and a
// directory of non-blog post pages ("pages") (e.g., /about), collects metadata
// and builds the output into a static website.
//
// We assume that the website is served from root (e.g., https://example.org/abc
// serves <output>/abc).
package site

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"

	pagetemplate "render/template"
	"render/typst"
)

// Entry pairs a compiled post's metadata with the URL it's served at.
type Entry struct {
	typst.Metadata
	URL string
}

// RenderBlog renders every Typst file under blogDir, and puts the result in
// outDir. It also returns a list of entries, with metadata about the post.
//
// It wraps each Typst document in templateDir/page.html, which should provide
// the necessary layout and navigation.
func RenderBlog(blogDir, outDir, templateDir string) ([]Entry, error) {
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

// RenderPages renders every Typst file under dir and returns the number of pages made.
//
// It wraps each Typst document in templateDir/page.html, which should provide
// the necessary layout and navigation.
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

// CopyStatic copies every file under dir to the same relative path under
// outDir, for assets that don't go through Typst (a favicon, downloadable
// files, images you don't want base64-inlined by Typst's HTML export, etc.).
func CopyStatic(dir, outDir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(outDir, rel)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, content, 0o644)
	})
}

// compile renders a Typst document (at typstPath) and wraps the document in
// a template. The result is stored at htmlPath (directories are created as
// necessary).
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
	content, err := typst.Body(htmlPath)
	if err != nil {
		return typst.Metadata{}, err
	}

	out, err := pagetemplate.Render(templateDir, templateName, pagetemplate.PageData{
		Title:       meta.Title,
		Description: meta.Description,
		Authors:     strings.Join(meta.Authors, ", "),
		Keywords:    strings.Join(meta.Keywords, ", "),
		Date:        meta.Date,
		Content:     content,
	})
	if err != nil {
		return typst.Metadata{}, err
	}
	if err := os.WriteFile(htmlPath, out, 0o644); err != nil {
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
//
// Keywords are compared under Unicode NFKC normalization.
func Keywords(entries []Entry) []string {
	seen := make(map[string]bool)
	var keywords []string
	for _, e := range entries {
		for _, k := range normalizeKeywords(e.Keywords) {
			if !seen[k] {
				seen[k] = true
				keywords = append(keywords, k)
			}
		}
	}
	sort.Strings(keywords)
	return keywords
}

// normalizeKeywords returns keywords under Unicode NFKC normalization, so
// e.g. combining-character, precomposed, and compatibility spellings of the
// same keyword compare equal.
func normalizeKeywords(keywords []string) []string {
	normalized := make([]string, len(keywords))
	for i, k := range keywords {
		normalized[i] = norm.NFKC.String(k)
	}
	return normalized
}

// recentPostsLimit caps how many posts the root index lists.
const recentPostsLimit = 5

// WriteIndexes renders the (root) index and index pages for each keyword.
// The index contains the recentPostsLimit-most recent posts.
//
// The index page is at index.html, and each keyword at keywords/<keyword>.html.
func WriteIndexes(outDir string, entries []Entry, templateDir string) error {
	root, err := renderIndex("index.html", "", entries, recentPostsLimit, templateDir)
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

// withKeyword returns the entries tagged with keyword. keyword is assumed to
// already be NFKC-normalized, as Keywords returns it.
func withKeyword(entries []Entry, keyword string) []Entry {
	var matched []Entry
	for _, e := range entries {
		for _, k := range e.Keywords {
			if norm.NFKC.String(k) == keyword {
				matched = append(matched, e)
				break
			}
		}
	}
	return matched
}

// WriteAllPosts renders a page listing every post, grouped by year (most
// recent year first, most recent post first within each year), to
// outDir/all-posts.html.
func WriteAllPosts(outDir string, entries []Entry, templateDir string) error {
	yearGroups := groupByYear(entries)

	years := make([]pagetemplate.YearGroup, len(yearGroups))
	for i, g := range yearGroups {
		list, err := pagetemplate.RenderEntries(templateDir, recentPostsFor(g.entries, 0))
		if err != nil {
			return err
		}
		years[i] = pagetemplate.YearGroup{Year: strconv.Itoa(g.year), Entries: list}
	}

	page, err := pagetemplate.RenderAllPosts(templateDir, "all-posts.html", pagetemplate.AllPostsData{
		Title: "All Posts",
		Years: years,
	})
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "all-posts.html"), page, 0o644)
}

// yearGroup is every entry published in year, used by groupByYear.
type yearGroup struct {
	year    int
	entries []Entry
}

// groupByYear buckets entries by their Date's year, returning groups sorted
// most-recent year first.
func groupByYear(entries []Entry) []yearGroup {
	byYear := make(map[int][]Entry)
	for _, e := range entries {
		year := e.Date.Year()
		byYear[year] = append(byYear[year], e)
	}

	years := make([]int, 0, len(byYear))
	for y := range byYear {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	groups := make([]yearGroup, len(years))
	for i, y := range years {
		groups[i] = yearGroup{year: y, entries: byYear[y]}
	}
	return groups
}

// renderIndex renders the recentPostsFor(entries, limit) list under
// templateName (e.g., "index.html" or "keyword.html") in templateDir.
func renderIndex(templateName, title string, entries []Entry, limit int, templateDir string) ([]byte, error) {
	list, err := pagetemplate.RenderEntries(templateDir, recentPostsFor(entries, limit))
	if err != nil {
		return nil, err
	}

	return pagetemplate.RenderIndex(templateDir, templateName, pagetemplate.IndexData{
		Title:   title,
		Entries: list,
	})
}

// recentPostsFor sorts entries most-recent-first, keeps at most limit of
// them (or all of them if limit is 0), and converts them into the
// []pagetemplate.Metadata RenderEntries expects.
func recentPostsFor(entries []Entry, limit int) []pagetemplate.Metadata {
	sorted := append([]Entry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})
	if limit > 0 && len(sorted) > limit {
		sorted = sorted[:limit]
	}

	metas := make([]pagetemplate.Metadata, len(sorted))
	for i, e := range sorted {
		meta := e.Metadata
		meta.Keywords = normalizeKeywords(meta.Keywords)
		metas[i] = pagetemplate.Metadata{Metadata: meta, URL: e.URL}
	}
	return metas
}

// recentPostsPlaceholder is a marker any compiled page can embed (e.g. via
// #html.elem("div", attrs: (id: "recent-posts-here")) in Typst) to have it
// replaced with the recentPostsLimit most recent posts.
const recentPostsPlaceholder = `<div id="recent-posts-here"></div>`

// InjectRecentPosts replaces the placeholder with a set of recent posts, only
// for index.html.
func InjectRecentPosts(outDir string, entries []Entry, templateDir string) error {
	list, err := pagetemplate.RenderEntries(templateDir, recentPostsFor(entries, recentPostsLimit))
	if err != nil {
		return err
	}

	return filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "index.html" {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Contains(content, []byte(recentPostsPlaceholder)) {
			return nil
		}
		replaced := bytes.Replace(content, []byte(recentPostsPlaceholder), []byte(list), 1)
		return os.WriteFile(path, replaced, 0o644)
	})
}
