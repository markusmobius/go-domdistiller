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
	nurl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/markusmobius/go-domdistiller/internal/pagination/parser"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

// PageNumberFinder parses the document to collect groups of adjacent plain text numbers and
// outlinks with digital anchor text.
type PageNumberFinder struct {
	wordCounter              stringutil.WordCounter
	adjacentNumberGroups     *info.MonotonicPageInfoGroups
	numForwardLinksProcessed int

	timingInfo *data.TimingInfo
	logger     logutil.Logger
}

func NewPageNumberFinder(wc stringutil.WordCounter, timingInfo *data.TimingInfo, logger logutil.Logger) *PageNumberFinder {
	return &PageNumberFinder{
		wordCounter:          wc,
		adjacentNumberGroups: &info.MonotonicPageInfoGroups{},

		timingInfo: timingInfo,
		logger:     logger,
	}
}

func (pnf *PageNumberFinder) FindPagination(root *html.Node, pageURL *nurl.URL) (pagination data.PaginationInfo) {
	url := *pageURL
	url.Path = strings.TrimSuffix(url.Path, "/")
	url.RawPath = url.Path
	strPageURL := stringutil.UnescapedString(&url)

	paramInfo := pnf.FindOutlink(root, &url)
	if paramInfo.Type != info.PageNumber {
		return
	}

	pagination.PrevPage = ""
	pagination.NextPage = paramInfo.NextPagingURL

	// If next page URL is empty but there are related page info, it means we are in
	// the last page, so the last page info is for previous page.
	nPageInfo := len(paramInfo.AllPageInfo)
	if pagination.NextPage == "" && nPageInfo > 0 {
		for i := nPageInfo - 1; i >= 0; i-- {
			currentInfo := paramInfo.AllPageInfo[i]
			if currentInfo.URL != strPageURL {
				pagination.PrevPage = currentInfo.URL
				break
			}
		}
		return
	}

	// If next page URL is not empty, find it in list of page info.
	// The page info before it will point to previous page.
	if pagination.NextPage != "" {
		nextPageIdx := -1
		for i, pageInfo := range paramInfo.AllPageInfo {
			if pageInfo.URL == pagination.NextPage {
				nextPageIdx = i
				break
			}
		}

		for i := nextPageIdx - 1; i >= 0; i-- {
			currentURL := paramInfo.AllPageInfo[i].URL
			if currentURL == "" || currentURL != strPageURL {
				pagination.PrevPage = currentURL
				break
			}
		}
	}

	return
}

// FindOutlink parses the document to collect outlinks with numeric anchor text and numeric text
// around them. Returns PageParamInfo, always (never null). If no page parameter is detected or
// determined to be best, its Type is info.Unset.
func (pnf *PageNumberFinder) FindOutlink(root *html.Node, pageURL *nurl.URL) *info.PageParamInfo {
	start := time.Now()

	idx := 0
	allLinks := dom.GetElementsByTagName(root, "a")
	for idx < len(allLinks) {
		link := allLinks[idx]
		pageInfo, _ := pnf.getPageInfoAndText(link, pageURL)
		if pageInfo == nil {
			idx++
			continue
		}

		// This link is a good candidate for pagination.

		// Close current group of adjacent numbers, add a new group if necessary.
		pnf.adjacentNumberGroups.AddGroup()

		// Before we append the link to the new group of adjacent numbers, check if it's
		// preceded by a text node with numeric text; if so, add it before the link.
		pnf.findAndAddClosestValidLeafNodes(link, false, true, pageURL)

		// Add the link to the current group of adjacent numbers.
		pnf.adjacentNumberGroups.AddPageInfo(pageInfo)

		// Add all following text nodes and links with numeric text.
		pnf.numForwardLinksProcessed = 0
		pnf.findAndAddClosestValidLeafNodes(link, false, false, pageURL)

		// Skip the current link and links already processed in the forward
		// findandAddClosestValidLeafNodes().
		idx += 1 + pnf.numForwardLinksProcessed
	}

	pnf.adjacentNumberGroups.CleanUp()
	pnf.timingInfo.AddEntry(start, "PageParameterParser")

	start = time.Now()
	paramInfo := parser.DetectParamInfo(pnf.adjacentNumberGroups, pageURL.String(), pnf.logger)
	pnf.timingInfo.AddEntry(start, "PageParameterDetector")

	return paramInfo
}

// getPageInfoAndText returns a populated PageInfoAndText if given link is to be added to
// adjacentNumbersGroups. Otherwise, returns null if link is to be ignored. "javascript:"
// links with numeric text are considered valid links to be added.
func (pnf *PageNumberFinder) getPageInfoAndText(link *html.Node, pageURL *nurl.URL) (*info.PageInfo, string) {
	// In original dom-distiller they ignore the invisible link. Unfortunately it's
	// impossible to do that here. NEED-COMPUTE-CSS.

	// Get visible text using innerText instead of textContent
	linkText := strings.TrimSpace(domutil.InnerText(link))
	number, err := pnf.linkTextToNumber(linkText)
	if err != nil || number < 0 || number > MaxNumForPageParam {
		return nil, ""
	}

	linkHref := dom.GetAttribute(link, "href")
	linkHref = stringutil.CreateAbsoluteURL(linkHref, pageURL)

	isEmptyHref := linkHref == ""
	isJavascriptLink := strings.HasPrefix(linkHref, "javascript:")

	var hrefURL *nurl.URL
	if !isEmptyHref && !isJavascriptLink {
		hrefURL, err = nurl.ParseRequestURI(linkHref)
		if err != nil || hrefURL.Host != pageURL.Host {
			return nil, ""
		}

		hrefURL, err = nurl.Parse(linkHref)
		if err != nil {
			return nil, ""
		}

		hrefURL.Path = strings.TrimSuffix(hrefURL.Path, "/")
		hrefURL.RawPath = hrefURL.Path
		hrefURL.Fragment = ""
		hrefURL.RawFragment = ""
	}

	if isEmptyHref || isJavascriptLink {
		return &info.PageInfo{
			PageNumber: number,
			URL:        linkHref,
		}, linkText
	}

	return &info.PageInfo{
		PageNumber: number,
		URL:        hrefURL.String(),
	}, linkText
}

// findAndAddClosestValidLeafNodes finds and adds the leaf node(s) closest to the given start node.
// This recurses and keeps finding and, if necessary, adding the numeric text of valid nodes, collecting
// the PageInfo for the current adjacency group. For backward search, i.e. nodes before start node,
// search terminates (i.e. recursion stops) once a text node or anchor is encountered. If the text
// node contains numeric text, it's added to the current adjacency group. Otherwise, a new group is
// created to break the adjacency.
// For forward search, i.e. nodes after start node, search continues (i.e. recursion continues) until
// a text node or anchor with non-numeric text is encountered. In the process, text nodes and anchors
// with numeric text are added to the current adjacency group. When a non-numeric text node or anchor
// is encountered, a new group is started to break the adjacency, and search ends.
//
// Returns true to continue search, false to stop.
func (pnf *PageNumberFinder) findAndAddClosestValidLeafNodes(start *html.Node, checkStart, backward bool, pageURL *nurl.URL) bool {
	var node *html.Node
	if checkStart {
		node = start
	} else {
		if backward {
			node = start.PrevSibling
		} else {
			node = start.NextSibling
		}
	}

	if node == nil {
		node = start.Parent
		if node == nil || rxInvalidParentWrapper.MatchString(domutil.NodeName(node)) {
			return false
		}
		return pnf.findAndAddClosestValidLeafNodes(node, false, backward, pageURL)
	}

	checkStart = false
	switch node.Type {
	case html.TextNode:
		// Text must contain words.
		text := node.Data
		if text == "" || pnf.wordCounter.Count(text) == 0 {
			break
		}

		added := pnf.addNonLinkTextIfValid(text)

		// For backward search, we're done regardless if text was added.
		// For forward search, we're done only if text was invalid, otherwise continue.
		if backward || !added {
			return false
		}

	case html.ElementNode:
		if dom.TagName(node) == "a" {
			// For backward search, we're done because we've already processed the anchor.
			if backward {
				return false
			}

			// For forward search, we're done only if link was invalid, otherwise continue.
			pnf.numForwardLinksProcessed++
			added := pnf.addLinkIfValid(node, pageURL)
			if !added {
				return false
			}
			break
		}

		// Intentionally fallthrough
		fallthrough

	default:
		// Check children nodes.
		if len(dom.ChildNodes(node)) == 0 {
			break
		}

		checkStart = true // We want to check the child node.
		if backward {
			// Start the backward search with the rightmost child i.e. last and closest to
			// given node.
			node = node.LastChild
		} else {
			// Start the forward search with the leftmost child i.e. first and closest to
			// given node.
			node = node.FirstChild
		}
	}

	return pnf.findAndAddClosestValidLeafNodes(node, checkStart, backward, pageURL)
}

// addNonLinkTextIfValid handles the text for a non-link node. Each numeric term in the text
// that is a valid plain page number adds a PageParamInfo.PageInfo into the current adjacent
// group. All other terms break the adjacency in the current group, adding a new group instead.
//
// Returns true if text was added to current group of adjacent numbers. Otherwise, false with
// a new group created to break the current adjacency.
func (pnf *PageNumberFinder) addNonLinkTextIfValid(text string) bool {
	// If the text does not contain valid number(s); if necessary, current group of adjacent
	// numbers should be closed, adding a new group if possible.
	if !containsNumber(text) {
		pnf.adjacentNumberGroups.AddGroup()
		return false
	}

	// Extract terms from the text, differentiating between those that contain only digits and
	// those that contain non-digits.
	added := false
	for _, term := range rxTerms.FindAllString(text, -1) {
		number := -1
		termWithDigits := rxSurroundingDigits.FindStringSubmatch(term)
		if len(termWithDigits) > 1 {
			number, _ = strconv.Atoi(termWithDigits[1])
		}

		if number >= 0 && number <= MaxNumForPageParam {
			// This text is a valid candidate of plain text page number, add it to last group of
			// adjacent numbers.
			pnf.adjacentNumberGroups.AddNumber(number, "")
			added = true
		} else {
			// The text is not a valid number, so current group of adjacent numbers should be
			// closed, adding a new group if possible.
			pnf.adjacentNumberGroups.AddGroup()
		}
	}

	return added
}

// addLinkIfValid adds PageInfo to the current adjacent group for a link if its text is numeric.
// Otherwise, add a new group to break the adjacency.
//
// Returns true if link was added, false otherwise.
func (pnf *PageNumberFinder) addLinkIfValid(link *html.Node, pageURL *nurl.URL) bool {
	pageInfo, _ := pnf.getPageInfoAndText(link, pageURL)
	if pageInfo != nil {
		pnf.adjacentNumberGroups.AddPageInfo(pageInfo)
		return true
	}

	pnf.adjacentNumberGroups.AddGroup()
	return false
}

func (pnf *PageNumberFinder) linkTextToNumber(linkText string) (int, error) {
	linkText = rxLinkNumberCleaner.ReplaceAllString(linkText, "")
	linkText = strings.TrimSpace(linkText)
	return strconv.Atoi(linkText)
}
