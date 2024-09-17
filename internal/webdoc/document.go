// ORIGINAL: java/webdocument/WebDocument.java

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

package webdoc

import (
	"strings"
)

// Document is a simplified view of the underlying webpage. It contains the
// logical elements (blocks of text, image + caption, video, etc).
type Document struct {
	Elements []Element
}

func NewDocument() *Document {
	return &Document{}
}

func (doc *Document) AddElements(elements ...Element) {
	doc.Elements = append(doc.Elements, elements...)
}

func (doc *Document) GenerateOutput(textOnly bool) string {
	var buffer strings.Builder
	for _, e := range doc.Elements {
		if !e.IsContent() {
			continue
		}

		buffer.WriteString(e.GenerateOutput(textOnly))
		if textOnly {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}

// CreateTextDocument generates a web document to be processed by distiller.
// Text groups have been introduced to help retain element order when adding
// images and embeds.
func (doc *Document) CreateTextDocument() *TextDocument {
	textBlocks := []*TextBlock{}
	firstTextIdx := doc.getNextTextIndex(0)
	if firstTextIdx == len(doc.Elements) {
		return NewTextDocument(nil)
	}

	currentText := doc.Elements[firstTextIdx].(*Text)
	currentGroup := currentText.GroupNumber
	previousGroup := currentGroup

	currentBlockTexts := []*Text{}
	for _, element := range doc.Elements {
		text, isText := element.(*Text)
		if !isText {
			continue
		}

		currentGroup = text.GroupNumber
		if currentGroup == previousGroup {
			currentBlockTexts = append(currentBlockTexts, text)
		} else {
			textBlocks = append(textBlocks, NewTextBlock(currentBlockTexts...))
			previousGroup = currentGroup
			currentBlockTexts = []*Text{text}
		}
	}

	textBlocks = append(textBlocks, NewTextBlock(currentBlockTexts...))
	return NewTextDocument(textBlocks)
}

// GetImageURLs returns list of source URLs of all image inside the document.
func (doc *Document) GetImageURLs() []string {
	imageURLs := []string{}
	for _, e := range doc.Elements {
		if !e.IsContent() {
			continue
		}

		// TODO: if we allow images in Text later, handle it here.
		switch element := e.(type) {
		case *Image:
			imageURLs = append(imageURLs, element.GetURLs()...)
		case *Figure:
			imageURLs = append(imageURLs, element.GetURLs()...)
		case *Table:
			imageURLs = append(imageURLs, element.GetImageURLs()...)
		}
	}

	return imageURLs
}

func (doc *Document) getNextTextIndex(startIndex int) int {
	for i := startIndex; i < len(doc.Elements); i++ {
		if _, isText := doc.Elements[i].(*Text); isText {
			return i
		}
	}

	return len(doc.Elements)
}
