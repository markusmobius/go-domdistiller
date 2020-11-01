// ORIGINAL: java/webdocument/WebDocumentBuilder.java and
//           java/webdocument/WebDocumentBuilderInterface.java

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
	nurl "net/url"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type DocumentBuilder interface {
	SkipNode(e *html.Node)
	StartNode(e *html.Node)
	EndNode()
	AddTextNode(textNode *html.Node)
	AddLineBreak(node *html.Node)
	AddDataTable(e *html.Node)
	AddTag(tag *Tag)
	AddEmbed(embed Element)
}

type WebDocumentBuilder struct {
	tagLevel      int
	nextTextIndex int
	groupNumber   int
	flush         bool
	document      *Document
	textBuilder   *TextBuilder
	actionStack   []ElementAction
	pageURL       *nurl.URL
}

func NewWebDocumentBuilder(wc stringutil.WordCounter, pageURL *nurl.URL) *WebDocumentBuilder {
	return &WebDocumentBuilder{
		document:    &Document{},
		textBuilder: NewTextBuilder(wc),
		pageURL:     pageURL,
	}
}

func (db *WebDocumentBuilder) SkipNode(e *html.Node) {
	db.flush = true
}

func (db *WebDocumentBuilder) StartNode(e *html.Node) {
	action := GetActionForElement(e)
	db.actionStack = append(db.actionStack, action)

	if action.ChangesTagLevel {
		db.tagLevel++
	}

	if action.IsAnchor {
		db.textBuilder.EnterAnchor()
	}

	db.flush = db.flush || action.Flush
}

func (db *WebDocumentBuilder) EndNode() {
	nActions := len(db.actionStack)
	if nActions == 0 {
		return
	}

	lastAction := db.actionStack[nActions-1]

	if lastAction.ChangesTagLevel {
		db.tagLevel--
	}

	if db.flush || lastAction.Flush {
		db.flushBlock(db.groupNumber)
		db.groupNumber++
	}

	if lastAction.IsAnchor {
		db.textBuilder.ExitAnchor()
	}

	db.actionStack = db.actionStack[:nActions-1]
}

func (db *WebDocumentBuilder) AddTextNode(textNode *html.Node) {
	if db.flush {
		db.flushBlock(db.groupNumber)
		db.groupNumber++
		db.flush = false
	}

	db.textBuilder.AddTextNode(textNode, db.tagLevel)
}

func (db *WebDocumentBuilder) AddLineBreak(br *html.Node) {
	if db.flush {
		db.flushBlock(db.groupNumber)
		db.groupNumber++
		db.flush = false
	}

	db.textBuilder.AddLineBreak(br)
}

func (db *WebDocumentBuilder) AddDataTable(table *html.Node) {
	db.flushBlock(db.groupNumber)
	db.document.AddElements(&Table{
		Element: table,
		PageURL: db.pageURL,
	})
}

func (db *WebDocumentBuilder) AddTag(tag *Tag) {
	db.flushBlock(db.groupNumber)
	db.document.AddElements(tag)
}

func (db *WebDocumentBuilder) AddEmbed(embed Element) {
	db.flushBlock(db.groupNumber)
	db.document.AddElements(embed)
}

func (db *WebDocumentBuilder) Build() *Document {
	db.flushBlock(db.groupNumber)
	return db.document
}

func (db *WebDocumentBuilder) flushBlock(group int) {
	if text := db.textBuilder.Build(db.nextTextIndex); text != nil {
		text.GroupNumber = group
		db.nextTextIndex++
		db.addText(*text)
	}
}

func (db *WebDocumentBuilder) addText(text Text) {
	for _, action := range db.actionStack {
		for _, label := range action.Labels {
			text.AddLabel(label)
		}
	}

	db.document.AddElements(&text)
}
