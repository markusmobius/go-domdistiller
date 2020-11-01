// ORIGINAL: javatest/ImageHeuristicsTest.java

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

package scorer_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter/scorer"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_Scorer_ImageHasFigureScorer(t *testing.T) {
	root := testutil.CreateDiv(0)
	fig := dom.CreateElement("figure")

	goodImage := dom.CreateElement("img")
	dom.SetAttribute(goodImage, "style", "width: 100px; height: 100px; display: block;")

	badImage := dom.CreateElement("img")
	dom.SetAttribute(badImage, "style", "width: 100px; height: 100px; display: block;")

	dom.AppendChild(fig, goodImage)
	dom.AppendChild(root, fig)
	dom.AppendChild(root, badImage)

	imgScorer := scorer.NewImageHasFigureScorer(50)

	assert.True(t, imgScorer.GetImageScore(goodImage) > 0)
	assert.Equal(t, 0, imgScorer.GetImageScore(badImage))
	assert.Equal(t, 0, imgScorer.GetImageScore(nil))
}
