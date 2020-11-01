// ORIGINAL: javatest/DomWalkerTest.java

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

package domutil_test

import (
	"strconv"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_DomUtil_WalkNodes_TopNodeHasNextSiblingAndParent(t *testing.T) {
	root := testutil.CreateDiv(0)
	firstChild := testutil.CreateDiv(1)
	secondChild := testutil.CreateDiv(2)

	dom.AppendChild(root, firstChild)
	dom.AppendChild(root, secondChild)

	runWalkerTest(t, firstChild, []walkVisitData{
		{1, false},
	})
}

func Test_DomUtil_WalkNodes_DivTree(t *testing.T) {
	root := testutil.CreateDivTree()[0]

	runWalkerTest(t, root, []walkVisitData{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, true},
		{6, true},
		{7, true},
		{8, true},
		{9, true},
		{10, true},
		{11, true},
		{12, true},
		{13, true},
		{14, true},
	}, "AllVisited")

	runWalkerTest(t, root, []walkVisitData{
		{0, false},
	}, "RootOnly")

	runWalkerTest(t, root, []walkVisitData{
		{0, true},
		{1, false},
		{8, false},
	}, "OnlyFirstLevel")

	runWalkerTest(t, root, []walkVisitData{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, true},
		{6, false},
		{7, false},
		{8, true},
		{9, false},
		// These are skipped because children of 9 aren't visited.
		// {10,false},
		// {11,false},
		{12, true},
		{13, false},
		{14, true},
	}, "SomeSkipped")
}

type walkVisitData struct {
	ExpectedID  int
	ShouldVisit bool
}

func runWalkerTest(t *testing.T, root *html.Node, listVisitData []walkVisitData, msgs ...interface{}) {
	visitIdx := -1
	nVisitData := len(listVisitData)
	visitedNodes := []*html.Node{}

	fnVisit := func(node *html.Node) bool {
		assert.Equal(t, html.ElementNode, node.Type, msgs...)
		assert.True(t, visitIdx < nVisitData-1, msgs...)

		visitIdx++
		visitData := listVisitData[visitIdx]

		nodeID, _ := strconv.Atoi(dom.GetAttribute(node, "id"))
		assert.Equal(t, visitData.ExpectedID, nodeID, msgs...)

		if visitData.ShouldVisit {
			visitedNodes = append(visitedNodes, node)
		}
		return visitData.ShouldVisit
	}

	fnExit := func(node *html.Node) {
		lastVisited := visitedNodes[len(visitedNodes)-1]
		visitedNodes = visitedNodes[:len(visitedNodes)-1]
		assert.Equal(t, lastVisited, node, msgs...)
	}

	domutil.WalkNodes(root, fnVisit, fnExit)
	assert.Equal(t, nVisitData-1, visitIdx, msgs...)
}
