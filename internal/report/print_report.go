package report

import (
	"fmt"
	"sort"
)

type ReportPage struct {
	URL   string
	Count int
}

func sortPages(pages map[string]int) []ReportPage {
	result := make([]ReportPage, 0, len(pages))
	for url, count := range pages {
		result = append(result, ReportPage{URL: url, Count: count})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].URL < result[j].URL
	})

	return result
}

func PrintReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)

	reportPages := sortPages(pages)
	for _, reportPage := range reportPages {
		fmt.Printf("Found %v internal links to %s\n", reportPage.Count, reportPage.URL)
	}
}
