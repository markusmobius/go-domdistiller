// ORIGINAL: java/webdocument/WebImage.java

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

package webdoc

import (
	"fmt"
	nurl "net/url"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type Image struct {
	BaseElement
	Element *html.Node // node for the image
	PageURL *nurl.URL  // url of page where image is placed

	cloned *html.Node
}

func (i *Image) ElementType() string {
	return "image"
}

func (i *Image) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	if i.cloned == nil {
		i.cloned = i.cloneAndProcessNode()
	}

	return dom.OuterHTML(i.cloned)
}

// GetURLs returns the list of source URLs of this image.
func (i *Image) GetURLs() []string {
	if i.cloned == nil {
		i.cloned = i.cloneAndProcessNode()
	}

	urls := []string{}
	src := dom.GetAttribute(i.cloned, "src")
	if src != "" {
		urls = append(urls, src)
	}

	urls = append(urls, domutil.GetAllSrcSetURLs(i.cloned)...)
	return urls
}

func (i *Image) getProcessedNode() *html.Node {
	if i.cloned == nil {
		i.cloned = i.cloneAndProcessNode()
	}
	return i.cloned
}

func (i *Image) cloneAndProcessNode() *html.Node {
	cloned := dom.Clone(i.Element, true)
	img := domutil.GetFirstElementByTagNameInc(cloned, "img")
	if img != nil {
		if src := dom.GetAttribute(img, "src"); src != "" {
			src = stringutil.CreateAbsoluteURL(src, i.PageURL)
			dom.SetAttribute(img, "src", src)
		}
	}

	domutil.MakeAllSrcAttributesAbsolute(cloned, i.PageURL)
	domutil.MakeAllSrcSetAbsolute(cloned, i.PageURL)
	domutil.StripAttributes(cloned)
	return cloned
}

func (i *Image) String() string {
	return fmt.Sprintf("ELEMENT %q: html=%q, is_content=%v",
		i.ElementType(), dom.OuterHTML(i.Element), i.isContent)
}
