package normalizer

import (
	"net/url"
	"regexp"
	"strings"
)

var multiSlashRegex = regexp.MustCompile(`/+`)

func NormalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	path := strings.Trim(parsedURL.Path, "/")
	path = multiSlashRegex.ReplaceAllString(path, "/")

	return strings.TrimSuffix(parsedURL.Host+"/"+path, "/"), nil
}
