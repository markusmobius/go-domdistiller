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
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/pattern"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Pattern_QPPP_IsPagingURL(t *testing.T) {
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryA=v1&queryB=4&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryA=v1&queryB=growl&queryB=5&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryA=v1&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryB=2&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryC=v3&queryC=v4",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3&queryC=v4")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?page=3",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b/",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b.htm",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b.html",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, false,
		"http://www.foo.com/a/b?queryA=v1&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, false,
		"http://www.foo.com/a/b?queryB=bar&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, false,
		"http://www.foo.com/a/b?queryA=v1",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
}

func Test_Pagination_Pattern_QPPP_IsPagePatternValid(t *testing.T) {
	qpppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12",
		"http://www.google.com/forum-12?page=[*!]")
	qpppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12?sid=12345",
		"http://www.google.com/forum-12?page=[*!]&sort=d")
	qpppAssertPagePatternValid(t, false,
		"http://www.google.com/a/forum-12?sid=12345",
		"http://www.google.com/b/forum-12?page=[*!]&sort=d")
	qpppAssertPagePatternValid(t, false,
		"http://www.google.com/forum-11?sid=12345",
		"http://www.google.com/forum-12?page=[*!]&sort=d")
}

func qpppAssertPagingURL(t *testing.T, expected bool, strURL string, strPattern string) {
	pattern := createQueryParamPagePattern(strPattern)
	assert.NotNil(t, pattern)
	assert.Equal(t, expected, pattern.IsPagingURL(strURL))
}

func qpppAssertPagePatternValid(t *testing.T, expected bool, strURL string, strPattern string) {
	parsedURL, _ := nurl.ParseRequestURI(strURL)
	assert.NotNil(t, parsedURL)

	pattern := createQueryParamPagePattern(strPattern)
	assert.NotNil(t, pattern)

	assert.Equal(t, expected, pattern.IsValidFor(parsedURL))
}

func createQueryParamPagePattern(strPattern string) pattern.PagePattern {
	url, err := nurl.ParseRequestURI(strPattern)
	if err != nil {
		return nil
	}

	for key, values := range url.Query() {
		for _, value := range values {
			if value == pattern.PageParamPlaceholder {
				pattern, _ := pattern.NewQueryParamPagePattern(url, key, "8")
				return pattern
			}
		}
	}

	return nil
}
