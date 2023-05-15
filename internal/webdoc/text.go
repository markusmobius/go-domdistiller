// ORIGINAL: java/webdocument/WebText.java

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

package webdoc

import (
	"fmt"
	nurl "net/url"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"golang.org/x/net/html"
)

type Text struct {
	BaseElement

	Text           string
	NumWords       int
	NumLinkedWords int
	Labels         map[string]struct{}
	TagLevel       int
	OffsetBlock    int
	GroupNumber    int
	PageURL        *nurl.URL

	TextNodes     []*html.Node
	Start         int
	End           int
	FirstWordNode int
	LastWordNode  int
}

func (t *Text) ElementType() string {
	return "text"
}

func (t *Text) GenerateOutput(textOnly bool) string {
	if t.HasLabel(label.Title) {
		return ""
	}

	// TODO: Instead of doing this next part, in the future track font size weight
	// and etc. and wrap the nodes in a "p" tag.
	clonedRoot := domutil.TreeClone(t.GetTextNodes())

	// To keep formatting/structure, at least one parent element should be in the output.
	// This is necessary because many times a WebText is only a single text node.
	if clonedRoot.Type != html.ElementNode {
		parentClone := dom.Clone(t.GetTextNodes()[0].Parent, false)
		dom.AppendChild(parentClone, clonedRoot)
		clonedRoot = parentClone
	}

	// The body element should not be used in the output.
	if dom.TagName(clonedRoot) == "body" {
		div := dom.CreateElement("div")
		dom.SetInnerHTML(div, dom.InnerHTML(clonedRoot))
		clonedRoot = div
	}

	// Retain parent tags until the root is not an inline element, to make sure the
	// style is display:block.
	var srcRoot *html.Node
	for {
		display := domutil.GetDisplayStyle(clonedRoot)
		if display != "inline" {
			break
		}

		if srcRoot == nil {
			srcRoot = domutil.GetNearestCommonAncestor(t.GetTextNodes()...)
			if srcRoot.Type != html.ElementNode {
				srcRoot = domutil.GetParentElement(srcRoot)
			}
		}

		srcRoot = domutil.GetParentElement(srcRoot)
		if dom.TagName(srcRoot) == "body" {
			break
		}

		parentClone := dom.Clone(srcRoot, false)
		dom.AppendChild(parentClone, clonedRoot)
		clonedRoot = parentClone
	}

	// Make sure links are absolute and IDs are gone.
	domutil.MakeAllLinksAbsolute(clonedRoot, t.PageURL)
	domutil.StripAttributes(clonedRoot)
	// TODO: if we allow images in WebText later, add StripImageElements().

	// Since there are tag elements that are being wrapped by a pair of Tags,
	// we only need to get the innerHTML, otherwise these tags would be duplicated.
	if textOnly {
		return domutil.InnerText(clonedRoot)
	}

	if CanBeNested(dom.TagName(clonedRoot)) {
		return dom.InnerHTML(clonedRoot)
	}

	return dom.OuterHTML(clonedRoot)
}

func (t *Text) AddLabel(s string) {
	if t.Labels == nil {
		t.Labels = make(map[string]struct{})
	}
	t.Labels[s] = struct{}{}
}

func (t *Text) HasLabel(s string) bool {
	_, exist := t.Labels[s]
	return exist
}

func (t *Text) TakeLabels() map[string]struct{} {
	res := t.Labels
	t.Labels = make(map[string]struct{})
	return res
}

func (t Text) FirstNonWhitespaceTextNode() *html.Node {
	return t.TextNodes[t.FirstWordNode]
}

func (t Text) LastNonWhitespaceTextNode() *html.Node {
	return t.TextNodes[t.LastWordNode]
}

func (t Text) GetTextNodes() []*html.Node {
	return t.TextNodes[t.Start:t.End]
}

func (t *Text) String() string {
	return fmt.Sprintf("ELEMENT %q: text=%q, labels=%v, is_content=%v",
		t.ElementType(), t.Text, t.Labels, t.isContent)
}
