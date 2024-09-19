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
	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"github.com/markusmobius/go-domdistiller/internal/tableclass"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

type ConverterFlag uint

const (
	Default        ConverterFlag = 0
	SkipUnlikelies ConverterFlag = 1 << iota
)

// DomConverter converts a node and its children into a Document.
type DomConverter struct {
	builder         webdoc.DocumentBuilder
	embedExtractors []embed.EmbedExtractor
	embedTagNames   map[string]struct{}
	pageURL         *nurl.URL
	logger          logutil.Logger
	tableClassifier *tableclass.Classifier
	flags           ConverterFlag
}

func NewDomConverter(flags ConverterFlag, builder webdoc.DocumentBuilder, pageURL *nurl.URL, logger logutil.Logger) *DomConverter {
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
		flags:           flags,
	}
}

func (dc *DomConverter) Convert(root *html.Node) {
	clone := dom.Clone(root, true)
	domutil.WalkNodes(clone, dc.visitNodeHandler, dc.exitNodeHandler)
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

	// Skip social and sharing elements.
	// See crbug.com/692553, crbug.com/696556, and crbug.com/674557
	className := dom.ClassName(node)
	component := dom.GetAttribute(node, "data-component")
	if className == "sharing" || className == "socialArea" || component == "share" {
		return false
	}

	// Skip byline (author)
	nodeData := dom.ClassName(node) + " " + dom.ID(node)
	if isByline(node, nodeData) {
		return false
	}

	// Skip unlikely candidates
	tagName := dom.TagName(node)
	if dc.hasFlag(SkipUnlikelies) {
		if re2go.IsUnlikelyCandidates(nodeData) && !re2go.MaybeItsACandidate(nodeData) &&
			!domutil.HasAncestor(node, "table") && tagName != "body" && tagName != "a" {
			return false
		}

		role := dom.GetAttribute(node, "role")
		if _, isUnlikely := unlikelyRoles[role]; isUnlikely {
			return false
		}
	}

	// Remove DIV, SECTION, and HEADER nodes without any
	// content(e.g. text, image, video, or iframe).
	switch tagName {
	case "div", "section", "header",
		"h1", "h2", "h3", "h4", "h5", "h6":
		if isElementWithoutContent(node) {
			return false
		}
	}

	// Node-type specific extractors check for elements they are interested in here.
	// Everything else will be filtered through the switch below.
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

		// If anchor has Javascript and only contains simple text content, we treat it as text node.
		if strings.HasPrefix(href, "javascript:") {
			linkChildNodes := dom.ChildNodes(node)
			if len(linkChildNodes) == 1 && linkChildNodes[0].Type == html.TextNode {
				textNode := linkChildNodes[0]

				// Replace node with the text node
				dom.DetachChild(textNode)
				node.Parent.InsertBefore(textNode, node)
				dom.DetachChild(node)

				// Save the text
				dc.builder.AddTextNode(textNode)
				return false
			}
		}

	case "span":
		if className == "mw-editsection" {
			// Skip "[edit]" on mediawiki desktop version.
			// See crbug.com/647667.
			return false
		}

	case "font":
		// Replace font element with span
		node.Attr = nil
		node.Data = "span"
		dc.builder.StartNode(node)
		return true

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
	case "option", "object", "embed", "applet",
		"input", "button", "form", "textarea", "select":
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
	if dc.logger == nil || dc.logger.InternallyNil() {
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

func (dc *DomConverter) hasFlag(flag ConverterFlag) bool {
	return dc.flags&flag != 0
}
