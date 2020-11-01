// ORIGINAL: java/OpenGraphProtocolParserAccessor.java

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

package opengraph

import (
	"strings"

	"github.com/markusmobius/go-domdistiller/data"
)

// Title returns the required "title" of the document.
func (ps *Parser) Title() string {
	return ps.propertyTable[TitleProp]
}

// Type returns the required "type" of the document if it's an
// article, empty string otherwise.
func (ps *Parser) Type() string {
	objType := ps.propertyTable[TypeProp]
	if strings.ToLower(objType) == ArticleObjtype {
		return "Article"
	}

	return ""
}

// URL returns the required "url" of the document.
func (ps *Parser) URL() string {
	return ps.propertyTable[URLProp]
}

// Images returns the structured properties of all "image"
// structures. Each "image" structure consists of image, image:url,
// image:secure_url, image:type, image:width, and image:height.
func (ps *Parser) Images() []data.MarkupImage {
	return ps.imageParser.ImageList
}

// Description returns the optional "description" of the document.
func (ps *Parser) Description() string {
	return ps.propertyTable[DescriptionProp]
}

// Publisher returns the optional "site_name" of the document.
func (ps *Parser) Publisher() string {
	return ps.propertyTable[SiteNameProp]
}

// Copyright returns empty since OpenGraph not support it.
func (ps *Parser) Copyright() string {
	return ""
}

// Author returns the concatenated first_name and last_name
// (delimited by a whitespace) of the "profile" object when
// value of "og:type" is "profile".
func (ps *Parser) Author() string {
	return ps.profileParser.GetFullName(ps.propertyTable)
}

// Article returns the properties of the "article" object when
// value of "og:type" is "article". The properties are published_time,
// modified_time and expiration_time, section, and a list of URLs
// to each author's profile.
func (ps *Parser) Article() *data.MarkupArticle {
	article := data.MarkupArticle{
		PublishedTime:  ps.propertyTable[ArticlePublishedTimeProp],
		ModifiedTime:   ps.propertyTable[ArticleModifiedTimeProp],
		ExpirationTime: ps.propertyTable[ArticleExpirationTimeProp],
		Section:        ps.propertyTable[ArticleSectionProp],
		Authors:        ps.articleParser.Authors,
	}

	if article.Section == "" &&
		article.PublishedTime == "" &&
		article.ModifiedTime == "" &&
		article.ExpirationTime == "" &&
		len(article.Authors) == 0 {
		return nil
	}

	return &article
}

// OptOut returns false since OpenGraph not support it. While
// this is not directly supported, the page owner can simply
// omit the required tags and init() will return a null
// OpenGraphProtocolParser.
func (ps *Parser) OptOut() bool {
	return false
}
