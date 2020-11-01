// ORIGINAL: javatest/webdocument/ElementActionTest.java

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
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_WebDoc_Flush, in our case Flush will always true because we cant compute CSS.
// - Test_WebDoc_ChangesTagLevel, in our case ChangeTagLevel will always true.

func Test_WebDoc_ElementAction_IsAnchor(t *testing.T) {
	assert.False(t, actForHtml(`<span></span>`).IsAnchor)
	assert.False(t, actForHtml(`<div></div>`).IsAnchor)
	assert.False(t, actForHtml(`<a></a>`).IsAnchor)
	assert.True(t, actForHtml(`<a href="http://example.com"></a>`).IsAnchor)
}

func Test_WebDoc_ElementAction_Labels(t *testing.T) {
	assert.Len(t, actForHtml("<span></span>").Labels, 0)
	assert.Len(t, actForHtml("<div></div>").Labels, 0)
	assert.Len(t, actForHtml("<p></p>").Labels, 0)
	assert.Len(t, actForHtml("<h1></h1>").Labels, 2)
	assert.Len(t, actForHtml("<h2></h2>").Labels, 2)
	assert.Len(t, actForHtml("<li></li>").Labels, 1)
	assert.Len(t, actForHtml("<nav></nav>").Labels, 1)
	assert.Len(t, actForHtml("<aside></aside>").Labels, 1)

	assert.True(t, actHasLabel(actForHtml("<h1></h1>"), label.H1))
	assert.True(t, actHasLabel(actForHtml("<h1></h1>"), label.Heading))
	assert.True(t, actHasLabel(actForHtml("<h4></h4>"), label.Heading))
	assert.True(t, actHasLabel(actForHtml("<h6></h6>"), label.Heading))
	assert.True(t, actHasLabel(actForHtml("<nav></nav>"), label.StrictlyNotContent))
	assert.True(t, actHasLabel(actForHtml("<aside></aside>"), label.StrictlyNotContent))
}

func Test_WebDoc_ElementAction_CommentsLabels(t *testing.T) {
	assert.False(t, actHasLabel(actForHtml(`<span></span>`), label.StrictlyNotContent))
	assert.False(t, actHasLabel(actForHtml(`<div></div>`), label.StrictlyNotContent))

	doc := testutil.CreateHTML()
	dom.SetAttribute(doc, "class", "comment")
	assert.False(t, actHasLabel(webdoc.GetActionForElement(doc), label.StrictlyNotContent))

	body := dom.QuerySelector(doc, "body")
	dom.SetAttribute(body, "class", "comment")
	assert.False(t, actHasLabel(webdoc.GetActionForElement(body), label.StrictlyNotContent))

	assert.True(t, actHasLabel(actForHtml(`<div class=" comment "></div>`), label.StrictlyNotContent))
	assert.True(t, actHasLabel(actForHtml(`<div class="foo.1 comment-thing"></div>`), label.StrictlyNotContent))
	assert.True(t, actHasLabel(actForHtml(`<div id="comments"></div>`), label.StrictlyNotContent))
	assert.True(t, actHasLabel(actForHtml(`<div class="user-comments"></div>`), label.StrictlyNotContent))
	assert.False(t, actHasLabel(actForHtml(`<article class="user-comments"></div>`), label.StrictlyNotContent))

	// Element.getClassName() returns SVGAnimatedString for SvgElement
	// https://code.google.com/p/google-web-toolkit/issues/detail?id=9195
	assert.False(t, actHasLabel(actForHtml("<svg></svg>"), label.StrictlyNotContent))

	assert.False(t, actHasLabel(actForHtml(
		`<div class="user-comments another-class lots-of-classes too-many-classes`+
			`class1 class2 class3 class4 class5 class6 class7 class8"></div>`),
		label.StrictlyNotContent))

	assert.True(t, actHasLabel(actForHtml(
		`<div class="     user-comments                         a          "></div>`),
		label.StrictlyNotContent))
}

func actForHtml(rawHTML string) webdoc.ElementAction {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, rawHTML)

	children := dom.Children(div)
	return webdoc.GetActionForElement(children[0])
}

func actHasLabel(a webdoc.ElementAction, wantedLabel string) bool {
	for _, label := range a.Labels {
		if label == wantedLabel {
			return true
		}
	}
	return false
}
