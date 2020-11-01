// ORIGINAL: java/webdocument/filters/images/AreaScorer.java

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

package scorer

import "golang.org/x/net/html"

// ImageAreaScorer uses image area (length*width) as its heuristic.
// Unfortunately to do that we need to compute CSS which is impossible
// in Go, so this scorer do nothing. NEED-COMPUTE-CSS.
type ImageAreaScorer struct {
	maxScore int
	minArea  int
	maxArea  int
}

// NewImageAreaScorer returns and initiates the ImageAreaScorer.
func NewImageAreaScorer(maxScore, minArea, maxArea int) *ImageAreaScorer {
	return &ImageAreaScorer{
		maxScore: maxScore,
		minArea:  minArea,
		maxArea:  maxArea,
	}
}

func (s *ImageAreaScorer) GetImageScore(_ *html.Node) int {
	return 0
}

func (s *ImageAreaScorer) GetMaxScore() int {
	return s.maxScore
}
