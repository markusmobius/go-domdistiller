// ORIGINAL: java/filters/english/NumWordsRulesClassifier.java

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

package english

import (
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// NumWordsRulesClassifier classifies several TextBlock as content or not-content through
// rules that have been determined using the C4.8 machine learning algorithm, as described
// in the paper "Boilerplate Detection using Shallow Text Features" (WSDM 2010), particularly
// using number of words per block and link density per block.
type NumWordsRulesClassifier struct{}

func NewNumWordsRulesClassifier() *NumWordsRulesClassifier {
	return &NumWordsRulesClassifier{}
}

func (f *NumWordsRulesClassifier) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) == 0 {
		return false
	}

	hasChanges := false
	for i, block := range textBlocks {
		var prevBlock, nextBlock *webdoc.TextBlock
		if i > 0 {
			prevBlock = textBlocks[i-1]
		}
		if i+1 < len(textBlocks) {
			nextBlock = textBlocks[i+1]
		}

		changed := f.classify(prevBlock, block, nextBlock)
		hasChanges = hasChanges || changed
	}

	return hasChanges
}

func (f *NumWordsRulesClassifier) classify(prev, current, next *webdoc.TextBlock) bool {
	isContent := false

	if current.LinkDensity <= 0.333333 {
		if prev == nil || prev.LinkDensity <= 0.555556 {
			if current.NumWords <= 16 {
				if next == nil || next.NumWords <= 15 {
					isContent = prev != nil && prev.NumWords > 4
				} else {
					isContent = true
				}
			} else {
				isContent = true
			}
		} else {
			if current.NumWords <= 40 {
				isContent = next != nil && next.NumWords > 17
			} else {
				isContent = true
			}
		}
	} else {
		isContent = false
	}

	return current.SetIsContent(isContent)
}
