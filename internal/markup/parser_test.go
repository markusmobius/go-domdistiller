// ORIGINAL: javatest/MarkupParserTest.java and javatest/MarkupParserProtoTest.java

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

package markup_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/markup"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_Markup_NullOpenGraphProtocolParser(t *testing.T) {
	expectedTitle := "Testing null OpenGraphProtocolParser."

	doc := testutil.CreateHTML()
	head := dom.QuerySelector(doc, "head")
	body := dom.QuerySelector(doc, "body")

	dom.AppendChild(head, testutil.CreateTitle(expectedTitle))
	dom.AppendChild(head, testutil.CreateMetaName("title", expectedTitle))
	dom.AppendChild(body, testutil.CreateHeading(1, expectedTitle))

	parser := markup.NewParser(doc, nil)
	assert.Equal(t, expectedTitle, parser.Title())
}

// TODO: write more tests if or when we determine:
// - which parser takes precedence
// - how we merge the different values retrieved from the different parsers.

func Test_Markup_CompleteInfoWithMultipleImages(t *testing.T) {
	// Create the required properties for OpenGraphProtocol except for "image".
	doc := testutil.CreateHTML()
	createDefaultOGTitle(doc)
	createDefaultOGType(doc)
	createDefaultOGUrl(doc)

	// Create the image properties and their structures.
	expectedURL1 := "http://test/image1.jpeg"
	expectedSecureURL1 := "https://test/image1.jpeg"
	expectedType1 := "image/jpeg"
	createMeta(doc, "og:image", expectedURL1)
	createMeta(doc, "og:image:url", expectedURL1)
	createMeta(doc, "og:image:secure_url", expectedSecureURL1)
	createMeta(doc, "og:image:type", expectedType1)
	createMeta(doc, "og:image:width", "600")
	createMeta(doc, "og:image:height", "400")

	expectedURL2 := "http://test/image2.jpeg"
	expectedSecureURL2 := "https://test/image2.jpeg"
	expectedType2 := "image/gif"
	createMeta(doc, "og:image", expectedURL2)
	createMeta(doc, "og:image:url", expectedURL2)
	createMeta(doc, "og:image:secure_url", expectedSecureURL2)
	createMeta(doc, "og:image:type", expectedType2)
	createMeta(doc, "og:image:width", "1024")
	createMeta(doc, "og:image:height", "900")

	parser := markup.NewParser(doc, nil)
	markupInfo := parser.MarkupInfo()
	assert.Equal(t, "dummy title", markupInfo.Title)
	assert.Equal(t, "", markupInfo.Type)
	assert.Equal(t, "http://dummy/url.html", markupInfo.URL)
	assert.Equal(t, 2, len(markupInfo.Images))

	markupImage := markupInfo.Images[0]
	assert.Equal(t, expectedURL1, markupImage.URL)
	assert.Equal(t, expectedSecureURL1, markupImage.SecureURL)
	assert.Equal(t, expectedType1, markupImage.Type)
	assert.Equal(t, 600, markupImage.Width)
	assert.Equal(t, 400, markupImage.Height)

	markupImage = markupInfo.Images[1]
	assert.Equal(t, expectedURL2, markupImage.URL)
	assert.Equal(t, expectedSecureURL2, markupImage.SecureURL)
	assert.Equal(t, expectedType2, markupImage.Type)
	assert.Equal(t, 1024, markupImage.Width)
	assert.Equal(t, 900, markupImage.Height)
}

func Test_Markup_Article(t *testing.T) {
	// Create the required properties for OpenGraphProtocol except for "type".
	doc := testutil.CreateHTML()
	createDefaultOGTitle(doc)
	createDefaultOGUrl(doc)
	createDefaultOGImage(doc)

	// Create the "article" object.
	expectedSection := "GWT Testing"
	expectedPublishedTime := "2014-04-01T01:23:59Z"
	expectedModifiedTime := "2014-04-02T02:23:59Z"
	expectedExpirationTime := "2014-04-03T03:23:59Z"
	expectedAuthor1 := "http://blah/author1.html"
	expectedAuthor2 := "http://blah/author2.html"
	createMeta(doc, "og:type", "article")
	createMeta(doc, "article:section", expectedSection)
	createMeta(doc, "article:published_time", expectedPublishedTime)
	createMeta(doc, "article:modified_time", expectedModifiedTime)
	createMeta(doc, "article:expiration_time", expectedExpirationTime)
	createMeta(doc, "article:author", expectedAuthor1)
	createMeta(doc, "article:author", expectedAuthor2)

	parser := markup.NewParser(doc, nil)
	markupInfo := parser.MarkupInfo()
	markupArticle := markupInfo.Article
	assert.Equal(t, "Article", markupInfo.Type)
	assert.Equal(t, expectedPublishedTime, markupArticle.PublishedTime)
	assert.Equal(t, expectedModifiedTime, markupArticle.ModifiedTime)
	assert.Equal(t, expectedExpirationTime, markupArticle.ExpirationTime)
	assert.Equal(t, expectedSection, markupArticle.Section)
	assert.Equal(t, 2, len(markupArticle.Authors))
	assert.Equal(t, expectedAuthor1, markupArticle.Authors[0])
	assert.Equal(t, expectedAuthor2, markupArticle.Authors[1])
}

func createDefaultOGTitle(doc *html.Node) {
	createMeta(doc, "og:title", "dummy title")
}

func createDefaultOGType(doc *html.Node) {
	createMeta(doc, "og:type", "website")
}

func createDefaultOGUrl(doc *html.Node) {
	createMeta(doc, "og:url", "http://dummy/url.html")
}

func createDefaultOGImage(doc *html.Node) {
	createMeta(doc, "og:image", "http://dummy/image.jpeg")
}

func createMeta(doc *html.Node, property, content string) {
	head := dom.QuerySelector(doc, "head")
	if head != nil {
		meta := testutil.CreateMetaProperty(property, content)
		dom.AppendChild(head, meta)
	}
}
