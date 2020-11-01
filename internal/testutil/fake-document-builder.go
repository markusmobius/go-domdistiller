// ORIGINAL: javatest/webdocument/FakeWebDocumentBuilder.java

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
	"bytes"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// FakeWebDocumentBuilder is a simple builder that just creates an html-like string
// from the calls. Only used for dom-converter test.
type FakeWebDocumentBuilder struct {
	buffer bytes.Buffer
	nodes  []*html.Node
}

func NewFakeWebDocumentBuilder() *FakeWebDocumentBuilder {
	return &FakeWebDocumentBuilder{}
}

func (db *FakeWebDocumentBuilder) Build() string {
	return db.buffer.String()
}

func (db *FakeWebDocumentBuilder) SkipNode(e *html.Node) {}

func (db *FakeWebDocumentBuilder) StartNode(e *html.Node) {
	db.nodes = append(db.nodes, e)
	db.buffer.WriteString("<")
	db.buffer.WriteString(dom.TagName(e))
	for _, attr := range e.Attr {
		db.buffer.WriteString(" ")
		db.buffer.WriteString(attr.Key)
		db.buffer.WriteString(`="`)
		db.buffer.WriteString(attr.Val)
		db.buffer.WriteString(`"`)
	}
	db.buffer.WriteString(">")
}

func (db *FakeWebDocumentBuilder) EndNode() {
	node := db.nodes[len(db.nodes)-1]
	db.nodes = db.nodes[:len(db.nodes)-1]
	db.buffer.WriteString("</" + dom.TagName(node) + ">")
}

func (db *FakeWebDocumentBuilder) AddTextNode(textNode *html.Node) {
	db.buffer.WriteString(textNode.Data)
}

func (db *FakeWebDocumentBuilder) AddLineBreak(node *html.Node) {
	db.buffer.WriteString("\n")
}

func (db *FakeWebDocumentBuilder) AddDataTable(e *html.Node) {
	db.buffer.WriteString("<datatable/>")
}

func (db *FakeWebDocumentBuilder) AddTag(tag *webdoc.Tag) {}

func (db *FakeWebDocumentBuilder) AddEmbed(embed webdoc.Element) {}
