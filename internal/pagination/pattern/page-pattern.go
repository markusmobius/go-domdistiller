// ORIGINAL: java/PageParameterDetector.java

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
	nurl "net/url"
)

// PagePattern is the interface that page pattern handlers must implement to detect
// page parameter from potential pagination URLs.
type PagePattern interface {
	// String returns the string of the URL page pattern.
	String() string

	// PageNumber returns the page number extracted from the URL during creation of
	// object that implements this interface.
	PageNumber() int

	// IsValidFor validates this page pattern according to the current document URL
	// through a pipeline of rules. Returns true if page pattern is valid.
	// docUrl is the current document URL.
	IsValidFor(docURL *nurl.URL) bool

	// IsPagingURL returns true if a URL matches this page pattern based on a pipeline of rules.
	// url is the URL to evaluate.
	IsPagingURL(url string) bool
}
