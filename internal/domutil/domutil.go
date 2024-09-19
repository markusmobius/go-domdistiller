// ORIGINAL: java/DomUtil.java

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

package domutil

import (
	nurl "net/url"
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

var (
	rxDisplay          = regexp.MustCompile(`(?i)display:\s*([\w-]+)\s*(?:;|$)`)
	rxVisibilityHidden = regexp.MustCompile(`(?i)visibility:\s*(:?hidden|collapse)`)
	rxSrcsetURL        = regexp.MustCompile(`(?i)(\S+)(\s+[\d.]+[xw])?(\s*(?:,|$))`)

	elementWithSizeAttr = map[string]struct{}{
		"table": {},
		"th":    {},
		"td":    {},
		"hr":    {},
		"pre":   {},
	}
)

// HasRootDomain checks if a provided URL has the specified root domain
// (ex. http://a.b.c/foo/bar has root domain of b.c).
func HasRootDomain(url string, root string) bool {
	if url == "" || root == "" {
		return false
	}

	if strings.HasPrefix(url, "//") {
		url = "http:" + url
	}

	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil {
		return false
	}

	return parsedURL.Host == root || strings.HasSuffix(parsedURL.Host, "."+root)
}

// GetFirstElementByTagNameInc returns the first element with `tagName` in the
// tree rooted at `root`, including root. Nil if none is found.
func GetFirstElementByTagNameInc(root *html.Node, tagName string) *html.Node {
	if dom.TagName(root) == tagName {
		return root
	}
	return GetFirstElementByTagName(root, tagName)
}

// GetFirstElementByTagName returns the first element with `tagName` in the
// tree rooted at `root`. Nil if none is found.
func GetFirstElementByTagName(root *html.Node, tagName string) *html.Node {
	nodes := dom.GetElementsByTagName(root, tagName)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

// GetNearestCommonAncestor returns the nearest common ancestor of nodes.
func GetNearestCommonAncestor(nodes ...*html.Node) *html.Node {
	_, nearestAncestor := GetAncestors(nodes...)
	return nearestAncestor
}

// GetAncestors returns all ancestor of the `nodes` and also the nearest common ancestor.
func GetAncestors(nodes ...*html.Node) (map[*html.Node]int, *html.Node) {
	// Find all ancestors
	ancestorList := []*html.Node{}
	ancestorCount := make(map[*html.Node]int)
	saveAncestor := func(ancestor *html.Node) {
		if _, exist := ancestorCount[ancestor]; !exist {
			ancestorList = append(ancestorList, ancestor)
		}
		ancestorCount[ancestor]++
	}

	for _, node := range nodes {
		// Include the node itself to list of ancestor
		saveAncestor(node)

		// Save parents of node to list ancestor
		parent := node.Parent
		for parent != nil {
			saveAncestor(parent)
			parent = parent.Parent
		}
	}

	// Find the nearest ancestor
	nNodes := len(nodes)
	var nearestAncestor *html.Node
	for _, node := range ancestorList {
		if ancestorCount[node] == nNodes {
			nearestAncestor = node
			break
		}
	}

	return ancestorCount, nearestAncestor
}

// MakeAllLinksAbsolute makes all anchors and video posters absolute.
func MakeAllLinksAbsolute(root *html.Node, pageURL *nurl.URL) {
	rootTagName := dom.TagName(root)

	if rootTagName == "a" {
		if href := dom.GetAttribute(root, "href"); href != "" {
			absHref := stringutil.CreateAbsoluteURL(href, pageURL)
			dom.SetAttribute(root, "href", absHref)
		}
	}

	if rootTagName == "video" {
		if poster := dom.GetAttribute(root, "poster"); poster != "" {
			absPoster := stringutil.CreateAbsoluteURL(poster, pageURL)
			dom.SetAttribute(root, "poster", absPoster)
		}
	}

	for _, link := range dom.GetElementsByTagName(root, "a") {
		if href := dom.GetAttribute(link, "href"); href != "" {
			absHref := stringutil.CreateAbsoluteURL(href, pageURL)
			dom.SetAttribute(link, "href", absHref)
		}
	}

	for _, video := range dom.GetElementsByTagName(root, "video") {
		if poster := dom.GetAttribute(video, "poster"); poster != "" {
			absPoster := stringutil.CreateAbsoluteURL(poster, pageURL)
			dom.SetAttribute(video, "poster", absPoster)
		}
	}

	MakeAllSrcAttributesAbsolute(root, pageURL)
	MakeAllSrcSetAbsolute(root, pageURL)
}

// MakeAllSrcAttributesAbsolute makes all "img", "source", "track", and "video"
// tags have an absolute "src" attribute.
func MakeAllSrcAttributesAbsolute(root *html.Node, pageURL *nurl.URL) {
	switch dom.TagName(root) {
	case "img", "source", "track", "video":
		if src := dom.GetAttribute(root, "src"); src != "" {
			absSrc := stringutil.CreateAbsoluteURL(src, pageURL)
			dom.SetAttribute(root, "src", absSrc)
		}
	}

	for _, element := range dom.QuerySelectorAll(root, "img,source,track,video") {
		if src := dom.GetAttribute(element, "src"); src != "" {
			absSrc := stringutil.CreateAbsoluteURL(src, pageURL)
			dom.SetAttribute(element, "src", absSrc)
		}
	}
}

// MakeAllSrcSetAbsolute makes all `srcset` within root absolute.
func MakeAllSrcSetAbsolute(root *html.Node, pageURL *nurl.URL) {
	if dom.HasAttribute(root, "srcset") {
		makeSrcSetAbsolute(root, pageURL)
	}

	for _, element := range dom.QuerySelectorAll(root, "[srcset]") {
		makeSrcSetAbsolute(element, pageURL)
	}
}

func GetSrcSetURLs(node *html.Node) []string {
	srcset := dom.GetAttribute(node, "srcset")
	if srcset == "" {
		return nil
	}

	matches := rxSrcsetURL.FindAllStringSubmatch(srcset, -1)
	urls := make([]string, len(matches))
	for i, group := range matches {
		urls[i] = group[1]
	}

	return urls
}

func GetAllSrcSetURLs(root *html.Node) []string {
	urls := GetSrcSetURLs(root)
	for _, node := range dom.QuerySelectorAll(root, "[srcset]") {
		urls = append(urls, GetSrcSetURLs(node)...)
	}

	return urls
}

func StripAttributes(node *html.Node) {
	elements := dom.GetElementsByTagName(node, "*")
	elements = append(elements, node)

	for _, elem := range elements {
		tagName := dom.TagName(elem)
		finalAttrs := []html.Attribute{}
		_, elementAllowedToHaveSize := elementWithSizeAttr[tagName]

		for _, attr := range elem.Attr {
			// Exclude identification and presentational attributes.
			switch attr.Key {
			case "id", "class", "align", "background", "bgcolor", "border", "cellpadding",
				"cellspacing", "frame", "hspace", "rules", "style", "valign", "vspace":
				continue

			case "width", "height":
				if !elementAllowedToHaveSize {
					continue
				}
			}

			// Exclude unsafe attributes
			if _, allowed := allowedAttributes[attr.Key]; !allowed {
				continue
			}

			finalAttrs = append(finalAttrs, attr)
		}

		elem.Attr = finalAttrs
	}
}

// CloneAndProcessList clones and process list of relevant nodes for output.
func CloneAndProcessList(outputNodes []*html.Node, pageURL *nurl.URL) *html.Node {
	if len(outputNodes) == 0 {
		return nil
	}

	clonedSubTree := TreeClone(outputNodes)
	if clonedSubTree == nil || clonedSubTree.Type != html.ElementNode {
		return nil
	}

	MakeAllLinksAbsolute(clonedSubTree, pageURL)
	StripAttributes(clonedSubTree)
	return clonedSubTree
}

// CloneAndProcessTree clone and process a given node tree/subtree.
// In original dom-distiller this will ignore hidden elements,
// unfortunately we can't do that here, so we will include hidden
// elements as well. NEED-COMPUTE-CSS.
func CloneAndProcessTree(root *html.Node, pageURL *nurl.URL) *html.Node {
	return CloneAndProcessList(GetOutputNodes(root), pageURL)
}

// GetOutputNodes returns list of relevant nodes for output from a subtree.
func GetOutputNodes(root *html.Node) []*html.Node {
	outputNodes := []*html.Node{}
	WalkNodes(root, func(node *html.Node) bool {
		switch node.Type {
		case html.TextNode:
			outputNodes = append(outputNodes, node)
			return false

		case html.ElementNode:
			outputNodes = append(outputNodes, node)
			return true

		default:
			return false
		}
	}, nil)

	return outputNodes
}

// GetParentNodes returns list of all the parents of this node starting with the node itself.
func GetParentNodes(node *html.Node) []*html.Node {
	result := []*html.Node{}
	current := node
	for current != nil {
		result = append(result, current)
		current = current.Parent
	}
	return result
}

// GetNodeDepth the depth of the given node in the DOM tree.
func GetNodeDepth(node *html.Node) int {
	return len(GetParentNodes(node)) - 1
}

// makeSrcSetAbsolute makes `srcset` for this node absolute.
func makeSrcSetAbsolute(node *html.Node, pageURL *nurl.URL) {
	srcset := dom.GetAttribute(node, "srcset")
	if srcset == "" {
		dom.RemoveAttribute(node, "srcset")
		return
	}

	newSrcset := rxSrcsetURL.ReplaceAllStringFunc(srcset, func(s string) string {
		p := rxSrcsetURL.FindStringSubmatch(s)
		return stringutil.CreateAbsoluteURL(p[1], pageURL) + p[2] + p[3]
	})

	dom.SetAttribute(node, "srcset", newSrcset)
}

// =================================================================================
// Functions below these point are functions that exist in original Dom-Distiller
// code but that can't be perfectly replicated by this package. This is because
// in original Dom-Distiller they uses GWT which able to compute stylesheet.
// Unfortunately, Go can't do this unless we are using some kind of headless
// browser, so here we only do some kind of workaround (which works but obviously
// not as good as GWT) or simply ignore it.
// =================================================================================

// InnerText in JS and GWT is used to capture text from an element while excluding
// text from hidden children. A child is hidden if it's computed width is 0, whether
// because its CSS (e.g `display: none`, `visibility: hidden`, etc), or if the child
// has `hidden` attribute. Since we can't compute stylesheet, we only look at `hidden`
// attribute and inline style.
//
// Besides excluding text from hidden children, difference between this function and
// `dom.TextContent` is the latter will skip <br> tag while this function will preserve
// <br> as whitespace. NEED-COMPUTE-CSS
func InnerText(node *html.Node) string {
	var buffer strings.Builder
	var finder func(*html.Node)

	finder = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			buffer.WriteString(" ")
			buffer.WriteString(n.Data)
			buffer.WriteString(" ")

		case html.ElementNode:
			if n.Data == "br" {
				buffer.WriteString(`|\/|`)
				return
			}

			if !IsProbablyVisible(n) {
				return
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}

	finder(node)
	text := buffer.String()
	text = strings.Join(strings.Fields(text), " ")
	text = re2go.TidyUpPunctuation(text)
	text = re2go.FixTempNewline(text)
	return text
}

// GetArea in original code returns area of a node by multiplying
// offsetWidth and offsetHeight. Since it's not possible in Go, we
// simply return 0. NEED-COMPUTE-CSS
func GetArea(node *html.Node) int {
	return 0
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for dom management.
// =================================================================================

// SomeNode iterates over a NodeList, return true if any of the
// provided iterate function calls returns true, false otherwise.
func SomeNode(nodeList []*html.Node, fn func(*html.Node) bool) bool {
	for i := 0; i < len(nodeList); i++ {
		if fn(nodeList[i]) {
			return true
		}
	}
	return false
}

// GetParentElement returns the nearest element parent.
func GetParentElement(node *html.Node) *html.Node {
	for parent := node.Parent; parent != nil; parent = parent.Parent {
		if parent.Type == html.ElementNode {
			return parent
		}
	}

	return nil
}

// Contains checks if child is inside node.
func Contains(node, child *html.Node) bool {
	if node == nil || child == nil {
		return false
	}

	if node == child {
		return true
	}

	childParent := child.Parent
	for childParent != nil {
		if childParent == node {
			return true
		}
		childParent = childParent.Parent
	}

	return false
}

// NodeName returns the name of the current node as a string.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeName
func NodeName(node *html.Node) string {
	switch node.Type {
	case html.TextNode:
		return "#text"
	case html.DocumentNode:
		return "#document"
	case html.CommentNode:
		return "#comment"
	case html.ElementNode, html.DoctypeNode:
		return node.Data
	default:
		return ""
	}
}

// HasAncestor check if node has ancestor with specified tag names.
func HasAncestor(node *html.Node, ancestorTagNames ...string) bool {
	// Convert array to map
	mapAncestorTags := make(map[string]struct{})
	for _, tag := range ancestorTagNames {
		mapAncestorTags[tag] = struct{}{}
	}

	// Check ancestors
	for parent := node.Parent; parent != nil; parent = parent.Parent {
		parentTagName := dom.TagName(parent)
		if _, exist := mapAncestorTags[parentTagName]; exist {
			return true
		}
	}

	return false
}

// IsProbablyVisible determines if a node is visible.
func IsProbablyVisible(node *html.Node) bool {
	displayStyle := GetDisplayStyle(node)
	styleAttr := dom.GetAttribute(node, "style")
	nodeAriaHidden := dom.GetAttribute(node, "aria-hidden")
	className := dom.GetAttribute(node, "class")

	// Have to null-check node.style and node.className.indexOf to deal
	// with SVG and MathML nodes. Also check for "fallback-image" so that
	// Wikimedia Math images are displayed
	return displayStyle != "none" &&
		!dom.HasAttribute(node, "hidden") &&
		!rxVisibilityHidden.MatchString(styleAttr) &&
		(nodeAriaHidden == "" || nodeAriaHidden != "true" || strings.Contains(className, "fallback-image"))
}

// GetDisplayStyle returns the default "display" in style property for the specified node.
func GetDisplayStyle(node *html.Node) string {
	// Check if display specified in inline style
	style := dom.GetAttribute(node, "style")
	parts := rxDisplay.FindStringSubmatch(style)
	if len(parts) >= 2 {
		return parts[1]
	}

	// Use default display
	switch dom.TagName(node) {
	case "address", "article", "blockquote", "body", "dd", "details", "dialog", "div",
		"dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "h1", "h2",
		"h3", "h4", "h5", "h6", "header", "hr", "html", "legend", "main", "nav", "ol",
		"p", "pre", "section", "ul":
		return "block"
	case "a", "abbr", "acronym", "audio", "b", "bdi", "bdo", "br", "canvas", "circle", "cite",
		"code", "data", "defs", "del", "dfn", "ellipse", "em", "embed", "font", "i", "iframe", "img",
		"ins", "kbd", "label", "lineargradient", "mark", "object", "output", "picture", "polygon",
		"q", "rect", "s", "source", "span", "stop", "strong", "sub", "sup", "svg", "tt", "text",
		"time", "track", "u", "var", "video", "wbr":
		return "inline"
	case "button", "input":
		return "inline-block"
	case "li", "summary":
		return "list-item"
	case "ruby":
		return "ruby"
	case "rt":
		return "ruby-text"
	case "table":
		return "table"
	case "caption":
		return "table-caption"
	case "td", "th":
		return "table-cell"
	case "col":
		return "table-column"
	case "colgroup":
		return "table-column-group"
	case "tfoot":
		return "table-footer-group"
	case "thead":
		return "table-header-group"
	case "tr":
		return "table-row"
	case "tbody":
		return "table-row-group"
	case "meta", "script", "style", "link":
		return "none"
	}

	// In case I forgot any tags, fallback to block
	return "block"
}
