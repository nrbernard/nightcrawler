package report

import (
	"reflect"
	"testing"
)

func TestSortPages(t *testing.T) {
	tests := []struct {
		name      string
		inputURLs map[string]int
		expected  []ReportPage
	}{
		{
			name:      "sorts by count",
			inputURLs: map[string]int{"https://blog.boot.dev/b-path": 1, "https://blog.boot.dev/a-path": 2},
			expected:  []ReportPage{{URL: "https://blog.boot.dev/a-path", Count: 2}, {URL: "https://blog.boot.dev/b-path", Count: 1}},
		},
		{
			name:      "sorts by count then alphabetically by path",
			inputURLs: map[string]int{"https://blog.boot.dev/b-path": 2, "https://blog.boot.dev/a-path": 2, "https://blog.boot.dev/c-path": 3},
			expected:  []ReportPage{{URL: "https://blog.boot.dev/c-path", Count: 3}, {URL: "https://blog.boot.dev/a-path", Count: 2}, {URL: "https://blog.boot.dev/b-path", Count: 2}},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := sortPages(tc.inputURLs)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected pages: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
