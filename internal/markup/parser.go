// ORIGINAL: java/MarkupParser.java

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

package markup

import (
	"time"

	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/markup/iereader"
	"github.com/markusmobius/go-domdistiller/internal/markup/opengraph"
	"github.com/markusmobius/go-domdistiller/internal/markup/schemaorg"
	"golang.org/x/net/html"
)

// Parser loads the different parsers that are based on different markup specifications, and
// allows retrieval of different distillation-related markup properties from a document. It retrieves
// the requested properties from one or more parsers.  If necessary, it may merge the information
// from multiple parsers.
//
// Currently, three markup format are supported: OpenGraphProtocol, IEReadingView and SchemaOrg.
// For now, OpenGraphProtocolParser takes precedence because it uses specific meta tags and hence
// extracts information the fastest; it also demands conformance to rules. If the rules are broken
// or the properties retrieved are null or empty, we try with SchemaOrg then IEReadingView.
//
// The properties that matter to distilled content are:
// - individual properties: title, page type, page url, description, publisher, author, copyright
// - dominant and inline images and their properties: url, secure_url, type, caption, width, height
// - article and its properties: section name, published time, modified time, expiration time,
//   authors.
//
// TODO: for some properties, e.g. dominant and inline images, we might want to retrieve from
// multiple parsers; IEReadingViewParser provides more information as it scans all images in the
// document.  If we do so, we would need to merge the multiple versions in a meaningful way.
type Parser struct {
	accessors []Accessor
}

func NewParser(root *html.Node, timingInfo *data.TimingInfo) *Parser {
	// Initiate parser
	ps := &Parser{}
	ps.accessors = make([]Accessor, 0)

	// Add accessors
	start := time.Now()
	ogParser, err := opengraph.NewParser(root, timingInfo)
	if err == nil && ogParser != nil {
		ps.accessors = append(ps.accessors, ogParser)
	}
	timingInfo.AddEntry(start, "OpenGraphProtocolParser")

	start = time.Now()
	ps.accessors = append(ps.accessors, schemaorg.NewParser(root, timingInfo))
	timingInfo.AddEntry(start, "SchemaOrgParserAccessor")

	start = time.Now()
	// TODO: Use eager evaluation in IEReadingViewParser, but only for profiling.
	ps.accessors = append(ps.accessors, iereader.NewParser(root))
	timingInfo.AddEntry(start, "SchemaOrgParserAccessor")

	return ps
}

func (ps *Parser) Title() string {
	for _, accessor := range ps.accessors {
		if title := accessor.Title(); title != "" {
			return title
		}
	}
	return ""
}

func (ps *Parser) Type() string {
	for _, accessor := range ps.accessors {
		tp := accessor.Type()
		if tp != "" {
			return tp
		}
	}
	return ""
}

func (ps *Parser) URL() string {
	for _, accessor := range ps.accessors {
		if url := accessor.URL(); url != "" {
			return url
		}
	}
	return ""
}

func (ps *Parser) Images() []data.MarkupImage {
	for _, accessor := range ps.accessors {
		if images := accessor.Images(); len(images) > 0 {
			return images
		}
	}
	return nil
}

func (ps *Parser) Description() string {
	for _, accessor := range ps.accessors {
		if description := accessor.Description(); description != "" {
			return description
		}
	}
	return ""
}

func (ps *Parser) Publisher() string {
	for _, accessor := range ps.accessors {
		if publisher := accessor.Publisher(); publisher != "" {
			return publisher
		}
	}
	return ""
}

func (ps *Parser) Copyright() string {
	for _, accessor := range ps.accessors {
		if copyright := accessor.Copyright(); copyright != "" {
			return copyright
		}
	}
	return ""
}

func (ps *Parser) Author() string {
	for _, accessor := range ps.accessors {
		if author := accessor.Author(); author != "" {
			return author
		}
	}
	return ""
}

func (ps *Parser) Article() *data.MarkupArticle {
	for _, accessor := range ps.accessors {
		if article := accessor.Article(); article != nil {
			return article
		}
	}
	return nil
}

func (ps *Parser) OptOut() bool {
	for _, accessor := range ps.accessors {
		if optOut := accessor.OptOut(); optOut {
			return true
		}
	}
	return false
}

func (ps *Parser) MarkupInfo() data.MarkupInfo {
	if ps.OptOut() {
		return data.MarkupInfo{}
	}

	info := data.MarkupInfo{
		Title:       ps.Title(),
		Type:        ps.Type(),
		URL:         ps.URL(),
		Description: ps.Description(),
		Publisher:   ps.Publisher(),
		Copyright:   ps.Copyright(),
		Author:      ps.Author(),
	}

	article := ps.Article()
	if article != nil {
		info.Article = data.MarkupArticle{
			PublishedTime:  article.PublishedTime,
			ModifiedTime:   article.ModifiedTime,
			ExpirationTime: article.ExpirationTime,
			Section:        article.Section,
			Authors:        append([]string{}, article.Authors...),
		}
	}

	for _, image := range ps.Images() {
		info.Images = append(info.Images, data.MarkupImage{
			URL:       image.URL,
			SecureURL: image.SecureURL,
			Type:      image.Type,
			Caption:   image.Caption,
			Width:     image.Width,
			Height:    image.Height,
		})
	}

	return info
}
