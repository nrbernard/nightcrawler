package main

import (
	"fmt"
	"net/url"
	"os"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	fmt.Printf("crawling %s\n", rawCurrentURL)
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println("error parsing base URL: ", err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("error parsing current URL: ", err)
		return
	}

	if currentURL.Host != baseURL.Host {
		fmt.Println("current URL is not on the same host as the base URL")
		return
	}

	normalizedCurrentURL, err := normalizeURL(currentURL.String())
	if err != nil {
		fmt.Println("error normalizing current URL: ", err)
		return
	}

	if pages[normalizedCurrentURL] > 0 {
		pages[normalizedCurrentURL]++
		return
	} else {
		pages[normalizedCurrentURL] = 1
	}

	html, err := getHTML(currentURL.String())
	if err != nil {
		fmt.Println("error getting HTML: ", err)
		return
	}

	fmt.Printf("html for %s\n", normalizedCurrentURL)
	fmt.Println(html)

	urls, err := getURLsFromHTML(html, currentURL.String())
	if err != nil {
		fmt.Println("error getting URLs from HTML: ", err)
		return
	}

	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	url := args[0]
	fmt.Printf("starting crawl of: %s\n", url)

	pages := map[string]int{}
	crawlPage(url, url, pages)

	for page, count := range pages {
		fmt.Printf("%s: %d\n", page, count)
	}

	os.Exit(0)
}
