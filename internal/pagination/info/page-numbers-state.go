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

// PageNumbersState is struct that returned by getPageNumbersState() after it
// has checked if the given list of PageLinkInfo's and PageInfo's are adjacent
// and consecutive, and if there's a gap in the list.
type PageNumbersState struct {
	IsAdjacent    bool
	IsConsecutive bool
	NextPagingURL string
}

func (pns *PageNumbersState) isPageNumberSequence(ascendingNumbers []*PageInfo) bool {
	if len(ascendingNumbers) <= 1 {
		return false
	}

	// The first one must have a URL unless it is the first page.
	firstPage := ascendingNumbers[0]
	if firstPage.PageNumber != 1 && firstPage.URL == "" {
		return false
	}

	// There's only one plain number without URL in ascending numbers group.
	hasPlainNum := false
	for _, page := range ascendingNumbers {
		if page.URL == "" {
			if hasPlainNum {
				return false
			}
			hasPlainNum = true
		} else if hasPlainNum && pns.NextPagingURL == "" {
			pns.NextPagingURL = page.URL
		}
	}

	// If there are only two pages, they must be siblings.
	if len(ascendingNumbers) == 2 {
		return firstPage.PageNumber+1 == ascendingNumbers[1].PageNumber
	}

	// Sometimes there is page where not all its page number is in sequence.
	// For example, ArsTechnica do its numbering like this :
	//   1, 2, 3, 4, 5, 32, 33
	// We only care about 1-5, so here we will look for the longest group of
	// consecutive page number.

	// Create group of consecutive page numbers sequence
	currentStart := 0
	mapSequenceEnd := make(map[int]int)
	for i := 0; i < len(ascendingNumbers)-1; i++ {
		currentPage := ascendingNumbers[i]
		nextPage := ascendingNumbers[i+1]

		if nextPage.PageNumber != currentPage.PageNumber+1 {
			mapSequenceEnd[currentStart] = i + 1
			currentStart = i + 1
		}
	}
	mapSequenceEnd[currentStart] = len(ascendingNumbers)

	// If there are too many group of consecutive numbers,
	// most likely the ascendingNumbers is not page number sequence
	// For example, this is page numbers :
	//   1, 2, 3, 4, 5, 32, 33
	// While this is not :
	//   1, 2, 3, 4, 5, 32, 33, 45, 46
	if len(mapSequenceEnd) > 2 {
		return false
	}

	// Find the longest group
	maxSequenceLength := 0
	sequenceStart := 0
	sequenceEnd := 0

	for start, end := range mapSequenceEnd {
		sequenceLength := end - start
		if sequenceLength > maxSequenceLength {
			maxSequenceLength = sequenceLength
			sequenceStart = start
			sequenceEnd = end
		}
	}

	// If the longest sequence is too small, stop
	if maxSequenceLength <= 1 {
		return false
	}

	// Make sure the longest group contains at most one page info without URL
	// (which might indicates the page info is for our current page)
	nEmptyURL := 0
	for _, page := range ascendingNumbers[sequenceStart:sequenceEnd] {
		if page.URL == "" {
			nEmptyURL++
		}
	}

	return nEmptyURL <= 1
}
