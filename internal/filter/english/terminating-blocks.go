// ORIGINAL: java/filters/english/TerminatingBlocksFinder.java

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

package english

import (
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// TerminatingBlocksFinder finds blocks which are potentially indicating the end of
// an article text and marks them with label.StrictlyNotContent.
type TerminatingBlocksFinder struct{}

func NewTerminatingBlocksFinder() *TerminatingBlocksFinder {
	return &TerminatingBlocksFinder{}
}

func (f *TerminatingBlocksFinder) Process(doc *webdoc.TextDocument) bool {
	changes := false

	for _, block := range doc.TextBlocks {
		if f.isTerminating(block) {
			block.AddLabels(label.StrictlyNotContent)
			changes = true
		}
	}

	return changes
}

func (f *TerminatingBlocksFinder) isTerminating(tb *webdoc.TextBlock) bool {
	if tb.NumWords > 14 {
		return false
	}

	text := strings.TrimSpace(tb.Text)
	if stringutil.CharCount(text) >= 8 {
		return re2go.IsTerminatingBlocks(text)
	} else if tb.LinkDensity == 1 {
		return text == "Comment"
	} else if text == "Shares" {
		// Skip social and sharing elements.
		// See crbug.com/692553
		return true
	}

	return false
}
