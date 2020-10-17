// ORIGINAL: javatest/PageParamInfoTest.java

package info_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_PPI_InsertFirstPage(t *testing.T) {
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
	contentInfo *PageParamContentInfo
	isPageInfo  bool
}

func ppciExNumberInPlainText(number int) *pageParamContentInfoEx {
	return &pageParamContentInfoEx{
		contentInfo: ppciNumberInPlainText(number),
	}
}

func ppciExNumericOutlink(targetURL string, number int, isPageInfo bool) *pageParamContentInfoEx {
	return &pageParamContentInfoEx{
		contentInfo: ppciNumericOutlink(targetURL, number),
		isPageInfo:  isPageInfo,
	}
}

func canInsertFirstPage(docURL string, allContentInfo []*pageParamContentInfoEx, pageParamInfo *info.PageParamInfo) bool {
	ascendingNumbers := []*info.PageInfo{}
	pageParamInfo.AllPageInfo = []*info.PageInfo{}

	for _, ex := range allContentInfo {
		switch ex.contentInfo.Type {
		case NumberInPlainText:
			ascendingNumbers = append(ascendingNumbers, &info.PageInfo{
				PageNumber: ex.contentInfo.Number,
			})

		case NumericOutlink:
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
