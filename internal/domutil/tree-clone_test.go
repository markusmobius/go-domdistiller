// ORIGINAL: javatest/TreeCloneBuilderTest.java

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

package domutil_test

import (
	"regexp"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_DomUtil_TreeClone_FullTreeBuilder(t *testing.T) {
	expectedHTML := `
<div id="0">
<div id="1">
<div id="2">
<div id="3"></div>
<div id="4"></div>
</div>
<div id="5"></div>
</div>
<div id="8">
<div id="12">
<div id="14"></div>
</div>
</div>
</div>`

	expectedHTML = regexp.MustCompile(`(?mi)^\s+`).ReplaceAllString(expectedHTML, "")
	expectedHTML = regexp.MustCompile(`(?i)\n`).ReplaceAllString(expectedHTML, "")

	divs := testutil.CreateDivTree()
	leafNodes := []*html.Node{divs[3], divs[4], divs[5], divs[14]}

	root := domutil.TreeClone(leafNodes)
	assert.Equal(t, expectedHTML, dom.OuterHTML(root))
}

func Test_DomUtil_TreeClone_SingleLeafNode(t *testing.T) {
	leafNodes := []*html.Node{dom.CreateTextNode("some content")}

	root := domutil.TreeClone(leafNodes)
	assert.Equal(t, dom.TextContent(leafNodes[0]), dom.TextContent(root))
}

func Test_DomUtil_TreeClone_NoCommonAncestors(t *testing.T) {
	divs := testutil.CreateDivTree()
	leafNodes := []*html.Node{divs[3], dom.CreateElement("div")}

	root := domutil.TreeClone(leafNodes)
	assert.Nil(t, root)
}
