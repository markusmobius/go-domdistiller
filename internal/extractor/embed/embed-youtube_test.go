// ORIGINAL: javatest/EmbedExtractorTest.java

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

package embed_test

import (
	nurl "net/url"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/extractor/embed"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_Embed_YouTube_Extract(t *testing.T) {
	youtube := dom.CreateElement("iframe")
	dom.SetAttribute(youtube, "src", "//www.youtube.com/embed/M7lc1UVf-VE?autoplay=1&hl=zh_TW")

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewYouTubeExtractor(pageURL, nil)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "M7lc1UVf-VE", result.ID)
	assert.Equal(t, "1", result.Params["autoplay"])
	assert.Equal(t, "zh_TW", result.Params["hl"])

	// Begin negative test
	notYoutube := dom.CreateElement("iframe")
	dom.SetAttribute(notYoutube, "src", "http://www.notyoutube.com/embed/M7lc1UVf-VE?autoplay=1")

	result, _ = (extractor.Extract(notYoutube)).(*webdoc.Embed)
	assert.Nil(t, result)
}

func Test_Embed_YouTube_ExtractID(t *testing.T) {
	youtube := dom.CreateElement("iframe")
	dom.SetAttribute(youtube, "src", "http://www.youtube.com/embed/M7lc1UVf-VE///?autoplay=1")

	extractor := embed.NewYouTubeExtractor(nil, nil)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "M7lc1UVf-VE", result.ID)

	// Begin negative test
	notYoutube := dom.CreateElement("iframe")
	dom.SetAttribute(notYoutube, "src", "http://www.youtube.com/embed")

	result, _ = (extractor.Extract(notYoutube)).(*webdoc.Embed)
	assert.Nil(t, result)
}

func Test_Embed_YouTube_Object(t *testing.T) {
	html := `<object>` +
		`<param name="movie" ` +
		`value="//www.youtube.com/v/DML2WUhn2M0&hl=en_US&fs=1&hd=1">` +
		`</param>` +
		`<param name="allowFullScreen" value="true">` +
		`</param>` +
		`</object>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	youtube := dom.FirstElementChild(div)
	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewYouTubeExtractor(pageURL, nil)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "DML2WUhn2M0", result.ID)
	assert.Equal(t, "en_US", result.Params["hl"])
	assert.Equal(t, "1", result.Params["fs"])
	assert.Equal(t, "1", result.Params["hd"])
}

func Test_Embed_YouTube_Object2(t *testing.T) {
	html := `<object type="application/x-shockwave-flash" ` +
		`data="http://www.youtube.com/v/ZuNNhOEzJGA&hl=fr&fs=1&rel=0&color1=0x006699&color2=0x54abd6&border=1">` +
		`</object>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	youtube := dom.FirstElementChild(div)
	extractor := embed.NewYouTubeExtractor(nil, nil)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "ZuNNhOEzJGA", result.ID)
	assert.Equal(t, "fr", result.Params["hl"])
	assert.Equal(t, "1", result.Params["fs"])
	assert.Equal(t, "0", result.Params["rel"])
	assert.Equal(t, "0x006699", result.Params["color1"])
	assert.Equal(t, "0x54abd6", result.Params["color2"])
	assert.Equal(t, "1", result.Params["border"])
}
