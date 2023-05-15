// ORIGINAL: java/webdocument/WebTable.java

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

package webdoc

import (
	"fmt"
	nurl "net/url"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

type Table struct {
	BaseElement

	Element *html.Node
	PageURL *nurl.URL

	cloned *html.Node
}

func (t *Table) ElementType() string {
	return "table"
}

func (t *Table) GenerateOutput(textOnly bool) string {
	if t.cloned == nil {
		t.cloned = domutil.CloneAndProcessTree(t.Element, t.PageURL)
	}

	if textOnly {
		return domutil.InnerText(t.cloned)
	}

	return dom.OuterHTML(t.cloned)
}

// GetImageURLs returns list of source URLs of all image inside the table.
func (t *Table) GetImageURLs() []string {
	if t.cloned == nil {
		t.cloned = domutil.CloneAndProcessTree(t.Element, t.PageURL)
	}

	imgURLs := []string{}
	for _, img := range dom.QuerySelectorAll(t.cloned, "img,source") {
		src := dom.GetAttribute(img, "src")
		if src != "" {
			imgURLs = append(imgURLs, src)
		}

		imgURLs = append(imgURLs, domutil.GetAllSrcSetURLs(img)...)
	}

	return imgURLs
}

func (t *Table) String() string {
	return fmt.Sprintf("ELEMENT %q: html=%q, is_content=%v",
		t.ElementType(), dom.OuterHTML(t.Element), t.isContent)
}
