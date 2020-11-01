// ORIGINAL: java/PageParamInfo.java

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

package info

import "fmt"

// PageParamInfo stores information about the page parameter detected from potential pagination
// URLs with numeric anchor text:
// - type of page parameter detected
// - URL pattern which contains a PageParameterDetector.PageParamPlaceholder to replace the page
//   parameter
// - list of pagination URLs with their page numbers
// - coefficient and delta values of the linear formula formed by the pagination URLs:
//      pageParamValue = coefficient * pageNum + delta
// - next paging URL.
type PageParamInfo struct {
	Type          ParamType
	PagePattern   string
	AllPageInfo   []*PageInfo
	Formula       *LinearFormula
	NextPagingURL string
}

// CompareTo evaluates which ParamInfo is better based on the properties of ParamInfo.
// We prefer the one if the list of PageLinkInfo forms a linear formula, see getLinearFormula().
// Returns 1 if this is better, -1 if other is better and 0 if they are equal.
func (pi *PageParamInfo) CompareTo(other *PageParamInfo) int {
	// We prefer the one where the LinkInfo array forms a linear formula, see isLinearFormula.
	if pi.Formula != nil && other.Formula == nil {
		return 1
	}

	if pi.Formula == nil && other.Formula != nil {
		return -1
	}

	if pi.Type == other.Type {
		return 0
	}

	// For different page param types, we prefer PageNumber.
	if pi.Type == PageNumber {
		return 1
	}

	if other.Type == PageNumber {
		return -1
	}

	// Can't decide as they have unknown page type.
	return 0
}

// CanInsertFirstPage returns true if the given URL, which is the current document URL,
// can be inserted as first page.
//
// Often times, the first page of paginated content does not have a page parameter to identify
// itself, so it is hard for us to cluster the first page into a paginated cluster. However,
// while parsing the first page, we do have some extra signals that can help decide if the
// current page is the first page of its cluster.
//
// docUrl is the current document URL that was parsed.
// ascendingNumbers the list of PageInfo's with ascending PageNumber.
func (pi *PageParamInfo) CanInsertFirstPage(docURL string, ascendingNumbers []*PageInfo) bool {
	// Not enough info to determine whether the URL is fit to be first page.
	if len(pi.AllPageInfo) < 2 {
		return false
	}

	// Already detected first page, no need to add another one.
	if pi.AllPageInfo[0].PageNumber == 1 {
		return false
	}

	// If the current document URL is not shorter than other paginated page URL, it should have
	// a page parameter to identify itself.  This means we could detect it as part of paginated
	// cluster while parsing other members.
	// On the other hand, it is still possible to be the last page of cluster.
	if len(docURL) >= len(pi.AllPageInfo[0].URL) {
		return false
	}

	// The other paginated members must be page 2 through last of paginated content, and the
	// current document isn't detected as any one of them.
	for i := 0; i < len(pi.AllPageInfo); i++ {
		if pi.AllPageInfo[i].PageNumber != i+2 {
			return false
		}

		if pi.AllPageInfo[i].URL == docURL {
			return false
		}
	}

	// If there is a digital outlink with anchor text "1", don't insert this URL as first page,
	// because the first page rarely has an outlink with anchor text "1" pointing to other
	// pages. Normally this is the last page of paginated cluster.
	for _, link := range ascendingNumbers {
		if link.PageNumber == 1 && link.URL != "" && link.URL != docURL {
			return false
		}
	}

	return true
}

// InsertFirstPage inserts the given URL, which is the current document URL, as first page.
// Only call this if canInsertFirstPage() returns true.
//
// docUrl is the current document URL that was parsed.
func (pi *PageParamInfo) InsertFirstPage(docURL string) {
	newInfo := []*PageInfo{{PageNumber: 1, URL: docURL}}
	pi.AllPageInfo = append(newInfo, pi.AllPageInfo...)
}

// determineNextPagingUrl determines the next paging URL for the given document URL.
func (pi *PageParamInfo) DetermineNextPagingURL(docURL string) {
	if pi.NextPagingURL != "" || len(pi.AllPageInfo) == 0 {
		return
	}

	// If document URL is among allPageInfo, the next page is the one after.
	hasDocURL := false
	for _, page := range pi.AllPageInfo {
		if hasDocURL {
			pi.NextPagingURL = page.URL
			return
		}

		if page.URL == docURL {
			hasDocURL = true
		}
	}
}

func (pi *PageParamInfo) String() string {
	str := fmt.Sprintf("Type: %d", pi.Type)
	str += fmt.Sprintf("\nPageInfo: %d", len(pi.AllPageInfo))
	str += fmt.Sprintf("\npattern: %s", pi.PagePattern)

	for _, page := range pi.AllPageInfo {
		str += "\n  " + page.String()
	}

	if pi.Formula == nil {
		str += "\nformula: nil"
	} else {
		str += fmt.Sprintf("\nformula: %s", pi.Formula.String())
	}

	str += fmt.Sprintf("\nnextPagingURL: %s", pi.NextPagingURL)
	return str
}
