// ORIGINAL: java/filters/heuristics/HeadingFusion.java

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
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// HeadingFusion fuses headings with the blocks after them. If the heading was
// marked as boilerplate, the fused block will be labeled to prevent
// BlockProximityFusion from merging through it.
type HeadingFusion struct{}

func NewHeadingFusion() *HeadingFusion {
	return &HeadingFusion{}
}

func (f *HeadingFusion) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	changes := false
	currentBlock := textBlocks[0]
	var prevBlock *webdoc.TextBlock

	for i := 1; i < len(textBlocks); i++ {
		prevBlock = currentBlock
		currentBlock = textBlocks[i]

		if !prevBlock.HasLabel(label.Heading) {
			continue
		}

		if prevBlock.HasLabel(label.StrictlyNotContent) || currentBlock.HasLabel(label.StrictlyNotContent) {
			continue
		}

		if prevBlock.HasLabel(label.Title) || currentBlock.HasLabel(label.Title) {
			continue
		}

		if currentBlock.IsContent() {
			changes = true

			headingWasContent := prevBlock.IsContent()
			prevBlock.MergeNext(currentBlock)
			currentBlock = prevBlock

			currentBlock.RemoveLabels(label.Heading)
			if !headingWasContent {
				currentBlock.AddLabels(label.BoilerplateHeadingFused)
			}

			// These lines is used to remove item from array.
			copy(textBlocks[i:], textBlocks[i+1:])
			textBlocks[len(textBlocks)-1] = nil
			textBlocks = textBlocks[:len(textBlocks)-1]
			i--
		} else if prevBlock.IsContent() {
			changes = true
			prevBlock.SetIsContent(false)
		}
	}

	doc.TextBlocks = textBlocks
	return changes
}
