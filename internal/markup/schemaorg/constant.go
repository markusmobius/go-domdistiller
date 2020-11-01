// ORIGINAL: java/SchemaOrgParser.java

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

package schemaorg

const (
	NameProp            = "name"
	URLProp             = "url"
	DescriptionProp     = "description"
	ImageProp           = "image"
	HeadlineProp        = "headline"
	PublisherProp       = "publisher"
	CopyrightHolderProp = "copyrightHolder"
	CopyrightYearProp   = "copyrightYear"
	ContentURLProp      = "contentUrl"
	EncodingFormatProp  = "encodingFormat"
	CaptionProp         = "caption"
	RepresentativeProp  = "representativeOfPage"
	WidthProp           = "width"
	HeightProp          = "height"
	DatePublishedProp   = "datePublished"
	DateModifiedProp    = "dateModified"
	AuthorProp          = "author"
	CreatorProp         = "creator"
	SectionProp         = "articleSection"
	AssociatedMediaProp = "associatedMedia"
	EncodingProp        = "encoding"
	FamilyNameProp      = "familyName"
	GivenNameProp       = "givenName"
	LegalNameProp       = "legalName"
	AuthorRel           = "author"
)

type SchemaType uint

const (
	Unsupported SchemaType = iota
	Image
	Article
	Person
	Organization
)

var schemaTypeURLs = map[string]SchemaType{
	"http://schema.org/ImageObject":             Image,
	"http://schema.org/Article":                 Article,
	"http://schema.org/BlogPosting":             Article,
	"http://schema.org/NewsArticle":             Article,
	"http://schema.org/ScholarlyArticle":        Article,
	"http://schema.org/TechArticle":             Article,
	"http://schema.org/Person":                  Person,
	"http://schema.org/Organization":            Organization,
	"http://schema.org/Corporation":             Organization,
	"http://schema.org/EducationalOrganization": Organization,
	"http://schema.org/GovernmentOrganization":  Organization,
	"http://schema.org/NGO":                     Organization,
}

// The key for `tagAttributeMap` is the tag name, while the entry value is an
// array of attributes in the specified tag from which to extract information:
// - 0th attribute: contains the value for the property specified in itemprop
// - 1st attribute: if available, contains the value for the author property.
var tagAttributeMap = map[string]string{
	"img":    "src",
	"audio":  "src",
	"embed":  "src",
	"iframe": "src",
	"source": "src",
	"track":  "src",
	"video":  "src",
	"a":      "href",
	"link":   "href",
	"area":   "href",
	"meta":   "content",
	"time":   "datetime",
	"object": "data",
	"data":   "value",
	"meter":  "value",
}
