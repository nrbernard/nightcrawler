package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if cfg.pages[normalizedURL] > 0 {
		cfg.pages[normalizedURL]++
		return false
	} else {
		cfg.pages[normalizedURL] = 1
		return true
	}
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	fmt.Printf("crawling %s\n", rawCurrentURL)

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("error parsing current URL: ", err)
		return
	}

	if currentURL.Host != cfg.baseURL.Host {
		fmt.Println("current URL is not on the same host as the base URL")
		return
	}

	normalizedCurrentURL, err := normalizeURL(currentURL.String())
	if err != nil {
		fmt.Println("error normalizing current URL: ", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizedCurrentURL)
	if !isFirst {
		return
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
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}

func main() {
	start := time.Now()

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("error parsing base URL: ", err)
		os.Exit(1)
	}

	cfg := config{
		pages:              map[string]int{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 10),
		wg:                 &sync.WaitGroup{},
	}

	fmt.Printf("starting crawl of: %s\n", baseURL.String())

	cfg.wg.Add(1)
	cfg.crawlPage(baseURL.String())

	cfg.wg.Wait()

	fmt.Printf("pages crawled: %v\n", len(cfg.pages))
	fmt.Printf("crawl completed in %v\n", time.Since(start))
	os.Exit(0)
}
