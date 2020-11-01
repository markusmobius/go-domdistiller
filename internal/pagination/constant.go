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

import "regexp"

const (
	// If the numeric value of a link's anchor text is greater than this number,
	// we don't think it represents the page number of the link.
	MaxNumForPageParam = 100
)

var (
	rxNumber        = regexp.MustCompile(`\d`)
	rxNumberAtStart = regexp.MustCompile(`^\d+`)

	// Regex for page number finder
	rxLinkNumberCleaner    = regexp.MustCompile(`[()\[\]{}]`)
	rxInvalidParentWrapper = regexp.MustCompile(`(?i)(body)|(html)`)
	rxTerms                = regexp.MustCompile(`(?i)(\S*[\w\x{00C0}-\x{1FFF}\x{2C00}-\x{D7FF}]\S*)`)
	rxSurroundingDigits    = regexp.MustCompile(`(?i)^[\W_]*(\d+)[\W_]*$`)

	// Regex for prev next finder
	rxNextLink       = regexp.MustCompile(`(?i)(next|weiter|continue|>([^\|]|$)|»([^\|]|$))`)
	rxPrevLink       = regexp.MustCompile(`(?i)(prev|early|old|new|<|«)`)
	rxPositive       = regexp.MustCompile(`(?i)article|body|content|entry|hentry|main|page|pagination|post|text|blog|story`)
	rxNegative       = regexp.MustCompile(`(?i)combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|outbrain|promo|related|shoutbox|sidebar|sponsor|shopping|tags|tool|widget`)
	rxExtraneous     = regexp.MustCompile(`(?i)print|archive|comment|discuss|e[\-]?mail|share|reply|all|login|sign|single|as one|article|post|篇`)
	rxPagination     = regexp.MustCompile(`(?i)pag(e|ing|inat)`)
	rxLinkPagination = regexp.MustCompile(`(?i)p(a|g|ag)?(e|ing|ination)?(=|\/)[0-9]{1,2}$`)
	rxFirstLast      = regexp.MustCompile(`(?i)(first|last)`)
)
