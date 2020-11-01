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
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_Embed_Twitter_ExtractNotRenderedBasic(t *testing.T) {
	tweetBlock := dom.CreateElement("blockquote")
	dom.SetAttribute(tweetBlock, "class", "twitter-tweet")

	p := dom.CreateElement("p")
	dom.AppendChild(p, testutil.CreateAnchor("//twitter.com/foo", "extra content"))
	dom.AppendChild(tweetBlock, p)
	dom.AppendChild(tweetBlock, testutil.CreateAnchor("//twitter.com/foo/bar/12345", "January 1, 1900"))

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewTwitterExtractor(pageURL, nil)
	result, _ := (extractor.Extract(tweetBlock)).(*webdoc.Embed)

	assert.NotNil(t, result)
	assert.Equal(t, "twitter", result.Type)
	assert.Equal(t, "12345", result.ID)

	// Test trailing slash
	tweetBlock = dom.CreateElement("blockquote")
	dom.SetAttribute(tweetBlock, "class", "twitter-tweet")

	p = dom.CreateElement("p")
	dom.AppendChild(p, testutil.CreateAnchor("http://twitter.com/foo", "extra content"))
	dom.AppendChild(tweetBlock, p)
	dom.AppendChild(tweetBlock, testutil.CreateAnchor("http://twitter.com/foo/bar/12345///", "January 1, 1900"))

	result, _ = (extractor.Extract(tweetBlock)).(*webdoc.Embed)

	assert.NotNil(t, result)
	assert.Equal(t, "twitter", result.Type)
	assert.Equal(t, "12345", result.ID)
}

func Test_Embed_Twitter_ExtractNotRenderedTrailingSlash(t *testing.T) {
	tweetBlock := dom.CreateElement("blockquote")
	dom.SetAttribute(tweetBlock, "class", "twitter-tweet")

	p := dom.CreateElement("p")
	dom.AppendChild(p, testutil.CreateAnchor("http://twitter.com/foo", "extra content"))
	dom.AppendChild(tweetBlock, p)
	dom.AppendChild(tweetBlock, testutil.CreateAnchor("http://twitter.com/foo/bar/12345///", "January 1, 1900"))

	extractor := embed.NewTwitterExtractor(nil, nil)
	result, _ := (extractor.Extract(tweetBlock)).(*webdoc.Embed)

	assert.NotNil(t, result)
	assert.Equal(t, "twitter", result.Type)
	assert.Equal(t, "12345", result.ID)
}

func Test_Embed_Twitter_ExtractNotRenderedBadTweet(t *testing.T) {
	tweetBlock := dom.CreateElement("blockquote")
	dom.SetAttribute(tweetBlock, "class", "random-class")

	p := dom.CreateElement("p")
	dom.AppendChild(p, testutil.CreateAnchor("http://nottwitter.com/foo", "extra content"))
	dom.AppendChild(tweetBlock, p)
	dom.AppendChild(tweetBlock, testutil.CreateAnchor("http://nottwitter.com/12345", "timestamp"))

	extractor := embed.NewTwitterExtractor(nil, nil)
	result, _ := (extractor.Extract(tweetBlock)).(*webdoc.Embed)

	assert.Nil(t, result)
}

func Test_Embed_Twitter_ExtractRenderedBasic(t *testing.T) {
	tweet := dom.CreateElement("iframe")
	dom.SetAttribute(tweet, "id", "twitter-widget")
	dom.SetAttribute(tweet, "title", "Twitter Tweet")
	dom.SetAttribute(tweet, "src", "https://platform.twitter.com/embed/index.html")
	dom.SetAttribute(tweet, "data-tweet-id", "12345")

	extractor := embed.NewTwitterExtractor(nil, nil)
	result, _ := (extractor.Extract(tweet)).(*webdoc.Embed)

	assert.NotNil(t, result)
	assert.Equal(t, "twitter", result.Type)
	assert.Equal(t, "12345", result.ID)
}

func Test_Embed_Twitter_ExtractRenderedBadTweet(t *testing.T) {
	tweet := dom.CreateElement("iframe")
	dom.SetAttribute(tweet, "id", "twitter-widget")
	dom.SetAttribute(tweet, "title", "Twitter Tweet")
	dom.SetAttribute(tweet, "src", "https://platform.not-twitter.com/embed/index.html")
	dom.SetAttribute(tweet, "data-bad-id", "12345")

	extractor := embed.NewTwitterExtractor(nil, nil)
	result, _ := (extractor.Extract(tweet)).(*webdoc.Embed)

	assert.Nil(t, result)
}
