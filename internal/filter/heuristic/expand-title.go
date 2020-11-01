// ORIGINAL: java/filters/heuristics/ExpandTitleToContentFilter.java

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

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// ExpandTitleToContent marks all TextBlocks "content" which are between the headline and the
// part that has already been marked content, if they are marked with label.MightBeContent.
// This filter is quite specific to the news domain.
type ExpandTitleToContent struct{}

func NewExpandTitleToContent() *ExpandTitleToContent {
	return &ExpandTitleToContent{}
}

func (f *ExpandTitleToContent) Process(doc *webdoc.TextDocument) bool {
	title := -1
	contentStart := -1
	for i, tb := range doc.TextBlocks {
		if contentStart == -1 && tb.HasLabel(label.Title) {
			title = i
			contentStart = -1
		}

		if contentStart == -1 && tb.IsContent() {
			contentStart = i
		}
	}

	if contentStart <= title || title == -1 {
		return false
	}

	changes := false
	for _, tb := range doc.TextBlocks[title:contentStart] {
		if tb.HasLabel(label.MightBeContent) {
			changed := tb.SetIsContent(true)
			changes = changes || changed
		}
	}

	return changes
}
