// ORIGINAL: java/PageParameterParser.java

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

package pagination

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	// If the numeric value of a link's anchor text is greater than this number,
	// we don't think it represents the page number of the link.
	MaxNumForPageParam = 100
)

var (
	// Regex for page number finder. If you are looking for regex for prev next finder,
	// they are compiled to re2go because it's quite slow.
	rxLinkNumberCleaner    = regexp.MustCompile(`[()\[\]{}]`)
	rxInvalidParentWrapper = regexp.MustCompile(`(?i)(body)|(html)`)
	rxTerms                = regexp.MustCompile(`(?i)(\S*[\w\x{00C0}-\x{1FFF}\x{2C00}-\x{D7FF}]\S*)`)
	rxSurroundingDigits    = regexp.MustCompile(`(?i)^[\W_]*(\d+)[\W_]*$`)
)

func containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func getStartingNumber(s string) (int, bool) {
	var b strings.Builder
	for _, r := range s {
		if !unicode.IsDigit(r) {
			break
		}
		b.WriteRune(r)
	}

	str := b.String()
	if str == "" {
		return 0, false
	}

	i, err := strconv.Atoi(b.String())
	return i, err == nil
}
