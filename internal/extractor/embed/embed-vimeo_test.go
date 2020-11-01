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

func Test_Embed_Vimeo_Extract(t *testing.T) {
	vimeo := dom.CreateElement("iframe")
	dom.SetAttribute(vimeo, "src", "//player.vimeo.com/video/12345?portrait=0")

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewVimeoExtractor(pageURL, nil)
	result, _ := (extractor.Extract(vimeo)).(*webdoc.Embed)

	// Check Vimeo specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "vimeo", result.Type)
	assert.Equal(t, "12345", result.ID)
	assert.Equal(t, "0", result.Params["portrait"])

	// Begin negative test
	wrongDomain := dom.CreateElement("iframe")
	dom.SetAttribute(wrongDomain, "src", "http://vimeo.com/video/09876?portrait=1")

	result, _ = (extractor.Extract(wrongDomain)).(*webdoc.Embed)
	assert.Nil(t, result)
}

func Test_Embed_Vimeo_ExtractID(t *testing.T) {
	vimeo := dom.CreateElement("iframe")
	dom.SetAttribute(vimeo, "src", "http://player.vimeo.com/video/12345?portrait=0")

	extractor := embed.NewVimeoExtractor(nil, nil)
	result, _ := (extractor.Extract(vimeo)).(*webdoc.Embed)

	// Check Vimeo specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "vimeo", result.Type)
	assert.Equal(t, "12345", result.ID)

	// Begin negative test
	wrongDomain := dom.CreateElement("iframe")
	dom.SetAttribute(wrongDomain, "src", "http://player.vimeo.com/video")

	result, _ = (extractor.Extract(wrongDomain)).(*webdoc.Embed)
	assert.Nil(t, result)
}
