// ORIGINAL: java/QueryParamPagePattern.java

package pattern

import (
	"errors"
	nurl "net/url"
	"strconv"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
)

// QueryParamPagePattern detects the page parameter in the query of a potential pagination URL. If
// detected, it replaces the page param value with PageParamPlaceholder, then creates and returns
// a new object. This object can then be called via PagePattern interface to:
// - validate the generated URL page pattern against the document URL
// - determine if a URL is a paging URL based on the page pattern.
// Example: if the original url is "http://www.foo.com/a/b/?page=2&query=a", the page pattern is
// "http://www.foo.com/a/b?page=[*!]&query=a"
type QueryParamPagePattern struct {
	url        *nurl.URL
	strURL     string
	pageNumber int
}

func NewQueryParamPagePattern(url *nurl.URL, queryName, queryValue string) (*QueryParamPagePattern, error) {
	if queryName == "" {
		return nil, errors.New("query name must not empty")
	}

	if queryValue == "" {
		return nil, errors.New("query value must not empty")
	}

	if !stringutil.IsStringAllDigit(queryValue) {
		return nil, errors.New("query value has non-digits: " + queryValue)
	}

	if _, isBad := badPageParamNames[queryName]; isBad {
		return nil, errors.New("query name is bad page param name: " + queryName)
	}

	value, err := strconv.Atoi(queryValue)
	if err != nil {
		return nil, errors.New("query value is invalid number: " + queryValue)
	}

	newURL := replaceUrlQueryValue(url, queryName, PageParamPlaceholder)
	return &QueryParamPagePattern{
		url:        newURL,
		strURL:     newURL.String(),
		pageNumber: value,
	}, nil
}

func (pp *QueryParamPagePattern) String() string {
	return pp.strURL
}

func (pp *QueryParamPagePattern) PageNumber() int {
	return pp.pageNumber
}

// IsValidFor returns true if page pattern and URL have the same path components.
func (pp *QueryParamPagePattern) IsValidFor(docURL *nurl.URL) bool {
	urlPath := rxTrailingSlashHTML.ReplaceAllString(pp.url.Path, "")
	docURLPath := rxTrailingSlashHTML.ReplaceAllString(docURL.Path, "")

	return pp.url.Scheme == docURL.Scheme &&
		pp.url.Host == docURL.Host &&
		urlPath == docURLPath
}

// IsPagingURL returns true if a URL matches this page pattern based on a pipeline of rules:
// - suffix (part of pattern after page param placeholder) must be same, and
// - scheme, host, and path must be same, and
// - query params, except that for page number, must be same in value, and
// - query value must be a plain number.
func (pp *QueryParamPagePattern) IsPagingURL(url string) bool {
	// Parse URL
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil {
		return false
	}

	// Make sure URL has same prefix as the pattern
	patternURLPath := rxTrailingSlashHTML.ReplaceAllString(pp.url.Path, "")
	parsedURLPath := rxTrailingSlashHTML.ReplaceAllString(parsedURL.Path, "")
	if pp.url.Scheme != parsedURL.Scheme ||
		pp.url.Host != parsedURL.Host ||
		patternURLPath != parsedURLPath {
		return false
	}

	// All queries in parsed URL must exist in the pattern URL
	patternURLQueries := pp.url.Query()
	parsedURLQueries := parsedURL.Query()

	for key := range parsedURLQueries {
		if _, exist := patternURLQueries[key]; !exist {
			return false
		}
	}

	// All queries (except page number) in pattern URL must exist in the parsed URL
	// If page number query exist in parsed URL, it must be number only.
	for key, patternParamValues := range patternURLQueries {
		// Check if this parameter is for page number
		isPageNumberParam := false
		for _, value := range patternParamValues {
			if value == PageParamPlaceholder {
				isPageNumberParam = true
				break
			}
		}

		// If it's not for page number, this parameter must exist in parsed URL
		parsedParamValues, parsedParamExist := parsedURLQueries[key]
		if !isPageNumberParam && !parsedParamExist {
			return false
		}

		// If it's for page number, the value must be number
		if isPageNumberParam {
			isNumberOnly := false
			for _, value := range parsedParamValues {
				if _, err := strconv.Atoi(value); err == nil {
					isNumberOnly = true
					break
				}
			}

			if parsedParamExist && !isNumberOnly {
				return false
			}

			continue
		}

		// Make sure the parameter values are the same
		if len(patternParamValues) != len(parsedParamValues) {
			return false
		}

		for i, patternParamValue := range patternParamValues {
			parsedParamValue := parsedParamValues[i]
			if patternParamValue != parsedParamValue {
				return false
			}
		}
	}

	return true
}
