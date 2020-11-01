// ORIGINAL: javatest/webdocument/filters/NestedElementRetainerTest.java

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
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_NER_OrderedListStructure(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	olStart := builder.AddTagStart("ol")
	liStart1 := builder.AddTagStart("li")
	builder.AddText("text 1").SetIsContent(false)
	liEnd1 := builder.AddTagEnd("li")
	liStart2 := builder.AddTagStart("li")
	builder.AddText("text 2").SetIsContent(false)
	liEnd2 := builder.AddTagEnd("li")
	liStart3 := builder.AddTagStart("li")
	builder.AddText("text 3").SetIsContent(true)
	liEnd3 := builder.AddTagEnd("li")
	olEnd := builder.AddTagEnd("ol")

	document := builder.Build()
	docfilter.NewNestedElementRetainer().Process(document)

	assert.True(t, olStart.IsContent())
	assert.False(t, liStart1.IsContent())
	assert.False(t, liEnd1.IsContent())
	assert.False(t, liStart2.IsContent())
	assert.False(t, liEnd2.IsContent())
	assert.True(t, liStart3.IsContent())
	assert.True(t, liEnd3.IsContent())
	assert.True(t, olEnd.IsContent())
}

func Test_Filter_DocFilter_NER_UnorderedListStructure(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	ulStart1 := builder.AddTagStart("ul")
	liStart1 := builder.AddTagStart("li")
	builder.AddText("text 1").SetIsContent(true)
	ulStart2 := builder.AddTagStart("ul")
	liStart2 := builder.AddTagStart("li")
	builder.AddText("text 2").SetIsContent(false)
	liEnd2 := builder.AddTagEnd("li")
	ulEnd2 := builder.AddTagEnd("ul")
	liEnd1 := builder.AddTagEnd("li")
	ulEnd1 := builder.AddTagEnd("ul")

	document := builder.Build()
	docfilter.NewNestedElementRetainer().Process(document)

	assert.True(t, ulStart1.IsContent())
	assert.True(t, liStart1.IsContent())
	assert.False(t, ulStart2.IsContent())
	assert.False(t, liStart2.IsContent())
	assert.False(t, liEnd2.IsContent())
	assert.False(t, ulEnd2.IsContent())
	assert.True(t, liEnd1.IsContent())
	assert.True(t, ulEnd1.IsContent())
}

func Test_Filter_DocFilter_NER_ContentFromListStrcture(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	olStartLevel1 := builder.AddTagStart("ol")
	olStartLevel2 := builder.AddTagStart("ol")
	liStart1 := builder.AddTagStart("li")
	builder.AddText("text 1").SetIsContent(false)
	liEnd1 := builder.AddTagEnd("li")
	olStartLevel3 := builder.AddTagStart("ol")
	liStart2 := builder.AddTagStart("li")
	builder.AddText("text 2").SetIsContent(true)
	liEnd2 := builder.AddTagEnd("li")
	olStartLevel4 := builder.AddTagStart("ol")
	liStart3 := builder.AddTagStart("li")
	builder.AddText("text 3").SetIsContent(false)
	liEnd3 := builder.AddTagEnd("li")
	liStart4 := builder.AddTagStart("li")
	builder.AddText("text 4").SetIsContent(false)
	liEnd4 := builder.AddTagEnd("li")
	liStart5 := builder.AddTagStart("li")
	builder.AddText("text 5").SetIsContent(false)
	liEnd5 := builder.AddTagEnd("li")
	liStart6 := builder.AddTagStart("li")
	builder.AddText("text 6").SetIsContent(false)
	liEnd6 := builder.AddTagEnd("li")
	olEndLevel4 := builder.AddTagEnd("ol")
	olEndLevel3 := builder.AddTagEnd("ol")
	liStart7 := builder.AddTagStart("li")
	builder.AddText("text 7").SetIsContent(true)
	liEnd7 := builder.AddTagEnd("li")
	olEndLevel2 := builder.AddTagEnd("ol")
	olEndLevel1 := builder.AddTagEnd("ol")

	document := builder.Build()
	docfilter.NewNestedElementRetainer().Process(document)

	assert.True(t, olStartLevel1.IsContent())
	assert.True(t, olStartLevel2.IsContent())
	assert.False(t, liStart1.IsContent())
	assert.False(t, liEnd1.IsContent())
	assert.True(t, olStartLevel3.IsContent())
	assert.True(t, liStart2.IsContent())
	assert.True(t, liEnd2.IsContent())
	assert.False(t, olStartLevel4.IsContent())
	assert.False(t, liStart3.IsContent())
	assert.False(t, liEnd3.IsContent())
	assert.False(t, liStart4.IsContent())
	assert.False(t, liEnd4.IsContent())
	assert.False(t, liStart5.IsContent())
	assert.False(t, liEnd5.IsContent())
	assert.False(t, liStart6.IsContent())
	assert.False(t, liEnd6.IsContent())
	assert.False(t, olEndLevel4.IsContent())
	assert.True(t, olEndLevel3.IsContent())
	assert.True(t, liStart7.IsContent())
	assert.True(t, liEnd7.IsContent())
	assert.True(t, olEndLevel2.IsContent())
	assert.True(t, olEndLevel1.IsContent())
}

func Test_Filter_DocFilter_NER_NoContentFromListStrcture(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	olStartLevel1 := builder.AddTagStart("ol")
	olStartLevel2 := builder.AddTagStart("ol")
	liStart1 := builder.AddTagStart("li")
	builder.AddText("text 1").SetIsContent(false)
	liEnd1 := builder.AddTagEnd("li")
	olStartLevel4 := builder.AddTagStart("ol")
	liStart3 := builder.AddTagStart("li")
	builder.AddText("text 3").SetIsContent(false)
	liEnd3 := builder.AddTagEnd("li")
	liStart4 := builder.AddTagStart("li")
	builder.AddText("text 4").SetIsContent(false)
	liEnd4 := builder.AddTagEnd("li")
	liStart5 := builder.AddTagStart("li")
	builder.AddText("text 5").SetIsContent(false)
	liEnd5 := builder.AddTagEnd("li")
	liStart6 := builder.AddTagStart("li")
	builder.AddText("text 6").SetIsContent(false)
	liEnd6 := builder.AddTagEnd("li")
	olEndLevel4 := builder.AddTagEnd("ol")
	olEndLevel2 := builder.AddTagEnd("ol")
	olEndLevel1 := builder.AddTagEnd("ol")

	document := builder.Build()
	docfilter.NewNestedElementRetainer().Process(document)

	assert.False(t, olStartLevel1.IsContent())
	assert.False(t, olStartLevel2.IsContent())
	assert.False(t, liStart1.IsContent())
	assert.False(t, liEnd1.IsContent())
	assert.False(t, olStartLevel4.IsContent())
	assert.False(t, liStart3.IsContent())
	assert.False(t, liEnd3.IsContent())
	assert.False(t, liStart4.IsContent())
	assert.False(t, liEnd4.IsContent())
	assert.False(t, liStart5.IsContent())
	assert.False(t, liEnd5.IsContent())
	assert.False(t, liStart6.IsContent())
	assert.False(t, liEnd6.IsContent())
	assert.False(t, olEndLevel4.IsContent())
	assert.False(t, olEndLevel2.IsContent())
	assert.False(t, olEndLevel1.IsContent())
}

func Test_Filter_DocFilter_NER_NestedListStructure(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	ulStart := builder.AddTagStart("ul")
	liStart := builder.AddTagStart("li")
	builder.AddText("text 1").SetIsContent(true)
	liEnd := builder.AddTagEnd("li")
	olStart := builder.AddTagStart("ol")
	liOLStart := builder.AddTagStart("li")
	builder.AddText("text 2").SetIsContent(true)
	liOLEnd := builder.AddTagEnd("li")
	olEnd := builder.AddTagEnd("ol")
	ulEnd := builder.AddTagEnd("ul")

	document := builder.Build()
	docfilter.NewNestedElementRetainer().Process(document)

	assert.True(t, ulStart.IsContent())
	assert.True(t, liStart.IsContent())
	assert.True(t, olStart.IsContent())
	assert.True(t, liOLStart.IsContent())
	assert.True(t, liOLEnd.IsContent())
	assert.True(t, olEnd.IsContent())
	assert.True(t, liEnd.IsContent())
	assert.True(t, ulEnd.IsContent())
}
