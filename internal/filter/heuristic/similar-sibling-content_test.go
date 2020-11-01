// ORIGINAL: javatest/SimilarSiblingContentExpansionTest.java

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

package heuristic_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_Heuristic_SSC_SimpleExpansion(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, "<div>"+
		"<div>text</div>"+
		"<div>text</div>"+
		"</div>")

	wc := stringutil.FastWordCounter{}
	doc := testutil.NewTextDocumentFromPage(body, wc, nil)
	assert.Len(t, doc.TextBlocks, 2)

	doc.TextBlocks[0].SetIsContent(true)
	assert.False(t, doc.TextBlocks[1].IsContent())

	classifier := heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.True(t, doc.TextBlocks[1].IsContent())
}

func Test_Filter_Heuristic_SSC_RequireSameTag(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, "<div>"+
		"<div>text</div>"+
		"<p>text</p>"+
		"</div>")

	wc := stringutil.FastWordCounter{}
	doc := testutil.NewTextDocumentFromPage(body, wc, nil)
	assert.Len(t, doc.TextBlocks, 2)

	doc.TextBlocks[0].SetIsContent(true)
	assert.False(t, doc.TextBlocks[1].IsContent())

	classifier := heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.False(t, doc.TextBlocks[1].IsContent())

	classifier = heuristic.NewSimilarSiblingContentExpansion()
	classifier.AllowMixedTags = true
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.True(t, doc.TextBlocks[1].IsContent())
}

func Test_Filter_Heuristic_SSC_DoNotCrossTitles(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, "<div>"+
		"<div>text</div>"+
		"<p>title</p>"+
		"<div>text</div>"+
		"</div>")

	wc := stringutil.FastWordCounter{}
	doc := testutil.NewTextDocumentFromPage(body, wc, nil)
	assert.Len(t, doc.TextBlocks, 3)

	doc.TextBlocks[1].AddLabels(label.Title)
	doc.TextBlocks[2].SetIsContent(true)
	assert.False(t, doc.TextBlocks[0].IsContent())

	classifier := heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.False(t, doc.TextBlocks[0].IsContent())

	classifier = heuristic.NewSimilarSiblingContentExpansion()
	classifier.AllowCrossTitles = true
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.True(t, doc.TextBlocks[0].IsContent())
}

func Test_Filter_Heuristic_SSC_DoNotCrossHeadings(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, "<div>"+
		"<div>text</div>"+
		"<p>heading</p>"+
		"<div>text</div>"+
		"</div>")

	wc := stringutil.FastWordCounter{}
	doc := testutil.NewTextDocumentFromPage(body, wc, nil)
	assert.Len(t, doc.TextBlocks, 3)

	doc.TextBlocks[1].AddLabels(label.Heading)
	doc.TextBlocks[2].SetIsContent(true)
	assert.False(t, doc.TextBlocks[0].IsContent())

	classifier := heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.False(t, doc.TextBlocks[0].IsContent())

	classifier = heuristic.NewSimilarSiblingContentExpansion()
	classifier.AllowCrossHeadings = true
	classifier.MaxBlockDistance = 3
	classifier.Process(doc)
	assert.True(t, doc.TextBlocks[0].IsContent())
}

func Test_Filter_Heuristic_SSC_ExpansionMaxDistance(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, "<div>"+
		"<div>text</div>"+
		"<p>text</p>"+
		"<div>text</div>"+
		"</div>")

	wc := stringutil.FastWordCounter{}
	doc := testutil.NewTextDocumentFromPage(body, wc, nil)
	assert.Len(t, doc.TextBlocks, 3)

	doc.TextBlocks[0].SetIsContent(true)
	assert.False(t, doc.TextBlocks[1].IsContent())
	assert.False(t, doc.TextBlocks[2].IsContent())

	classifier := heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxBlockDistance = 1
	classifier.Process(doc)
	assert.False(t, doc.TextBlocks[1].IsContent())
	assert.False(t, doc.TextBlocks[2].IsContent())

	classifier = heuristic.NewSimilarSiblingContentExpansion()
	classifier.AllowCrossHeadings = true
	classifier.MaxBlockDistance = 2
	classifier.Process(doc)
	assert.False(t, doc.TextBlocks[1].IsContent())
	assert.True(t, doc.TextBlocks[2].IsContent())
}

func Test_Filter_Heuristic_SSC_MaxLinkDensity(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, "<div>"+
		"<div>text</div>"+
		"<div>text <a href='http://example.com'>link</a></div>"+
		"</div>")

	wc := stringutil.FastWordCounter{}
	doc := testutil.NewTextDocumentFromPage(body, wc, nil)
	assert.Len(t, doc.TextBlocks, 2)

	doc.TextBlocks[0].SetIsContent(true)
	assert.True(t, doc.TextBlocks[1].LinkDensity > 0.4)
	assert.True(t, doc.TextBlocks[1].LinkDensity < 0.6)

	classifier := heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxLinkDensity = 0.4
	classifier.MaxBlockDistance = 1
	classifier.Process(doc)
	assert.False(t, doc.TextBlocks[1].IsContent())

	classifier = heuristic.NewSimilarSiblingContentExpansion()
	classifier.MaxLinkDensity = 0.6
	classifier.MaxBlockDistance = 1
	classifier.Process(doc)
	assert.True(t, doc.TextBlocks[1].IsContent())
}
