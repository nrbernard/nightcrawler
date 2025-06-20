package parser

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ParseHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	rootNode, err := html.Parse(reader)
	if err != nil {
		return []string{}, err
	}

	return extractURLs(rootNode, rawBaseURL), nil
}

func extractURLs(node *html.Node, rawBaseURL string) []string {
	urls := []string{}

	// Process current node if it's an anchor tag
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				parsedURL, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}

				base, err := url.Parse(rawBaseURL)
				if err != nil {
					continue
				}

				resolved := base.ResolveReference(parsedURL)
				urls = append(urls, resolved.String())
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		urls = append(urls, extractURLs(child, rawBaseURL)...)
	}

	return urls
}
