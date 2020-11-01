// ORIGINAL: javatest/TestUtil.java

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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

var (
	rxDirAttributes = regexp.MustCompile(`(?i) dir="(ltr|rtl|inherit|auto)"`)
)

func CreateDivTree() []*html.Node {
	divs := []*html.Node{CreateDiv(0)}
	createDivTree(divs[0], 0, &divs)
	return divs
}

// CreateDiv creates a div with the integer id as its id.
func CreateDiv(id int) *html.Node {
	div := dom.CreateElement("div")
	dom.SetAttribute(div, "id", strconv.Itoa(id))
	return div
}

func CreateTitle(value string) *html.Node {
	title := dom.CreateElement("title")
	dom.SetInnerHTML(title, value)
	return title
}

func CreateHeading(n int, value string) *html.Node {
	h := dom.CreateElement(fmt.Sprintf("h%d", n))
	dom.SetInnerHTML(h, value)
	return h
}

func CreateAnchor(href, text string) *html.Node {
	anchor := dom.CreateElement("a")
	dom.SetAttribute(anchor, "href", href)
	dom.SetInnerHTML(anchor, text)
	return anchor
}

func CreateMetaProperty(property string, content string) *html.Node {
	meta := dom.CreateElement("meta")
	dom.SetAttribute(meta, "property", property)
	dom.SetAttribute(meta, "content", content)
	return meta
}

func CreateMetaName(name string, content string) *html.Node {
	meta := dom.CreateElement("meta")
	dom.SetAttribute(meta, "name", name)
	dom.SetAttribute(meta, "content", content)
	return meta
}

func CreateSpan(text string) *html.Node {
	span := dom.CreateElement("span")
	dom.SetInnerHTML(span, text)
	return span
}

func CreateParagraph(text string) *html.Node {
	p := dom.CreateElement("p")
	dom.SetInnerHTML(p, text)
	return p
}

func CreateListItem(text string) *html.Node {
	li := dom.CreateElement("li")
	dom.SetTextContent(li, text)
	return li
}

func RemoveAllDirAttributes(str string) string {
	return rxDirAttributes.ReplaceAllString(str, "")
}

func createDivTree(e *html.Node, depth int, divs *[]*html.Node) {
	if depth > 2 {
		return
	}

	for i := 0; i < 2; i++ {
		child := CreateDiv(len(*divs))
		*divs = append(*divs, child)
		dom.AppendChild(e, child)
		createDivTree(child, depth+1, divs)
	}
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for testing.
// =================================================================================

// CreateHTML returns an <html> that consist of empty <head> and <body>.
// This is an additional method and doesn't exist in original Java code.
func CreateHTML() *html.Node {
	rawHTML := `
<!DOCTYPE html>
<html lang="en">
<head></head>
<body></body>
</html>`

	root, _ := html.Parse(strings.NewReader(rawHTML))
	return dom.GetElementsByTagName(root, "html")[0]
}

// GetPrettyHTML returns formatted outer HTML of the node.
func GetPrettyHTML(node *html.Node) string {
	return gohtml.Format(dom.OuterHTML(node))
}
