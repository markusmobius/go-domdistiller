// ORIGINAL: javatest/webdocument/filters/LeadImageFinderTest.java

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

package docfilter_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_LIF_EmptyDocument(t *testing.T) {
	document := webdoc.NewDocument()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.Len(t, document.Elements, 0)
}

func Test_Filter_DocFilter_LIF_RelevantLeadImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddLeadImage()
	builder.AddText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.True(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, wi.IsContent())
}

func Test_Filter_DocFilter_LIF_NoContent(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddLeadImage()
	builder.AddText("text 1").SetIsContent(false)

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.False(t, wi.IsContent())
}

func Test_Filter_DocFilter_LIF_MultipleLeadImageCandidates(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	priority := builder.AddLeadImage()
	ignored := builder.AddLeadImage()
	builder.AddText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.True(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, priority.IsContent())
	assert.False(t, ignored.IsContent())
}

func Test_Filter_DocFilter_LIF_IrrelevantLeadImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	builder.AddText("text 2").SetIsContent(true)
	priority := builder.AddLeadImage()

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.False(t, priority.IsContent())
}

func Test_Filter_DocFilter_LIF_MultipleLeadImageCandidatesWithinTexts(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	priority := builder.AddLeadImage()
	builder.AddText("text 2").SetIsContent(true)
	ignored := builder.AddLeadImage()

	document := builder.Build()
	assert.True(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, priority.IsContent())
	assert.False(t, ignored.IsContent())
}

func Test_Filter_DocFilter_LIF_IrrelevantLeadImageWithContentImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	smallImage := builder.AddImage()
	smallImage.SetIsContent(true)
	largeImage := builder.AddLeadImage()
	builder.AddNestedText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, smallImage.IsContent())
	assert.False(t, largeImage.IsContent())
}

func Test_Filter_DocFilter_LIF_IrrelevantLeadImageAfterSingleText(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddImage()
	builder.AddNestedText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.False(t, wi.IsContent())
}
