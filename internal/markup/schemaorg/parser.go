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
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"golang.org/x/net/html"
)

// Parser recognizes and parses schema.org markup tags, and returns the properties that matter
// to distilled content. Schema.org markup (http://schema.org) is based on the microdata format
// (http://www.whatwg.org/specs/web-apps/current-work/multipage/microdata.html).
//
// For the basic Schema.org Thing type, the basic properties are: name, url, description, image.
// In addition, for each type that we support, we also parse more specific properties:
// - Article: headline (i.e. title), publisher, copyright year, copyright holder, date published,
//            date modified, author, article section
// - ImageObject: headline (i.e. title), publisher, copyright year, copyright holder, content url,
//                encoding format, caption, representative of page, width, height
// - Person: family name, given name
// - Organization: legal name.
//
// The value of a Schema.Org property can be a Schema.Org type, i.e. embedded. E.g., the author or
// publisher of article or publisher of image could be a Schema.Org Person or Organization type;
// in fact, this is the reason we support Person and Organization types.
type Parser struct {
	itemScopes    []ThingItem
	itemElement   map[*html.Node]ThingItem
	authorFromRel string
}

func NewParser(root *html.Node, timingInfo *data.TimingInfo) *Parser {
	// Initiate parser
	ps := &Parser{}
	ps.itemElement = make(map[*html.Node]ThingItem)

	start := time.Now()
	ps.parse(root)
	timingInfo.AddEntry(start, "SchemaOrg.parse")

	return ps
}

func (ps *Parser) parse(root *html.Node) {
	allProp := dom.QuerySelectorAll(root, "[itemprop],[itemscope]")

	// Root node (html) is not included in the result of querySelectorAll, so need to
	// handle it explicitly here.
	ps.parseElement(root, nil)
	for _, prop := range allProp {
		ps.parseElement(prop, ps.getItemScopeParent(prop))
	}

	// As per http://schema.org/author (or http://schema.org/Article and search for
	// "author" property), if <a> or <link> tags specify rel="author", extract it.
	allProp = dom.QuerySelectorAll(root, "a[rel=author],link[rel=author]")
	for _, prop := range allProp {
		if ps.authorFromRel == "" {
			ps.authorFromRel = ps.getAuthorFromRelAttribute(prop)
		}
	}
}

func (ps *Parser) parseElement(element *html.Node, parentItem ThingItem) {
	var newItem ThingItem
	var propertyNames []string
	if parentItem != nil {
		propertyNames = ps.getItemProp(element)
	}

	if ps.isItemScope(element) {
		// The "itemscope" and "itemtype" attributes of |e| indicate the start of an item.
		// Create the corresponding extended-ThingItem, and add it to the list if:
		// 1) its type is supported, and
		// 2) if the parent is an unsupported type, it's not an "itemprop" attribute of the
		//    parent, based on the rule that an item is a top-level item if its element doesn't
		//    have an itemprop attribute.
		newItem = ps.createItemForElement(element)
		if newItem != nil && newItem.isSupported() &&
			(parentItem == nil || parentItem.isSupported() || len(propertyNames) == 0) {
			ps.itemScopes = append(ps.itemScopes, newItem)
			ps.itemElement[element] = newItem
		}
	}

	// If parent is a supported type, parse the element for >= 1 properties in "itemprop"
	// attribute.
	if len(propertyNames) > 0 && (parentItem != nil && parentItem.isSupported()) &&
		(newItem == nil || newItem.isSupported()) {
		for _, prop := range propertyNames {
			if newItem != nil {
				// If a new item was created above, the property value of this "itemprop"
				// attribute is an embedded item, so add it to the parent item.
				parentItem.putItemValue(prop, newItem)
			} else {
				// Otherwise, extract the property value from the tag itself, and add
				// it to the parent item.
				parentItem.putStringValue(prop, ps.getPropertyValue(element))
			}
		}
	}
}

// getItemScopeParent is assumed the ItemScope parent of Element e is already parsed.
// For querySelectorAll(), parent nodes are guaranteed to appear before child nodes,
// so this assumption is met.
func (ps *Parser) getItemScopeParent(element *html.Node) ThingItem {
	parentElement := element
	for parentElement != nil {
		parentElement = parentElement.Parent
		if parentElement == nil {
			return nil
		}

		if ps.isItemScope(parentElement) {
			return ps.itemElement[parentElement]
		}
	}

	return nil
}

func (ps *Parser) createItemForElement(element *html.Node) ThingItem {
	switch ps.getItemType(element) {
	case Image:
		return NewImageItem(element)
	case Article:
		return NewArticleItem(element)
	case Person:
		return NewPersonItem(element)
	case Organization:
		return NewOrganizationItem(element)
	case Unsupported:
		return NewUnsupportedItem(element)
	default:
		return nil
	}
}

func (ps *Parser) isItemScope(element *html.Node) bool {
	return dom.HasAttribute(element, "itemscope") &&
		dom.HasAttribute(element, "itemtype")
}

func (ps *Parser) getItemProp(element *html.Node) []string {
	itemProp := dom.GetAttribute(element, "itemprop")
	if itemProp == "" {
		return nil
	}

	return strings.Fields(itemProp)
}

func (ps *Parser) getItemType(element *html.Node) SchemaType {
	schemaType := dom.GetAttribute(element, "itemtype")
	return schemaTypeURLs[schemaType]
}

// Extracts the property value from `element`. For some tags, the value
// is a specific attribute, while for others, it's the text between
// the start and end tags.
func (ps *Parser) getPropertyValue(element *html.Node) string {
	var value string
	tagName := dom.TagName(element)
	if attrName, exist := tagAttributeMap[tagName]; exist {
		value = dom.GetAttribute(element, attrName)
	}

	if value == "" {
		value = dom.TextContent(element)
		value = strings.TrimSpace(value)
	}

	return value
}

// Extracts the author property from the "rel=author" attribute of an
// anchor or a link element.
func (ps *Parser) getAuthorFromRelAttribute(element *html.Node) string {
	rel := dom.GetAttribute(element, "rel")
	tagName := strings.ToLower(dom.TagName(element))
	if (tagName == "a" || tagName == "link") && strings.ToLower(rel) == AuthorRel {
		author := dom.TextContent(element)
		return strings.TrimSpace(author)
	}

	return ""
}

func (ps *Parser) getArticleItems() []*ArticleItem {
	articles := []*ArticleItem{}
	for _, item := range ps.itemScopes {
		articleItem, ok := item.(*ArticleItem)
		if item.getType() == Article && ok {
			articles = append(articles, articleItem)
		}
	}
	return articles
}

func (ps *Parser) getImageItems() []*ImageItem {
	images := []*ImageItem{}
	for _, item := range ps.itemScopes {
		imageItem, ok := item.(*ImageItem)
		if item.getType() == Image && ok {
			images = append(images, imageItem)
		}
	}
	return images
}
