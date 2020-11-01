// ORIGINAL: java/webdocument/WebText.java, java/webdocument/WebTag.java,
//           java/webdocument/WebImage.java

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
	"regexp"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

type TagType uint

const (
	TagStart TagType = iota
	TagEnd
)

var lazyImageAttrs = map[string]string{
	"data-srcset": "srcset",
}

func CanBeNested(tagName string) bool {
	switch tagName {
	case "ul", "ol", "li", "blockquote", "pre":
		return true

	default:
		return false
	}
}

// All inline elements except for impossible tags: br, object, and script.
// Please refer to DomConverter.visitElement() for skipped tags.
// Reference: https://developer.mozilla.org/en-US/docs/HTML/Inline_elements
var inlineTagNames = map[string]struct{}{}

var rxDisplay = regexp.MustCompile(`(?i)display:\s*([\w-]+)\s*(?:;|$)`)

func getDisplayStyle(node *html.Node) string {
	// Check if display specified in inline style
	style := dom.GetAttribute(node, "style")
	parts := rxDisplay.FindStringSubmatch(style)
	if len(parts) >= 2 {
		return parts[1]
	}

	// Use default display
	switch dom.TagName(node) {
	case "address", "article", "blockquote", "body", "dd", "details", "dialog", "div",
		"dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "h1", "h2",
		"h3", "h4", "h5", "h6", "header", "hr", "html", "legend", "main", "nav", "ol",
		"p", "pre", "section", "ul":
		return "block"
	case "a", "abbr", "acronym", "audio", "b", "bdi", "bdo", "br", "canvas", "circle", "cite",
		"code", "data", "defs", "del", "dfn", "ellipse", "em", "embed", "font", "i", "iframe", "img",
		"ins", "kbd", "label", "lineargradient", "mark", "object", "output", "picture", "polygon",
		"q", "rect", "s", "source", "span", "stop", "strong", "sub", "sup", "svg", "tt", "text",
		"time", "track", "u", "var", "video", "wbr":
		return "inline"
	case "button", "input":
		return "inline-block"
	case "li", "summary":
		return "list-item"
	case "ruby":
		return "ruby"
	case "rt":
		return "ruby-text"
	case "table":
		return "table"
	case "caption":
		return "table-caption"
	case "td", "th":
		return "table-cell"
	case "col":
		return "table-column"
	case "colgroup":
		return "table-column-group"
	case "tfoot":
		return "table-footer-group"
	case "thead":
		return "table-header-group"
	case "tr":
		return "table-row"
	case "tbody":
		return "table-row-group"
	default:
		return "none"
	}
}
