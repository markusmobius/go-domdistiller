// ORIGINAL: Part of javatest/StringUtilTest.java

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

// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package stringutil_test

import (
	nurl "net/url"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/stretchr/testify/assert"
)

func Test_StringUtil_IsStringAllWhitespace(t *testing.T) {
	assert.True(t, stringutil.IsStringAllWhitespace(""))
	assert.True(t, stringutil.IsStringAllWhitespace(" \t\r\n"))
	assert.True(t, stringutil.IsStringAllWhitespace(" \u00a0     \t\t\t"))
	assert.False(t, stringutil.IsStringAllWhitespace("a"))
	assert.False(t, stringutil.IsStringAllWhitespace("     a  "))
	assert.False(t, stringutil.IsStringAllWhitespace("\u00a0\u0460"))
	assert.False(t, stringutil.IsStringAllWhitespace("\n\t_ "))
}

// =================================================================================
// Tests below these point are test for function that doesn't exist in original code
// =================================================================================

func Test_StringUtil_CreateAbsoluteURL(t *testing.T) {
	relURL1 := "#here"
	relURL2 := "/test/123"
	relURL3 := "test/123"
	relURL4 := "//www.google.com"
	relURL5 := "https://www.google.com"
	relURL6 := "ftp://ftp.server.com"
	relURL7 := "www.google.com"
	relURL8 := "http//www.google.com"
	relURL9 := "../hello/relative"
	relURL10 := "image.png"

	absURL1 := "#here"
	absURL2 := "http://example.com/test/123"
	absURL3 := "http://example.com/page/test/123"
	absURL4 := "http://www.google.com"
	absURL5 := "https://www.google.com"
	absURL6 := "ftp://ftp.server.com"
	absURL7 := "http://example.com/page/www.google.com"
	absURL8 := "http://example.com/page/http//www.google.com"
	absURL9 := "http://example.com/hello/relative"
	absURL10 := "http://example.com/page/image.png"

	baseURL, _ := nurl.ParseRequestURI("http://example.com/page/")
	baseURL2, _ := nurl.ParseRequestURI("http://example.com/page/doc.html")
	assert.Equal(t, absURL1, stringutil.CreateAbsoluteURL(relURL1, baseURL))
	assert.Equal(t, absURL2, stringutil.CreateAbsoluteURL(relURL2, baseURL))
	assert.Equal(t, absURL3, stringutil.CreateAbsoluteURL(relURL3, baseURL))
	assert.Equal(t, absURL4, stringutil.CreateAbsoluteURL(relURL4, baseURL))
	assert.Equal(t, absURL5, stringutil.CreateAbsoluteURL(relURL5, baseURL))
	assert.Equal(t, absURL6, stringutil.CreateAbsoluteURL(relURL6, baseURL))
	assert.Equal(t, absURL7, stringutil.CreateAbsoluteURL(relURL7, baseURL))
	assert.Equal(t, absURL8, stringutil.CreateAbsoluteURL(relURL8, baseURL))
	assert.Equal(t, absURL9, stringutil.CreateAbsoluteURL(relURL9, baseURL))
	assert.Equal(t, absURL10, stringutil.CreateAbsoluteURL(relURL10, baseURL2))
}
