// ORIGINAL: java/filters/heuristics/BlockProximityFusion.java

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
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// BlockProximityFusion fuses adjacent blocks if their distance (in blocks) does not
// exceed a certain limit. This probably makes sense only in cases where an upstream
// filter already has removed some blocks.
type BlockProximityFusion struct {
	postFiltering bool
}

func NewBlockProximityFusion(postFiltering bool) *BlockProximityFusion {
	return &BlockProximityFusion{
		postFiltering: postFiltering,
	}
}

func (f *BlockProximityFusion) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	changes := false
	prevBlock := textBlocks[0]

	for i := 1; i < len(textBlocks); i++ {
		block := textBlocks[i]
		if !block.IsContent() || !prevBlock.IsContent() {
			prevBlock = block
			continue
		}

		diffBlocks := block.OffsetBlocksStart() - prevBlock.OffsetBlocksEnd() - 1
		if diffBlocks <= 1 {
			ok := true
			if f.postFiltering {
				if prevBlock.TagLevel != block.TagLevel {
					ok = false
				}
			} else {
				if block.HasLabel(label.BoilerplateHeadingFused) {
					ok = false
				}
			}

			if prevBlock.HasLabel(label.StrictlyNotContent) != block.HasLabel(label.StrictlyNotContent) {
				ok = false
			}

			if prevBlock.HasLabel(label.Title) != block.HasLabel(label.Title) {
				ok = false
			}

			if (!prevBlock.IsContent() && prevBlock.HasLabel(label.Li)) && !block.HasLabel(label.Li) {
				ok = false
			}

			if ok {
				changes = true
				prevBlock.MergeNext(block)

				// These lines is used to remove item from array.
				copy(textBlocks[i:], textBlocks[i+1:])
				textBlocks[len(textBlocks)-1] = nil
				textBlocks = textBlocks[:len(textBlocks)-1]
				i--
			}
		} else {
			prevBlock = block
		}
	}

	doc.TextBlocks = textBlocks
	return changes
}
