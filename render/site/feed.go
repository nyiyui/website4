package site

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	nethtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// atomFeed, atomPerson, atomLink, atomEntry, and atomContent mirror the
// subset of RFC 4287 (Atom) this site's feed needs.
type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Xmlns   string      `xml:"xmlns,attr"`
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	Updated string      `xml:"updated"`
	Author  atomPerson  `xml:"author"`
	Links   []atomLink  `xml:"link"`
	Entries []atomEntry `xml:"entry"`
}

type atomPerson struct {
	Name string `xml:"name"`
}

type atomLink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr,omitempty"`
}

type atomEntry struct {
	ID        string      `xml:"id"`
	Title     string      `xml:"title"`
	Updated   string      `xml:"updated"`
	Published string      `xml:"published"`
	Links     []atomLink  `xml:"link"`
	Summary   string      `xml:"summary,omitempty"`
	Content   atomContent `xml:"content"`
}

// atomContent holds a full post's body as escaped HTML character data
// (type="html"), which is what RFC 4287 section 4.1.3.3 requires for inline
// HTML content: encoding/xml's chardata escaping of this string is exactly
// the entity-escaping the spec calls for.
type atomContent struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

// WriteFeed writes a full-content Atom feed of entries, most-recent-first,
// to outDir/feed.xml.
//
// baseURL must be an absolute URL ending in "/" (e.g. "https://example.org/"
// or "https://example.org/ja/"): Atom entries require absolute ids/links,
// but Entry.URL is generated relative to the locale's own root (so it
// resolves correctly under that locale's <base href> when served directly),
// and Entry.Content's root-relative ("/"-prefixed) links point at
// site-wide-root assets shared across locales.
func WriteFeed(outDir string, entries []Entry, baseURL, siteName string) error {
	sorted := append([]Entry(nil), entries...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})

	updated := time.Now().UTC()
	if len(sorted) > 0 {
		updated = sorted[0].Date.UTC()
	}

	feed := atomFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		ID:      baseURL,
		Title:   siteName,
		Updated: updated.Format(time.RFC3339),
		Author:  atomPerson{Name: siteName},
		Links: []atomLink{
			{Rel: "alternate", Href: baseURL, Type: "text/html"},
			{Rel: "self", Href: baseURL + "feed.xml", Type: "application/atom+xml"},
		},
	}

	for _, e := range sorted {
		content, err := absolutizeContent(e.Content, baseURL)
		if err != nil {
			return fmt.Errorf("absolutize content for %s: %w", e.URL, err)
		}
		url := baseURL + e.URL
		date := e.Date.UTC().Format(time.RFC3339)
		feed.Entries = append(feed.Entries, atomEntry{
			ID:        url,
			Title:     e.Title,
			Updated:   date,
			Published: date,
			Links:     []atomLink{{Rel: "alternate", Href: url, Type: "text/html"}},
			Summary:   e.Description,
			Content:   atomContent{Type: "html", Body: string(content)},
		})
	}

	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(outDir, "feed.xml"))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.WriteString(f, xml.Header); err != nil {
		return err
	}
	_, err = f.Write(out)
	return err
}

// absolutizeContent rewrites root-relative href/src/srcset attributes (those
// starting with "/") in content to be absolute under baseURL, so links and
// images still resolve when a post is read in a feed reader outside the
// site's own origin. Fragment-only links (e.g. footnote references and
// backlinks, which start with "#") are left untouched, since they point
// within the same entry.
//
// Typst's native #image() can base64-inline images directly into the HTML;
// none of the current posts do this (they go through static files copied by
// CopyStatic instead), so there's no feed-size concern today, but a post
// that switches to #image() would embed that data verbatim in the feed too.
func absolutizeContent(content template.HTML, baseURL string) (template.HTML, error) {
	nodes, err := nethtml.ParseFragment(strings.NewReader(string(content)), &nethtml.Node{
		Type:     nethtml.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return "", err
	}

	for _, n := range nodes {
		absolutizeNode(n, baseURL)
	}

	var buf strings.Builder
	for _, n := range nodes {
		if err := nethtml.Render(&buf, n); err != nil {
			return "", err
		}
	}
	return template.HTML(buf.String()), nil
}

// absolutizeNode rewrites n's href/src/srcset attributes (if any) in place,
// then recurses into its children.
func absolutizeNode(n *nethtml.Node, baseURL string) {
	if n.Type == nethtml.ElementNode {
		for i, a := range n.Attr {
			switch a.Key {
			case "href", "src":
				if strings.HasPrefix(a.Val, "/") {
					n.Attr[i].Val = baseURL + strings.TrimPrefix(a.Val, "/")
				}
			case "srcset":
				n.Attr[i].Val = absolutizeSrcset(a.Val, baseURL)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		absolutizeNode(c, baseURL)
	}
}

// absolutizeSrcset rewrites a srcset attribute's root-relative URL
// candidates to be absolute under baseURL, preserving any width/density
// descriptor after each URL.
func absolutizeSrcset(srcset, baseURL string) string {
	candidates := strings.Split(srcset, ",")
	for i, c := range candidates {
		parts := strings.SplitN(strings.TrimSpace(c), " ", 2)
		if strings.HasPrefix(parts[0], "/") {
			parts[0] = baseURL + strings.TrimPrefix(parts[0], "/")
		}
		candidates[i] = strings.Join(parts, " ")
	}
	return strings.Join(candidates, ", ")
}
