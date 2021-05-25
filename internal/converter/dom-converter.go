// ORIGINAL: java/webdocument/DomConverter.java

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

package converter

import (
	nurl "net/url"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/extractor/embed"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/tableclass"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// DomConverter converts a node and its children into a Document.
type DomConverter struct {
	builder         webdoc.DocumentBuilder
	embedExtractors []embed.EmbedExtractor
	embedTagNames   map[string]struct{}
	pageURL         *nurl.URL
	logger          logutil.Logger
	tableClassifier *tableclass.Classifier
}

func NewDomConverter(builder webdoc.DocumentBuilder, pageURL *nurl.URL, logger logutil.Logger) *DomConverter {
	extractors := []embed.EmbedExtractor{
		embed.NewImageExtractor(pageURL, logger),
		embed.NewTwitterExtractor(pageURL, logger),
		embed.NewVimeoExtractor(pageURL, logger),
		embed.NewYouTubeExtractor(pageURL, logger),
	}

	embedTagNames := make(map[string]struct{})
	for _, extractor := range extractors {
		for _, tagName := range extractor.RelevantTagNames() {
			embedTagNames[tagName] = struct{}{}
		}
	}

	return &DomConverter{
		builder:         builder,
		embedExtractors: extractors,
		embedTagNames:   embedTagNames,
		pageURL:         pageURL,
		logger:          logger,
		tableClassifier: tableclass.NewClassifier(logger),
	}
}

func (dc *DomConverter) Convert(root *html.Node) {
	domutil.WalkNodes(root, dc.visitNodeHandler, dc.exitNodeHandler)
}

func (dc *DomConverter) visitNodeHandler(node *html.Node) bool {
	switch node.Type {
	case html.TextNode:
		dc.builder.AddTextNode(node)
		return false

	case html.ElementNode:
		return dc.visitElementNodeHandler(node)

	default:
		return false
	}
}

func (dc *DomConverter) exitNodeHandler(node *html.Node) {
	if node.Type == html.ElementNode {
		if tagName := dom.TagName(node); webdoc.CanBeNested(tagName) {
			dc.builder.AddTag(webdoc.NewTag(tagName, webdoc.TagEnd))
		}
	}

	dc.builder.EndNode()
}

func (dc *DomConverter) visitElementNodeHandler(node *html.Node) bool {
	// In original dom-distiller they skip invisible or uninteresting elements.
	// Unfortunately it's impossible to do that perfectly here (NEED-COMPUTE-CSS).
	if !domutil.IsProbablyVisible(node) {
		return false
	}

	// Node-type specific extractors check for elements they are interested in here.
	// Everything else will be filtered through the switch below.
	tagName := dom.TagName(node)
	if _, isEmbed := dc.embedTagNames[tagName]; isEmbed {
		// If the tag is marked as interesting, check the extractors.
		for _, extractor := range dc.embedExtractors {
			embed := extractor.Extract(node)
			if embed != nil {
				dc.builder.AddEmbed(embed)
				return false
			}
		}
	}

	// Skip social and sharing elements.
	// See crbug.com/692553, crbug.com/696556, and crbug.com/674557
	className := dom.GetAttribute(node, "class")
	component := dom.GetAttribute(node, "data-component")
	if className == "sharing" || className == "socialArea" || component == "share" {
		return false
	}

	// Create a placeholder for the elements we want to preserve.
	if webdoc.CanBeNested(tagName) {
		dc.builder.AddTag(webdoc.NewTag(tagName, webdoc.TagStart))
	}

	switch tagName {
	case "a":
		// The "section" parameter is to differentiate with "redlinks".
		// Ref: https://en.wikipedia.org/wiki/Wikipedia:Red_link
		href := dom.GetAttribute(node, "href")
		if strings.Contains(href, "action=edit&section=") {
			// Skip "edit section" on mediawiki.
			// See crbug.com/647667.
			return false
		}

	case "span":
		if className == "mw-editsection" {
			// Skip "[edit]" on mediawiki desktop version.
			// See crbug.com/647667.
			return false
		}

	case "br":
		dc.builder.AddLineBreak(node)
		return false

	case "table":
		tableType, _ := dc.tableClassifier.Classify(node)
		dc.logTableInfo(node, tableType)
		if tableType == tableclass.Data {
			dc.builder.AddDataTable(node)
			return false
		}

	case "video":
		dc.builder.AddEmbed(webdoc.NewVideo(node, dc.pageURL, 0, 0))
		return false

	// These element types are all skipped (but may affect document construction).
	case "option", "object", "embed", "applet":
		dc.builder.SkipNode(node)
		return false

	// These types are skipped and don't affect document construction.
	case "head", "style", "script", "link", "noscript", "iframe", "svg":
		return false
	}

	dc.builder.StartNode(node)
	return true
}

func (dc *DomConverter) logTableInfo(table *html.Node, tableType tableclass.Type) {
	if dc.logger == nil {
		return
	}

	if !dc.logger.IsLogVisibility() {
		return
	}

	id := dom.GetAttribute(table, "id")
	class := dom.GetAttribute(table, "class")

	logMsg := "Table: " + tableType.String()
	if id != "" {
		logMsg += " #" + id
	}
	if class != "" {
		logMsg += " class=" + class
	}

	parent := domutil.GetParentElement(table)
	if parent != nil {
		tagName := dom.TagName(parent)
		id := dom.GetAttribute(parent, "id")
		class := dom.GetAttribute(parent, "class")

		logMsg += ", parent=[" + tagName
		if id != "" {
			logMsg += " #" + id
		}
		if class != "" {
			logMsg += " class=" + class
		}
		logMsg += "]"
	}

	dc.logger.PrintVisibilityInfo(logMsg)
}
