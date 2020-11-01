// ORIGINAL: javatest/webdocument/TestWebDocumentBuilder.java

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
)

// WebDocumentBuilder is a simple builder for testing.
type WebDocumentBuilder struct {
	document    *webdoc.Document
	textBuilder *TextBuilder
}

func NewWebDocumentBuilder() *WebDocumentBuilder {
	return &WebDocumentBuilder{
		document:    webdoc.NewDocument(),
		textBuilder: NewTextBuilder(stringutil.FastWordCounter{}),
	}
}

func (db *WebDocumentBuilder) AddText(text string) *webdoc.Text {
	wt := db.textBuilder.CreateForText(text)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddNestedText(text string) *webdoc.Text {
	wt := db.textBuilder.CreateNestedText(text, 5)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddAnchorText(text string) *webdoc.Text {
	wt := db.textBuilder.CreateForAnchorText(text)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddTable(innerHTML string) *webdoc.Table {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<table>"+innerHTML+"</table>")

	table := dom.QuerySelector(div, "table")
	wt := &webdoc.Table{Element: table}
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddImage() *webdoc.Image {
	image := dom.CreateElement("img")
	dom.SetAttribute(image, "src", "http://www.example.com/foo.jpg")

	wi := &webdoc.Image{
		Element: image,
		Width:   100,
		Height:  100,
	}

	db.document.AddElements(wi)
	return wi
}

func (db *WebDocumentBuilder) AddLeadImage() *webdoc.Image {
	image := dom.CreateElement("img")
	dom.SetAttribute(image, "width", "600")
	dom.SetAttribute(image, "height", "400")
	dom.SetAttribute(image, "src", "http://www.example.com/lead.bmp")

	wi := &webdoc.Image{
		Element: image,
		Width:   100,
		Height:  100,
	}

	db.document.AddElements(wi)
	return wi
}

func (db *WebDocumentBuilder) AddTagStart(tagName string) *webdoc.Tag {
	wt := webdoc.NewTag(tagName, webdoc.TagStart)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddTagEnd(tagName string) *webdoc.Tag {
	wt := webdoc.NewTag(tagName, webdoc.TagEnd)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) Build() *webdoc.Document {
	return db.document
}
