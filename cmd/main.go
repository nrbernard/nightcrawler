package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	crawler "github.com/nrbernard/nightcrawler/internal/crawler"
	report "github.com/nrbernard/nightcrawler/internal/report"
)

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

	cfg := crawler.NewConfig(baseURL, maxPages, maxConcurrency)

	fmt.Printf("starting crawl of: %s\n", baseURL.String())

	cfg.AddWaitGroup(1)
	cfg.CrawlPage(baseURL.String())

	cfg.Wait()

	report.PrintReport(cfg.Pages, baseURL.String())
	fmt.Printf("crawl completed in %v\n", time.Since(start))
	os.Exit(0)
}
