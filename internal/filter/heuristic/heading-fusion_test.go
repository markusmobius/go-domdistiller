// ORIGINAL: javatest/HeadingFusionTest.java

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
	"strings"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_Heuristic_HF_HeadingFused(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(headingText, label.Heading)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.True(t, changed)

	textBlocks := textDocument.TextBlocks
	assert.Len(t, textBlocks, 2)
	assert.False(t, textBlocks[0].HasLabel(label.Heading))
	assert.False(t, textBlocks[0].HasLabel(label.BoilerplateHeadingFused))

	docContent := testutil.GetContentFromTextDocument(textDocument)
	assert.True(t, strings.Contains(docContent, headingText))
	assert.True(t, strings.Contains(docContent, longText))
	assert.True(t, strings.Contains(docContent, shortText))
}

func Test_Filter_Heuristic_HF_BoilerplateHeadingFused(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddNonContentBlock(headingText, label.Heading)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.True(t, changed)

	textBlocks := textDocument.TextBlocks
	assert.Len(t, textBlocks, 2)
	assert.False(t, textBlocks[0].HasLabel(label.Heading))
	assert.True(t, textBlocks[0].HasLabel(label.BoilerplateHeadingFused))

	docContent := testutil.GetContentFromTextDocument(textDocument)
	assert.True(t, strings.Contains(docContent, headingText))
	assert.True(t, strings.Contains(docContent, longText))
	assert.True(t, strings.Contains(docContent, shortText))
}

func Test_Filter_Heuristic_HF_BeforeBoilerplate(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(headingText, label.Heading)
	builder.AddNonContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.True(t, changed)

	textBlocks := textDocument.TextBlocks
	assert.Len(t, textBlocks, 3)
	assert.False(t, textBlocks[0].IsContent())
}

func Test_Filter_Heuristic_HF_TitleNotFused(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(headingText, label.Heading, label.Title)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.False(t, changed)
	assert.Len(t, textDocument.TextBlocks, 3)
}
