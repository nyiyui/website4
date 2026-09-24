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
	flag.Parse()

	entries, err := site.RenderBlog(*blogDir, *outDir, *templateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render-once: build posts: %v\n", err)
		os.Exit(1)
	}

	if err := site.WriteIndexes(*outDir, entries, *templateDir); err != nil {
		fmt.Fprintf(os.Stderr, "render-once: write indexes: %v\n", err)
		os.Exit(1)
	}

	if err := site.WriteAllPosts(*outDir, entries, *templateDir); err != nil {
		fmt.Fprintf(os.Stderr, "render-once: write all-posts: %v\n", err)
		os.Exit(1)
	}

	pageCount := 0
	if _, err := os.Stat(*pagesDir); err == nil {
		pageCount, err = site.RenderPages(*pagesDir, *outDir, *templateDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render-once: render pages: %v\n", err)
			os.Exit(1)
		}
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "render-once: stat pages dir: %v\n", err)
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

	if err := site.InjectRecentPosts(*outDir, entries, *templateDir); err != nil {
		fmt.Fprintf(os.Stderr, "render-once: inject recent posts: %v\n", err)
		os.Exit(1)
	}

	if err := site.InjectHeadStyles(*outDir); err != nil {
		fmt.Fprintf(os.Stderr, "render-once: inject head styles: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("rendered %d post(s), %d page(s) into %s\n", len(entries), pageCount, *outDir)
}
