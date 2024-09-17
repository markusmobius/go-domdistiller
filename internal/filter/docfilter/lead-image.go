// ORIGINAL: java/webdocument/filters/LeadImageFinder.java

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

package docfilter

import (
	"fmt"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter/scorer"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// In original dom-distiller they specified this value as 26. However since
// we only use two heuristics instead of four I decided to halved it.
const imageMinimumAcceptedScore = 13

// LeadImageFinder is used to identify a lead image for an article (sometimes known as a mast).
// Each candidate image is scored based on several heuristics including:
// - If an ancestor node of the image is the "figure" tag.
// - If the image is close with content.
//
// In original dom-distiller they have two more heuristics which cannot be implemented here
// because it would require us to compute the stylesheet (NEED-COMPUTE-CSS):
// - The ratio of width/height.
// - The area of the image (width * height) relative to its container.
type LeadImageFinder struct {
	logger logutil.Logger
}

func NewLeadImageFinder(logger logutil.Logger) *LeadImageFinder {
	return &LeadImageFinder{
		logger: logger,
	}
}

func (f *LeadImageFinder) Process(doc *webdoc.Document) bool {
	candidates := []webdoc.Element{}
	var firstContent, lastContent *webdoc.Text

	// TODO: WebDocument should have a separate list that point to all images
	// in the document.
	for _, e := range doc.Elements {
		wt, isText := e.(*webdoc.Text)
		if !isText || !e.IsContent() {
			continue
		}

		if firstContent == nil {
			firstContent = wt
		}

		lastContent = wt
	}

	if lastContent == nil {
		return false
	}

	for _, e := range doc.Elements {
		// If the element is an image and not already considered content.
		webImage, isImage := e.(*webdoc.Image)
		webFigure, isFigure := e.(*webdoc.Figure)
		if ((isImage || isFigure) && e.IsContent()) || e == lastContent {
			// If we hit the last content or a image that is
			// content, then we are no longer searching for a
			// "lead" image.
			break
		} else if isImage {
			candidates = append(candidates, webImage)
		} else if isFigure {
			candidates = append(candidates, webFigure)
		}
	}

	return f.findLeadImage(candidates, firstContent)
}

func (f *LeadImageFinder) findLeadImage(candidates []webdoc.Element, firstContent *webdoc.Text) bool {
	if len(candidates) == 0 {
		return false
	}

	var contentElement *html.Node
	if firstContent != nil {
		contentElement = firstContent.FirstNonWhitespaceTextNode()
	}

	bestScore := 0
	heuristics := f.getLeadHeuristics(contentElement)
	var bestImage webdoc.Element

	for _, candidate := range candidates {
		currentScore := f.getImageScore(candidate, heuristics)
		if currentScore > imageMinimumAcceptedScore {
			if bestImage == nil || bestScore < currentScore {
				bestImage = candidate
				bestScore = currentScore
			}
		}
	}

	if bestImage == nil {
		return false
	}

	bestImage.SetIsContent(true)
	return true
}

func (f *LeadImageFinder) getImageScore(e webdoc.Element, heuristics []scorer.ImageScorer) int {
	if e == nil || len(heuristics) == 0 {
		return 0
	}

	var imgNode *html.Node
	webImage, isImage := e.(*webdoc.Image)
	webFigure, isFigure := e.(*webdoc.Figure)
	if !isImage && !isFigure {
		return 0
	} else if isImage {
		imgNode = webImage.Element
	} else {
		imgNode = webFigure.Element
	}

	score := 0
	for _, ir := range heuristics {
		score += ir.GetImageScore(imgNode)
	}

	f.logFinalScore(imgNode, score)
	return score
}

func (f *LeadImageFinder) getLeadHeuristics(firstContent *html.Node) []scorer.ImageScorer {
	return []scorer.ImageScorer{
		// scorer.NewImageAreaScorer(25, 75_000, 200_000),
		// scorer.NewImageRatioScorer(25),
		scorer.NewImageDomDistanceScorer(25, firstContent),
		scorer.NewImageHasFigureScorer(15),
	}
}

func (f *LeadImageFinder) logFinalScore(node *html.Node, score int) {
	if f.logger == nil || f.logger.InternallyNil() {
		return
	}

	logMsg := "null image can't be scored"
	if node != nil {
		src := dom.GetAttribute(node, "src")
		logMsg = fmt.Sprintf("Final image score: %d : %s", score, src)
	}

	f.logger.PrintVisibilityInfo(logMsg)
}
