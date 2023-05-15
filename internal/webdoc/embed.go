// ORIGINAL: java/webdocument/WebEmbed.java

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
	"fmt"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

// Embed is the base class for many site-specific embedded
// elements (Twitter, YouTube, etc.).
type Embed struct {
	BaseElement

	Element *html.Node
	ID      string
	Type    string
	Params  map[string]string
}

func (e *Embed) ElementType() string {
	return "embed"
}

func (e *Embed) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	embed := dom.CreateElement("div")
	dom.SetAttribute(embed, "class", "embed-placeholder")
	dom.SetAttribute(embed, "data-type", e.Type)
	dom.SetAttribute(embed, "data-id", e.ID)

	// Radhi:
	// I just realize the embed element never used in original dom-distiller. No wonder
	// Chromium doesn't render any embedded element. To be fair Readability.js doesn't
	// render some embed  as well citing security concerns. In my opinion since dom-
	// distiller usually only used in page that we already visit, the embedded iframe
	// should automatically be trustworthy enough.
	// TODO: Maybe just to be save we should sanitize it.
	tagName := dom.TagName(e.Element)
	if tagName == "blockquote" || tagName == "iframe" {
		domutil.StripAttributes(e.Element)
		dom.AppendChild(embed, e.Element)
	}

	return dom.OuterHTML(embed)
}

func (e *Embed) String() string {
	return fmt.Sprintf("ELEMENT %q: type=%q id=%q, is_content=%v",
		e.ElementType(), e.Type, e.ID, e.isContent)
}
