package parser

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "empty HTML",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
			</body>
		</html>
		`,
			expected: []string{},
		},
		{
			name:     "links with query parameters and fragments",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a href="/search?q=test">Search</a>
				<a href="/page#section">Section</a>
				<a href="https://other.com/page?param=value#frag">External</a>
			</body>
		</html>
		`,
			expected: []string{
				"https://example.com/search?q=test",
				"https://example.com/page#section",
				"https://other.com/page?param=value#frag",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URLs: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
