package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	appleDeveloperHost = "https://developer.apple.com"
	appleDocsDataHost  = "https://developer.apple.com/tutorials/data"

	rootDocPath = "/documentation/apple_ads"
)

type docPage struct {
	References    map[string]docRef   `json:"references"`
	TopicSections []docTopicSection   `json:"topicSections"`
	SeeAlso       []docSeeAlsoSection `json:"seeAlsoSections"`
}

type docTopicSection struct {
	Title       string   `json:"title"`
	Anchor      string   `json:"anchor"`
	Identifiers []string `json:"identifiers"`
}

type docSeeAlsoSection struct {
	Title       string   `json:"title"`
	Identifiers []string `json:"identifiers"`
}

type docRef struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Role  string `json:"role"`
	Kind  string `json:"kind"`
}

func main() {
	var outPath string
	flag.StringVar(&outPath, "out", "docs/OFFICIAL_APPLE_ADS_DOC_INDEX.md", "Output markdown path")
	flag.Parse()

	client := &http.Client{Timeout: 30 * time.Second}

	pagesByDocPath := map[string]*docPage{}
	allPages := map[string]string{} // docPath -> title

	// Root must be present; fetch it explicitly so we can fail with a useful error.
	root, err := fetchDocPage(client, rootDocPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch root docs JSON (%s): %v\n", rootDocPath, err)
		os.Exit(1)
	}
	pagesByDocPath[rootDocPath] = root

	visited := map[string]bool{rootDocPath: true}
	queue := []string{}

	for _, ref := range root.References {
		if !strings.HasPrefix(ref.URL, rootDocPath) {
			continue
		}
		if ref.Title != "" {
			allPages[ref.URL] = ref.Title
		}
		if shouldFetch(ref) {
			queue = append(queue, ref.URL)
		}
	}

	for len(queue) > 0 {
		docPath := queue[0]
		queue = queue[1:]
		if visited[docPath] {
			continue
		}
		visited[docPath] = true

		page, err := fetchDocPage(client, docPath)
		if err != nil {
			// If Apple doesn't expose JSON for some doc paths, ignore them.
			// The index remains link-only, so missing topic section expansion is acceptable.
			continue
		}
		pagesByDocPath[docPath] = page

		for _, ref := range page.References {
			if !strings.HasPrefix(ref.URL, rootDocPath) {
				continue
			}
			if ref.Title != "" {
				allPages[ref.URL] = ref.Title
			}
			if shouldFetch(ref) && !visited[ref.URL] {
				queue = append(queue, ref.URL)
			}
		}
	}

	out, err := buildMarkdownIndex(root, pagesByDocPath, allPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build markdown: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outPath, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
		os.Exit(1)
	}
}

func shouldFetch(ref docRef) bool {
	// "collectionGroup" pages contain topicSections that enumerate endpoints/objects.
	// Some deprecated areas (like Creative Sets) are marked as "article" but still
	// include topicSections, so we fetch those too.
	if ref.Role == "collectionGroup" {
		return true
	}
	if ref.Kind == "article" && ref.Role == "article" {
		return true
	}
	return false
}

func fetchDocPage(client *http.Client, docPath string) (*docPage, error) {
	// Apple Developer docs are rendered from DocC JSON under /tutorials/data.
	// The mapping is stable: {docPath}.json
	url := appleDocsDataHost + docPath + ".json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	var p docPage
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func buildMarkdownIndex(root *docPage, pagesByDocPath map[string]*docPage, allPages map[string]string) ([]byte, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	var b bytes.Buffer
	b.WriteString("# Official Apple Ads Docs (Full Index)\n\n")
	b.WriteString("This file is **link-only**: it indexes Apple's official Apple Ads documentation without copying the content.\n\n")
	b.WriteString(fmt.Sprintf("Generated: `%s` from `%s%s`.\n\n", now, appleDeveloperHost, rootDocPath))
	b.WriteString("Regenerate:\n")
	b.WriteString("- `go run ./tools/appleadsdocsindex`\n\n")

	b.WriteString("Related:\n")
	b.WriteString("- CLI mapping: `docs/OFFICIAL_APPLE_ADS_DOCS.md`\n\n")

	// Top-level groups from the root doc page.
	for _, ts := range root.TopicSections {
		if ts.Title == "" {
			continue
		}
		b.WriteString("## " + ts.Title + "\n\n")

		for _, ident := range ts.Identifiers {
			ref, ok := root.References[ident]
			if !ok || ref.URL == "" || ref.Title == "" {
				continue
			}
			b.WriteString(fmt.Sprintf("- [%s](%s%s)\n", ref.Title, appleDeveloperHost, ref.URL))
		}
		b.WriteString("\n")

		// Expand each top-level collection under this group if we fetched it.
		for _, ident := range ts.Identifiers {
			ref, ok := root.References[ident]
			if !ok || ref.URL == "" || ref.Title == "" {
				continue
			}
			p := pagesByDocPath[ref.URL]
			if p == nil || len(p.TopicSections) == 0 {
				continue
			}

			b.WriteString("### " + ref.Title + "\n\n")
			writeTopicSections(&b, p)
		}
	}

	// Include additional collections we discovered that aren't top-level (e.g. "Usability and Errors").
	extraCollections := make([]string, 0, len(pagesByDocPath))
	topLevelSet := map[string]bool{rootDocPath: true}
	for _, ts := range root.TopicSections {
		for _, ident := range ts.Identifiers {
			if ref, ok := root.References[ident]; ok && ref.URL != "" {
				topLevelSet[ref.URL] = true
			}
		}
	}
	for docPath, p := range pagesByDocPath {
		if topLevelSet[docPath] {
			continue
		}
		if p == nil || len(p.TopicSections) == 0 {
			continue
		}
		extraCollections = append(extraCollections, docPath)
	}
	sort.Strings(extraCollections)
	if len(extraCollections) > 0 {
		b.WriteString("## Additional Collections\n\n")
		for _, docPath := range extraCollections {
			p := pagesByDocPath[docPath]
			title := allPages[docPath]
			if title == "" {
				title = docPath
			}
			b.WriteString(fmt.Sprintf("### [%s](%s%s)\n\n", title, appleDeveloperHost, docPath))
			writeTopicSections(&b, p)
		}
	}

	// Alphabetical list of everything discovered.
	type item struct {
		Title string
		URL   string
	}
	var items []item
	for url, title := range allPages {
		if url == "" || title == "" {
			continue
		}
		items = append(items, item{Title: title, URL: url})
	}
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})

	b.WriteString("## All Pages (Alphabetical)\n\n")
	for _, it := range items {
		b.WriteString(fmt.Sprintf("- [%s](%s%s)\n", it.Title, appleDeveloperHost, it.URL))
	}
	b.WriteString("\n")

	return b.Bytes(), nil
}

func writeTopicSections(b *bytes.Buffer, p *docPage) {
	for _, ts := range p.TopicSections {
		if ts.Title == "" {
			continue
		}
		b.WriteString("#### " + ts.Title + "\n\n")
		for _, ident := range ts.Identifiers {
			ref, ok := p.References[ident]
			if !ok || ref.URL == "" || ref.Title == "" {
				continue
			}
			b.WriteString(fmt.Sprintf("- [%s](%s%s)\n", ref.Title, appleDeveloperHost, ref.URL))
		}
		b.WriteString("\n")
	}
}
