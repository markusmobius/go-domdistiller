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

const imageBase64 = "data:image/png;base64,iVBORw0KGgo" +
	"AAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/" +
	"w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="

const shortImageBase64 = "data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_Embed_Image_WithSettingDimension
// - Test_Embed_Image_WithHeightCSS
// - Test_Embed_Image_WithWidthHeightCSSPx
// - Test_Embed_Image_WithWidthAttributeHeightCSSPx
// - Test_Embed_Image_WithWidthAttributeHeightCSS
// - Test_Embed_Image_WithAttributeCSS
// - Test_Embed_Image_WithAttributesCSSHeightCMAndWidthAttrb
// - Test_Embed_Image_WithAttributesCSSHeightCM

func Test_Embed_Image_HasWidthHeightAttributes(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", imageBase64)
	dom.SetAttribute(img, "width", "32")
	dom.SetAttribute(img, "height", "32")

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(img)).(*webdoc.Image)

	assert.NotNil(t, result)
	assert.Equal(t, 32, result.Width)
	assert.Equal(t, 32, result.Height)
}

func Test_Embed_Image_HasNoAttributes(t *testing.T) {
	img := dom.CreateElement("img")

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(img)).(*webdoc.Image)

	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Width)
	assert.Equal(t, 0, result.Height)
}

func Test_Embed_Image_HasWidthAttribute(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", imageBase64)
	dom.SetAttribute(img, "width", "32")

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(img)).(*webdoc.Image)

	assert.NotNil(t, result)
	assert.Equal(t, 32, result.Width)
	assert.Equal(t, 0, result.Height)
}

func Test_Embed_Image_LazyLoadedImage(t *testing.T) {
	// Common lazy attributes
	extractLazyLoadedImage(t, "data-src", "src")
	extractLazyLoadedImage(t, "datasrc", "src")
	extractLazyLoadedImage(t, "data-original", "src")
	extractLazyLoadedImage(t, "data-url", "src")
	extractLazyLoadedImage(t, "data-srcset", "srcset")
	extractLazyLoadedImage(t, "datasrcset", "srcset")

	// Custom lazy attributes
	extractLazyLoadedImage(t, "lazy-src", "src")
	extractLazyLoadedImage(t, "lazysrc", "src")
	extractLazyLoadedImage(t, "lazy-srcset", "srcset")
	extractLazyLoadedImage(t, "lazysrcset", "srcset")

	// Image with small base64 image source
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", shortImageBase64)
	dom.SetAttribute(img, "lazy-srcset", "image.png 1x")

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL, nil)

	result, _ := (extractor.Extract(img)).(*webdoc.Image)
	assert.NotNil(t, result)
	assert.Equal(t, `<img srcset="http://example.com/image.png 1x"/>`, result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image.png"}, result.GetURLs())
}

func Test_Embed_Image_LazyLoadedFigure(t *testing.T) {
	// Common lazy attributes
	extractLazyLoadedFigure(t, "data-src", "src")
	extractLazyLoadedFigure(t, "datasrc", "src")
	extractLazyLoadedFigure(t, "data-original", "src")
	extractLazyLoadedFigure(t, "data-url", "src")
	extractLazyLoadedFigure(t, "data-srcset", "srcset")
	extractLazyLoadedFigure(t, "datasrcset", "srcset")

	// Custom lazy attributes
	extractLazyLoadedFigure(t, "lazy-src", "src")
	extractLazyLoadedFigure(t, "lazysrc", "src")
	extractLazyLoadedFigure(t, "lazy-srcset", "srcset")
	extractLazyLoadedFigure(t, "lazysrcset", "srcset")

	// Figure with image with small base64 image source
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", shortImageBase64)
	dom.SetAttribute(img, "lazy-srcset", "image.png 1x")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL, nil)
	expected := `<figure><img srcset="http://example.com/image.png 1x"/></figure>`

	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image.png"}, result.GetURLs())

	// Figure with noscript image
	lazyImg := dom.CreateElement("img")
	dom.SetAttribute(lazyImg, "src", "image-lq.png")

	realImg := dom.CreateElement("img")
	dom.SetAttribute(realImg, "data-src", "image-hq.png")

	noscript := dom.CreateElement("noscript")
	dom.SetInnerHTML(noscript, dom.OuterHTML(realImg))

	figure = dom.CreateElement("figure")
	dom.AppendChild(figure, lazyImg)
	dom.AppendChild(figure, noscript)

	extractor = embed.NewImageExtractor(pageURL, nil)
	expected = `<figure><img src="http://example.com/image-hq.png"/></figure>`

	result, _ = (extractor.Extract(figure)).(*webdoc.Figure)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image-hq.png"}, result.GetURLs())
}

func Test_Embed_Image_FigureWithoutCaptionWithNoscript(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	noscript := dom.CreateElement("noscript")
	dom.SetInnerHTML(noscript, "<span>text</span>")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, noscript)

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	// In original dom-distiller the text is not included because by default
	// noscript is hidden element so inner text wouldn't capture it. However
	// we works in server side so we will catch it as well.
	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>text</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, 100, result.Width)
	assert.Equal(t, 100, result.Height)
	assert.Equal(t, expected, result.GenerateOutput(false))
}

func Test_Embed_Image_FigureWithoutImageAndCaption(t *testing.T) {
	figure := dom.CreateElement("figure")

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	assert.Nil(t, result)
}

func Test_Embed_Image_FigureWithPictureWithoutImg(t *testing.T) {
	source := dom.CreateElement("source")
	dom.SetAttribute(source, "srcset", "http://www.example.com/image-240-200.jpg")
	dom.SetAttribute(source, "media", "(min-width: 800px)")

	picture := dom.CreateElement("picture")
	dom.AppendChild(picture, source)

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, picture)

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	expected := `<figure><picture>` +
		`<source srcset="http://www.example.com/image-240-200.jpg" media="(min-width: 800px)"/>` +
		`<img srcset="http://www.example.com/image-240-200.jpg"/>` +
		`</picture></figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
}

func Test_Embed_Image_FigureWithPictureWithoutSourceAndImg(t *testing.T) {
	picture := dom.CreateElement("picture")
	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, picture)

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)
	expected := `<figure><picture></picture></figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
}

func Test_Embed_Image_FigureCaptionTextOnly(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	figcaption := dom.CreateElement("figcaption")
	dom.SetTextContent(figcaption, "This is a caption")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, figcaption)

	extractor := embed.NewImageExtractor(nil, nil)
	result := extractor.Extract(figure)

	assert.NotNil(t, result)
	assert.Equal(t, "This is a caption", result.GenerateOutput(true))
}

func Test_Embed_Image_FigureCaptionWithAnchor(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	anchor := dom.CreateElement("a")
	dom.SetAttribute(anchor, "href", "link_page.html")
	dom.SetInnerHTML(anchor, "caption<br>link")

	figcaption := dom.CreateElement("figcaption")
	dom.AppendChild(figcaption, dom.CreateTextNode("This is a "))
	dom.AppendChild(figcaption, anchor)

	figure := dom.CreateElement("figure")
	dom.SetAttribute(figure, "attribute", "value")
	dom.SetAttribute(figure, "class", "test")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, figcaption)

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL, nil)
	result := extractor.Extract(figure)

	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>` +
		`This is a <a href="http://example.com/link_page.html">caption<br/>link</a>` +
		`</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, "This is a caption\nlink", result.GenerateOutput(true))
}

func Test_Embed_Image_FigureCaptionWithoutAnchor(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	figcaption := dom.CreateElement("figcaption")
	dom.SetInnerHTML(figcaption, "<div><span>This is a caption</span><a></a></div>")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, figcaption)

	extractor := embed.NewImageExtractor(nil, nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>This is a caption</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, 100, result.Width)
	assert.Equal(t, 100, result.Height)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, "This is a caption", result.GenerateOutput(true))
}

func Test_Embed_Image_FigureDivCaption(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	div := dom.CreateElement("div")
	dom.SetTextContent(div, "This is a caption")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, div)

	extractor := embed.NewImageExtractor(nil, nil)
	result := extractor.Extract(figure)

	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>This is a caption</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, "This is a caption", result.GenerateOutput(true))
}

func extractLazyLoadedImage(t *testing.T, lazyAttr, expectedAttr string) {
	// Prepare test image
	testURL := "image.png"
	if expectedAttr == "srcset" {
		testURL += " 1x"
	}

	testImg := dom.CreateElement("img")
	dom.SetAttribute(testImg, lazyAttr, testURL)

	// Prepare expected image
	expectedURL := "http://example.com/image.png"
	if expectedAttr == "srcset" {
		expectedURL += " 1x"
	}

	expectedImg := dom.CreateElement("img")
	dom.SetAttribute(expectedImg, expectedAttr, expectedURL)

	// Test extractor
	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL, nil)

	result, _ := (extractor.Extract(testImg)).(*webdoc.Image)
	assert.NotNil(t, result)
	assert.Equal(t, dom.OuterHTML(expectedImg), result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image.png"}, result.GetURLs())
}

func extractLazyLoadedFigure(t *testing.T, lazyAttr, expectedAttr string) {
	// Prepare test figure
	testURL := "image.png"
	if expectedAttr == "srcset" {
		testURL += " 1x"
	}

	testImg := dom.CreateElement("img")
	dom.SetAttribute(testImg, lazyAttr, testURL)

	testFigure := dom.CreateElement("figure")
	dom.AppendChild(testFigure, testImg)

	// Prepare expected image
	expectedURL := "http://example.com/image.png"
	if expectedAttr == "srcset" {
		expectedURL += " 1x"
	}

	expectedImg := dom.CreateElement("img")
	dom.SetAttribute(expectedImg, expectedAttr, expectedURL)

	expectedFigure := dom.CreateElement("figure")
	dom.AppendChild(expectedFigure, expectedImg)

	// Test extractor
	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL, nil)

	result, _ := (extractor.Extract(testFigure)).(*webdoc.Figure)
	assert.NotNil(t, result)
	assert.Equal(t, dom.OuterHTML(expectedFigure), result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image.png"}, result.GetURLs())
}
