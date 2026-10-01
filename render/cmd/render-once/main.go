// Command render-once builds the site a single time and exits.
package main

import (
	"flag"
	"fmt"
	"os"

	"render/site"
)

func main() {
	blogDir := flag.String("blog", "./blog", "directory of dated posts")
	pagesDir := flag.String("pages", "./pages", "directory of non-dated posts")
	staticDir := flag.String("static", "./static", "directory of static assets copied as-is")
	outDir := flag.String("out", "./public", "output directory")
	templateDir := flag.String("template", "./render/template", "template directory")
	blogJaDir := flag.String("blog-ja", "./blog/ja", "directory of dated Japanese posts")
	pagesJaDir := flag.String("pages-ja", "./pages/ja", "directory of non-dated Japanese pages")
	outJaDir := flag.String("out-ja", "./public/ja", "Japanese output directory")
	templateJaDir := flag.String("template-ja", "./render/template-ja", "Japanese template directory")
	flag.Parse()

	enPaths, err := site.ContentPaths(*blogDir, *pagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render-once: scan English content: %v\n", err)
		os.Exit(1)
	}
	jaPaths, err := site.ContentPaths(*blogJaDir, *pagesJaDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render-once: scan Japanese content: %v\n", err)
		os.Exit(1)
	}

	entries, pageCount, err := buildLocale(*blogDir, *pagesDir, *outDir, *templateDir, site.AltLang{Prefix: "/ja", Paths: jaPaths})
	if err != nil {
		fmt.Fprintf(os.Stderr, "render-once: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(*staticDir); err == nil {
		if err := site.CopyStatic(*staticDir, *outDir); err != nil {
			fmt.Fprintf(os.Stderr, "render-once: copy static: %v\n", err)
			os.Exit(1)
		}
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "render-once: stat static dir: %v\n", err)
		os.Exit(1)
	}

	jaEntries, jaPageCount, err := buildLocale(*blogJaDir, *pagesJaDir, *outJaDir, *templateJaDir, site.AltLang{Prefix: "", Paths: enPaths})
	if err != nil {
		fmt.Fprintf(os.Stderr, "render-once: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("rendered %d post(s), %d page(s) into %s\n", len(entries), pageCount, *outDir)
	fmt.Printf("rendered %d post(s), %d page(s) into %s\n", len(jaEntries), jaPageCount, *outJaDir)
}

// buildLocale runs the full build pipeline (posts, indexes, all-posts,
// pages, recent-posts injection, head-style injection) for one locale.
// blogDir and pagesDir are each optional: if a directory doesn't exist, that
// step is skipped, so a locale can start out as pages-only (or even empty).
func buildLocale(blogDir, pagesDir, outDir, templateDir string, alt site.AltLang) ([]site.Entry, int, error) {
	var entries []site.Entry
	if _, statErr := os.Stat(blogDir); statErr == nil {
		var err error
		entries, err = site.RenderBlog(blogDir, outDir, templateDir, alt)
		if err != nil {
			return nil, 0, fmt.Errorf("build posts: %w", err)
		}
	} else if !os.IsNotExist(statErr) {
		return nil, 0, fmt.Errorf("stat blog dir: %w", statErr)
	}

	if err := site.WriteIndexes(outDir, entries, templateDir, alt); err != nil {
		return nil, 0, fmt.Errorf("write indexes: %w", err)
	}

	if err := site.WriteAllPosts(outDir, entries, templateDir, alt); err != nil {
		return nil, 0, fmt.Errorf("write all-posts: %w", err)
	}

	pageCount := 0
	if _, statErr := os.Stat(pagesDir); statErr == nil {
		var err error
		pageCount, err = site.RenderPages(pagesDir, outDir, templateDir, alt)
		if err != nil {
			return nil, 0, fmt.Errorf("render pages: %w", err)
		}
	} else if !os.IsNotExist(statErr) {
		return nil, 0, fmt.Errorf("stat pages dir: %w", statErr)
	}

	if err := site.InjectRecentPosts(outDir, entries, templateDir); err != nil {
		return nil, 0, fmt.Errorf("inject recent posts: %w", err)
	}

	if err := site.InjectHeadStyles(outDir); err != nil {
		return nil, 0, fmt.Errorf("inject head styles: %w", err)
	}

	return entries, pageCount, nil
}
