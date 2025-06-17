package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	pages              map[string]int
	maxPages           int
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

func (cfg *config) shouldCrawlPage() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.pages) < cfg.maxPages
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if !cfg.shouldCrawlPage() {
		return
	}

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

	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("error parsing base URL: ", err)
		os.Exit(1)
	}

	maxConcurrency := 3
	maxPages := 10

	if len(args) > 1 {
		concurrencyArg, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("error parsing max concurrency: ", err)
			os.Exit(1)
		}
		maxConcurrency = concurrencyArg
	}

	if len(args) > 2 {
		pagesArg, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("error parsing max pages: ", err)
			os.Exit(1)
		}
		maxPages = pagesArg
	}

	cfg := config{
		pages:              map[string]int{},
		maxPages:           maxPages,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	fmt.Printf("starting crawl of: %s\n", baseURL.String())

	cfg.wg.Add(1)
	cfg.crawlPage(baseURL.String())

	cfg.wg.Wait()

	fmt.Printf("pages crawled: %v\n", len(cfg.pages))
	for page, count := range cfg.pages {
		fmt.Printf("%s: %d\n", page, count)
	}
	fmt.Printf("crawl completed in %v\n", time.Since(start))
	os.Exit(0)
}
