// ORIGINAL: java/TreeCloneBuilder.java

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

package domutil

import (
	"golang.org/x/net/html"
)

// TreeClone takes a list of nodes and returns a clone of the minimum tree in the
// DOM that contains all of them. This is done by going through each node, cloning its
// parent and adding children to that parent until the next node is not contained in
// that parent (originally). The list cannot contain a parent of any of the other nodes.
// Children of the nodes in the provided list are excluded.
//
// This implementation doesn't come from the original dom-distiller code. Instead I
// created it from scratch to make it simpler and more Go idiomatic.
func TreeClone(nodes []*html.Node) *html.Node {
	// Get the nearest ancestor
	allAncestors, nearestAncestor := GetAncestors(nodes...)
	if nearestAncestor == nil {
		return nil
	}

	// Clone the ancestor and childrens that required to reach specified nodes
	var fnClone func(src *html.Node) *html.Node
	fnClone = func(src *html.Node) *html.Node {
		clone := &html.Node{
			Type:     src.Type,
			DataAtom: src.DataAtom,
			Data:     src.Data,
			Attr:     append([]html.Attribute{}, src.Attr...),
		}

		for child := src.FirstChild; child != nil; child = child.NextSibling {
			if _, exist := allAncestors[child]; exist {
				clone.AppendChild(fnClone(child))
			}
		}

		return clone
	}

	return fnClone(nearestAncestor)
}
