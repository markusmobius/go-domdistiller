// ORIGINAL: java/webdocument/filters/NestedElementRetainer.java

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

package docfilter

import (
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

type NestedElementRetainer struct{}

func NewNestedElementRetainer() *NestedElementRetainer {
	return &NestedElementRetainer{}
}

func (f *NestedElementRetainer) Process(doc *webdoc.Document) bool {
	isContent := false
	stackMark := -1
	stack := []*webdoc.Tag{}

	for _, e := range doc.Elements {
		if webTag, isTag := e.(*webdoc.Tag); !isTag {
			if !isContent {
				isContent = e.IsContent()
			}
		} else {
			if webTag.Type == webdoc.TagStart {
				webTag.SetIsContent(isContent)
				stack = append(stack, webTag)
				isContent = false
			} else {
				startWebTag := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				isContent = isContent || stackMark >= len(stack)
				if isContent {
					stackMark = len(stack) - 1
				}

				wasContent := startWebTag.IsContent()
				startWebTag.SetIsContent(isContent)
				webTag.SetIsContent(isContent)
				isContent = wasContent
			}
		}
	}

	return true
}
