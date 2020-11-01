// ORIGINAL: javatest/QueryParamPagePatternTest.java

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

package pattern_test

import (
	nurl "net/url"
	"strings"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/pattern"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Pattern_PCPP_IsPagingURL(t *testing.T) {
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/abc-2.html", "http://www.foo.com/a/abc-[*!].html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/abc.html", "http://www.foo.com/a/abc-[*!].html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/abc", "http://www.foo.com/a/abc-[*!]")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/abc-2", "http://www.foo.com/a/abc-[*!]")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/b-c-3", "http://www.foo.com/a/b-[*!]-c-3")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a-c-3", "http://www.foo.com/a-[*!]-c-3")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a-p-1-c-3", "http://www.foo.com/a-p-[*!]-c-3")
	pcppAssertPagingURL(t, false, "http://www.foo.com/a/abc-page", "http://www.foo.com/a/abc-[*!]")
	pcppAssertPagingURL(t, false, "http://www.foo.com/a/2", "http://www.foo.com/a/abc-[*!]")

	pcppAssertPagingURL(t, true, "http://www.foo.com/a/page/2", "http://www.foo.com/a/page/[*!]")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a", "http://www.foo.com/a/page/[*!]")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/page/2/abc.html", "http://www.foo.com/a/page/[*!]/abc.html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/abc.html", "http://www.foo.com/a/page/[*!]/abc.html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/abc.html", "http://www.foo.com/a/[*!]/abc.html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/2/abc.html", "http://www.foo.com/a/[*!]/abc.html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/abc.html", "http://www.foo.com/a/[*!]/abc.html")
	pcppAssertPagingURL(t, true, "http://www.foo.com/a/page/2page", "http://www.foo.com/a/page/[*!]page")
	pcppAssertPagingURL(t, false, "http://www.foo.com/a/page/2", "http://www.foo.com/a/page/[*!]page")
	pcppAssertPagingURL(t, false, "http://www.foo.com/a/page/b", "http://www.foo.com/a/page/[*!]")
	pcppAssertPagingURL(t, false, "http://www.foo.com/m/page/2", "http://www.foo.com/p/page/[*!]")
}

func Test_Pagination_Pattern_PCPP_IsPagePatternValid(t *testing.T) {
	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12",
		"http://www.google.com/forum-12/page/[*!]")
	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12",
		"http://www.google.com/forum-12/[*!]")
	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12",
		"http://www.google.com/forum-12/page-[*!]")

	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12/food",
		"http://www.google.com/forum-12/food/for/bar/[*!]")
	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12-food",
		"http://www.google.com/forum-12-food-[*!]")

	pcppAssertPagePatternValid(t, false,
		"http://www.google.com/forum-12/food",
		"http://www.google.com/forum-12/food/2012/01/[*!]")
	pcppAssertPagePatternValid(t, false,
		"http://www.google.com/forum-12/food/2012/01/01",
		"http://www.google.com/forum-12/food/2012/01/[*!]")

	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/thread/12",
		"http://www.google.com/thread/12/page/[*!]")
	pcppAssertPagePatternValid(t, false,
		"http://www.google.com/thread/12/foo",
		"http://www.google.com/thread/12/page/[*!]/foo")
	pcppAssertPagePatternValid(t, true,
		"http://www.google.com/thread/12/foo",
		"http://www.google.com/thread/12/[*!]/foo")
}

func Test_Pagination_Pattern_PCPP_IsLastNumericPathComponentBad(t *testing.T) {
	// Path component is not numeric i.e. contains non-digits.
	url, _ := nurl.Parse("http://www.foo.com/a2")
	digitStart := strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, false, url.Path, digitStart)

	// Numeric path component is first.
	url, _ = nurl.Parse("http://www.foo.com/2")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, false, url.Path, digitStart)

	// Numeric path component follows a path component that is not a bad page param name.
	url, _ = nurl.Parse("http://www.foo.com/good/2")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, false, url.Path, digitStart)

	// Numeric path component follows a path component that is a bad page param name.
	url, _ = nurl.Parse("http://www.foo.com/wiki/2")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, true, url.Path, digitStart)

	// (s)htm(l) extension doesn't follow digit.
	url, _ = nurl.Parse("http://www.foo.com/2a")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, false, url.Path, digitStart)

	// .htm follows digit, previous path component is not a bad page param name.
	url, _ = nurl.Parse("http://www.foo.com/good/2.htm")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, false, url.Path, digitStart)

	// .html follows digit, previous path component is a bad page param name.
	url, _ = nurl.Parse("http://www.foo.com/wiki/2.html")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, true, url.Path, digitStart)

	// .shtml follows digit, previous path component is not a bad page param name, but the one
	// before that is.
	url, _ = nurl.Parse("http://www.foo.com/wiki/good/2.shtml")
	digitStart = strings.Index(url.Path, "2")
	pcppAssertLastNumericPathComponentBad(t, false, url.Path, digitStart)
}

func pcppAssertPagingURL(t *testing.T, expected bool, strURL string, strPattern string) {
	pattern := createPathComponentPagePattern(strPattern)
	assert.NotNil(t, pattern)
	assert.Equal(t, expected, pattern.IsPagingURL(strURL))
}

func pcppAssertPagePatternValid(t *testing.T, expected bool, strURL string, strPattern string) {
	parsedURL, _ := nurl.ParseRequestURI(strURL)
	assert.NotNil(t, parsedURL)

	pattern := createPathComponentPagePattern(strPattern)
	assert.NotNil(t, pattern)

	assert.Equal(t, expected, pattern.IsValidFor(parsedURL))
}

func pcppAssertLastNumericPathComponentBad(t *testing.T, expected bool, urlPath string, digitStart int) {
	isBad := pattern.IsLastNumericPathComponentBad(urlPath, digitStart, digitStart+1)
	assert.Equal(t, expected, isBad)
}

func createPathComponentPagePattern(strPattern string) pattern.PagePattern {
	// Parse pattern
	url, err := nurl.ParseRequestURI(strPattern)
	if err != nil {
		return nil
	}

	// Get digit location
	digitStart := strings.Index(url.Path, pattern.PageParamPlaceholder)
	digitEnd := digitStart + 1

	// Convert pattern placholder to number
	url.Path = strings.Replace(url.Path, pattern.PageParamPlaceholder, "8", 1)
	url.RawPath = url.Path

	pattern, _ := pattern.NewPathComponentPagePattern(url, digitStart, digitEnd)
	return pattern
}
