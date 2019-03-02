package types

import (
	"fmt"
	"path"
	"strings"
)

// URL string type value
type URL string

// String returns string value
func (url URL) String() string {
	return fmt.Sprintf("%s", string(url))
}

// Base returns the end of the url + "~/"
func (url URL) Base() string {
	return "~/" + path.Base(url.String())
}

// ReplaceContext replaces all context {param} in the URL
func (url URL) ReplaceContext(mapKeysValues map[string]string) URL {
	var newURL = url.String()
	for key, value := range mapKeysValues {
		if strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}") {
			newURL = strings.Replace(newURL, key, value, -1)
		}
	}
	return URL(newURL)
}
