// ORIGINAL: java/DomWalker.java

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
	"golang.org/x/net/html"
)

// WalkNodes used to walk the subtree of the DOM rooted at a particular root. It has two
// function parameters, i.e. fnVisit and fnExit :
// - fnVisit is called when we reach a node during the walk. If it returns false, children
//   of the node will be skipped and fnExit won't be called for this node.
// - fnExit is called when exiting a node, after visiting all of its children.
func WalkNodes(root *html.Node, fnVisit func(*html.Node) bool, fnExit func(*html.Node)) {
	if root == nil {
		return
	}

	visitChildren := false
	if fnVisit != nil {
		visitChildren = fnVisit(root)
	}

	if !visitChildren {
		return
	}

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		WalkNodes(child, fnVisit, fnExit)
	}

	if fnExit != nil {
		fnExit(root)
	}
}
