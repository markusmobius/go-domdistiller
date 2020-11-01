// ORIGINAL: javatest/webdocument/DomConverterTest.java

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

package converter_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/converter"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_Converter_DisplayNone
// - Test_Converter_VisibilityHidden
// - Test_Converter_InvisibleInVisible
// - Test_Converter_VisibleInInvisible
// - Test_Converter_VisibleInInvisible2
// - Test_Converter_InvisibleInInvisible
// - Test_Converter_DifferentChildrenInVisible
// - Test_Converter_DifferentChildrenInInvisible
// - Test_Converter_KeepHidden
// - Test_Converter_KeepHiddenNested
// - Test_Converter_KeepContinue
// - Test_Converter_KeepContinueNested
// - Test_Converter_WikipediaFoldedSections

func Test_Converter_VisibleText(t *testing.T) {
	html := "visible text"
	assertEqual(t, html, html)
}

func Test_Converter_VisibleElement(t *testing.T) {
	html := "<div>visible element</div>"
	assertEqual(t, html, html)
}

func Test_Converter_VisibleInVisible(t *testing.T) {
	html := "<div>visible parent" +
		"<div>visible child</div>" +
		"</div>"
	assertEqual(t, html, html)
}

func Test_Converter_DataTable(t *testing.T) {
	html := `<table align="left" role="grid">` + // role=grid make this a data table.
		`<tbody align="left">` +
		`<tr>` +
		`<td>row1col1</td>` +
		`<td>row1col2</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	assertEqual(t, html, "<datatable/>")
}

func Test_Converter_NonDataTable(t *testing.T) {
	html := `<table align="left">` +
		`<tbody align="left">` +
		`<tr>` +
		`<td>row1col1</td>` +
		`<td>row1col2</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	assertEqual(t, html, html)
}

func Test_Converter_IgnorableElements(t *testing.T) {
	assertEmpty(t, "<head></head>")
	assertEmpty(t, "<style></style>")
	assertEmpty(t, "<script></script>")
	assertEmpty(t, "<link></link>")
	assertEmpty(t, "<noscript></noscript>")

	assertEqual(t, "<iframe></iframe>", "")
	assertEqual(t, "<svg></svg>", "")
	assertEqual(t, "<option></option>", "")
	assertEqual(t, "<object></object>", "")
	assertEqual(t, "<embed></embed>", "")
	assertEqual(t, "<applet></applet>", "")
}

func Test_Converter_SvgTagNameCase(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<SVG></SVG>")
	assert.Equal(t, "svg", dom.TagName(dom.FirstElementChild(div)))
}

func Test_Converter_ElementOrder(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, `Text content <img src="http://example.com/1.jpg"> more content`)

	builder := webdoc.NewWebDocumentBuilder(stringutil.FastWordCounter{}, nil)
	converter.NewDomConverter(converter.Default, builder, nil, nil).Convert(div)

	doc := builder.Build()
	elements := doc.Elements

	assert.Equal(t, 3, len(elements))
	assert.IsType(t, &webdoc.Text{}, elements[0])
	assert.IsType(t, &webdoc.Image{}, elements[1])
	assert.IsType(t, &webdoc.Text{}, elements[2])
}

func Test_Converter_LineBreak(t *testing.T) {
	html := "text<br>split<br/>with<br/>lines"
	assertEqual(t, html, "text\nsplit\nwith\nlines")
}

func Test_Converter_List(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<ol><li>some text1</li><li>some text2</li></ol>")

	builder := webdoc.NewWebDocumentBuilder(stringutil.FastWordCounter{}, nil)
	converter.NewDomConverter(converter.Default, builder, nil, nil).Convert(div)

	doc := builder.Build()
	elements := doc.Elements

	assert.Equal(t, 8, len(elements))
	assert.IsType(t, &webdoc.Tag{}, elements[0])
	assert.Equal(t, webdoc.TagStart, (elements[0].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Tag{}, elements[1])
	assert.Equal(t, webdoc.TagStart, (elements[1].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Text{}, elements[2])

	assert.IsType(t, &webdoc.Tag{}, elements[3])
	assert.Equal(t, webdoc.TagEnd, (elements[3].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Tag{}, elements[4])
	assert.Equal(t, webdoc.TagStart, (elements[4].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Text{}, elements[5])

	assert.IsType(t, &webdoc.Tag{}, elements[6])
	assert.Equal(t, webdoc.TagEnd, (elements[6].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Tag{}, elements[7])
	assert.Equal(t, webdoc.TagEnd, (elements[7].(*webdoc.Tag).Type))
}

func Test_Converter_SocialElements(t *testing.T) {
	assertEqual(t, `<div></div>`, ``)
	assertEqual(t, `<div data-component="share"></div>`, ``)
	assertEqual(t, `<div class="socialArea"></div>`, ``)
	assertEqual(t, `<li></li>`, `<li></li>`)
	assertEqual(t, `<li class="sharing"></li>`, ``)
}

func Test_Converter_WikiEditLinks(t *testing.T) {
	html := `<a href="index.php?action=edit&redlink=1"></a>`
	assertEqual(t, html, html)

	html = `<a href="index.php?action=edit&section=3" class="mw-ui-icon"></a>`
	assertEqual(t, html, "")
}

func Test_Converter_WikiEditSection(t *testing.T) {
	html := `<span class="mw-headline"></span>`
	assertEqual(t, html, html)

	html = `<span class="mw-editsection"></span>`
	assertEqual(t, html, "")
}

func assertEqual(t *testing.T, innerHTML, expectedHTML string) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, innerHTML)

	builder := testutil.NewFakeWebDocumentBuilder()
	converter.NewDomConverter(converter.Default, builder, nil, nil).Convert(div)

	expected := "<div>" + expectedHTML + "</div>"
	assert.Equal(t, expected, builder.Build())
}

func assertEmpty(t *testing.T, innerHTML string) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, innerHTML)

	builder := testutil.NewFakeWebDocumentBuilder()
	converter.NewDomConverter(converter.Default, builder, nil, nil).Convert(div)

	assert.Equal(t, "", builder.Build())
}
