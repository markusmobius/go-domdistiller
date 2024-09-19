// ORIGINAL: java/document/TextDocument.java and
//           java/document/TextDocumentStatistics.java

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

package webdoc

import (
	"strings"
)

// TextDocument is a text document, consisting of one or more TextBlock.
type TextDocument struct {
	TextBlocks []*TextBlock
}

func NewTextDocument(textBlocks []*TextBlock) *TextDocument {
	return &TextDocument{textBlocks}
}

func (td *TextDocument) ApplyToModel() {
	for _, tb := range td.TextBlocks {
		tb.ApplyToModel()
	}
}

// CountWordsInContent returns the sum of number of words in content blocks.
func (td *TextDocument) CountWordsInContent() int {
	numWords := 0
	for _, tb := range td.TextBlocks {
		if tb.IsContent() {
			numWords += tb.NumWords
		}
	}
	return numWords
}

// DebugString returns detailed debugging information about the contained TextBlocks.
func (td *TextDocument) DebugString() string {
	var buffer strings.Builder
	for _, tb := range td.TextBlocks {
		buffer.WriteString(tb.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
