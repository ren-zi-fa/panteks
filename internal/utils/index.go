package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"golang.org/x/net/html"
)

// Scrape fetches a web page and returns the cleaned/minified <body> content
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

// find node <body>
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

// removeTags delete all node with certain name tags
func removeTags(n *html.Node, tags ...string) {
	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.ElementNode {
			if _, found := tagSet[c.Data]; found {
				// dedelete node
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

// removeAttrs 
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

// removeComments 
func removeComments(n *html.Node) {
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.CommentNode {
			//delete comment
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

// ExtractBody get all <body>, remove script/footer/style, atribut style="", komentar, and minify to byte
func ExtractBody(htmlBytes []byte) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	body := findBody(doc)
	if body == nil {
		return nil, fmt.Errorf("no <body> found")
	}

	removeTags(body, "script", "footer", "style")


	removeAttrs(body, "style")

	removeComments(body)

	var buf bytes.Buffer
	if err := html.Render(&buf, body); err != nil {
		return nil, fmt.Errorf("failed to render body: %w", err)
	}

	// minify and remove space lane
	minified := strings.Join(strings.Fields(buf.String()), " ")

	return []byte(minified), nil
}
