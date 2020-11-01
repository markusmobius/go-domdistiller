// ORIGINAL: javatest/webdocument/WebImageTest.java

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

// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package webdoc_test

import (
	nurl "net/url"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_WebDoc_Image_GenerateOutput(t *testing.T) {
	html := `<picture>` +
		`<source srcset="image"/>` +
		`<img dirty-attributes/>` +
		`</picture>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	picture := dom.QuerySelector(div, "picture")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{Element: picture, PageURL: baseURL}

	expected := `<picture><source srcset="http://example.com/image"/><img/></picture>`
	assert.Equal(t, expected, webImage.GenerateOutput(false))
}

func Test_WebDoc_Image_GetSrcList(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", "image")
	dom.SetAttribute(img, "srcset", "image200 200w, image400 400w")

	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{
		Element: img,
		PageURL: baseURL,
	}

	urls := webImage.GetURLs()
	assert.Equal(t, 3, len(urls))
	assert.Equal(t, "http://example.com/image", urls[0])
	assert.Equal(t, "http://example.com/image200", urls[1])
	assert.Equal(t, "http://example.com/image400", urls[2])
}

func Test_WebDoc_Image_GetSrcListInPicture(t *testing.T) {
	html := `<picture>` +
		`<source srcset="image100 100w, //example.org/image300 300w"/>` +
		`<img/>` +
		`</picture>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	picture := dom.QuerySelector(div, "picture")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{Element: picture, PageURL: baseURL}

	urls := webImage.GetURLs()
	assert.Equal(t, 2, len(urls))
	assert.Equal(t, "http://example.com/image100", urls[0])
	assert.Equal(t, "http://example.org/image300", urls[1])
}

func Test_WebDoc_Image_PictureWithoutImg(t *testing.T) {
	html := `<picture>` +
		`<source srcset="image"/>` +
		`</picture>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	picture := dom.QuerySelector(div, "picture")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{Element: picture, PageURL: baseURL}

	expected := `<picture><source srcset="http://example.com/image"/></picture>`
	assert.Equal(t, expected, webImage.GenerateOutput(false))
}
