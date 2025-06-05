package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "https scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "http scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "trailing slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "root URL",
			inputURL: "https://blog.boot.dev/",
			expected: "blog.boot.dev",
		},
		{
			name:     "query parameters",
			inputURL: "https://blog.boot.dev/path?param=value",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "fragment",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "multiple slashes",
			inputURL: "https://blog.boot.dev/path//to//file",
			expected: "blog.boot.dev/path/to/file",
		},
		{
			name:     "subdomain",
			inputURL: "https://sub.blog.boot.dev/path",
			expected: "sub.blog.boot.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
