// ORIGINAL: java/PathComponentPagePattern.java

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
	"path"
	"strconv"
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
)

// PathComponentPagePattern detects the page parameter in the path of a potential pagination URL.
// If detected, it replaces the page param value with PageParamPlaceholder, then creates and returns
// a new object. This object can then be accessed via PagePattern interface to:
// - validate the generated URL page pattern against the document URL
// - determine if a URL is a paging URL based on the page pattern.
// Example: if the original url is "http://www.foo.com/a/b-3.html", the page pattern is
// "http://www.foo.com/a/b-[*!].html".
type PathComponentPagePattern struct {
	url        *nurl.URL
	strURL     string
	pageNumber int

	paramIndex              int
	placeholderStart        int
	placeholderSegmentStart int
	prefix                  string
	suffix                  string
}

func NewPathComponentPagePattern(url *nurl.URL, digitStart, digitEnd int) (*PathComponentPagePattern, error) {
	// Clone URL to make sure we don't mutate the original
	// Since we assume original URL is already absolute, here we only parse
	// without checking error.
	clonedURL, err := nurl.Parse(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to clone url: %w", err)
	}

	// Clean all fragment and queries from URL
	clonedURL.RawQuery = ""
	clonedURL.Fragment = ""
	clonedURL.RawFragment = ""

	// Make sure last numeric path is good
	if IsLastNumericPathComponentBad(clonedURL.Path, digitStart, digitEnd) {
		return nil, errors.New("bad last numeric path component")
	}

	// Parse page number
	valueStr := clonedURL.Path[digitStart:digitEnd]
	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 0 {
		return nil, errors.New("value in path component is invalid number")
	}

	// Create path for URL pattern
	patternPath := clonedURL.Path[:digitStart] + PageParamPlaceholder + clonedURL.Path[digitEnd:]
	clonedURL.Path = patternPath
	clonedURL.RawPath = patternPath
	strURL := stringutil.UnescapedString(clonedURL)

	// Calculate placeholder location
	placeholderStart := strings.Index(strURL, PageParamPlaceholder)
	placeholderSegmentStart := strings.LastIndex(strURL[:placeholderStart], "/")

	// Create prefix
	prefix := strURL[:placeholderSegmentStart]

	// Create suffix, if available
	lenURL := len(strURL)
	lenSuffix := lenURL - placeholderStart - len(PageParamPlaceholder)

	var suffix string
	if lenSuffix > 0 {
		suffix = strURL[lenURL-lenSuffix:]
	}

	// Determine placeholder param index
	paramIndex := -1
	pathComponents := strings.Split(patternPath, "/")
	for i, pathComponent := range pathComponents {
		if strings.Contains(pathComponent, PageParamPlaceholder) {
			paramIndex = i
			break
		}
	}

	// Make sure the placeholder is not the first path component if:
	// - pattern doesn't have suffix; and
	// - placeholder is located at the start of path component.
	if paramIndex == 1 && len(pathComponents) == 2 && suffix == "" &&
		placeholderStart-1 == placeholderSegmentStart {
		return nil, errors.New("page number is in the first path component")
	}

	return &PathComponentPagePattern{
		url:        clonedURL,
		strURL:     strURL,
		pageNumber: value,

		paramIndex:              paramIndex,
		placeholderStart:        placeholderStart,
		placeholderSegmentStart: placeholderSegmentStart,
		prefix:                  prefix,
		suffix:                  suffix,
	}, nil
}

func PathComponentPagePatternsFromURL(url *nurl.URL) []PagePattern {
	urlPath := strings.Trim(url.Path, "/")
	if urlPath == "" || !rxNumber.MatchString(urlPath) {
		return nil
	}

	patterns := []PagePattern{}
	for _, indexes := range rxNumber.FindAllStringIndex(url.Path, -1) {
		start, end := indexes[0], indexes[1]
		pattern, err := NewPathComponentPagePattern(url, start, end)
		if err == nil && pattern != nil {
			patterns = append(patterns, pattern)
		}
	}

	// Sometimes there are URL without number, eg:
	// http://example.com/foo/bar
	// There is possibility that it's the first page, so we'll add it to pattern as well.
	if len(patterns) == 0 {
		newPath := path.Join(url.Path, PageParamPlaceholder)
		start := strings.Index(newPath, PageParamPlaceholder)

		newPath = strings.Replace(newPath, PageParamPlaceholder, "1", 1)
		end := start + 1

		cloned := *url
		cloned.Path = newPath
		cloned.RawPath = newPath

		pattern, err := NewPathComponentPagePattern(&cloned, start, end)
		if err == nil && pattern != nil {
			patterns = append(patterns, pattern)
		}
	}

	return patterns
}

func (pp *PathComponentPagePattern) String() string {
	return pp.strURL
}

func (pp *PathComponentPagePattern) PageNumber() int {
	return pp.pageNumber
}

// IsValidFor returns true if pattern and URL are sufficiently similar and the pattern's
// components are not calendar digits.
func (pp *PathComponentPagePattern) IsValidFor(docURL *nurl.URL) bool {
	urlComponents := strings.Split(docURL.Path, "/")
	patternComponents := strings.Split(pp.url.Path, "/")

	urlComponentsLen := len(urlComponents)
	patternComponentsLen := len(patternComponents)

	// Both the pattern and doc URL must have the similar path.
	if urlComponentsLen > patternComponentsLen {
		return false
	}

	// If both doc URL and page pattern have only 1 component, their common prefix+suffix must
	// be at least half of the entire component in doc URL, e.g doc URL is
	// "foo.com/foo-bar-threads-132" and pattern is "foo.com/foo-bar-threads-132-[*!]".
	if urlComponentsLen == 1 && patternComponentsLen == 1 {
		urlComponent := urlComponents[0]
		patternComponent := patternComponents[0]
		commonPrefixLen := pp.getLongestCommonPrefixLength(urlComponent, patternComponent)
		commonSuffixLen := pp.getLongestCommonSuffixLength(urlComponent, patternComponent, commonPrefixLen)
		return (commonSuffixLen+commonPrefixLen)*2 >= len(urlComponent)
	}

	if !pp.hasSamePathComponentsAs(docURL) {
		return false
	}

	if pp.isCalendarPage() {
		return false
	}

	return true
}

// IsPagingURL returns true if a URL matches this page pattern based on a pipeline of rules:
// - suffix (part of pattern after page param placeholder) must be same, and
// - different set of rules depending on if page param is at start of path component or not.
func (pp *PathComponentPagePattern) IsPagingURL(url string) bool {
	// Both url and pattern must have the same suffix, if available.
	if pp.suffix != "" && !strings.HasSuffix(url, pp.suffix) {
		return false
	}

	if pp.strURL[pp.placeholderStart-1] == '/' {
		return pp.isPagingUrlForStartOfPathComponent(url)
	}

	return pp.isPagingUrlForNotStartOfPathComponent(url)
}

func (pp *PathComponentPagePattern) getLongestCommonPrefixLength(str1, str2 string) int {
	if str1 == "" || str2 == "" {
		return 0
	}

	limit := len(str1)
	if lenStr2 := len(str2); lenStr2 < limit {
		limit = lenStr2
	}

	var i int
	for i = 0; i < limit; i++ {
		if str1[i] != str2[i] {
			break
		}
	}

	return i
}

func (pp *PathComponentPagePattern) getLongestCommonSuffixLength(str1, str2 string, startIndex int) int {
	var commonSuffixLen int
	for i, j := len(str1)-1, len(str2)-1; i > startIndex && j > startIndex; i, j = i-1, j-1 {
		if str1[i] != str2[j] {
			break
		}
		commonSuffixLen++
	}
	return commonSuffixLen
}

// hasSamePathComponentsAs returns true if, except for the path component containing the page param, the
// other path components of doc URL are the same as pattern's. But pattern may have more components, e.g.:
// - doc URL is /thread/12, pattern is /thread/12/page/[*!]
//   returns true because "thread" and "12" in doc URL match those in pattern
// - doc URL is /thread/12/foo, pattern is /thread/12/page/[*!]/foo
//   returns false because "foo" in doc URL doesn't match "page" in pattern whose page param
//   path component comes after.
// - doc URL is /thread/12/foo, pattern is /thread/12/[*!]/foo
//   returns true because "foo" in doc URL would match "foo" in pattern whose page param path
//   component is skipped when matching.
// - doc URL is /thread/foo.html, pattern is /thread/foo/[*!]
//   returns true because "foo.html" in doc URL would match "foo" in pattern whose trailing
//   html extension is skipped when matching.
func (pp *PathComponentPagePattern) hasSamePathComponentsAs(parsedURL *nurl.URL) bool {
	// Trim trailing shtml extension from doc URL path
	parsedURLPath := rxEndOrHasSHTML.ReplaceAllString(parsedURL.Path, "")

	urlComponents := strings.Split(parsedURLPath, "/")
	patternComponents := strings.Split(pp.url.Path, "/")
	passedParamComponent := false

	for i, j := 0, 0; i < len(urlComponents) && j < len(patternComponents); i, j = i+1, j+1 {
		if i == pp.paramIndex && !passedParamComponent {
			passedParamComponent = true

			// Repeat current path component if doc URL has less components (as per comments
			// just above, doc URL may have less components).
			if len(urlComponents) < len(patternComponents) {
				i--
			}
			continue
		}

		if !stringutil.EqualsIgnoreCase(urlComponents[i], patternComponents[j]) {
			return false
		}
	}

	return true
}

// isCalendarPage returns true if pattern is for a calendar page, e.g. 2012/01/[*!], which
// would be a false-positive.
func (pp *PathComponentPagePattern) isCalendarPage() bool {
	if pp.paramIndex < 2 {
		return false
	}

	// Only if param is the entire path component. This handles some cases erroneously
	// considered false-positives e.g. first page is
	// http://www.politico.com/story/2014/07/barack-obama-immigration-legal-questions-109467.html,
	// and second page is
	// http://www.politico.com/story/2014/07/barack-obama-immigration-legal-questions-109467_Page2.html,
	// would be considered false-positives otherwise because of "2014" and "07".
	patternComponents := strings.Split(pp.url.Path, "/")
	if len(patternComponents[pp.paramIndex]) != len(PageParamPlaceholder) {
		return false
	}

	month, _ := strconv.Atoi(patternComponents[pp.paramIndex-1])
	if month >= 1 && month <= 12 {
		year, _ := strconv.Atoi(patternComponents[pp.paramIndex-2])
		if year > 1970 && year < 3000 {
			return true
		}
	}

	return false
}

// isPagingUrlForStartOfPathComponent returns true if url is a paging URL based on the page pattern
// where the page param is at the start of a path component.
// If the page pattern is www.foo.com/a/[*!]/abc.html, expected doc URL is:
// - www.foo.com/a/2/abc.html
// - www.foo.com/a/abc.html
// - www.foo.com/abc.html.
func (pp *PathComponentPagePattern) isPagingUrlForStartOfPathComponent(url string) bool {
	urlLen := len(url)
	suffixLen := len(pp.suffix)
	suffixStart := urlLen - suffixLen

	urlOrigin := strings.Index(pp.strURL, pp.url.Path)
	prevComponentPos := strings.LastIndex(pp.strURL[:pp.placeholderSegmentStart], "/")
	if prevComponentPos >= urlOrigin {
		if prevComponentPos+suffixLen == urlLen {
			// The url doesn't have page number param and previous path component, like
			// www.foo.com/abc.html.
			return url[:prevComponentPos] == pp.strURL[:prevComponentPos]
		}
	}

	// If both url and pattern have the same prefix, url must have nothing else.
	if strings.HasPrefix(url, pp.prefix) {
		acceptLen := pp.placeholderSegmentStart + suffixLen
		// The url doesn't have page number parameter, like www.foo.com/a/abc.html.
		if acceptLen == urlLen {
			return true
		}
		if acceptLen > urlLen {
			return false
		}

		// While we are here, the url must have page number param, so the url must have a '/'
		// at the pattern's path component start position.
		if url[pp.placeholderSegmentStart] != '/' {
			return false
		}

		val, err := strconv.Atoi(url[pp.placeholderSegmentStart+1 : suffixStart])
		return err == nil && val >= 0
	}

	return false
}

// isPagingUrlForNotStartOfPathComponent returns true if url is a paging URL based on the page
// pattern where the page param is not at the start of a path component.
// If the page pattern is www.foo.com/a/abc-[*!].html, expected doc URL is:
// - www.foo.com/a/abc-2.html
// - www.foo.com/a/abc.html.
func (pp *PathComponentPagePattern) isPagingUrlForNotStartOfPathComponent(url string) bool {
	urlLen := len(url)
	suffixLen := len(pp.suffix)
	suffixStart := urlLen - suffixLen

	// The page param path component of both url and pattern must have the same prefix.
	if !strings.HasPrefix(url, pp.prefix) {
		return false
	}

	// Find the first different character in page param path component just before
	// placeholder or suffix, then check if it's acceptable.
	maxPos := pp.placeholderStart
	if suffixStart < maxPos {
		maxPos = suffixStart
	}

	firstDiffPos := pp.placeholderSegmentStart
	for ; firstDiffPos < maxPos; firstDiffPos++ {
		if url[firstDiffPos] != pp.strURL[firstDiffPos] {
			break
		}
	}

	if firstDiffPos == suffixStart { // First different character is the suffix.
		if firstDiffPos+1 == pp.placeholderStart &&
			rxPageParamSeparator.MatchString(string(pp.strURL[firstDiffPos])) {
			return true
		}

		// If the url doesn't have page parameter, it is fine.
		if firstDiffPos+suffixLen == urlLen {
			return true
		}
	} else if firstDiffPos == pp.placeholderStart { // First different character is page param.
		val, err := strconv.Atoi(url[firstDiffPos:suffixStart])
		if err == nil && val >= 0 {
			return true
		}
	}

	return false
}

// IsLastNumericPathComponentBad returns true if :
// - the digitStart to digitEnd of urlStr is the last path component, and
// - the entire path component is numeric, and
// - the previous path component is a bad page param name.
// E.g. "www.foo.com/tag/2" will return true because of the above reasons and "tag"
// is a bad page param.
func IsLastNumericPathComponentBad(urlPath string, digitStart, digitEnd int) bool {
	postMatch := urlPath[digitEnd:]

	// Checks that this is the last path component, and trailing characters, if available,
	// are (s)htm(l) extensions.
	if rxEndOrHasSHTML.MatchString(postMatch) {
		// Entire component is numeric, get previous path component.
		prevPathComponent := rxLastPathComponent.FindStringSubmatch(urlPath[:digitStart])
		if len(prevPathComponent) > 0 {
			if _, isBad := badPageParamNames[prevPathComponent[1]]; isBad {
				return true
			}
		}
	}

	return false
}
