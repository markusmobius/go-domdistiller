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

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"golang.org/x/net/html"
)

var (
	rxOgpNsPrefix         = regexp.MustCompile(`(?i)((\w+):\s+(http:\/\/ogp.me\/ns(\/\w+)*#))\s*`)
	rxOgpNsNonPrefixName  = regexp.MustCompile(`(?i)^xmlns:(\w+)`)
	rxOgpNsNonPrefixValue = regexp.MustCompile(`(?i)^http:\/\/ogp.me\/ns(\/\w+)*#`)
)

// Parser recognizes and parses the Open Graph Protocol markup tags and returns the properties
// that matter to distilled content.
//
// First, it extracts the prefix and/or xmlns attributes from the HTML or HEAD tags to determine the
// prefixes that will be used for the protocol. If no prefix is specified, we fall back to the
// commonly used ones, e.g. "og". Then, it scans the OpenGraph Protocol <meta> elements that we
// care about, extracts their content, and stores them semantically, i.e. taking into consideration
// arrays, structures, and object types. Callers call get* to access these properties.
//
// The properties we care about are:
// - 4 required properties: title, type, image, url.
// - 2 optional properties: description, site_name.
// - image structured properties: image:url, image:secure_url, image:type, image:width, image:height
// - profile object properties: first_name, last_name
// - article object properties: section, published_time, modified_time, expiration_time, author;
//                              each author is a URL to the author's profile.
type Parser struct {
	prefixes      PrefixNameList
	propertyTable map[string]string
	imageParser   ImagePropParser
	profileParser ProfilePropParser
	articleParser ArticlePropParser
}

func NewParser(root *html.Node, timingInfo *data.TimingInfo) (*Parser, error) {
	// Initiate parser
	ps := &Parser{}
	ps.prefixes = make(PrefixNameList)
	ps.propertyTable = make(map[string]string)

	start := time.Now()
	timingInfo.AddEntry(start, "OpenGraphProtocolParser.parse")

	start = time.Now()
	ps.findPrefixes(root)
	timingInfo.AddEntry(start, "OpenGraphProtocolParser.findPrefixes")

	start = time.Now()
	ps.parseMetaTags(root)
	timingInfo.AddEntry(start, "OpenGraphProtocolParser.parseMetaTags")

	start = time.Now()
	ps.imageParser.Verify()
	timingInfo.AddEntry(start, "OpenGraphProtocolParser.imageParser.verify")

	prefix := ps.prefixes[OG] + ":"
	switch {
	case ps.propertyTable[TitleProp] == "":
		return nil, fmt.Errorf("required \"%s:title\" property is missing", prefix)

	case ps.propertyTable[TypeProp] == "":
		return nil, fmt.Errorf("required \"%s:type\" property is missing", prefix)

	case ps.propertyTable[URLProp] == "":
		return nil, fmt.Errorf("required \"%s:url\" property is missing", prefix)

	case len(ps.imageParser.ImageList) == 0:
		return nil, fmt.Errorf("required \"%s:image\" property is missing", prefix)
	}

	return ps, nil
}

func (ps *Parser) findPrefixes(root *html.Node) {
	strPrefixes := ""

	// See if HTML tag has "prefix" attribute.
	if dom.TagName(root) == "html" {
		strPrefixes = dom.GetAttribute(root, "prefix")
	}

	// Otherwise, see if HEAD tag has "prefix" attribute.
	if strPrefixes == "" {
		head := dom.QuerySelector(root, "head")
		if head != nil {
			strPrefixes = dom.GetAttribute(head, "prefix")
		}
	}

	// If there's "prefix" attribute, its value is something like
	// "og: http://ogp.me/ns# profile: http://ogp.me/ns/profile# article: http://ogp.me/ns/article#".
	if strPrefixes != "" {
		matches := rxOgpNsPrefix.FindAllStringSubmatch(strPrefixes, -1)
		for _, groups := range matches {
			ps.prefixes.addObjectType(groups[2], groups[4])
		}
	} else {
		// Still no "prefix" attribute, see if HTMl tag has "xmlns" attributes e.g.:
		// - "xmlns:og="http://ogp.me/ns#"
		// - "xmlns:profile="http://ogp.me/ns/profile#"
		// - "xmlns:article="http://ogp.me/ns/article#".
		for _, attr := range root.Attr {
			attrName := strings.ToLower(attr.Key)
			nameMatch := rxOgpNsNonPrefixName.FindStringSubmatch(attrName)
			if nameMatch == nil {
				continue
			}

			valueMatch := rxOgpNsNonPrefixValue.FindStringSubmatch(attr.Val)
			if valueMatch != nil {
				ps.prefixes.addObjectType(nameMatch[1], valueMatch[1])
			}
		}
	}

	ps.prefixes.setDefault()
}

func (ps *Parser) parseMetaTags(root *html.Node) {
	// Fetch meta nodes
	var metaNodes []*html.Node
	if doPrefixFiltering {
		// Attribute selectors with prefix
		// https://developer.mozilla.org/en-US/docs/Web/CSS/Attribute_selectors
		query := ""
		for _, prefix := range ps.prefixes {
			query += `meta[property^=` + prefix + `],`
		}

		query = strings.TrimSuffix(query, ",")
		metaNodes = dom.QuerySelectorAll(root, query)
	} else {
		metaNodes = dom.QuerySelectorAll(root, "meta[property]")
	}

	// Parse property
	for _, meta := range metaNodes {
		content := dom.GetAttribute(meta, "content")
		property := dom.GetAttribute(meta, "property")
		property = strings.ToLower(property)

		// Only store properties that we care about for distillation.
		for _, importantProperty := range importantProperties {
			prefixWithColon := ps.prefixes[importantProperty.Prefix] + ":"

			// Note that `==` won't work here because importantProperties uses "image:"
			// (ImageStructPropPfx) for all image structured properties, so as to prevent
			// repetitive property name comparison - here and then again in ImageParser.
			if !strings.HasPrefix(property, prefixWithColon+importantProperty.Name) {
				continue
			}

			addProperty := true
			property = strings.TrimPrefix(property, prefixWithColon)
			switch importantProperty.Type {
			case "image":
				addProperty = ps.imageParser.Parse(property, content, ps.propertyTable)
			case "profile":
				addProperty = ps.profileParser.Parse(property, content, ps.propertyTable)
			case "article":
				addProperty = ps.articleParser.Parse(property, content, ps.propertyTable)
			}

			if addProperty {
				ps.propertyTable[importantProperty.Name] = content
			}
		}
	}
}
