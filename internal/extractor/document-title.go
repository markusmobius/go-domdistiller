// ORIGINAL: java/DocumentTitleGetter.java

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
// and licensed under the Apache License, Version 2.0.

package extractor

import (
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

var (
	rxTitleSeparator       = regexp.MustCompile(`(?i) [\|\-\\/>»] `)
	rxTitleHierarchySep    = regexp.MustCompile(`(?i) [\\/>»] `)
	rxTitleRemoveFinalPart = regexp.MustCompile(`(?i)(.*)[\|\-\\/>»] .*`)
	rxTitleRemove1stPart   = regexp.MustCompile(`(?i)[^\|\-\\/>»]*[\|\-\\/>»](.*)`)
	rxTitleAnySeparator    = regexp.MustCompile(`(?i)[\|\-\\/>»]+`)
)

// getDocumentTitle attempt to returns the title for the distilled document, whose functionality
// is migrated from Mozilla's Readability.js. It starts with the document's <title> element and
// extracts parts of the text based on delimiters '|', '-' or ':'. If resulting title is too short
// or long, it uses the document's first H1 element. If the resulting trimmed title is still too
// short, it reverts back to the original title in the document's <title> element.
//
// The implementation of this function is a bit different compared to the one in original
// dom-distiller. Here we use implementation from `go-readability`, the port of Readability.js
// in Go language.
func getDocumentTitle(root *html.Node, wc stringutil.WordCounter) string {
	curTitle := ""
	origTitle := ""
	titleHadHierarchicalSeparators := false

	// If document is empty, there are nothing to do
	if root == nil || wc == nil {
		return ""
	}

	// If they had an element with tag "title" in their HTML
	titleNode := dom.QuerySelector(root, "title")
	if titleNode != nil {
		origTitle = domutil.InnerText(titleNode)
		curTitle = origTitle
	}

	// If there's a separator in the title, first remove the final part
	if rxTitleSeparator.MatchString(curTitle) {
		titleHadHierarchicalSeparators = rxTitleHierarchySep.MatchString(curTitle)
		curTitle = rxTitleRemoveFinalPart.ReplaceAllString(origTitle, "$1")

		// If the resulting title is too short (3 words or fewer), remove
		// the first part instead:
		if wc.Count(curTitle) < 3 {
			curTitle = rxTitleRemove1stPart.ReplaceAllString(origTitle, "$1")
		}
	} else if strings.Contains(curTitle, ": ") {
		// Check if we have an heading containing this exact string, so
		// we could assume it's the full title.
		headings := []*html.Node{}
		headings = append(headings, dom.GetElementsByTagName(root, "h1")...)
		headings = append(headings, dom.GetElementsByTagName(root, "h2")...)

		trimmedTitle := strings.TrimSpace(curTitle)
		match := domutil.SomeNode(headings, func(heading *html.Node) bool {
			return strings.TrimSpace(dom.TextContent(heading)) == trimmedTitle
		})

		// If we don't, let's extract the title out of the original
		// title string.
		if !match {
			curTitle = origTitle[strings.LastIndex(origTitle, ":")+1:]

			// If the title is now too short, try the first colon instead:
			if wc.Count(curTitle) < 3 {
				curTitle = origTitle[strings.Index(origTitle, ":")+1:]
				// But if we have too many words before the colon there's
				// something weird with the titles and the H tags so let's
				// just use the original title instead
			} else if wc.Count(origTitle[:strings.Index(origTitle, ":")]) > 5 {
				curTitle = origTitle
			}
		}
	} else if stringutil.CharCount(curTitle) > 150 || stringutil.CharCount(curTitle) < 15 {
		if h1 := dom.QuerySelector(root, "h1"); h1 != nil {
			curTitle = domutil.InnerText(h1)
		}
	}

	curTitle = strings.TrimSpace(curTitle)
	curTitle = strings.Join(strings.Fields(curTitle), " ")
	// If we now have 4 words or fewer as our title, and either no
	// 'hierarchical' separators (\, /, > or ») were found in the original
	// title or we decreased the number of words by more than 1 word, use
	// the original title.
	curTitleWordCount := wc.Count(curTitle)
	tmpOrigTitle := rxTitleAnySeparator.ReplaceAllString(origTitle, "")

	if curTitleWordCount <= 4 && origTitle != "" &&
		(!titleHadHierarchicalSeparators ||
			curTitleWordCount != wc.Count(tmpOrigTitle)-1) {
		curTitle = origTitle
	}

	return curTitle
}
