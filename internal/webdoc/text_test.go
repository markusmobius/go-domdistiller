// ORIGINAL: javatest/webdocument/WebTextTest.java

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

package webdoc_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_WebDoc_Text_GenerateOutputMultipleContentNodes(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	container := dom.CreateElement("div")
	dom.AppendChild(body, container)

	content1 := dom.CreateElement("p")
	dom.AppendChild(content1, dom.CreateTextNode("Some text content 1."))
	dom.AppendChild(container, content1)

	content2 := dom.CreateElement("p")
	dom.AppendChild(content2, dom.CreateTextNode("Some text content 2."))
	dom.AppendChild(container, content2)

	wc := stringutil.SelectWordCounter(dom.TextContent(doc))
	builder := webdoc.NewTextBuilder(wc)
	builder.AddTextNode(content1.FirstChild, 0)
	builder.AddTextNode(content2.FirstChild, 0)

	text := builder.Build(0)
	got := text.GenerateOutput(false)
	want := "<div><p>Some text content 1.</p><p>Some text content 2.</p></div>"
	assert.Equal(t, want, testutil.RemoveAllDirAttributes(got))
}

func Test_WebDoc_Text_GenerateOutputSingleContentNode(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	container := dom.CreateElement("div")
	dom.AppendChild(body, container)

	content1 := dom.CreateElement("p")
	dom.AppendChild(content1, dom.CreateTextNode("Some text content 1."))
	dom.AppendChild(container, content1)

	content2 := dom.CreateElement("p")
	dom.AppendChild(content2, dom.CreateTextNode("Some text content 2."))
	dom.AppendChild(container, content2)

	wc := stringutil.SelectWordCounter(dom.TextContent(container))
	builder := webdoc.NewTextBuilder(wc)
	builder.AddTextNode(content1.FirstChild, 0)

	text := builder.Build(0)
	got := text.GenerateOutput(false)
	want := "<p>Some text content 1.</p>"
	assert.Equal(t, want, testutil.RemoveAllDirAttributes(got))
}

func Test_WebDoc_Text_GenerateOutputBrElements(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	container := dom.CreateElement("div")
	dom.AppendChild(body, container)

	content1 := dom.CreateElement("p")
	dom.AppendChild(content1, dom.CreateTextNode("Words"))
	dom.AppendChild(content1, dom.CreateElement("br"))
	dom.AppendChild(content1, dom.CreateTextNode("split"))
	dom.AppendChild(content1, dom.CreateElement("br"))
	dom.AppendChild(content1, dom.CreateTextNode("with"))
	dom.AppendChild(content1, dom.CreateElement("br"))
	dom.AppendChild(content1, dom.CreateTextNode("lines"))
	dom.AppendChild(container, content1)

	children := dom.ChildNodes(content1)
	wc := stringutil.SelectWordCounter(dom.TextContent(container))
	builder := webdoc.NewTextBuilder(wc)
	builder.AddTextNode(children[0], 0)
	builder.AddLineBreak(children[1])
	builder.AddTextNode(children[2], 0)
	builder.AddLineBreak(children[3])
	builder.AddTextNode(children[4], 0)
	builder.AddLineBreak(children[5])
	builder.AddTextNode(children[6], 0)

	text := builder.Build(0)
	got := text.GenerateOutput(false)
	want := "<p>Words<br/>split<br/>with<br/>lines</p>"
	assert.Equal(t, want, testutil.RemoveAllDirAttributes(got))

	got = text.GenerateOutput(true)
	want = "Words\nsplit\nwith\nlines"
	assert.Equal(t, want, got)
}

func Test_WebDoc_Text_StripUnsafeAttributes(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	container := dom.CreateElement("div")
	dom.AppendChild(body, container)

	content1 := dom.CreateElement("p")
	dom.SetAttribute(content1, "allowfullscreen", "true") // This should be passed through
	dom.SetAttribute(content1, "onclick", "alert(1)")     // This should be stripped
	dom.AppendChild(content1, dom.CreateTextNode("Text"))
	dom.AppendChild(container, content1)

	wc := stringutil.SelectWordCounter(dom.TextContent(container))
	builder := webdoc.NewTextBuilder(wc)
	builder.AddTextNode(content1.FirstChild, 0)

	text := builder.Build(0)
	got := text.GenerateOutput(false)
	want := `<p allowfullscreen="true">Text</p>`
	assert.Equal(t, want, testutil.RemoveAllDirAttributes(got))
}

func Test_WebDoc_Text_GenerateOutputLiElements(t *testing.T) {
	container := dom.CreateElement("li")
	dom.AppendChild(container, dom.CreateTextNode("Some text content 1."))

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, container)

	wc := stringutil.SelectWordCounter(dom.TextContent(container))
	builder := webdoc.NewTextBuilder(wc)
	builder.AddTextNode(container.FirstChild, 0)

	text := builder.Build(0)
	got := text.GenerateOutput(false)
	want := "Some text content 1."
	assert.Equal(t, want, testutil.RemoveAllDirAttributes(got))
}
