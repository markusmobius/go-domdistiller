// ORIGINAL: javatest/TestTextDocumentBuilder.java

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

package testutil

import (
	"net/url"

	"github.com/markusmobius/go-domdistiller/internal/converter"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

type TextDocumentBuilder struct {
	textBlocks  []*webdoc.TextBlock
	textBuilder *TextBuilder
}

func NewTextDocumentBuilder(wc stringutil.WordCounter) *TextDocumentBuilder {
	return &TextDocumentBuilder{
		textBuilder: NewTextBuilder(wc),
	}
}

func (tdb *TextDocumentBuilder) AddContentBlock(str string, labels ...string) *webdoc.TextBlock {
	tb := tdb.addBlock(str, labels...)
	tb.SetIsContent(true)
	return tb
}

func (tdb *TextDocumentBuilder) AddNonContentBlock(str string, labels ...string) *webdoc.TextBlock {
	tb := tdb.addBlock(str, labels...)
	tb.SetIsContent(false)
	return tb
}

func (tdb *TextDocumentBuilder) Build() *webdoc.TextDocument {
	return webdoc.NewTextDocument(tdb.textBlocks)
}

func (tdb *TextDocumentBuilder) addBlock(str string, labels ...string) *webdoc.TextBlock {
	wt := tdb.textBuilder.CreateForText(str)
	for _, label := range labels {
		wt.AddLabel(label)
	}

	tdb.textBlocks = append(tdb.textBlocks, webdoc.NewTextBlock(wt))
	return tdb.textBlocks[len(tdb.textBlocks)-1]
}

func NewTextDocumentFromPage(doc *html.Node, wc stringutil.WordCounter, pageURL *url.URL) *webdoc.TextDocument {
	builder := webdoc.NewWebDocumentBuilder(wc, pageURL)
	converter.NewDomConverter(builder, pageURL, nil).Convert(doc)
	return builder.Build().CreateTextDocument()
}
