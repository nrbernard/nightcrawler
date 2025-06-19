package reporter

import (
	"fmt"
)

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
