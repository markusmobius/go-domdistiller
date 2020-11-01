// ORIGINAL: javatest/webdocument/filters/RelevantElementsTest.java

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

func Test_Filter_DocFilter_RE_EmptyDocument(t *testing.T) {
	document := webdoc.NewDocument()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.True(t, len(document.Elements) == 0)
}

func Test_Filter_DocFilter_RE_NoContent(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1")
	builder.AddText("text 2")
	builder.AddTable("<tbody><tr><td>t1</td></tr></tbody>")

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	for _, e := range document.Elements {
		assert.False(t, e.IsContent())
	}
}

func Test_Filter_DocFilter_RE_RelevantTable(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	wt := builder.AddTable("<tbody><tr><td>t1</td></tr></tbody>")

	document := builder.Build()
	assert.True(t, docfilter.NewRelevantElements().Process(document))
	assert.True(t, wt.IsContent())
}

func Test_Filter_DocFilter_RE_NonRelevantTable(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	builder.AddText("text 2")
	wt := builder.AddTable("<tbody><tr><td>t1</td></tr></tbody>")

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.False(t, wt.IsContent())
}

func Test_Filter_DocFilter_RE_RelevantImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	wi := builder.AddImage()

	document := builder.Build()
	assert.True(t, docfilter.NewRelevantElements().Process(document))
	assert.True(t, wi.IsContent())
}

func Test_Filter_DocFilter_RE_NonRelevantImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddImage()
	builder.AddText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.False(t, wi.IsContent())
}

func Test_Filter_DocFilter_RE_ImageAfterNonContent(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	builder.AddText("text 2").SetIsContent(false)
	wi := builder.AddImage()

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.False(t, wi.IsContent())
}
