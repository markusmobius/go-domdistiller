// ORIGINAL: javatest/PageParamInfoTest.java

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

package info_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	tu "github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Info_PPI_InsertFirstPage(t *testing.T) {
	paramInfo := &info.PageParamInfo{}
	paramInfo.Type = info.PageNumber

	{
		testURL := "http://www.google.com/article/bar"
		allContentInfo := []*pageParamContentInfoEx{
			ppciExNumberInPlainText(1),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=2", 2, true),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=3", 3, true),
		}

		canInsert := canInsertFirstPage(testURL, allContentInfo, paramInfo)
		assert.True(t, canInsert)
		assert.Len(t, paramInfo.AllPageInfo, 2)

		// The current document is inserted as first page.
		paramInfo.InsertFirstPage(testURL)
		assert.Len(t, paramInfo.AllPageInfo, 3)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 1, page.PageNumber)
		assert.Equal(t, "http://www.google.com/article/bar", page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, "http://www.google.com/article/bar?page=2", page.URL)

		page = paramInfo.AllPageInfo[2]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, "http://www.google.com/article/bar?page=3", page.URL)
	}

	// Current document url has same length as other paginated pages, so it shouldn't be
	// inserted as first page.
	{
		testURL := "http://www.google.com/article/bar?page=1"
		allContentInfo := []*pageParamContentInfoEx{
			ppciExNumberInPlainText(1),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=2", 2, true),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=3", 3, true),
		}

		canInsert := canInsertFirstPage(testURL, allContentInfo, paramInfo)
		assert.False(t, canInsert)
		assert.Len(t, paramInfo.AllPageInfo, 2)
	}

	// Current document url is the last page, so shouldn't be inserted as first page.
	{
		testURL := "http://www.google.com/article/bar?page=4"
		allContentInfo := []*pageParamContentInfoEx{
			ppciExNumericOutlink("http://www.google.com/article/bar", 1, false),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=2", 2, true),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=3", 3, true),
			ppciExNumberInPlainText(4),
		}

		canInsert := canInsertFirstPage(testURL, allContentInfo, paramInfo)
		assert.False(t, canInsert)
		assert.Len(t, paramInfo.AllPageInfo, 2)
	}

	// Page one has an outlink to itself, should be inserted as first page.
	{
		testURL := "http://www.google.com/article/bar"
		allContentInfo := []*pageParamContentInfoEx{
			ppciExNumericOutlink("http://www.google.com/article/bar", 1, false),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=2", 2, true),
			ppciExNumericOutlink("http://www.google.com/article/bar?page=3", 3, true),
		}

		canInsert := canInsertFirstPage(testURL, allContentInfo, paramInfo)
		assert.True(t, canInsert)
		assert.Len(t, paramInfo.AllPageInfo, 2)

		paramInfo.InsertFirstPage(testURL)
		assert.Len(t, paramInfo.AllPageInfo, 3)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 1, page.PageNumber)
		assert.Equal(t, "http://www.google.com/article/bar", page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, "http://www.google.com/article/bar?page=2", page.URL)

		page = paramInfo.AllPageInfo[2]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, "http://www.google.com/article/bar?page=3", page.URL)
	}
}

type pageParamContentInfoEx struct {
	contentInfo *tu.PageParamContentInfo
	isPageInfo  bool
}

func ppciExNumberInPlainText(number int) *pageParamContentInfoEx {
	return &pageParamContentInfoEx{
		contentInfo: tu.PPCINumberInPlainText(number),
	}
}

func ppciExNumericOutlink(targetURL string, number int, isPageInfo bool) *pageParamContentInfoEx {
	return &pageParamContentInfoEx{
		contentInfo: tu.PPCINumericOutlink(targetURL, number),
		isPageInfo:  isPageInfo,
	}
}

func canInsertFirstPage(docURL string, allContentInfo []*pageParamContentInfoEx, pageParamInfo *info.PageParamInfo) bool {
	ascendingNumbers := []*info.PageInfo{}
	pageParamInfo.AllPageInfo = []*info.PageInfo{}

	for _, ex := range allContentInfo {
		switch ex.contentInfo.Type {
		case tu.NumberInPlainText:
			ascendingNumbers = append(ascendingNumbers, &info.PageInfo{
				PageNumber: ex.contentInfo.Number,
			})

		case tu.NumericOutlink:
			page := &info.PageInfo{
				PageNumber: ex.contentInfo.Number,
				URL:        ex.contentInfo.TargetURL,
			}

			ascendingNumbers = append(ascendingNumbers, page)
			if ex.isPageInfo {
				pageParamInfo.AllPageInfo = append(pageParamInfo.AllPageInfo, page)
			}
		}
	}

	return pageParamInfo.CanInsertFirstPage(docURL, ascendingNumbers)
}
