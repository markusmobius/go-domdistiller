// ORIGINAL: java/PageLinkInfo.java and java/PageParamInfo.java

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

import "github.com/markusmobius/go-domdistiller/internal/pagination/pattern"

// PageLinkInfo stores information about the link (anchor) after PageParameterDetector
// has detected the page parameter:
// - the page number (as represented by the original plain text) for the link
// - the original page parameter numeric component in the URL (this component would be replaced
//   by pattern.PageParamPlaceholder in the URL pattern)
// - the position of this link in the list of ascending numbers.
type PageLinkInfo struct {
	PageNumber         int
	PageParamValue     int
	PosInAscendingList int
}

type ListLinkInfo []*PageLinkInfo

func (allLinkInfo ListLinkInfo) Evaluate(pagePattern pattern.PagePattern, ascendingNumbers []*PageInfo, firstPageURL string) *PageParamInfo {
	if len(allLinkInfo) >= minLinksToJustifyLinearMap {
		state := allLinkInfo.PageNumbersState(ascendingNumbers)
		if !state.IsAdjacent {
			return nil
		}

		// Type.PageNumber: ascending numbers must be consecutive and form a page number
		// sequence.
		if !state.IsConsecutive {
			return nil
		}

		if !state.isPageNumberSequence(ascendingNumbers) {
			return nil
		}

		allPageInfo := []*PageInfo{}
		for _, link := range allLinkInfo {
			allPageInfo = append(allPageInfo, &PageInfo{
				PageNumber: link.PageNumber,
				URL:        ascendingNumbers[link.PosInAscendingList].URL,
			})
		}

		linearFormula := allLinkInfo.LinearFormula()
		return &PageParamInfo{
			Type:          PageNumber,
			PagePattern:   pagePattern.String(),
			AllPageInfo:   allPageInfo,
			Formula:       linearFormula,
			NextPagingURL: state.NextPagingURL,
		}
	}

	// Most of news article have no more than 3 pages and the first page probably doesn't have
	// any page parameter. If the first page URL matches the the page pattern, we treat it as
	// the first page of this pattern.
	if len(allLinkInfo) == 1 && firstPageURL != "" {
		onlyLink := allLinkInfo[0]
		secondPageIsOutLink := onlyLink.PageNumber == 2 && onlyLink.PosInAscendingList == 1
		thirdPageIsOutLink := onlyLink.PageNumber == 3 && onlyLink.PosInAscendingList == 2

		// onlyLink's pos is 2 (evaluated right before), so ascendingNumbers has >= 3
		// elements; check if previous element is previous page.
		ascendingNumbers[1].PageNumber = 2

		// 1 PageLinkInfo means ascendingNumbers has >= 1 element.
		if ascendingNumbers[0].PageNumber == 1 && (secondPageIsOutLink || thirdPageIsOutLink) &&
			pagePattern.IsPagingURL(firstPageURL) {
			// Has valid PageParamInfo, create and populate it.
			var coefficient int
			delta := onlyLink.PageParamValue - onlyLink.PageNumber
			if delta == 0 || delta == 1 {
				coefficient = 1
			} else {
				coefficient = onlyLink.PageParamValue
				delta = 0
			}

			allPageInfo := []*PageInfo{}
			allPageInfo = append(allPageInfo,
				&PageInfo{
					PageNumber: 1,
					URL:        firstPageURL,
				},
				&PageInfo{
					PageNumber: onlyLink.PageNumber,
					URL:        ascendingNumbers[onlyLink.PosInAscendingList].URL,
				},
			)

			nextPagingURL := ""
			if thirdPageIsOutLink {
				nextPagingURL = allPageInfo[1].URL
			}

			return &PageParamInfo{
				Type:          PageNumber,
				PagePattern:   pagePattern.String(),
				AllPageInfo:   allPageInfo,
				Formula:       NewLinearFormula(coefficient, delta),
				NextPagingURL: nextPagingURL,
			}
		}
	}

	return nil
}

// LinearFormula determines if the list of PageLinkInfo's form a linear formula:
// pageParamValue = coefficient * pageNum + delta (delta == -coefficient or delta == 0).
//
// The coefficient and delta are calculated from the page parameter values and page numbers of 2
// PageLinkInfo's, and then validated against the remaining PageLinkInfo's.
// The order of page numbers doesn't matter.
//
// Returns LinearFormula, containing the coefficient and delta, if the page
// parameter formula could be determined. Otherwise, returns null.
//
// allLinkInfo is the list of PageLinkInfo's to evaluate
//
// TODO: As this gets rolled out, reassess the necessity of non-1 coefficient support.
func (allLinkInfo ListLinkInfo) LinearFormula() *LinearFormula {
	if len(allLinkInfo) < minLinksToJustifyLinearMap {
		return nil
	}

	firstLink := allLinkInfo[0]
	secondLink := allLinkInfo[1]

	if len(allLinkInfo) == 2 && maxInt(firstLink.PageNumber, secondLink.PageNumber) > 4 {
		return nil
	}

	deltaX := secondLink.PageNumber - firstLink.PageNumber
	if deltaX == 0 {
		return nil
	}

	deltaY := secondLink.PageParamValue - firstLink.PageParamValue
	coefficient := deltaY / deltaX
	if coefficient == 0 {
		return nil
	}

	delta := firstLink.PageParamValue - coefficient*firstLink.PageNumber
	if delta != 0 && delta != -coefficient {
		return nil
	}

	// Check if the remaining elements are on the same linear map.
	for i := 2; i < len(allLinkInfo); i++ {
		link := allLinkInfo[i]
		if link.PageParamValue != coefficient*link.PageNumber+delta {
			return nil
		}
	}

	return NewLinearFormula(coefficient, delta)
}

// PageNumbersState detects if page numbers in list of PageLinkInfo's are adjacent, if page
// numbers in list of PageInfo's are consecutive, and if there's a gap in the list.
//
// For adjacency, the page numbers in list of PageLinkInfo's must either be adjacent, or
// separated by at most 1 plain text number which must represent the current page number in one
// of the PageInfo's.
//
// For consecutiveness, there must be at least one pair of consecutive number values in the list
// of PageLinkInfo's, or between a PageLinkInfo and a plain text number.
//
// Returns a populated PageNumbersState.
//
// allLinkInfo is the list of PageLinkInfo's to evaluate
// ascendingNumbers is list of PageInfo's with ascending PageNum's
func (allLinkInfo ListLinkInfo) PageNumbersState(ascendingNumbers []*PageInfo) *PageNumbersState {
	state := &PageNumbersState{}

	// Check if elements in allLinkInfo are adjacent or there's only 1 gap i.e. the gap is
	// current page number represented in plain text.
	firstPos := -1
	lastPos := -1
	gapPos := -1

	pageParamSet := make(map[int]struct{})
	for _, linkInfo := range allLinkInfo {
		currentPos := linkInfo.PosInAscendingList
		if lastPos == -1 {
			firstPos = currentPos
		} else if currentPos != lastPos+1 {
			// If position is not strictly ascending, or the gap size is > 1 (e.g. "[3] [4] 5 6
			// [7]"), or there's more than 1 gap (e.g. "[3] 4 [5] 6 [7]"), allLinkInfo is not
			// adjacent.
			if currentPos <= lastPos || currentPos != lastPos+2 || gapPos != -1 {
				return state
			}

			gapPos = currentPos - 1
		}

		// Make sure page param value, i.e. page number represented in plain text, is unique.
		if _, exist := pageParamSet[linkInfo.PageParamValue]; exist {
			return state
		}

		pageParamSet[linkInfo.PageParamValue] = struct{}{}
		lastPos = currentPos
	} // for all LinkInfo's

	state.IsAdjacent = true

	// Now, determine if page numbers in ascendingNumbers are consecutive.

	// First, handle the gap.
	if gapPos != -1 {
		if gapPos <= 0 || gapPos >= len(ascendingNumbers)-1 {
			return state
		}

		// Check if its adjacent page numbers are consecutive.
		// e.g. "[1] [5] 6 [7] [12]" is accepted; "[4] 8 [16]" is rejected.
		// This can eliminate links affecting the number of items on a page.
		currentPageNum := ascendingNumbers[gapPos].PageNumber
		if ascendingNumbers[gapPos-1].PageNumber == currentPageNum-1 &&
			ascendingNumbers[gapPos+1].PageNumber == currentPageNum+1 {
			state.IsConsecutive = true
			state.NextPagingURL = ascendingNumbers[gapPos+1].URL
		}

		return state
	}

	// There is no gap.  Check if at least one of the following cases is satisfied:
	// Case #1: "[1] [2] ..." or "1 [2] ... ".
	if (firstPos == 0 || firstPos == 1) && ascendingNumbers[0].PageNumber == 1 &&
		ascendingNumbers[1].PageNumber == 2 {
		state.IsConsecutive = true
		return state
	}

	// Case #2: "[1] 2 [3] ..." where [1] doesn't belong to current pattern.
	if firstPos == 2 && ascendingNumbers[2].PageNumber == 3 &&
		ascendingNumbers[1].URL == "" && ascendingNumbers[0].URL != "" {
		state.IsConsecutive = true
		return state
	}

	// Case #3: "... [n-1] [n]" or "... [n - 1] n".
	numberSize := len(ascendingNumbers)
	if (lastPos == numberSize-1 || lastPos == numberSize-2) &&
		ascendingNumbers[numberSize-2].PageNumber+1 == ascendingNumbers[numberSize-1].PageNumber {
		state.IsConsecutive = true
		return state
	}

	// Case #4: "... [i-1] [i] [i+1] ...".
	for i := firstPos + 1; i < lastPos; i++ {
		if ascendingNumbers[i-1].PageNumber+2 == ascendingNumbers[i+1].PageNumber {
			state.IsConsecutive = true
			return state
		}
	}

	// Otherwise, there's no pair of consecutive values.
	return state
}
