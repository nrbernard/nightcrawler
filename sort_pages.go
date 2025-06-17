package main

import "sort"

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
