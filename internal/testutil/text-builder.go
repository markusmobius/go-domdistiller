// ORIGINAL: javatest/webdocument/TestWebTextBuilder.java

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

package testutil

import (
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

type TextBuilder struct {
	wordCounter stringutil.WordCounter
	textNodes   []*html.Node
}

func NewTextBuilder(wc stringutil.WordCounter) *TextBuilder {
	return &TextBuilder{wordCounter: wc}
}

func (tb *TextBuilder) CreateForText(str string) *webdoc.Text {
	return tb.create(str, false)
}

func (tb *TextBuilder) CreateForAnchorText(str string) *webdoc.Text {
	return tb.create(str, true)
}

func (tb *TextBuilder) CreateNestedText(str string, levels int) *webdoc.Text {
	div := dom.CreateElement("div")
	tmp := div

	for i := 0; i < levels-1; i++ {
		dom.AppendChild(tmp, dom.CreateElement("div"))
		tmp = dom.FirstElementChild(tmp)
	}

	dom.AppendChild(tmp, dom.CreateTextNode(str))
	tb.textNodes = append(tb.textNodes, tmp.FirstChild)

	idx := len(tb.textNodes) - 1
	numWords := tb.wordCounter.Count(str)

	return &webdoc.Text{
		Text:           str,
		TextNodes:      tb.textNodes,
		Start:          idx,
		End:            idx + 1,
		FirstWordNode:  idx,
		LastWordNode:   idx,
		NumWords:       numWords,
		NumLinkedWords: 0,
		TagLevel:       0,
		OffsetBlock:    idx,
	}
}

func (tb *TextBuilder) create(str string, isAnchor bool) *webdoc.Text {
	tb.textNodes = append(tb.textNodes, dom.CreateTextNode(str))

	idx := len(tb.textNodes) - 1
	numWords := tb.wordCounter.Count(str)
	numLinkedWords := numWords
	if !isAnchor {
		numLinkedWords = 0
	}

	return &webdoc.Text{
		Text:           str,
		TextNodes:      tb.textNodes,
		Start:          idx,
		End:            idx + 1,
		FirstWordNode:  idx,
		LastWordNode:   idx,
		NumWords:       numWords,
		NumLinkedWords: numLinkedWords,
		TagLevel:       0,
		OffsetBlock:    idx,
	}
}
