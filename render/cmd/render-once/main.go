// Command render-once builds the site a single time: it compiles every
// blog post and arbitrary page into the output directory, and writes the
// root and per-keyword indexes, then exits.
package main

import (
	"flag"
	"fmt"
	"os"

	"render/site"
)

func main() {
	blogDir := flag.String("blog", "./blog", "directory of dated posts, nested by year")
	pagesDir := flag.String("pages", "./pages", "directory of arbitrary, non-indexed pages")
	outDir := flag.String("out", "./public", "output directory served as the site's root")
	templateDir := flag.String("template", "", "template directory (default: render/template's DefaultDir)")
	flag.Parse()

	entries, err := site.Build(*blogDir, *outDir, *templateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render-once: build posts: %v\n", err)
		os.Exit(1)
	}

	if err := site.WriteIndexes(*outDir, entries, *templateDir); err != nil {
		fmt.Fprintf(os.Stderr, "render-once: write indexes: %v\n", err)
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

	fmt.Printf("rendered %d post(s), %d page(s) into %s\n", len(entries), pageCount, *outDir)
}
