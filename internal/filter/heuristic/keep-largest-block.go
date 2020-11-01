// ORIGINAL: java/filters/heuristics/KeepLargestBlockFilter.java

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

// boilerpipe
//
// Copyright (c) 2009 Christian Kohlschütter
//
// The author licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// KeepLargestBlock keeps the largest TextBlock only (by the number of words). In case of
// more than one block with the same number of words, the first block is chosen. All
// discarded blocks are marked "not content" and flagged as `label.MightBeContent`. Note
// that, by default, only TextBlocks marked as "content" are taken into consideration.
type KeepLargestBlock struct {
	expandToSiblings bool
}

func NewKeepLargestBlock(expandToSiblings bool) *KeepLargestBlock {
	return &KeepLargestBlock{
		expandToSiblings: expandToSiblings,
	}
}

func (f *KeepLargestBlock) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	maxNumWords := -1
	largestBlockIndex := -1
	var largestBlock *webdoc.TextBlock

	for i, tb := range textBlocks {
		if tb.IsContent() {
			if tb.NumWords > maxNumWords {
				largestBlock = tb
				maxNumWords = tb.NumWords
				largestBlockIndex = i
			}
		}
	}

	for _, tb := range textBlocks {
		if tb == largestBlock {
			tb.SetIsContent(true)
			tb.AddLabels(label.VeryLikelyContent)
		} else {
			tb.SetIsContent(false)
			tb.AddLabels(label.MightBeContent)
		}
	}

	if f.expandToSiblings && largestBlockIndex != -1 {
		f.maybeExpandContentToEarlierTextBlocks(textBlocks, largestBlock, largestBlockIndex)
		f.maybeExpandContentToLaterTextBlocks(textBlocks, largestBlock, largestBlockIndex)
	}

	return true
}

func (f *KeepLargestBlock) maybeExpandContentToEarlierTextBlocks(textBlocks []*webdoc.TextBlock, largestBlock *webdoc.TextBlock, largestBlockIndex int) {
	firstTextElement := domutil.GetParentElement(largestBlock.FirstNonWhitespaceTextNode())
	for i := largestBlockIndex - 1; i >= 0; i-- {
		candidate := textBlocks[i]
		candidateLastTextElement := domutil.GetParentElement(candidate.LastNonWhitespaceTextNode())
		if f.isSibling(firstTextElement, candidateLastTextElement) {
			candidate.SetIsContent(true)
			candidate.AddLabels(label.SiblingOfMainContent)
			firstTextElement = domutil.GetParentElement(candidate.FirstNonWhitespaceTextNode())
		}
	}
}

func (f *KeepLargestBlock) maybeExpandContentToLaterTextBlocks(textBlocks []*webdoc.TextBlock, largestBlock *webdoc.TextBlock, largestBlockIndex int) {
	lastTextElement := domutil.GetParentElement(largestBlock.LastNonWhitespaceTextNode())
	for i := largestBlockIndex + 1; i < len(textBlocks); i++ {
		candidate := textBlocks[i]
		candidateFirstTextElement := domutil.GetParentElement(candidate.FirstNonWhitespaceTextNode())
		if f.isSibling(lastTextElement, candidateFirstTextElement) {
			candidate.SetIsContent(true)
			candidate.AddLabels(label.SiblingOfMainContent)
			lastTextElement = domutil.GetParentElement(candidate.LastNonWhitespaceTextNode())
		}
	}
}

func (f *KeepLargestBlock) isSibling(e1, e2 *html.Node) bool {
	return domutil.GetParentElement(e1) == domutil.GetParentElement(e2)
}
