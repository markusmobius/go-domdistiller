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
	"bytes"
	nurl "net/url"
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

var (
	rxPunctuation      = regexp.MustCompile(`\s+([.?!,;])\s*(\S*)`)
	rxTempNewline      = regexp.MustCompile(`\s*\|\\/\|\s*`)
	rxDisplayNone      = regexp.MustCompile(`(?i)display:\s*none`)
	rxVisibilityHidden = regexp.MustCompile(`(?i)visibility:\s*(:?hidden|collapse)`)
	rxSrcsetURL        = regexp.MustCompile(`(?i)(\S+)(\s+[\d.]+[xw])?(\s*(?:,|$))`)
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
	ancestors := make(map[*html.Node]int)
	for _, node := range nodes {
		// Include the node itself to list of ancestor
		ancestors[node]++

		// Save parents of node to list ancestor
		parent := node.Parent
		for parent != nil {
			ancestors[parent]++
			parent = parent.Parent
		}
	}

	// Find common ancestor
	nNodes := len(nodes)
	commonAncestors := make(map[*html.Node]struct{})
	for node, count := range ancestors {
		if count == nNodes {
			commonAncestors[node] = struct{}{}
		}
	}

	// If there are no common ancestor found, stop
	if len(commonAncestors) == 0 {
		return nil, nil
	}

	// Find the nearest ancestor
	var nearestAncestor *html.Node
	for node := range commonAncestors {
		childIsAncestor := false
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if _, exist := commonAncestors[child]; exist {
				childIsAncestor = true
				break
			}
		}

		if !childIsAncestor {
			nearestAncestor = node
		}
	}

	return ancestors, nearestAncestor
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

// StripImageElement removes unnecessary attributes for image elements.
func StripImageElement(img *html.Node) {
	importantAttrs := []html.Attribute{}
	for _, attr := range img.Attr {
		switch attr.Key {
		case "src", "alt", "srcset", "dir", "width", "height", "title":
			importantAttrs = append(importantAttrs, attr)
		default:
			continue
		}
	}
	img.Attr = importantAttrs
}

func StripImageElements(root *html.Node) {
	if dom.TagName(root) == "img" {
		StripImageElement(root)
	}

	for _, img := range dom.QuerySelectorAll(root, "img") {
		StripImageElement(img)
	}
}

// StripAttributeFromTagss trips some attribute from certain tags in the tree
// rooted at `root`, including root itself.
func StripAttributeFromTags(root *html.Node, attr string, tagNames ...string) {
	rootTagName := dom.TagName(root)
	for _, tag := range tagNames {
		if rootTagName == tag || tag == "*" {
			dom.RemoveAttribute(root, attr)
			break
		}
	}

	for i, tag := range tagNames {
		tagNames[i] = tag + "[" + attr + "]"
	}

	selectors := strings.Join(tagNames, ",")
	for _, elem := range dom.QuerySelectorAll(root, selectors) {
		dom.RemoveAttribute(elem, attr)
	}
}

// StripIDs strips all "id" attributes from all nodes in the tree rooted at `root`
func StripIDs(root *html.Node) {
	StripAttributeFromTags(root, "id", "*")
}

// StripFontColorAttributes strips all "color" attributes from "font" nodes in the
// tree rooted at `root`
func StripFontColorAttributes(root *html.Node) {
	StripAttributeFromTags(root, "color", "font")
}

// StripTableBackgroundColorAttributes strips all "bgcolor" attributes from table
// nodes in the tree rooted at `root`
func StripTableBackgroundColorAttributes(root *html.Node) {
	StripAttributeFromTags(root, "bgcolor", "table", "tr", "td", "th")
}

// StripStyleAttributes strips all "style" attributes from all nodes in the tree
// rooted at `root`
func StripStyleAttributes(root *html.Node) {
	StripAttributeFromTags(root, "style", "*")
}

// StripTargetAttributes strips all "target" attributes from anchor nodes in the
// tree rooted at `root`
func StripTargetAttributes(root *html.Node) {
	StripAttributeFromTags(root, "target", "a")
}

// StripUnwantedClassNames strips unwanted classNames from all nodes in the tree
// rooted at `root`.
func StripUnwantedClassNames(root *html.Node) {
	if dom.HasAttribute(root, "class") {
		stripUnwantedClassNames(root)
	}

	for _, element := range dom.QuerySelectorAll(root, "[class]") {
		stripUnwantedClassNames(element)
	}
}

// StripAllUnsafeAttributes strips all attributes from nodes other than
// ones in the list of allowedAttributes.
func StripAllUnsafeAttributes(root *html.Node) {
	if root.Type == html.ElementNode {
		stripAllUnsafeAttributes(root)
	}

	for _, element := range dom.QuerySelectorAll(root, "*") {
		stripAllUnsafeAttributes(element)
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

	StripIDs(clonedSubTree)
	MakeAllLinksAbsolute(clonedSubTree, pageURL)
	StripTargetAttributes(clonedSubTree)
	StripFontColorAttributes(clonedSubTree)
	StripTableBackgroundColorAttributes(clonedSubTree)
	StripStyleAttributes(clonedSubTree)
	StripImageElements(clonedSubTree)
	StripAllUnsafeAttributes(clonedSubTree)
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

func stripUnwantedClassNames(node *html.Node) {
	class := dom.GetAttribute(node, "class")
	if strings.Contains(class, "caption") {
		dom.SetAttribute(node, "class", "caption")
	} else {
		dom.RemoveAttribute(node, "class")
	}
}

func stripAllUnsafeAttributes(node *html.Node) {
	allowedAttrs := []html.Attribute{}
	for _, attr := range node.Attr {
		if _, allowed := allowedAttributes[attr.Key]; allowed {
			allowedAttrs = append(allowedAttrs, attr)
		}
	}

	node.Attr = allowedAttrs
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
	var buffer bytes.Buffer
	var finder func(*html.Node)

	finder = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			buffer.WriteString(" " + n.Data + " ")

		case html.ElementNode:
			if n.Data == "br" {
				buffer.WriteString(`|\/|`)
				return
			}

			if dom.HasAttribute(n, "hidden") {
				return
			}

			styleAttr := dom.GetAttribute(n, "style")
			if rxDisplayNone.MatchString(styleAttr) || rxVisibilityHidden.MatchString(styleAttr) {
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
	text = rxPunctuation.ReplaceAllString(text, "$1 $2")
	text = rxTempNewline.ReplaceAllString(text, "\n")
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
