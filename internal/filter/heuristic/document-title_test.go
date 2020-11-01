// ORIGINAL: javatest/DocumentTitleMatchClassifierTest.java

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

	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_Heuristic_DTM_LabelsTitle(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := testutil.NewTextBlockBuilder(wc)

	titleBlock := builder.CreateForText(titleText)
	contentBlock := builder.CreateForText(contentText)
	document := webdoc.NewTextDocument([]*webdoc.TextBlock{
		titleBlock,
		contentBlock,
	})

	classifier := heuristic.NewDocumentTitleMatch(wc, titleText)
	changed := classifier.Process(document)

	assert.True(t, changed)
	assert.True(t, titleBlock.HasLabel(label.Title))
	assert.False(t, contentBlock.HasLabel(label.Title))
}

func Test_Filter_Heuristic_DTM_LabelsMultipleTitle(t *testing.T) {
	// This test mimics leading and trailing breadcrumbs containing the title.
	wc := stringutil.FastWordCounter{}
	builder := testutil.NewTextBlockBuilder(wc)

	titleBlockAsLi := builder.CreateForText(titleText)
	titleBlockAsLi.AddLabels(label.Li)

	titleBlock := builder.CreateForText(titleText)
	contentBlock := builder.CreateForText(contentText)

	trailingTitleBlockAsLi := builder.CreateForText(titleText)
	trailingTitleBlockAsLi.AddLabels(label.Li)

	document := webdoc.NewTextDocument([]*webdoc.TextBlock{
		titleBlockAsLi,
		titleBlock,
		contentBlock,
		trailingTitleBlockAsLi,
	})

	classifier := heuristic.NewDocumentTitleMatch(wc, titleText)
	changed := classifier.Process(document)

	assert.True(t, changed)
	assert.True(t, titleBlockAsLi.HasLabel(label.Title))
	assert.True(t, titleBlock.HasLabel(label.Title))
	assert.False(t, contentBlock.HasLabel(label.Title))
	assert.True(t, trailingTitleBlockAsLi.HasLabel(label.Title))
}

func Test_Filter_Heuristic_DTM_DoesNotLabelTitleInContent(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := testutil.NewTextBlockBuilder(wc)

	titleAndContentBlock := builder.CreateForText(titleText + " " + contentText)
	document := webdoc.NewTextDocument([]*webdoc.TextBlock{titleAndContentBlock})

	classifier := heuristic.NewDocumentTitleMatch(wc, titleText)
	changed := classifier.Process(document)

	assert.False(t, changed)
	assert.False(t, titleAndContentBlock.HasLabel(label.Title))
}

func Test_Filter_Heuristic_DTM_LabelsPartialTitleMatch(t *testing.T) {
	// Non-exhaustive test for the type of partial-matches that Boilerpipe performs.
	wc := stringutil.FastWordCounter{}
	builder := testutil.NewTextBlockBuilder(wc)

	titleBlock := builder.CreateForText(titleText)
	contentBlock := builder.CreateForText(contentText)
	document := webdoc.NewTextDocument([]*webdoc.TextBlock{
		titleBlock,
		contentBlock,
	})

	classifier := heuristic.NewDocumentTitleMatch(wc, "BreakingNews » "+titleText)
	changed := classifier.Process(document)

	assert.True(t, changed)
	assert.True(t, titleBlock.HasLabel(label.Title))
	assert.False(t, contentBlock.HasLabel(label.Title))
}

func Test_Filter_Heuristic_DTM_MatchesMultipleTitles(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := testutil.NewTextBlockBuilder(wc)
	secondTitleText := "I am another document title"

	titleBlock1 := builder.CreateForText(titleText)
	titleBlock2 := builder.CreateForText(secondTitleText)
	contentBlock := builder.CreateForText(contentText)
	document := webdoc.NewTextDocument([]*webdoc.TextBlock{
		titleBlock1,
		titleBlock2,
		contentBlock,
	})

	classifier := heuristic.NewDocumentTitleMatch(wc, titleText, secondTitleText)
	changed := classifier.Process(document)

	assert.True(t, changed)
	assert.True(t, titleBlock1.HasLabel(label.Title))
	assert.True(t, titleBlock2.HasLabel(label.Title))
	assert.False(t, contentBlock.HasLabel(label.Title))
}

func Test_Filter_Heuristic_DTM_TitleWithExtraCharacters(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := testutil.NewTextBlockBuilder(wc)
	text := "title:?! :?!text"

	titleBlock1 := builder.CreateForText(text)
	titleBlock2 := builder.CreateForText(text)
	document := webdoc.NewTextDocument([]*webdoc.TextBlock{
		titleBlock1,
		titleBlock2,
	})

	classifier := heuristic.NewDocumentTitleMatch(wc, "title text")
	changed := classifier.Process(document)

	assert.True(t, changed)
	assert.True(t, titleBlock1.HasLabel(label.Title))
	assert.True(t, titleBlock2.HasLabel(label.Title))
}
