// ORIGINAL: java/filters/heuristics/SimilarSiblingContentExpansion.java

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

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// SimilarSiblingContent marks "siblings" of content as content if they are "similar" enough.
//
// This calculates "siblings" by finding a "canonical" DOM node for each TextBlock. This node is the
// highest ancestor of the TextBlock's first contained node that does not contain (in its subtree)
// the last node of the previous TextBlock or the first node of the next TextBlock.
//
// If a content block and a non-content block are siblings and are "similar" enough, then the non-
// content block is marked as content. The "similarity" test is configurable in various ways.
type SimilarSiblingContent struct {
	AllowCrossTitles   bool
	AllowCrossHeadings bool
	AllowMixedTags     bool
	MaxLinkDensity     float64
	MaxBlockDistance   int
}

func NewSimilarSiblingContentExpansion() *SimilarSiblingContent {
	return &SimilarSiblingContent{}
}

func (f *SimilarSiblingContent) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	canonicalReps := f.findCanonicalReps(textBlocks)

	// After processing a block, it will be added to either the list of good or the list of bad
	// blocks. The good list contains blocks that are content, and bad contains non-content.
	// The range [goodBegin, goodEnd] is a set of blocks that could be potential siblings (and
	// similar for bad).
	bad := make([]int, len(textBlocks))
	good := make([]int, len(textBlocks))
	badBegin, badEnd := 0, 0
	goodBegin, goodEnd := 0, 0

	changes := false
	for i := 0; i < len(textBlocks); i++ {
		if (!f.AllowCrossTitles && textBlocks[i].HasLabel(label.Title)) ||
			(!f.AllowCrossHeadings && textBlocks[i].HasLabel(label.Heading)) {
			// Clear the sets of potential siblings (since expansion is not allowed to cross
			// this block).
			goodBegin = goodEnd
			badBegin = badEnd
			continue
		}

		if f.allowExpandFrom(textBlocks, i) {
			good[goodEnd] = i
			goodEnd++

			// Check the potential bad siblings and set any matches to content.
			for j := badBegin; j < badEnd; j++ {
				b := bad[j]
				if i-b > f.MaxBlockDistance {
					if j == badBegin {
						badBegin++
					}
					continue
				}
				if f.isSimilarIndex(canonicalReps, i, b) {
					changes = true
					textBlocks[b].SetIsContent(true)
					// Remove bad[j] from the "bad" potential set. There is no need to add it to
					// the good set since any further sibling of it will also be a sibling of
					// the current block which has already been added to the list.
					bad[j] = bad[badBegin]
					badBegin++
				}
			}
		} else if f.allowExpandTo(textBlocks, i) {
			j := 0

			// Check the potential good siblings. If there's a match, this block becomes
			// content.
			for j = goodBegin; j < goodEnd; j++ {
				g := good[j]
				if i-g > f.MaxBlockDistance {
					if j == goodBegin {
						goodBegin++
					}
					continue
				}
				if f.isSimilarIndex(canonicalReps, i, g) {
					changes = true
					textBlocks[i].SetIsContent(true)
					// Remove good[j] from the potential set. This can be done because any
					// sibling of it will also be a sibling of the current block, and the
					// current block will be added to the potential set below.
					good[j] = good[goodBegin]
					goodBegin++
					break
				}
			}
			if j == goodEnd {
				bad[badEnd] = i
				badEnd++
			} else {
				good[goodEnd] = i
				goodEnd++
			}
		}
	}

	return changes
}

func (f *SimilarSiblingContent) allowExpandFrom(textBlocks []*webdoc.TextBlock, i int) bool {
	return textBlocks[i].IsContent() &&
		!textBlocks[i].HasLabel(label.StrictlyNotContent) &&
		!textBlocks[i].HasLabel(label.Title)
}

func (f *SimilarSiblingContent) allowExpandTo(textBlocks []*webdoc.TextBlock, i int) bool {
	return textBlocks[i].LinkDensity <= f.MaxLinkDensity &&
		!textBlocks[i].IsContent() &&
		!textBlocks[i].HasLabel(label.StrictlyNotContent) &&
		!textBlocks[i].HasLabel(label.Title)
}

func (f *SimilarSiblingContent) findCanonicalReps(textBlocks []*webdoc.TextBlock) []*html.Node {
	reps := []*html.Node{}

	for i := 0; i < len(textBlocks); i++ {
		var nextNode *html.Node
		if i+1 < len(textBlocks) {
			nextNode = textBlocks[i+1].FirstNonWhitespaceTextNode()
		}

		var prevNode *html.Node
		if i > 0 {
			prevNode = textBlocks[i-1].LastNonWhitespaceTextNode()
		}

		currentNode := textBlocks[i].FirstNonWhitespaceTextNode()

		// Find the highest ancestor of currNode that is not also an ancestor of
		// one of prevNode or nextNode
		currentParent := currentNode.Parent
		for !domutil.Contains(currentParent, prevNode) && !domutil.Contains(currentParent, nextNode) {
			currentNode = currentParent
			currentParent = currentNode.Parent
		}

		reps = append(reps, currentNode)
	}

	return reps
}

func (f *SimilarSiblingContent) isSimilarIndex(canonicalReps []*html.Node, i, j int) bool {
	leftNode, rightNode := canonicalReps[i], canonicalReps[j]
	if !f.AllowMixedTags && !f.areSameTag(leftNode, rightNode) {
		return false
	}

	return leftNode.Parent == rightNode.Parent
}

func (f *SimilarSiblingContent) areSameTag(left, right *html.Node) bool {
	switch {
	case left.Type != right.Type:
		return false

	case left.Type != html.ElementNode:
		return true

	default:
		return domutil.NodeName(left) == domutil.NodeName(right)
	}
}
