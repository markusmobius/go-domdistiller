// ORIGINAL: javatest/webdocument/WebTableTest.java

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
	nurl "net/url"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_WebDoc_Table_GenerateOutput(t *testing.T) {
	html := `<table><tbody>` +
		`<tr>` +
		`<td>row1col1</td>` +
		`<td><img src="http://example.com/table.png"/></td>` +
		`<td><picture><img/></picture></td>` +
		`</tr>` +
		`</tbody></table>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	table := dom.QuerySelector(div, "table")
	webTable := webdoc.Table{Element: table}

	// Output should be the same as the input in this case.
	got := webTable.GenerateOutput(false)
	assert.Equal(t, html, testutil.RemoveAllDirAttributes(got))

	// Test GetImageURLs as well.
	imgURLs := webTable.GetImageURLs()
	assert.Equal(t, 1, len(imgURLs))
	assert.Equal(t, "http://example.com/table.png", imgURLs[0])
}

func Test_WebDoc_Table_GetImageURLs(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, `
<table>
<tbody>
<tr>
<td>
<img src="http://example.com/table.png" srcset="image100 100w, //example.org/image300 300w"/>
</td>
<td>
<picture>
<source srcset="image200 200w, //example.org/image400 400w"/>
<img/>
</picture>
</td>
</tr>
</tbody>
</table>`)

	table := dom.QuerySelector(div, "table")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webTable := webdoc.Table{Element: table, PageURL: baseURL}

	urls := webTable.GetImageURLs()
	assert.Equal(t, 5, len(urls))
	assert.Equal(t, "http://example.com/table.png", urls[0])
	assert.Equal(t, "http://example.com/image100", urls[1])
	assert.Equal(t, "http://example.org/image300", urls[2])
	assert.Equal(t, "http://example.com/image200", urls[3])
	assert.Equal(t, "http://example.org/image400", urls[4])
}
