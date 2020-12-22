// ORIGINAL: java/QueryParamPagePattern.java

// Copyright (c) 2020 Markus Mobius
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pattern

import (
	"errors"
	"fmt"
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
	// Clone URL to make sure we don't mutate the original
	// Since we assume original URL is already absolute, here we only parse
	// without checking error.
	clonedURL, err := nurl.Parse(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to clone url: %w", err)
	}

	// Validate function parameters
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

	// Replace URL queries to PageParamPlaceholder
	queries := clonedURL.Query()
	queries.Set(queryName, PageParamPlaceholder)
	clonedURL.RawQuery = queries.Encode()

	return &QueryParamPagePattern{
		url:        clonedURL,
		strURL:     stringutil.UnescapedString(clonedURL),
		pageNumber: value,
	}, nil
}

func QueryParamPagePatternsFromURL(url *nurl.URL) []PagePattern {
	patterns := []PagePattern{}
	for key, values := range url.Query() {
		for _, value := range values {
			pattern, err := NewQueryParamPagePattern(url, key, value)
			if err == nil && pattern != nil {
				patterns = append(patterns, pattern)
			}
		}
	}

	return patterns
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
