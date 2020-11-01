// ORIGINAL: java/IEReadingViewParser.java

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

package iereader

import (
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

// Parser recognizes and parses the IE Reading View markup tags according to
// http://ie.microsoft.com/testdrive/browser/readingview, and returns the properties that matter to
// distilled content - title, date, author, publisher, dominant image, inline images, caption,
// copyright, and opt-out. Some properties require the scanning and parsing of a lot of nodes, so
// each property is scanned for and verified for legitimacy lazily, i.e. only upon request.
// It implements markup.Accessor.
type Parser struct {
	root            *html.Node
	allMeta         []*html.Node
	determinedProps map[string]struct{}

	title     string
	date      string
	author    string
	publisher string
	copyright string
	optOut    bool
	images    []data.MarkupImage
}

func NewParser(root *html.Node) *Parser {
	return &Parser{
		root:            root,
		allMeta:         dom.GetElementsByTagName(root, "meta"),
		determinedProps: make(map[string]struct{}),
	}
}

func (p *Parser) Title() string {
	if _, determined := p.determinedProps["title"]; !determined {
		p.findTitle()
	}
	return p.title
}

func (p *Parser) Type() string {
	return ""
}

func (p *Parser) URL() string {
	return ""
}

func (p *Parser) Images() []data.MarkupImage {
	if _, determined := p.determinedProps["images"]; !determined {
		p.findImages()
	}
	return p.images
}

func (p *Parser) Description() string {
	return ""
}

func (p *Parser) Publisher() string {
	// The primary indicator is the OpenGraph Protocol "site_name" meta tag, which is
	// handled by opengraph.Parser. The secondary indicator is any html tag with the
	// "publisher" or "source_organization" attribute.
	if _, determined := p.determinedProps["publisher"]; !determined {
		p.findPublisher()
	}
	return p.publisher
}

func (p *Parser) Copyright() string {
	if _, determined := p.determinedProps["copyright"]; !determined {
		p.findCopyright()
	}
	return p.copyright
}

func (p *Parser) Author() string {
	if _, determined := p.determinedProps["author"]; !determined {
		p.findAuthor()
	}
	return p.author
}

func (p *Parser) Article() *data.MarkupArticle {
	if _, determined := p.determinedProps["date"]; !determined {
		p.findDate()
	}

	author := p.Author()
	article := &data.MarkupArticle{}
	article.PublishedTime = p.date
	if author != "" {
		article.Authors = []string{author}
	}

	return article
}

func (p *Parser) OptOut() bool {
	if _, determined := p.determinedProps["optout"]; !determined {
		p.findOptOut()
	}
	return p.optOut
}

func (p *Parser) findTitle() {
	// Mark this property as determined
	p.determinedProps["title"] = struct{}{}

	if len(p.allMeta) == 0 {
		return
	}

	// Make sure there's a <title> element.
	titles := dom.GetElementsByTagName(p.root, "title")
	if len(titles) == 0 {
		return
	}

	// Extract title text from meta tag with "title" as name.
	for _, meta := range p.allMeta {
		name := dom.GetAttribute(meta, "name")
		if strings.ToLower(name) == "title" {
			p.title = dom.GetAttribute(meta, "content")
			break
		}
	}
}

func (p *Parser) findImages() {
	// Mark this property as determined
	p.determinedProps["images"] = struct{}{}

	allImages := dom.GetElementsByTagName(p.root, "img")
	for _, img := range allImages {
		// As long as the image has a caption, it's relevant regardless of size;
		// otherwise, it's relevant if its size is good.
		caption := p.getImageCaption(img)
		if caption != "" || p.isImageRelevantBySize(img) {
			width, _ := strconv.Atoi(dom.GetAttribute(img, "width"))
			height, _ := strconv.Atoi(dom.GetAttribute(img, "height"))

			p.images = append(p.images, data.MarkupImage{
				URL:     dom.GetAttribute(img, "src"),
				Caption: caption,
				Width:   width,
				Height:  height,
			})
		}
	}
}

func (p *Parser) findPublisher() {
	// Mark this property as determined
	p.determinedProps["publisher"] = struct{}{}

	// Look for "publisher" or "source_organization" attribute in any html tag.
	allElems := dom.GetElementsByTagName(p.root, "*")
	for _, elem := range allElems {
		publisher := dom.GetAttribute(elem, "publisher")
		if publisher == "" {
			publisher = dom.GetAttribute(elem, "source_organization")
		}

		if publisher != "" {
			p.publisher = publisher
			break
		}
	}
}

func (p *Parser) findCopyright() {
	// Mark this property as determined
	p.determinedProps["copyright"] = struct{}{}

	// Get copyright from meta tag with "copyright" as name.
	for _, meta := range p.allMeta {
		name := dom.GetAttribute(meta, "name")
		if strings.ToLower(name) == "copyright" {
			p.copyright = dom.GetAttribute(meta, "content")
			break
		}
	}
}

func (p *Parser) findAuthor() {
	// Mark this property as determined
	p.determinedProps["author"] = struct{}{}

	// Get author from the first element that includes the "byline-name" class.
	// Note that we ignore the order of this element for now.
	elem := dom.QuerySelector(p.root, ".byline-name")
	if elem != nil {
		p.author = strings.TrimSpace(dom.TextContent(elem))
	}
}

func (p *Parser) findDate() {
	// Mark this property as determined
	p.determinedProps["date"] = struct{}{}

	// Get date from any element that includes the "dateline" class.
	elem := dom.QuerySelector(p.root, ".dateline")
	if elem != nil {
		p.date = strings.TrimSpace(dom.TextContent(elem))
		return
	}

	// Otherwise, get date from meta tag with "displaydate" as name.
	for _, meta := range p.allMeta {
		name := dom.GetAttribute(meta, "name")
		if strings.ToLower(name) == "displaydate" {
			p.date = dom.GetAttribute(meta, "content")
			break
		}
	}
}

func (p *Parser) findOptOut() {
	// Mark this property as determined
	p.determinedProps["optout"] = struct{}{}

	// Get optout from meta tag with "IE_RM_OFF" as name.
	for _, meta := range p.allMeta {
		name := dom.GetAttribute(meta, "name")
		if strings.ToUpper(name) == "IE_RM_OFF" {
			content := dom.GetAttribute(meta, "content")
			p.optOut = strings.ToLower(content) == "true"
			break
		}
	}
}

// isImageRelevantBySize specifies whether an image is relevant or not. The image is
// relevant if its width >= 400 and its aspect ratio between 1.3 and 3.0 inclusively.
// In real world this is done by computing stylesheet. However, since it's not possible
// to do that with Go, we just check width and height in image attribute.
// NEED-COMPUTE-CSS
func (p *Parser) isImageRelevantBySize(img *html.Node) bool {
	width, errWidth := strconv.Atoi(dom.GetAttribute(img, "width"))
	height, errHeight := strconv.Atoi(dom.GetAttribute(img, "height"))
	if errWidth != nil || errHeight != nil {
		return false
	}

	if width < 400 || height <= 0 {
		return false
	}

	aspectRatio := float64(width) / float64(height)
	return aspectRatio >= 1.3 && aspectRatio <= 3.0
}

func (p *Parser) getImageCaption(img *html.Node) string {
	// If <image> is a child of <figure>, then get the <figcaption> elements.
	if dom.TagName(img.Parent) != "figure" {
		return ""
	}

	caption := ""
	captionNodes := dom.GetElementsByTagName(img.Parent, "figcaption")
	if nCaption := len(captionNodes); nCaption > 0 && nCaption <= 2 {
		// Use innerText (instead of textContent) to get only visible captions.
		for _, captionNode := range captionNodes {
			caption = domutil.InnerText(captionNode)
			if caption != "" {
				break
			}
		}
	}

	return caption
}
