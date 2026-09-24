package typst

import (
	"bytes"
	"testing"

	"golang.org/x/net/html"
)

func renderInner(t *testing.T, n *html.Node) string {
	t.Helper()
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&buf, c); err != nil {
			t.Fatalf("render: %v", err)
		}
	}
	return buf.String()
}

func TestGetTOC(t *testing.T) {
	entries, err := GetTOC("testdata/sample.html")
	if err != nil {
		t.Fatalf("GetTOC: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("want 2 top-level entries, got %d: %+v", len(entries), entries)
	}

	intro := entries[0]
	if intro.Level != 2 || intro.Node.Data != "h2" {
		t.Errorf("intro: want level=2 h2, got level=%d node=%s", intro.Level, intro.Node.Data)
	}
	if got := renderInner(t, intro.Node); got != "Introduction" {
		t.Errorf("intro text: got %q", got)
	}
	if len(intro.Children) != 3 {
		t.Fatalf("intro children: want 3, got %d: %+v", len(intro.Children), intro.Children)
	}

	background := intro.Children[0]
	if background.Level != 3 {
		t.Errorf("background level: want 3, got %d", background.Level)
	}
	if got, want := renderInner(t, background.Node), "Background <em>emph</em> and <code>code</code>"; got != want {
		t.Errorf("background rich text: got %q, want %q", got, want)
	}

	methods := intro.Children[2]
	if methods.Level != 3 || len(methods.Children) != 2 {
		t.Fatalf("methods: want level=3 with 2 children, got level=%d children=%d", methods.Level, len(methods.Children))
	}
	if methods.Children[0].Level != 4 {
		t.Errorf("data collection level: want 4, got %d", methods.Children[0].Level)
	}

	conclusion := entries[1]
	if conclusion.Level != 2 || len(conclusion.Children) != 0 {
		t.Errorf("conclusion: want level=2 with 0 children, got level=%d children=%d", conclusion.Level, len(conclusion.Children))
	}
}
