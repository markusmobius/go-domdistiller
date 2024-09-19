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

package parser

import (
	nurl "net/url"
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/markusmobius/go-domdistiller/internal/pagination/pattern"
)

type PageCandidate struct {
	pagePattern pattern.PagePattern
	links       []*info.PageLinkInfo
}

// PageCandidatesMap stores a map of URL pattern to its associated list of PageLinkInfo's.
type PageCandidatesMap map[string]PageCandidate

func (pcm PageCandidatesMap) add(pagePattern pattern.PagePattern, link *info.PageLinkInfo) {
	strPattern := pagePattern.String()
	if entry, exist := pcm[strPattern]; exist {
		entry.links = append(entry.links, link)
		pcm[strPattern] = entry
	} else {
		pcm[strPattern] = PageCandidate{
			pagePattern: pagePattern,
			links:       []*info.PageLinkInfo{link},
		}
	}
}

// DetectionState keeps track of the detection state:
// - best PageParamInfo detected so far
// - if multiple page patterns have been found.
type DetectionState struct {
	bestPageParamInfo    *info.PageParamInfo
	hasMultiPagePatterns bool
}

func newDetectionStateFromMonotonicNumbers(monotonicNumbers []*info.PageInfo, isDescending bool,
	parsedDocURL *nurl.URL, acceptedPagePattern string) *DetectionState {
	// Count number of outlinks.
	outlinks := 0
	for _, pageInfo := range monotonicNumbers {
		if pageInfo.URL != "" {
			outlinks++
		}
	}

	if outlinks == 0 {
		return nil
	}

	if isDescending {
		// Reverse the monotonic numbers
		for i, j := 0, len(monotonicNumbers)-1; i < j; i, j = i+1, j-1 {
			monotonicNumbers[i], monotonicNumbers[j] = monotonicNumbers[j], monotonicNumbers[i]
		}
	}

	// Some documents only have two partial pages, where each has a digital outlink to the other.
	// But we need at least 2 URLs to extract page parameters. To handle this case, use current
	// doc URL as the URL for the original plain text number. Note we only do this when the known
	// digital link's page number is 1 or 2.
	if len(monotonicNumbers) == 2 && outlinks == 1 &&
		monotonicNumbers[0].PageNumber == 1 && monotonicNumbers[1].PageNumber == 2 {
		if monotonicNumbers[0].URL == "" {
			monotonicNumbers[0] = &info.PageInfo{PageNumber: 1, URL: parsedDocURL.String()}
		} else {
			monotonicNumbers[1] = &info.PageInfo{PageNumber: 2, URL: parsedDocURL.String()}
		}

		// Increment outlinks to include current document URL.
		outlinks++
	}

	// If there are too little outlinks, just stop
	if outlinks < 2 {
		return nil
	}

	// Now, extract the the page parameter.
	// Eliminate calendar date links.
	possibleDateNum := 0
	for _, page := range monotonicNumbers {
		if page.PageNumber == possibleDateNum+1 {
			possibleDateNum++
		}
	}

	if possibleDateNum >= 28 && possibleDateNum <= 31 {
		return nil
	}

	// Prepare candidates map
	firstPageURL := ""
	pageCandidates := make(PageCandidatesMap)
	parsedURLs := make([]*nurl.URL, len(monotonicNumbers))

	// First, try query components of URLs, looking out for first page URL.
	for i, page := range monotonicNumbers {
		if page.URL == "" {
			continue
		}

		_, err := nurl.ParseRequestURI(page.URL)
		if err != nil {
			parsedURLs[i] = nil
			continue
		}

		url, err := nurl.Parse(page.URL)
		if err != nil {
			parsedURLs[i] = nil
			continue
		}

		url.User = nil
		url.Fragment = ""
		url.RawFragment = ""
		parsedURLs[i] = url

		queryPatterns := pattern.QueryParamPagePatternsFromURL(url)
		for _, queryPattern := range queryPatterns {
			pageCandidates.add(queryPattern, &info.PageLinkInfo{
				PageNumber:         page.PageNumber,
				PageParamValue:     queryPattern.PageNumber(),
				PosInAscendingList: i,
			})
		}

		if page.PageNumber == 1 {
			firstPageURL = page.URL
		}
	}

	// If query components yield nothing, try paths of URLs.
	if len(pageCandidates) == 0 {
		for i, page := range monotonicNumbers {
			url := parsedURLs[i]
			if url == nil {
				continue
			}

			pathPatterns := pattern.PathComponentPagePatternsFromURL(url)
			for _, pathPattern := range pathPatterns {
				pageCandidates.add(pathPattern, &info.PageLinkInfo{
					PageNumber:         page.PageNumber,
					PageParamValue:     pathPattern.PageNumber(),
					PosInAscendingList: i,
				})
			}
		}
	}

	// Determine which URL page pattern is valid with a valid, and the best, PageParamInfo.
	state := &DetectionState{}
	for strPattern, candidate := range pageCandidates {
		if strPattern == acceptedPagePattern || len(candidate.links) > MaxPagingDocs ||
			!candidate.pagePattern.IsValidFor(parsedDocURL) {
			continue
		}

		pageParamInfo := info.ListLinkInfo(candidate.links).
			Evaluate(candidate.pagePattern, monotonicNumbers, firstPageURL)
		if pageParamInfo == nil {
			continue
		}

		// If feasible, insert current document URL as first page.
		// Otherwise, we enhance the heuristic: if current document URL fits the paging pattern
		// of the potential pagination URLs, consider it as first page too.
		docURL := strings.TrimSuffix(parsedDocURL.String(), "/")
		if pageParamInfo.CanInsertFirstPage(docURL, monotonicNumbers) {
			pageParamInfo.InsertFirstPage(docURL)
		} else if candidate.pagePattern.IsPagingURL(docURL) {
			firstPage := pageParamInfo.AllPageInfo[0]
			if firstPage.PageNumber == 2 && firstPage.URL != docURL && len(docURL) < len(firstPage.URL) {
				pageParamInfo.InsertFirstPage(docURL)
			}
		}

		state.compareAndUpdate(&DetectionState{
			bestPageParamInfo: pageParamInfo,
		})
	}

	if state.isEmpty() {
		return nil
	}

	return state
}

func (ds *DetectionState) compareAndUpdate(state *DetectionState) {
	if ds.isEmpty() {
		ds.bestPageParamInfo = state.bestPageParamInfo
		ds.hasMultiPagePatterns = state.hasMultiPagePatterns
		return
	}

	// Compare both PageParamInfo's.
	ret := ds.bestPageParamInfo.CompareTo(state.bestPageParamInfo)
	if ret == -1 { // The formal one is better.
		ds.bestPageParamInfo = state.bestPageParamInfo
		ds.hasMultiPagePatterns = state.hasMultiPagePatterns
	} else if ret == 0 { // Can't decide which one is better.
		ds.hasMultiPagePatterns = true
	}
}

func (ds *DetectionState) isEmpty() bool {
	return ds.bestPageParamInfo == nil
}
