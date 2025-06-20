package crawler

import (
	"fmt"
	"net/url"
	"sync"

	fetcher "github.com/nrbernard/nightcrawler/internal/fetch"
	normalizer "github.com/nrbernard/nightcrawler/internal/normalize"
	parser "github.com/nrbernard/nightcrawler/internal/parse"
)

type Config struct {
	Pages              map[string]int
	MaxPages           int
	BaseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func NewConfig(baseURL *url.URL, maxPages int, maxConcurrency int) *Config {
	return &Config{
		Pages:              map[string]int{},
		MaxPages:           maxPages,
		BaseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}
}

func (cfg *Config) AddWaitGroup(delta int) {
	cfg.wg.Add(delta)
}

func (cfg *Config) Wait() {
	cfg.wg.Wait()
}

func (cfg *Config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if cfg.Pages[normalizedURL] > 0 {
		cfg.Pages[normalizedURL]++
		return false
	} else {
		cfg.Pages[normalizedURL] = 1
		return true
	}
}

func (cfg *Config) shouldCrawlPage() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.Pages) < cfg.MaxPages
}

func (cfg *Config) CrawlPage(rawCurrentURL string) {
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

	if currentURL.Host != cfg.BaseURL.Host {
		fmt.Println("current URL is not on the same host as the base URL")
		return
	}

	normalizedCurrentURL, err := normalizer.NormalizeURL(currentURL.String())
	if err != nil {
		fmt.Println("error normalizing current URL: ", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizedCurrentURL)
	if !isFirst {
		return
	}

	html, err := fetcher.GetHTML(currentURL.String())
	if err != nil {
		fmt.Println("error getting HTML: ", err)
		return
	}

	urls, err := parser.GetURLsFromHTML(html, currentURL.String())
	if err != nil {
		fmt.Println("error getting URLs from HTML: ", err)
		return
	}

	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.CrawlPage(url)
	}
}
