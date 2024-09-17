// ORIGINAL: java/document/TextBlock.java

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
	"fmt"
	"sort"
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/label"
	"golang.org/x/net/html"
)

// TextBlock describes a block of text. A block can be an "atomic" text node (i.e., a sequence
// of text that is not interrupted by any HTML markup) or a compound of such atomic elements.
type TextBlock struct {
	TextElements     []*Text
	Text             string
	Labels           map[string]struct{}
	NumWords         int
	NumWordsInAnchor int
	LinkDensity      float64
	TagLevel         int

	isContent bool
}

func NewTextBlock(textElements ...*Text) *TextBlock {
	tb := &TextBlock{TextElements: textElements, TagLevel: -1}

	for _, wt := range textElements {
		tb.Text += wt.Text
		tb.NumWords += wt.NumWords
		tb.NumWordsInAnchor += wt.NumLinkedWords

		for label := range wt.TakeLabels() {
			tb.AddLabels(label)
		}

		if tb.TagLevel == -1 {
			tb.TagLevel = wt.TagLevel
		}
	}

	tb.LinkDensity = tb.calcLinkDensity()
	return tb
}

func (tb *TextBlock) IsContent() bool {
	return tb.isContent
}

// SetIsContent set the value of isContent.
// Returns true if isContent value changed.
func (tb *TextBlock) SetIsContent(isContent bool) bool {
	if isContent == tb.isContent {
		return false
	}

	tb.isContent = isContent
	return true
}

func (tb *TextBlock) MergeNext(other *TextBlock) {
	tb.Text += "\n" + other.Text
	tb.NumWords += other.NumWords
	tb.NumWordsInAnchor += other.NumWordsInAnchor
	tb.LinkDensity = tb.calcLinkDensity()
	tb.isContent = tb.isContent || other.isContent
	tb.TextElements = append(tb.TextElements, other.TextElements...)

	for label := range other.Labels {
		tb.AddLabels(label)
	}

	if other.TagLevel < tb.TagLevel {
		tb.TagLevel = other.TagLevel
	}
}

func (tb *TextBlock) AddLabels(labels ...string) {
	if tb.Labels == nil {
		tb.Labels = make(map[string]struct{})
	}

	for _, label := range labels {
		tb.Labels[label] = struct{}{}
	}
}

func (tb *TextBlock) RemoveLabels(labels ...string) {
	if tb.Labels == nil {
		tb.Labels = make(map[string]struct{})
	}

	for _, label := range labels {
		delete(tb.Labels, label)
	}
}

func (tb *TextBlock) HasLabel(label string) bool {
	_, exist := tb.Labels[label]
	return exist
}

func (tb *TextBlock) OffsetBlocksStart() int {
	firstText := tb.firstText()
	if firstText == nil {
		return -1
	}

	return firstText.OffsetBlock
}

func (tb *TextBlock) OffsetBlocksEnd() int {
	lastText := tb.lastText()
	if lastText == nil {
		return -1
	}

	return lastText.OffsetBlock
}

func (tb *TextBlock) FirstNonWhitespaceTextNode() *html.Node {
	return tb.firstText().FirstNonWhitespaceTextNode()
}

func (tb *TextBlock) LastNonWhitespaceTextNode() *html.Node {
	return tb.firstText().LastNonWhitespaceTextNode()
}

func (tb *TextBlock) ApplyToModel() {
	if !tb.isContent {
		return
	}

	for _, wt := range tb.TextElements {
		wt.SetIsContent(true)
		if tb.HasLabel(label.Title) {
			wt.AddLabel(label.Title)
		}
	}
}

func (tb *TextBlock) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(fmt.Sprintf("%d/%d;", tb.OffsetBlocksStart(), tb.OffsetBlocksEnd()))
	sb.WriteString(fmt.Sprintf("tl=%d;", tb.TagLevel))
	sb.WriteString(fmt.Sprintf("nw=%d;", tb.NumWords))
	sb.WriteString(fmt.Sprintf("ld=%.3f;", tb.LinkDensity))
	sb.WriteString("]\t")

	if tb.isContent {
		sb.WriteString("CONTENT,")
	} else {
		sb.WriteString("boilerplate,")
	}

	sb.WriteString(tb.labelsDebugString())
	sb.WriteString("\n")
	sb.WriteString(tb.Text)
	return sb.String()
}

func (tb *TextBlock) calcLinkDensity() float64 {
	if tb.NumWords == 0 {
		return 0
	}

	return float64(tb.NumWordsInAnchor) / float64(tb.NumWords)
}

func (tb TextBlock) firstText() *Text {
	if len(tb.TextElements) == 0 {
		return nil
	}

	return tb.TextElements[0]
}

func (tb TextBlock) lastText() *Text {
	if len(tb.TextElements) == 0 {
		return nil
	}

	return tb.TextElements[len(tb.TextElements)-1]
}

func (tb *TextBlock) labelsDebugString() string {
	var i int
	labels := make([]string, len(tb.Labels))
	for label := range tb.Labels {
		labels[i] = label
		i++
	}

	sort.Strings(labels)
	return strings.Join(labels, ",")
}
