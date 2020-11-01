// ORIGINAL: javatest/webdocument/WebVideoTest.java

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
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_WebDoc_Video_GenerateOutput(t *testing.T) {
	video := dom.CreateElement("video")
	dom.SetAttribute(video, "onfocus", "new XMLHttpRequest();") // should be stripped

	child := dom.CreateElement("source")
	dom.SetAttribute(child, "src", "http://example.com/foo.ogg")
	dom.AppendChild(video, child)

	child = dom.CreateElement("track")
	dom.SetAttribute(child, "src", "http://example.com/foo.vtt")
	dom.SetAttribute(child, "onclick", "alert(1)") // should be stripped
	dom.AppendChild(video, child)

	expected := `<video>` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>`

	webVideo := webdoc.NewVideo(video, nil, 0, 0)
	assert.Equal(t, expected, webVideo.GenerateOutput(false))
}

func Test_WebDoc_Video_GenerateOutputInvalidChildren(t *testing.T) {
	video := dom.CreateElement("video")
	dom.SetAttribute(video, "onfocus", "new XMLHttpRequest();") // should be stripped

	child := dom.CreateElement("source")
	dom.SetAttribute(child, "src", "http://example.com/foo.ogg")
	dom.AppendChild(video, child)

	child = dom.CreateElement("track")
	dom.SetAttribute(child, "src", "http://example.com/foo.vtt")
	dom.SetAttribute(child, "onclick", "alert(1)") // should be stripped
	dom.AppendChild(video, child)

	child = dom.CreateElement("div")
	dom.SetTextContent(child, "We do not use custom error messages!")
	dom.AppendChild(video, child)

	// Output should ignore anything other than "track" and "source" tags.
	expected := `<video>` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>`

	webVideo := webdoc.NewVideo(video, nil, 0, 0)
	assert.Equal(t, expected, webVideo.GenerateOutput(false))
}

func Test_WebDoc_Video_GenerateOutputRelativeURL(t *testing.T) {
	video := dom.CreateElement("video")
	dom.SetAttribute(video, "poster", "bar.png")                // should be stripped
	dom.SetAttribute(video, "onfocus", "new XMLHttpRequest();") // should be stripped

	child := dom.CreateElement("source")
	dom.SetAttribute(child, "src", "foo.ogg")
	dom.AppendChild(video, child)

	child = dom.CreateElement("track")
	dom.SetAttribute(child, "src", "foo.vtt")
	dom.SetAttribute(child, "onclick", "alert(1)") // should be stripped
	dom.AppendChild(video, child)

	expected := `<video poster="http://example.com/bar.png">` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>`

	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webVideo := webdoc.NewVideo(video, baseURL, 0, 0)
	assert.Equal(t, expected, webVideo.GenerateOutput(false))
}

func Test_WebDoc_Video_PosterEmpty(t *testing.T) {
	video := dom.CreateElement("video")
	webVideo := webdoc.NewVideo(video, nil, 400, 300)
	assert.Equal(t, "<video></video>", webVideo.GenerateOutput(false))
}
