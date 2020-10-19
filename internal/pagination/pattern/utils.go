package pattern

import (
	nurl "net/url"
)

// replaceUrlQueryValue replaces query value of the specified URL. The original URL
// is preserved and not changed. Returns the mutated URL after its query changed.
func replaceUrlQueryValue(url *nurl.URL, queryName string, queryValue string) *nurl.URL {
	clonedURL := *url
	queries := clonedURL.Query()
	queries.Set(queryName, PageParamPlaceholder)
	clonedURL.RawQuery = queries.Encode()
	return &clonedURL
}
