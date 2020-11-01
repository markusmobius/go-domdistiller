// ORIGINAL: javatest/webdocument/WebDocumentBuilderTest.java

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

package webdoc_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

const (
	WdbText1 = "Some really long text which should be content."
	WdbText2 = "Another really long text thing which should be content."
	WdbText3 = "And again a third long text for testing."
)

func Test_WebDoc_WebDocumentBuilder_SpansAsInline(t *testing.T) {
	// <span>
	//   TEXT1
	//   <span>
	//     TEXT2
	//   </span>
	//   TEXT3
	// </span>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)

	outerSpan := wdbAddElement(body, "span")
	docBuilder.StartNode(outerSpan)

	text := wdbAddTextNode(outerSpan, WdbText1)
	docBuilder.AddTextNode(text)

	innerSpan := wdbAddElement(outerSpan, "span")
	docBuilder.StartNode(innerSpan)

	text = wdbAddTextNode(innerSpan, WdbText2)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of inner span

	text = wdbAddTextNode(outerSpan, WdbText3)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of outer span
	docBuilder.EndNode() // end of body

	wdbAssertInline(t, docBuilder)
}

func Test_WebDoc_WebDocumentBuilder_DivsAsInline(t *testing.T) {
	// <span>
	//   TEXT1
	//   <div style="display: inline;">
	//     TEXT2
	//   </div>
	//   TEXT3
	// </span>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)

	span := wdbAddElement(body, "span")
	docBuilder.StartNode(span)

	text := wdbAddTextNode(span, WdbText1)
	docBuilder.AddTextNode(text)

	div := wdbAddElement(span, "div")
	dom.SetAttribute(div, "style", "display: inline;")
	docBuilder.StartNode(div)

	text = wdbAddTextNode(div, WdbText2)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of div

	text = wdbAddTextNode(span, WdbText3)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of span
	docBuilder.EndNode() // end of body

	wdbAssertInline(t, docBuilder)
}

func Test_WebDoc_WebDocumentBuilder_DivsAsBlock(t *testing.T) {
	// <div>
	//   TEXT1
	//   <div>
	//     TEXT2
	//   </div>
	//   TEXT3
	// </div>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)

	outerDiv := wdbAddElement(body, "div")
	docBuilder.StartNode(outerDiv)

	text := wdbAddTextNode(outerDiv, WdbText1)
	docBuilder.AddTextNode(text)

	innerDiv := wdbAddElement(outerDiv, "div")
	docBuilder.StartNode(innerDiv)

	text = wdbAddTextNode(innerDiv, WdbText2)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of inner div

	text = wdbAddTextNode(outerDiv, WdbText3)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of outer div
	docBuilder.EndNode() // end of body

	wdbAssertBlock(t, docBuilder)
}

func Test_WebDoc_WebDocumentBuilder_SpansAsBlock(t *testing.T) {
	// <div>
	//   TEXT1
	//   <span style="display: block;">
	//     TEXT2
	//   </span>
	//   TEXT3
	// </div>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)

	div := wdbAddElement(body, "div")
	docBuilder.StartNode(div)

	text := wdbAddTextNode(div, WdbText1)
	docBuilder.AddTextNode(text)

	span := wdbAddElement(div, "span")
	dom.SetAttribute(span, "style", "display: block;")
	docBuilder.StartNode(span)

	text = wdbAddTextNode(span, WdbText2)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of span

	text = wdbAddTextNode(div, WdbText3)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of div
	docBuilder.EndNode() // end of body

	wdbAssertBlock(t, docBuilder)
}

func Test_WebDoc_WebDocumentBuilder_HeadingsAsBlock(t *testing.T) {
	// <div>
	//   TEXT1
	//   <h1>
	//     TEXT2
	//   </h1>
	//   TEXT3
	// </div>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)

	div := wdbAddElement(body, "div")
	docBuilder.StartNode(div)

	text := wdbAddTextNode(div, WdbText1)
	docBuilder.AddTextNode(text)

	h1 := wdbAddElement(div, "h1")
	docBuilder.StartNode(h1)

	text = wdbAddTextNode(h1, WdbText2)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of h1

	text = wdbAddTextNode(div, WdbText3)
	docBuilder.AddTextNode(text)
	docBuilder.EndNode() // end of div
	docBuilder.EndNode() // end of body

	wdbAssertBlock(t, docBuilder)
}

func Test_WebDoc_WebDocumentBuilder_KeepsWhitespaceWithinTextBlock(t *testing.T) {
	//
	// <div>
	//   TEXT1
	//
	//   <span>
	//     TEXT2
	//   </span>
	//   TEXT3
	// </div>
	//
	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)
	docBuilder.AddTextNode(wdbAddTextNode(body, "\n")) // will be ignored

	outerDiv := wdbAddElement(body, "div")
	docBuilder.StartNode(outerDiv)
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, WdbText1))
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))

	span := wdbAddElement(outerDiv, "span")
	docBuilder.StartNode(span)
	docBuilder.AddTextNode(wdbAddTextNode(span, "\n"))
	docBuilder.AddTextNode(wdbAddTextNode(span, WdbText2))
	docBuilder.AddTextNode(wdbAddTextNode(span, "\n"))
	docBuilder.EndNode() // end of span

	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, WdbText3))
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))
	docBuilder.EndNode() // end of outer div

	docBuilder.AddTextNode(wdbAddTextNode(body, "\n")) // will be ignored
	docBuilder.EndNode()                               // end of body

	textBlocks := wdbGetBuilderTextBlocks(docBuilder)
	assert.Equal(t, 1, len(textBlocks))
	assert.Equal(t, "\n"+WdbText1+"\n\n"+WdbText2+"\n\n"+WdbText3+"\n", textBlocks[0].Text)
}

func Test_WebDoc_WebDocumentBuilder_NonWordCharacterNotMergedWithNextBlockLevelTextBlock(t *testing.T) {
	//
	// <div>
	//   -
	//   <div>TEXT1</div>
	//   <span>TEXT2</span>
	// </div>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)
	docBuilder.AddTextNode(wdbAddTextNode(body, "\n")) // will be ignored

	outerDiv := wdbAddElement(body, "div")
	docBuilder.StartNode(outerDiv)
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "-"))
	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))

	innerDiv := wdbAddElement(outerDiv, "div")
	docBuilder.StartNode(innerDiv)
	docBuilder.AddTextNode(wdbAddTextNode(innerDiv, WdbText1))
	docBuilder.EndNode() // end of inner div

	docBuilder.AddTextNode(wdbAddTextNode(outerDiv, "\n"))

	span := wdbAddElement(outerDiv, "span")
	docBuilder.StartNode(span)
	docBuilder.AddTextNode(wdbAddTextNode(span, WdbText2))
	docBuilder.EndNode() // end of span

	docBuilder.AddTextNode(wdbAddTextNode(body, "\n")) // will be ignored
	docBuilder.EndNode()                               // end of outer div

	docBuilder.EndNode() // end of body

	textBlocks := wdbGetBuilderTextBlocks(docBuilder)
	assert.Equal(t, 3, len(textBlocks))
	assert.Equal(t, "\n-\n", textBlocks[0].Text)
	assert.Equal(t, WdbText1, textBlocks[1].Text)
	assert.Equal(t, "\n"+WdbText2+"\n", textBlocks[2].Text)
}

func Test_WebDoc_WebDocumentBuilder_EmptyBlockNotMergedWithNextBlock(t *testing.T) {
	// This test simulates many social-bar/leading-link type UIs where
	// lists are used for laying out images.
	// <ul>
	//   <li><a href="foo.html> </a>
	//   </li>
	//   <li>TEXT1
	//   </li>
	// </ul>

	wc := stringutil.FastWordCounter{}
	docBuilder := webdoc.NewWebDocumentBuilder(wc, nil)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	docBuilder.StartNode(body)

	ul := wdbAddElement(body, "ul")
	docBuilder.StartNode(ul)
	docBuilder.AddTextNode(wdbAddTextNode(ul, "\n"))

	li1 := wdbAddElement(ul, "li")
	docBuilder.StartNode(li1)

	anchor := wdbAddElement(li1, "a")
	dom.SetAttribute(anchor, "href", "foo.html")
	docBuilder.StartNode(anchor)
	docBuilder.AddTextNode(wdbAddTextNode(anchor, " "))
	docBuilder.EndNode() // end of anchor

	docBuilder.AddTextNode(wdbAddTextNode(li1, "\n"))
	docBuilder.EndNode() // end of li1

	docBuilder.AddTextNode(wdbAddTextNode(ul, "\n"))

	li2 := wdbAddElement(ul, "li")
	docBuilder.StartNode(li2)
	docBuilder.AddTextNode(wdbAddTextNode(li2, WdbText1))
	docBuilder.AddTextNode(wdbAddTextNode(li2, "\n"))
	docBuilder.EndNode() // end of li2

	docBuilder.EndNode() // end of ul
	docBuilder.EndNode() // end of body

	textBlocks := wdbGetBuilderTextBlocks(docBuilder)
	assert.Equal(t, 1, len(textBlocks))
	assert.Equal(t, WdbText1+"\n", textBlocks[0].Text)
}

func Test_WebDoc_WebDocumentBuilder_Regression0(t *testing.T) {
	html := `<blockquote><p>“There are plenty of instances where provocation comes into` +
		` consideration, instigation comes into consideration, and I will be on the record` +
		` right here on national television and say that I am sick and tired of men` +
		` constantly being vilified and accused of things and we stop there,”` +
		` <a href="http://deadspin.com/i-do-not-believe-women-provoke-violence-says-stephen` +
		`-a-1611060016" target="_blank">Smith said.</a>  “I’m saying, “Can we go a step` +
		` further?” Since we want to dig all deeper into Chad Johnson, can we dig in deep` +
		` to her?”</p></blockquote>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.AppendChild(body, div)

	document := testutil.NewTextDocumentFromPage(div, stringutil.FastWordCounter{}, nil)
	textBlocks := document.TextBlocks
	assert.Len(t, textBlocks, 1)

	tb := textBlocks[0]
	assert.Equal(t, 74, tb.NumWords)
	assert.True(t, 0.1 > tb.LinkDensity)
}

func Test_WebDoc_WebDocumentBuilder_Regression1(t *testing.T) {
	html := "<p>\n" +
		"<a href=\"example\" target=\"_top\"><u>More news</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Search</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Features</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Blogs</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Horse Health</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Ask the Experts</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Horse Breeding</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Forms</u></a> | \n" +
		"<a href=\"example\" target=\"_top\"><u>Home</u></a> </p>\n"

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.AppendChild(body, div)

	document := testutil.NewTextDocumentFromPage(div, stringutil.FastWordCounter{}, nil)
	textBlocks := document.TextBlocks
	assert.Len(t, textBlocks, 1)

	tb := textBlocks[0]
	assert.Equal(t, 14, tb.NumWords)
	assert.Equal(t, 1.0, tb.LinkDensity)
}

func wdbAssertInline(t *testing.T, builder *webdoc.WebDocumentBuilder) {
	textBlocks := wdbGetBuilderTextBlocks(builder)
	assert.Equal(t, 1, len(textBlocks))
	assert.Equal(t, 1, textBlocks[0].TagLevel)
}

func wdbAssertBlock(t *testing.T, builder *webdoc.WebDocumentBuilder) {
	textBlocks := wdbGetBuilderTextBlocks(builder)
	assert.Equal(t, 3, len(textBlocks))
	assert.Equal(t, 2, textBlocks[0].TagLevel)
	assert.Equal(t, 3, textBlocks[1].TagLevel)
	assert.Equal(t, 2, textBlocks[2].TagLevel)
}

func wdbAddElement(parent *html.Node, tagName string) *html.Node {
	e := dom.CreateElement(tagName)
	dom.AppendChild(parent, e)
	return e
}

func wdbAddTextNode(parent *html.Node, data string) *html.Node {
	e := dom.CreateTextNode(data)
	dom.AppendChild(parent, e)
	return e
}

func wdbGetBuilderTextBlocks(builder *webdoc.WebDocumentBuilder) []*webdoc.TextBlock {
	return builder.Build().CreateTextDocument().TextBlocks
}
