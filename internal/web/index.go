package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Scrape fetch halaman web, lalu mengembalikan isi <body> bersih/minify
func Scrape(ctx context.Context, targetURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return ExtractBody(body)
}

// cari node <body>
func findBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findBody(c); result != nil {
			return result
		}
	}
	return nil
}

// removeTags hapus semua node dengan nama tertentu
func removeTags(n *html.Node, tags ...string) {
	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.ElementNode {
			if _, found := tagSet[c.Data]; found {
				// hapus node
				if c.PrevSibling != nil {
					c.PrevSibling.NextSibling = c.NextSibling
				} else {
					n.FirstChild = c.NextSibling
				}
				if c.NextSibling != nil {
					c.NextSibling.PrevSibling = c.PrevSibling
				} else {
					n.LastChild = c.PrevSibling
				}
				c = next
				continue
			}
		}
		removeTags(c, tags...)
		c = next
	}
}

// removeAttrs hapus atribut tertentu
func removeAttrs(n *html.Node, attrs ...string) {
	attrSet := make(map[string]struct{}, len(attrs))
	for _, a := range attrs {
		attrSet[a] = struct{}{}
	}

	if n.Type == html.ElementNode {
		filtered := make([]html.Attribute, 0, len(n.Attr))
		for _, a := range n.Attr {
			if _, found := attrSet[a.Key]; !found {
				filtered = append(filtered, a)
			}
		}
		n.Attr = filtered
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeAttrs(c, attrs...)
	}
}

// removeComments hapus semua komentar <!-- -->
func removeComments(n *html.Node) {
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.CommentNode {
			// hapus node komentar
			if c.PrevSibling != nil {
				c.PrevSibling.NextSibling = c.NextSibling
			} else {
				n.FirstChild = c.NextSibling
			}
			if c.NextSibling != nil {
				c.NextSibling.PrevSibling = c.PrevSibling
			} else {
				n.LastChild = c.PrevSibling
			}
		} else {
			removeComments(c)
		}
		c = next
	}
}

// ExtractBody ambil isi <body>, hapus script/footer/style, atribut style="", komentar, lalu minify
func ExtractBody(htmlBytes []byte) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	body := findBody(doc)
	if body == nil {
		return nil, fmt.Errorf("no <body> found")
	}

	// hapus tag
	removeTags(body, "script", "footer", "style")

	// hapus atribut style
	removeAttrs(body, "style")

	// hapus komentar
	removeComments(body)

	// render ulang
	var buf bytes.Buffer
	if err := html.Render(&buf, body); err != nil {
		return nil, fmt.Errorf("failed to render body: %w", err)
	}

	// minify sederhana: buang newline dan spasi berlebih
	minified := strings.Join(strings.Fields(buf.String()), " ")

	return []byte(minified), nil
}
