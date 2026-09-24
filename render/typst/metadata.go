package typst

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Metadata is the extractes metadata from Typst.
// Although title and description support rich text in Typst, since Typst's HTML
// exporter doesn't support rich text for these (or rather, HTML doesn't).
type Metadata struct {
	Title       string
	Authors     []string
	Description string
	Keywords    []string
	Date        time.Time
}

func GetMetadata(typstPath string, htmlPath string) (Metadata, error) {
	m, err := getOtherMetadata(htmlPath)
	if err != nil {
		return Metadata{}, err
	}
	date, err := getDate(typstPath)
	if err != nil {
		return Metadata{}, err
	}
	m.Date = date
	return m, nil
}

// pattern captures Typst's date formats:
// 1. datetime(year: 2026, month: 9, day: 23)
// 2. "2026-09-23"
// 3. auto or none
const pattern = `date\s*:\s*(datetime\([^)]+\)|"[^"]+"|auto|none)`

var unixEpoch = time.Unix(0, 0)
var errNoAutoDate = errors.New("typst: no auto date supported")
var re = regexp.MustCompile(pattern)

// getDate returns the document date from a Typst file using a Regexp.
// "auto" is not supported since i's not useful for blogs.
// "none" will be returned as the Unix epoch (Jan 1, 1970).
func getDate(typstPath string) (time.Time, error) {
	content, err := os.ReadFile(typstPath)
	if err != nil {
		return time.Time{}, err
	}

	// FindSubmatch returns the full match at [0] and the capture group at [1]
	matches := re.FindSubmatch(content)
	if len(matches) < 2 {
		return time.Time{}, fmt.Errorf("date metadata not found in %s", typstPath)
	}
	rawDate := string(matches[1])
	if rawDate == "auto" {
		return time.Time{}, errNoAutoDate
	}
	if rawDate == "none" {
		return unixEpoch, nil
	}
	// Handle format: "YYYY-MM-DD"
	if strings.HasPrefix(rawDate, `"`) {
		cleanString := strings.Trim(rawDate, `"`)
		parsedTime, err := time.Parse("2006-01-02", cleanString)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse date string %q: %w", cleanString, err)
		}
		return parsedTime, nil
	}

	// Handle format: datetime(year: 2026, month: 9, day: 23)
	if strings.HasPrefix(rawDate, "datetime(") {
		yearRe := regexp.MustCompile(`year\s*:\s*(\d+)`)
		monthRe := regexp.MustCompile(`month\s*:\s*(\d+)`)
		dayRe := regexp.MustCompile(`day\s*:\s*(\d+)`)

		yearMatch := yearRe.FindStringSubmatch(rawDate)
		monthMatch := monthRe.FindStringSubmatch(rawDate)
		dayMatch := dayRe.FindStringSubmatch(rawDate)

		if yearMatch == nil || monthMatch == nil || dayMatch == nil {
			return time.Time{}, fmt.Errorf("missing year, month, or day in datetime definition: %s", rawDate)
		}

		// Errors ignored because the \d+ regex guarantees valid integers
		year, _ := strconv.Atoi(yearMatch[1])
		month, _ := strconv.Atoi(monthMatch[1])
		day, _ := strconv.Atoi(dayMatch[1])

		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, fmt.Errorf("unrecognized date format: %s", rawDate)
}

func getOtherMetadata(htmlPath string) (Metadata, error) {
	f, err := os.Open(htmlPath)
	if err != nil {
		return Metadata{}, err
	}
	doc, err := html.Parse(f)
	if err != nil {
		return Metadata{}, err
	}
	var meta Metadata
	extract(doc, &meta)
	return meta, nil
}

// splitAndTrim splits a comma-delimited string and removes surrounding whitespace
func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}

	rawSlice := strings.Split(s, ",")
	var result []string

	for _, item := range rawSlice {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// extract traverses the HTML node tree to populate the Metadata struct
func extract(n *html.Node, meta *Metadata) {
	if n.Type == html.ElementNode {
		// OPTIMIZATION: Skip the body entirely.
		// Return immediately so we don't traverse its children.
		if n.Data == "body" {
			return
		}

		// Extract <title> tag text
		if n.Data == "title" && n.FirstChild != nil {
			meta.Title = n.FirstChild.Data
		}

		// Extract <meta> tags
		if n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "name":
					name = strings.ToLower(attr.Val)
				case "content":
					content = attr.Val
				}
			}

			switch name {
			case "description":
				meta.Description = content
			case "authors":
				meta.Authors = splitAndTrim(content)
			case "keywords":
				meta.Keywords = splitAndTrim(content)
			}
		}
	}

	// Traverse to the next child
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extract(c, meta)
	}
}
