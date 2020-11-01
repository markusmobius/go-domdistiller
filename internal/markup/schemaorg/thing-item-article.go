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

import (
	"github.com/markusmobius/go-domdistiller/data"
	"golang.org/x/net/html"
)

type ArticleItem struct {
	BaseThingItem
}

func NewArticleItem(element *html.Node) *ArticleItem {
	item := &ArticleItem{}
	item.init(Article, element)

	item.addStringPropertyName(HeadlineProp)
	item.addStringPropertyName(PublisherProp)
	item.addStringPropertyName(CopyrightHolderProp)
	item.addStringPropertyName(CopyrightYearProp)
	item.addStringPropertyName(DateModifiedProp)
	item.addStringPropertyName(DatePublishedProp)
	item.addStringPropertyName(AuthorProp)
	item.addStringPropertyName(CreatorProp)
	item.addStringPropertyName(SectionProp)

	item.addItemPropertyName(PublisherProp)
	item.addItemPropertyName(CopyrightHolderProp)
	item.addItemPropertyName(AuthorProp)
	item.addItemPropertyName(CreatorProp)
	item.addItemPropertyName(AssociatedMediaProp)
	item.addItemPropertyName(EncodingProp)

	return item
}

func (ai *ArticleItem) getArticle() *data.MarkupArticle {
	author := ai.getPersonOrOrganizationName(AuthorProp)
	if author == "" {
		author = ai.getPersonOrOrganizationName(CreatorProp)
	}

	var authors []string
	if author != "" {
		authors = []string{author}
	}

	return &data.MarkupArticle{
		PublishedTime: ai.getStringProperty(DatePublishedProp),
		ModifiedTime:  ai.getStringProperty(DateModifiedProp),
		Section:       ai.getStringProperty(SectionProp),
		Authors:       authors,
	}
}

func (ai *ArticleItem) getCopyright() string {
	copyright := ai.getStringProperty(CopyrightYearProp)
	copyrightHolder := ai.getPersonOrOrganizationName(CopyrightHolderProp)
	if copyright != "" && copyrightHolder != "" {
		copyright += " "
	}
	copyright += copyrightHolder

	if copyright != "" {
		return "Copyright " + copyright
	}
	return ""
}

func (ai *ArticleItem) getPersonOrOrganizationName(propertyName string) string {
	// Returns either the string value of `propertyName` or the value
	// returned by getName() of PersonItem or OrganizationItem.
	value := ai.getStringProperty(propertyName)
	if value != "" {
		return value
	}

	valueItem := ai.getItemProperty(propertyName)
	if valueItem != nil {
		switch valueItem.getType() {
		case Person:
			if personItem, ok := valueItem.(*PersonItem); ok {
				return personItem.getName()
			}

		case Organization:
			if orgItem, ok := valueItem.(*OrganizationItem); ok {
				return orgItem.getName()
			}
		}
	}

	return ""
}

func (ai *ArticleItem) getRepresentativeImageItem() *ImageItem {
	// Returns the corresponding ImageItem for "associatedMedia" or "encoding" property.
	item := ai.getItemProperty(AssociatedMediaProp)
	if item == nil {
		item = ai.getItemProperty(EncodingProp)
	}

	if item != nil && item.getType() == Image {
		if imageItem, ok := item.(*ImageItem); ok {
			return imageItem
		}
	}

	return nil
}

func (ai *ArticleItem) getImage() *data.MarkupImage {
	// Use value of "image" property to create a MarkupParser.Image.
	imageURL := ai.getStringProperty(ImageProp)
	if imageURL == "" {
		return nil
	}

	return &data.MarkupImage{URL: imageURL}
}
