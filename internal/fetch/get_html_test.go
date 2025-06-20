package fetch

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHTML(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedError  bool
		expectedHTML   string
	}{
		{
			name: "successful HTML response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte("<html><body>Test</body></html>"))
			},
			expectedError: false,
			expectedHTML:  "<html><body>Test</body></html>",
		},
		{
			name: "successful HTML response with charset in header",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte("<html><body>Test</body></html>"))
			},
			expectedError: false,
			expectedHTML:  "<html><body>Test</body></html>",
		},
		{
			name: "error status code",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Not Found"))
			},
			expectedError: true,
			expectedHTML:  "",
		},
		{
			name: "non-HTML content type",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"message": "test"}`))
			},
			expectedError: true,
			expectedHTML:  "",
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Server Error"))
			},
			expectedError: true,
			expectedHTML:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.serverResponse))
			defer server.Close()

			html, err := FetchHTML(server.URL)

			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if html != tc.expectedHTML {
				t.Errorf("expected HTML: %q, got: %q", tc.expectedHTML, html)
			}
		})
	}
}
