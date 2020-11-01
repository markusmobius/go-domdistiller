// ORIGINAL: javatest/OpenGraphProtocolParserAccessorTest.java

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

package opengraph_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/markup/opengraph"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_OpenGraph_RequiredPropertiesAndDescriptionAndSiteName(t *testing.T) {
	expectedTitle := "Testing required OpenGraph Protocol properties and optional description of the document."
	expectedImage := "http://test/image.jpeg"
	expectedURL := "http://test/test.html"
	expectedDescr := "This test expects to retrieve the required OpenGraph Protocol properties and optional description of the document."
	expectedSiteName := "Google"

	root := testutil.CreateHTML()
	createMeta(root, "og:title", expectedTitle)
	createMeta(root, "og:type", "video.movie")
	createMeta(root, "og:image", expectedImage)
	createMeta(root, "og:url", expectedURL)
	createMeta(root, "og:description", expectedDescr)
	createMeta(root, "og:site_name", expectedSiteName)

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	images := parser.Images()
	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, expectedURL, parser.URL())
	assert.Equal(t, expectedDescr, parser.Description())
	assert.Equal(t, expectedSiteName, parser.Publisher())
	assert.Equal(t, 1, len(images))
	assert.Equal(t, expectedImage, images[0].URL)
	assert.Equal(t, "", images[0].SecureURL)
	assert.Equal(t, "", images[0].Type)
	assert.Equal(t, 0, images[0].Width)
	assert.Equal(t, 0, images[0].Height)
}

func Test_OpenGraph_NoRequiredImage(t *testing.T) {
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultType(root)
	createDefaultURL(root)

	// Set an image structured property but not the root property.
	createMeta(root, "og:image:url", "http://test/image.jpeg")

	parser, _ := opengraph.NewParser(root, nil)
	assert.Nil(t, parser)
}

func Test_OpenGraph_OneImage(t *testing.T) {
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultType(root)
	createDefaultURL(root)

	expectedURL := "http://test/image.jpeg"
	expectedSecureURL := "https://test/image.jpeg"
	expectedType := "image/jpeg"

	createMeta(root, "og:image", expectedURL)
	createMeta(root, "og:image:url", expectedURL)
	createMeta(root, "og:image:secure_url", expectedSecureURL)
	createMeta(root, "og:image:type", expectedType)
	createMeta(root, "og:image:width", "600")
	createMeta(root, "og:image:height", "400")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	images := parser.Images()
	assert.Equal(t, 1, len(images))
	assert.Equal(t, expectedURL, images[0].URL)
	assert.Equal(t, expectedSecureURL, images[0].SecureURL)
	assert.Equal(t, expectedType, images[0].Type)
	assert.Equal(t, 600, images[0].Width)
	assert.Equal(t, 400, images[0].Height)
}

func Test_OpenGraph_CompleteMultipleImages(t *testing.T) {
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultType(root)
	createDefaultURL(root)

	// Create the image properties and their structures.
	expectedURL1 := "http://test/image1.jpeg"
	expectedSecureURL1 := "https://test/image1.jpeg"
	expectedType1 := "image/jpeg"
	expectedURL2 := "http://test/image2.jpeg"
	expectedSecureURL2 := "https://test/image2.jpeg"
	expectedType2 := "image/gif"

	createMeta(root, "og:image", expectedURL1)
	createMeta(root, "og:image:url", expectedURL1)
	createMeta(root, "og:image:secure_url", expectedSecureURL1)
	createMeta(root, "og:image:type", expectedType1)
	createMeta(root, "og:image:width", "600")
	createMeta(root, "og:image:height", "400")
	createMeta(root, "og:image", expectedURL2)
	createMeta(root, "og:image:url", expectedURL2)
	createMeta(root, "og:image:secure_url", expectedSecureURL2)
	createMeta(root, "og:image:type", expectedType2)
	createMeta(root, "og:image:width", "1024")
	createMeta(root, "og:image:height", "900")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	images := parser.Images()
	assert.Equal(t, 2, len(images))

	image := images[0]
	assert.Equal(t, expectedURL1, image.URL)
	assert.Equal(t, expectedSecureURL1, image.SecureURL)
	assert.Equal(t, expectedType1, image.Type)
	assert.Equal(t, 600, image.Width)
	assert.Equal(t, 400, image.Height)

	image = images[1]
	assert.Equal(t, expectedURL2, image.URL)
	assert.Equal(t, expectedSecureURL2, image.SecureURL)
	assert.Equal(t, expectedType2, image.Type)
	assert.Equal(t, 1024, image.Width)
	assert.Equal(t, 900, image.Height)
}

func Test_OpenGraph_IncompleteMultipleImages(t *testing.T) {
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultType(root)
	createDefaultURL(root)

	// Create the image properties and their structures.
	expectedURL1 := "http://test/image1.jpeg"
	expectedSecureURL1 := "https://test/image1.jpeg"
	expectedType1 := "image/jpeg"
	createMeta(root, "og:image", expectedURL1)
	createMeta(root, "og:image:url", expectedURL1)
	createMeta(root, "og:image:secure_url", expectedSecureURL1)
	createMeta(root, "og:image:type", expectedType1)

	// Intentionally insert a root image tag before the width and height
	// tags, so the width and height should belong to the 2nd image.
	expectedURL2 := "http://test/image2.jpeg"
	createMeta(root, "og:image", expectedURL2)
	createMeta(root, "og:image:width", "600")
	createMeta(root, "og:image:height", "400")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	images := parser.Images()
	assert.Equal(t, 2, len(images))

	image := images[0]
	assert.Equal(t, expectedURL1, image.URL)
	assert.Equal(t, expectedSecureURL1, image.SecureURL)
	assert.Equal(t, expectedType1, image.Type)
	assert.Equal(t, 0, image.Width)
	assert.Equal(t, 0, image.Height)

	image = images[1]
	assert.Equal(t, expectedURL2, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, 600, image.Width)
	assert.Equal(t, 400, image.Height)
}

func Test_OpenGraph_NoObjects(t *testing.T) {
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultType(root)
	createDefaultURL(root)
	createDefaultImage(root)

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, "", parser.Author())
	assert.Nil(t, parser.Article())
}

func Test_OpenGraph_Profile(t *testing.T) {
	// Create the required properties except for "type".
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultURL(root)
	createDefaultImage(root)

	// Create the "profile" object.
	createMeta(root, "og:type", "profile")
	createMeta(root, "profile:first_name", "Jane")
	createMeta(root, "profile:last_name", "Doe")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, "Jane Doe", parser.Author())
}

func Test_OpenGraph_Article(t *testing.T) {
	// Create the required properties except for "type".
	root := testutil.CreateHTML()
	createDefaultTitle(root)
	createDefaultURL(root)
	createDefaultImage(root)

	// Create the "article" object.
	expectedSection := "GWT Testing"
	expectedPublishedTime := "2014-04-01T01:23:59Z"
	expectedModifiedTime := "2014-04-02T02:23:59Z"
	expectedExpirationTime := "2014-04-03T03:23:59Z"
	expectedAuthor1 := "http://blah/author1.html"
	expectedAuthor2 := "http://blah/author2.html"
	createMeta(root, "og:type", "article")
	createMeta(root, "article:section", expectedSection)
	createMeta(root, "article:published_time", expectedPublishedTime)
	createMeta(root, "article:modified_time", expectedModifiedTime)
	createMeta(root, "article:expiration_time", expectedExpirationTime)
	createMeta(root, "article:author", expectedAuthor1)
	createMeta(root, "article:author", expectedAuthor2)

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, expectedPublishedTime, article.PublishedTime)
	assert.Equal(t, expectedModifiedTime, article.ModifiedTime)
	assert.Equal(t, expectedExpirationTime, article.ExpirationTime)
	assert.Equal(t, expectedSection, article.Section)
	assert.Equal(t, 2, len(article.Authors))
	assert.Equal(t, expectedAuthor1, article.Authors[0])
	assert.Equal(t, expectedAuthor2, article.Authors[1])
}

func Test_OpenGraph_OGAndProfilePrefixesInHtmlTag(t *testing.T) {
	// Set prefix attribute in HTML tag.
	root := testutil.CreateHTML()
	dom.SetAttribute(root, "prefix", "tstog: http://ogp.me/ns# tstpf: http://ogp.me/ns/profile#")

	// Create the required properties and description.
	expectedTitle := "Testing customized OG and profile prefixes"
	expectedURL := "http://test/url.html"
	expectedImage := "http://test/image.jpeg"
	expectedDescr := "This tests the use of customized OG and profile prefixes"
	expectedSiteName := "Google"
	createMeta(root, "tstog:title", expectedTitle)
	createMeta(root, "tstog:type", "profile")
	createMeta(root, "tstog:url", expectedURL)
	createMeta(root, "tstog:image", expectedImage)
	createMeta(root, "tstog:description", expectedDescr)
	createMeta(root, "tstog:site_name", expectedSiteName)

	// Create the image structure.
	expectedSecureURL := "https://test/image.jpeg"
	expectedImageType := "image/jpeg"
	createMeta(root, "tstog:image:url", expectedImage)
	createMeta(root, "tstog:image:secure_url", expectedSecureURL)
	createMeta(root, "tstog:image:type", expectedImageType)
	createMeta(root, "tstog:image:width", "600")
	createMeta(root, "tstog:image:height", "400")

	// Create the "profile" object.
	createMeta(root, "tstpf:first_name", "Jane")
	createMeta(root, "tstpf:last_name", "Doe")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	images := parser.Images()
	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, expectedURL, parser.URL())
	assert.Equal(t, expectedDescr, parser.Description())
	assert.Equal(t, expectedSiteName, parser.Publisher())
	assert.Equal(t, 1, len(images))
	assert.Equal(t, expectedImage, images[0].URL)
	assert.Equal(t, expectedSecureURL, images[0].SecureURL)
	assert.Equal(t, expectedImageType, images[0].Type)
	assert.Equal(t, 600, images[0].Width)
	assert.Equal(t, 400, images[0].Height)
	assert.Equal(t, "Jane Doe", parser.Author())
}

func Test_OpenGraph_ArticlePrefixInHeadTag(t *testing.T) {
	// Set prefix attribute in head tag.
	root := testutil.CreateHTML()
	headNode := dom.QuerySelector(root, "head")
	dom.SetAttribute(headNode, "prefix", "tstog: http://ogp.me/ns# tsta: http://ogp.me/ns/article#")

	// Create the required properties.
	createCustomizedTitle(root)
	createCustomizedURL(root)
	createCustomizedImage(root)
	createMeta(root, "tstog:type", "article")

	// Create the "article" object.
	expectedSection := "GWT Testing"
	expectedPublishedTime := "2014-04-01T01:23:59Z"
	expectedModifiedTime := "2014-04-02T02:23:59Z"
	expectedExpirationTime := "2014-04-03T03:23:59Z"
	expectedAuthor1 := "http://blah/author1.html"
	expectedAuthor2 := "http://blah/author2.html"
	createMeta(root, "tsta:section", expectedSection)
	createMeta(root, "tsta:published_time", expectedPublishedTime)
	createMeta(root, "tsta:modified_time", expectedModifiedTime)
	createMeta(root, "tsta:expiration_time", expectedExpirationTime)
	createMeta(root, "tsta:author", expectedAuthor1)
	createMeta(root, "tsta:author", expectedAuthor2)

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, "Article", parser.Type())

	article := parser.Article()
	assert.Equal(t, expectedPublishedTime, article.PublishedTime)
	assert.Equal(t, expectedModifiedTime, article.ModifiedTime)
	assert.Equal(t, expectedExpirationTime, article.ExpirationTime)
	assert.Equal(t, expectedSection, article.Section)
	assert.Equal(t, 2, len(article.Authors))
	assert.Equal(t, expectedAuthor1, article.Authors[0])
	assert.Equal(t, expectedAuthor2, article.Authors[1])
}

func Test_OpenGraph_IncorrectPrefix(t *testing.T) {
	// Set prefix attribute in HTML tag.
	root := testutil.CreateHTML()
	dom.SetAttribute(root, "prefix", "tstog: http://ogp.me/ns#")

	// Create the required properties.
	createCustomizedTitle(root)
	createCustomizedType(root)
	createCustomizedURL(root)
	createCustomizedImage(root)

	// Create the description property with the common "og" prefix, instead
	// of the customized "tstog" prefix.
	createMeta(root, "og:description", "this description should be ignored")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, "", parser.Description())
}

func Test_OpenGraph_OGAndProfileXmlns(t *testing.T) {
	// Set xmlns attribute in HTML tag.
	root := testutil.CreateHTML()
	dom.SetAttribute(root, "xmlns:tstog", "http://ogp.me/ns#")
	dom.SetAttribute(root, "xmlns:tstpf", "http://ogp.me/ns/profile#")

	// Create the required properties and description.
	expectedTitle := "Testing customized OG and profile xmlns"
	expectedURL := "http://test/url.html"
	expectedImage := "http://test/image.jpeg"
	expectedDescr := "This tests the use of customized OG and profile xmlns"
	expectedSiteName := "Google"
	createMeta(root, "tstog:title", expectedTitle)
	createMeta(root, "tstog:type", "profile")
	createMeta(root, "tstog:url", expectedURL)
	createMeta(root, "tstog:image", expectedImage)
	createMeta(root, "tstog:description", expectedDescr)
	createMeta(root, "tstog:site_name", expectedSiteName)

	// Create the image structure.
	expectedSecureURL := "https://test/image.jpeg"
	expectedImageType := "image/jpeg"
	createMeta(root, "tstog:image:url", expectedImage)
	createMeta(root, "tstog:image:secure_url", expectedSecureURL)
	createMeta(root, "tstog:image:type", expectedImageType)
	createMeta(root, "tstog:image:width", "600")
	createMeta(root, "tstog:image:height", "400")

	// Create the "profile" object.
	createMeta(root, "tstpf:first_name", "Jane")
	createMeta(root, "tstpf:last_name", "Doe")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)

	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, expectedURL, parser.URL())
	assert.Equal(t, expectedDescr, parser.Description())
	assert.Equal(t, expectedSiteName, parser.Publisher())
	assert.Equal(t, "Jane Doe", parser.Author())

	images := parser.Images()
	assert.Equal(t, 1, len(images))

	image := images[0]
	assert.Equal(t, expectedImage, image.URL)
	assert.Equal(t, expectedSecureURL, image.SecureURL)
	assert.Equal(t, expectedImageType, image.Type)
	assert.Equal(t, 600, image.Width)
	assert.Equal(t, 400, image.Height)
}

func Test_OpenGraph_ArticleXmlns(t *testing.T) {
	// Set xmlns attribute in HTML tag.
	root := testutil.CreateHTML()
	dom.SetAttribute(root, "xmlns:tstog", "http://ogp.me/ns#")
	dom.SetAttribute(root, "xmlns:tsta", "http://ogp.me/ns/article#")

	// Create the required properties.
	createCustomizedTitle(root)
	createCustomizedURL(root)
	createCustomizedImage(root)
	createMeta(root, "tstog:type", "article")

	// Create the "article" object.
	expectedSection := "GWT Testing"
	expectedPublishedTime := "2014-04-01T01:23:59Z"
	expectedModifiedTime := "2014-04-02T02:23:59Z"
	expectedExpirationTime := "2014-04-03T03:23:59Z"
	expectedAuthor1 := "http://blah/author1.html"
	expectedAuthor2 := "http://blah/author2.html"
	createMeta(root, "tsta:section", expectedSection)
	createMeta(root, "tsta:published_time", expectedPublishedTime)
	createMeta(root, "tsta:modified_time", expectedModifiedTime)
	createMeta(root, "tsta:expiration_time", expectedExpirationTime)
	createMeta(root, "tsta:author", expectedAuthor1)
	createMeta(root, "tsta:author", expectedAuthor2)

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, "Article", parser.Type())

	article := parser.Article()
	assert.Equal(t, expectedPublishedTime, article.PublishedTime)
	assert.Equal(t, expectedModifiedTime, article.ModifiedTime)
	assert.Equal(t, expectedExpirationTime, article.ExpirationTime)
	assert.Equal(t, expectedSection, article.Section)
	assert.Equal(t, 2, len(article.Authors))
	assert.Equal(t, expectedAuthor1, article.Authors[0])
	assert.Equal(t, expectedAuthor2, article.Authors[1])
}

func Test_OpenGraph_IncorrectXmlns(t *testing.T) {
	// Set prefix attribute in HTML tag.
	root := testutil.CreateHTML()
	dom.SetAttribute(root, "xmlns:tstog", "http://ogp.me/ns#")

	// Create the required properties.
	createCustomizedTitle(root)
	createCustomizedType(root)
	createCustomizedURL(root)
	createCustomizedImage(root)

	// Create the description property with the common "og" prefix, instead
	// of the customized "tstog" prefix.
	createMeta(root, "og:description", "this description should be ignored")

	parser, _ := opengraph.NewParser(root, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, "", parser.Description())
}

func createDefaultTitle(root *html.Node) {
	createMeta(root, "og:title", "dummy title")
}

func createCustomizedTitle(root *html.Node) {
	createMeta(root, "tstog:title", "dummy title")
}

func createDefaultType(root *html.Node) {
	createMeta(root, "og:type", "website")
}

func createCustomizedType(root *html.Node) {
	createMeta(root, "tstog:type", "website")
}

func createDefaultURL(root *html.Node) {
	createMeta(root, "og:url", "http://dummy/url.html")
}

func createCustomizedURL(root *html.Node) {
	createMeta(root, "tstog:url", "http://dummy/url.html")
}

func createDefaultImage(root *html.Node) {
	createMeta(root, "og:image", "http://dummy/image.jpeg")
}

func createCustomizedImage(root *html.Node) {
	createMeta(root, "tstog:image", "http://dummy/image.jpeg")
}

func createDescription(root *html.Node, description, prefix string) {
	createMeta(root, prefix+":description", description)
}

func createMeta(root *html.Node, property, content string) {
	head := dom.QuerySelector(root, "head")
	if head == nil {
		return
	}

	meta := testutil.CreateMetaProperty(property, content)
	dom.AppendChild(head, meta)
}
