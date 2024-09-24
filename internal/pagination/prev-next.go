// ORIGINAL: java/PagingLinksFinder.java

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

// Parts of this file are adapted from Readability.
//
// Readability is Copyright (c) 2010 Src90 Inc
// and licenced under the Apache License, Version 2.0.

package pagination

import (
	"fmt"
	nurl "net/url"
	"path"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type pagingLinkScore struct {
	linkText  string
	linkHref  string
	linkIndex int
	score     int
}

// PrevNextFinder finds the next and previous page links for the distilled document. The functionality
// for next page links is migrated from readability.getArticleTitle() in chromium codebase's
// third_party/readability/js/readability.js, and then expanded for previous page links; boilerpipe
// doesn't have such capability.
// First, it determines the prefix URL of the document. Then, for each anchor in the document, its
// href and text are compared to the prefix URL and examined for next- or previous-paging-related
// information. If it passes, its score is then determined by applying various heuristics on its
// href, text, class name and ID. Lastly, the page link with the highest score of at least 50 is
// considered to have enough confidence as the next or previous page link.
type PrevNextFinder struct {
	linkDebugInfo     map[*html.Node]string
	linkDebugMessages map[*html.Node]map[string]struct{}
	logger            logutil.Logger
}

func NewPrevNextFinder(logger logutil.Logger) *PrevNextFinder {
	return &PrevNextFinder{
		linkDebugInfo:     make(map[*html.Node]string),
		linkDebugMessages: make(map[*html.Node]map[string]struct{}),
		logger:            logger,
	}
}

func (pnf *PrevNextFinder) FindPagination(root *html.Node, pageURL *nurl.URL) data.PaginationInfo {
	return data.PaginationInfo{
		PrevPage: pnf.FindOutlink(root, pageURL, false),
		NextPage: pnf.FindOutlink(root, pageURL, true),
	}
}

func (pnf *PrevNextFinder) FindOutlink(root *html.Node, pageURL *nurl.URL, findNext bool) string {
	// Clean up URL
	tmp, err := nurl.Parse(pageURL.String())
	if err != nil {
		return ""
	}

	tmp.Fragment = ""
	tmp.RawFragment = ""

	// Remove trailing '/' from window location href, because it'll be used to compare with
	// other href's whose trailing '/' are also removed.
	tmp.Path = strings.TrimSuffix(tmp.Path, "/")
	tmp.RawPath = tmp.Path
	currentURL := stringutil.UnescapedString(tmp)
	pnf.printLog("Current URL:", currentURL)

	// Create folder URL
	tmp.Path = strings.TrimSuffix(path.Dir(tmp.Path), "/")
	tmp.RawPath = tmp.Path
	tmp.RawQuery = ""
	folderURL := stringutil.UnescapedString(tmp)
	pnf.printLog("Folder URL:", folderURL)

	// Create allowed prefix
	// The trailing "/" is essential to ensure the whole hostname is matched, and not just the
	// prefix of the hostname. It also maintains the requirement of having a "path" in the URL.
	tmp.Path = "/"
	tmp.RawPath = tmp.Path
	allowedPrefix := stringutil.UnescapedString(tmp)
	lenPrefix := len(allowedPrefix)
	pnf.printLog("Allowed prefix:", allowedPrefix)

	// Loop through all links, looking for hints that they may be next- or previous- page links.
	// Things like having "page" in their textContent, className or id, or being a child of a
	// node with a page-y className or id.
	// Also possible: levenshtein distance? longest common subsequence?
	// After we do that, assign each page a score.
	bannedURLs := make(map[string]struct{})
	candidates := make([]pagingLinkScore, 0)
	allLinks := dom.GetElementsByTagName(root, "a")
	for i, link := range allLinks {
		// Try to convert relative URL in link href to absolute URL
		linkHref := dom.GetAttribute(link, "href")
		linkHref = stringutil.CreateAbsoluteURL(linkHref, pageURL)

		// Make sure the link href is absolute
		_, err := nurl.ParseRequestURI(linkHref)
		if err != nil {
			pnf.appendDebugStrForLink(link, "ignored: can't converted to abs url")
			continue
		}

		// Make sure the href is related with current page
		if !stringutil.HasPrefixIgnoreCase(linkHref, allowedPrefix) {
			pnf.appendDebugStrForLink(link, "ignored: not prefix")
			continue
		}

		if findNext && !containsNumber(linkHref[lenPrefix:]) {
			pnf.appendDebugStrForLink(link, "ignored: not prefix + number")
			continue
		}

		// In original dom-distiller they skip invisible links, but we can't do that here
		// since Go can't compute stylesheet. NEED-COMPUTE-CSS.

		// Remove url anchor and then trailing '/' from link's href.
		tmp, err := nurl.Parse(linkHref)
		if err != nil {
			pnf.appendDebugStrForLink(link, "ignored: can't be cleaned")
			continue
		}

		tmp.Fragment = ""
		tmp.RawFragment = ""
		tmp.Path = strings.TrimSuffix(tmp.Path, "/")
		tmp.RawPath = tmp.Path
		linkHref = stringutil.UnescapedString(tmp)

		// Ignore page link that is the same as current window location.
		// If the page link is same as the folder URL:
		// - next page link: ignore it, since we would already have seen it.
		// - previous page link: don't ignore it, since some sites will simply have the same
		//                       folder URL for the first page.
		if stringutil.EqualsIgnoreCase(linkHref, currentURL) ||
			(findNext && stringutil.EqualsIgnoreCase(linkHref, folderURL)) {
			pnf.appendDebugStrForLink(link, "ignored: same as current or folder url "+folderURL)
			continue
		}

		// Get link text using inner text
		linkText := domutil.InnerText(link)
		linkText = strings.TrimSpace(linkText)

		// If the linkText looks like it's not the next or previous page, skip it.
		if len(linkText) > 25 {
			pnf.appendDebugStrForLink(link, "ignored: link text too long")
			continue
		}

		// If the linkText contains banned text, skip it, and also ban other anchors with the
		// same link URL.
		if re2go.IsExtraneousToPagination(linkText) {
			pnf.appendDebugStrForLink(link, "ignored: one of extra")
			bannedURLs[linkHref] = struct{}{}
			continue
		}

		// For next page link, if the initial part of the URL is identical to the folder URL, but
		// the rest of it doesn't contain any digits, it's certainly not a next page link.
		// However, this doesn't apply to previous page link, because most sites will just have
		// the folder URL for the first page.
		// TODO(kuan): do we need to apply this heuristic to previous page links if current page
		// number is not 2?
		if findNext {
			remainingLinkHref := linkHref
			if strings.HasPrefix(linkHref, folderURL) {
				remainingLinkHref = linkHref[len(folderURL):]
			}

			if !containsNumber(remainingLinkHref) {
				pnf.appendDebugStrForLink(link, "ignored: no number beyond folder url "+folderURL)
				continue
			}
		}

		// Prepare link score
		linkObj := pagingLinkScore{
			linkIndex: i,
			linkText:  linkText,
			linkHref:  linkHref,
			score:     0,
		}

		// If the folder URL isn't part of this URL, penalize this link.  It could still be the
		// link, but the odds are lower.
		// Example: http://www.actionscript.org/resources/articles/745/1/JavaScript-and-VBScript-Injection-in-ActionScript-3/Page1.html.
		if !strings.HasPrefix(linkHref, folderURL) {
			linkObj.score -= 25

			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, not part of folder url %s",
				linkObj.score, folderURL))
		}

		// Concatenate the link text with class name and id, and determine the score based on
		// existence of various paging-related words.
		linkData := linkText + " " + dom.GetAttribute(link, "class") + " " + dom.GetAttribute(link, "id")

		if (findNext && re2go.IsNextPaginationLink(linkData)) ||
			(!findNext && re2go.IsPrevPaginationLink(linkData)) {
			linkObj.score += 50

			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, has %s",
				linkObj.score, pnf.rxDebugName(findNext)))
		}

		if re2go.IsMightBePagination(linkData) {
			linkObj.score += 25

			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, has pag* word",
				linkObj.score))
		}

		if re2go.ContainsFirstLast(linkData) {
			// -65 is enough to negate any bonuses gotten from a > or » in the text.
			// If we already matched on "next", last is probably fine.
			// If we didn't, then it's bad.  Penalize.
			// Same for "prev".
			if (findNext && !re2go.IsNextPaginationLink(linkText)) ||
				(!findNext && !re2go.IsPrevPaginationLink(linkText)) {
				linkObj.score -= 65

				pnf.appendDebugStrForLink(link, fmt.Sprintf(
					"score %d, has first|last but no %s",
					linkObj.score, pnf.rxDebugName(findNext)))
			}
		}

		if re2go.IsNegativeToBePagination(linkData) || re2go.IsExtraneousToPagination(linkData) {
			linkObj.score -= 50
			pnf.appendDebugStrForLink(link,
				fmt.Sprintf("score %d, has negative or extra regex", linkObj.score))
		}

		if (findNext && re2go.IsPrevPaginationLink(linkData)) ||
			(!findNext && re2go.IsNextPaginationLink(linkData)) {
			linkObj.score -= 200

			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, has opp of %s",
				linkObj.score, pnf.rxDebugName(findNext)))
		}

		// Check if a parent element contains page or paging or paginate.
		positiveMatch := false
		negativeMatch := false
		parent := domutil.GetParentElement(link)
		for parent != nil && (!positiveMatch || !negativeMatch) {
			parentData := dom.GetAttribute(parent, "class") + " " + dom.GetAttribute(parent, "id")
			if !positiveMatch && re2go.IsMightBePagination(parentData) {
				linkObj.score += 25
				positiveMatch = true
				pnf.appendDebugStrForLink(link, fmt.Sprintf(
					"score %d, positive parent - %s",
					linkObj.score, parentData))
			}

			// TODO(kuan): to get 1st page for prev page link, this can't be applied; however,
			// the non-application might be the cause of recursive prev page being returned,
			// i.e. for page 1, it may incorrectly return page 3 for prev page link.
			if !negativeMatch && re2go.IsNegativeToBePagination(parentData) {
				// If this is just something like "footer", give it a negative.
				// If it's something like "body-and-footer", leave it be.
				if !re2go.IsPositiveToBePagination(parentData) {
					linkObj.score -= 25
					negativeMatch = true
					pnf.appendDebugStrForLink(link, fmt.Sprintf(
						"score %d, negative parent - %s",
						linkObj.score, parentData))
				}
			}

			parent = domutil.GetParentElement(parent)
		}

		// If the URL looks like it has paging in it, add to the score.
		// Things like /page/2/, /pagenum/2, ?p=3, ?page=11, ?pagination=34.
		if re2go.IsMightBeLinkPagination(linkHref) || re2go.IsMightBePagination(linkHref) {
			linkObj.score += 25
			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, has paging info", linkObj.score))
		}

		// If the URL contains negative values, give a slight decrease.
		if re2go.IsExtraneousToPagination(linkHref) {
			linkObj.score -= 15
			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, has extra regex", linkObj.score))
		}

		// If the link text is too long, penalize the link.
		if len(linkText) > 10 {
			linkObj.score -= len(linkText)
			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, text too long", linkObj.score))
		}

		// If the link text can be parsed as a number, give it a minor bonus, with a slight bias
		// towards lower numbered pages.  This is so that pages that might not have 'next' in
		// their text can still get scored, and sorted properly by score.
		// TODO(kuan): it might be wrong to assume that it knows about other pages in the
		// document and that it starts on the first page.
		linkTextAsNumber, _ := strconv.Atoi(linkText)
		if linkTextAsNumber > 0 {
			// Punish 1 since we're either already there, or it's probably before what we
			// want anyway.
			if findNext && linkTextAsNumber == 1 {
				linkObj.score -= 10
			} else {
				additionalScore := 10 - linkTextAsNumber
				if additionalScore < 0 {
					additionalScore = 0
				}
				linkObj.score += additionalScore
			}

			pnf.appendDebugStrForLink(link, fmt.Sprintf(
				"score %d, link text is a number (%d)",
				linkObj.score, linkTextAsNumber))
		}

		// Check difference between between page number in link href and current URL.
		// If the difference is exactly 1 (or -1 for previous) increase the score.
		diff, valid := pnf.getPageDiff(currentURL, linkHref, len(allowedPrefix))
		if valid {
			if (findNext && diff == 1) || (!findNext && diff == -1) {
				linkObj.score += 25

				pnf.appendDebugStrForLink(link, fmt.Sprintf(
					"score %d, diff (%d)", linkObj.score, diff))
			}
		}

		// Add final score to candidates
		candidates = append(candidates, linkObj)
	} // loop for all links

	// Loop through all of the possible pages from above and find the top candidate for the next
	// page URL. Require at least a score of 50, which is a relatively high confidence that
	// this page is the next link.
	var topPage *pagingLinkScore
	for i, pageObj := range candidates {
		if _, exist := bannedURLs[pageObj.linkHref]; exist {
			continue
		}

		if pageObj.score >= 50 && (topPage == nil || topPage.score < pageObj.score) {
			topPage = &candidates[i]
		}
	}

	pagingHref := ""
	if topPage != nil {
		pagingHref = topPage.linkHref
		pnf.appendDebugStrForLink(allLinks[topPage.linkIndex], fmt.Sprintf(
			"found: score %d, text=[%s] %s",
			topPage.score, topPage.linkText, topPage.linkHref))
	}

	pnf.printDebugInfo(findNext, pagingHref, allLinks)
	return pagingHref
}

func (pnf *PrevNextFinder) getPageDiff(pageURL, linkHref string, skip int) (int, bool) {
	maxLimit := len(pageURL)
	if lenHref := len(linkHref); lenHref < maxLimit {
		maxLimit = lenHref
	}

	var commonLen int
	for i := skip; i < maxLimit; i++ {
		if pageURL[i] != linkHref[i] {
			commonLen = i
			break
		}
	}

	var urlAsNumber int
	if num, ok := getStartingNumber(pageURL[commonLen:]); ok {
		urlAsNumber = num
	}

	var linkAsNumber int
	if num, ok := getStartingNumber(linkHref[commonLen:]); ok {
		linkAsNumber = num
	}

	if urlAsNumber > 0 && linkAsNumber > 0 {
		return linkAsNumber - urlAsNumber, true
	}

	return 0, false
}

func (pnf *PrevNextFinder) appendDebugStrForLink(link *html.Node, message string) {
	if pnf.logger == nil || pnf.logger.InternallyNil() || !pnf.logger.IsLogPagination() {
		return
	}

	// Check if this message already used for this link
	messageHasBeenUsed := false
	currentMessages, exist := pnf.linkDebugMessages[link]
	if exist && currentMessages != nil {
		_, messageHasBeenUsed = currentMessages[message]
	} else {
		pnf.linkDebugMessages[link] = make(map[string]struct{})
	}

	if messageHasBeenUsed {
		return
	}

	// Combine existing debug message with the new one
	strDebug := ""
	if str, exist := pnf.linkDebugInfo[link]; exist {
		strDebug = str
	}

	if strDebug != "" {
		strDebug += "; "
	}

	strDebug += message
	pnf.linkDebugInfo[link] = strDebug
	pnf.linkDebugMessages[link][message] = struct{}{}
}

func (pnf *PrevNextFinder) printLog(args ...interface{}) {
	if pnf.logger != nil && !pnf.logger.InternallyNil() {
		pnf.logger.PrintPaginationInfo(args...)
	}
}

func (pnf *PrevNextFinder) printDebugInfo(findNext bool, pagingHref string, allLinks []*html.Node) {
	if pnf.logger == nil || pnf.logger.InternallyNil() || !pnf.logger.IsLogPagination() {
		return
	}

	// This logs the following to the console:
	// - number of links processed
	// - the next or previous page link found
	// - for each link: its href, text, concatenated debug string.
	direction := "next"
	if !findNext {
		direction = "prev"
	}

	pnf.logger.PrintPaginationInfo(fmt.Sprintf(
		"nLinks=%d, found %s: %s",
		len(allLinks), direction, pagingHref))

	for i, link := range allLinks {
		text := domutil.InnerText(link)
		text = strings.Join(strings.Fields(text), " ")
		href := dom.GetAttribute(link, "href")
		debugMsg := pnf.linkDebugInfo[link]

		pnf.logger.PrintPaginationInfo(fmt.Sprintf(
			"%d) %s, txt=[%s], dbg=[%s]",
			i, href, text, debugMsg))
	}
}

func (pnf *PrevNextFinder) rxDebugName(findNext bool) string {
	if findNext {
		return "next regex"
	}

	return "prev regex"
}
