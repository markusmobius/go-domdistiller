// ORIGINAL: java/OpenGraphProtocolParser.java

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

package opengraph

const (
	TitleProp                 = "title"
	TypeProp                  = "type"
	ImageProp                 = "image"
	URLProp                   = "url"
	DescriptionProp           = "description"
	SiteNameProp              = "site_name"
	ImageStructPropPfx        = "image:"
	ImageURLProp              = "image:url"
	ImageSecureURLProp        = "image:secure_url"
	ImageTypeProp             = "image:type"
	ImageWidthProp            = "image:width"
	ImageHeightProp           = "image:height"
	ProfileFirstnameProp      = "first_name"
	ProfileLastnameProp       = "last_name"
	ArticleSectionProp        = "section"
	ArticlePublishedTimeProp  = "published_time"
	ArticleModifiedTimeProp   = "modified_time"
	ArticleExpirationTimeProp = "expiration_time"
	ArticleAuthorProp         = "author"
	ProfileObjtype            = "profile"
	ArticleObjtype            = "article"

	doPrefixFiltering = true
)

var importantProperties = []struct {
	Name   string
	Prefix Prefix
	Type   string
}{
	{TitleProp, OG, ""},
	{TypeProp, OG, ""},
	{URLProp, OG, ""},
	{DescriptionProp, OG, ""},
	{SiteNameProp, OG, ""},
	{ImageProp, OG, "image"},
	{ImageStructPropPfx, OG, "image"},
	{ProfileFirstnameProp, Profile, "profile"},
	{ProfileLastnameProp, Profile, "profile"},
	{ArticleSectionProp, Article, "article"},
	{ArticlePublishedTimeProp, Article, "article"},
	{ArticleModifiedTimeProp, Article, "article"},
	{ArticleExpirationTimeProp, Article, "article"},
	{ArticleAuthorProp, Article, "article"},
}
