// ORIGINAL: java/filters/simple/BoilerplateBlockFilter.java

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

package simple

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// BoilerplateBlock removes TextBlocks which have explicitly been
// marked as "not content".
type BoilerplateBlock struct {
	labelToKeep string
}

func NewBoilerplateBlock(labelToKeep string) *BoilerplateBlock {
	return &BoilerplateBlock{labelToKeep: labelToKeep}
}

func (f *BoilerplateBlock) Process(doc *webdoc.TextDocument) bool {
	hasChanges := false
	textBlocks := doc.TextBlocks

	for i := 0; i < len(textBlocks); i++ {
		tb := textBlocks[i]
		if !tb.IsContent() && (f.labelToKeep != "" || !tb.HasLabel(label.Title)) {
			hasChanges = true

			// These lines is used to remove item from array.
			copy(textBlocks[i:], textBlocks[i+1:])
			textBlocks[len(textBlocks)-1] = nil
			textBlocks = textBlocks[:len(textBlocks)-1]
			i--
		}
	}

	doc.TextBlocks = textBlocks
	return hasChanges
}
