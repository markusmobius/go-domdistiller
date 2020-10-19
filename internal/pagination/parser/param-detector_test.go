// ORIGINAL: java/PageParameterDetector.java

package parser_test

import (
	"fmt"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/markusmobius/go-domdistiller/internal/pagination/parser"
	tu "github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Parser_Detector_DynamicPara(t *testing.T) {
	testUrl := "http://bbs.globalimporter.net/bbslist-N11101107-33-0.htm"
	outlink1 := "http://bbs.globalimporter.net/bbslist-N11101107-32-0.htm"
	outlink2 := "http://bbs.globalimporter.net/bbslist-N11101107-34-0.htm"
	outlink3 := "http://bbs.globalimporter.net/bbslist-N11101107-35-0.htm"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 32),
		tu.PPCINumberInPlainText(33),
		tu.PPCINumericOutlink(outlink2, 34),
		tu.PPCINumericOutlink(outlink3, 35),
	}

	paramInfo := detectPageParameter(testUrl, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 32, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 34, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 35, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink2, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_DynamicParaForComma(t *testing.T) {
	testUrl := "http://forum.interia.pl/forum/praca,1094,2,2515,0,0"
	outlink1 := "http://forum.interia.pl/forum/praca,1094,2,2515,0,0"
	outlink2 := "http://forum.interia.pl/forum/praca,1094,3,2515,0,0"
	outlink3 := "http://forum.interia.pl/forum/praca,1094,4,2515,0,0"
	outlink4 := "http://forum.interia.pl/forum/praca,1094,5,2515,0,0"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCINumericOutlink(outlink2, 3),
		tu.PPCINumericOutlink(outlink3, 4),
		tu.PPCINumericOutlink(outlink4, 5),
	}

	paramInfo := detectPageParameter(testUrl, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 4)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 5, page.PageNumber)
	assert.Equal(t, outlink4, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink2, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_StaticPara(t *testing.T) {
	testUrl := "http://www.google.com/forum"
	outlink1 := "http://www.google.com/forum?page=2"
	outlink2 := "http://www.google.com/forum?page=3"
	outlink3 := "http://www.google.com/forum?page=4"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCINumericOutlink(outlink2, 3),
		tu.PPCINumericOutlink(outlink3, 4),
	}

	paramInfo := detectPageParameter(testUrl, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 4)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, testUrl, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink1, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_StaticParaPageAtSuffix(t *testing.T) {
	testURL := "http://www.google.com/forum/0"
	outlink1 := "http://www.google.com/forum/1"
	outlink2 := "http://www.google.com/forum/2"
	outlink3 := "http://www.google.com/forum/3"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCINumericOutlink(outlink2, 3),
		tu.PPCINumericOutlink(outlink3, 4),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, -1, paramInfo.Formula.Delta)
	assert.Equal(t, outlink1, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_HandleOnlyHasPreviousPage(t *testing.T) {
	// The current doc is the 2nd page and has only digital outlink which points to 1st page.
	testURL := "http://www.google.com/forum/thread-20"
	outlink1 := "http://www.google.com/forum/thread-0"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
		tu.PPCINumberInPlainText(2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 20, paramInfo.Formula.Coefficient)
	assert.Equal(t, -20, paramInfo.Formula.Delta)
	assert.True(t, paramInfo.NextPagingURL == "")
}

func Test_Pagination_Parser_Detector_RejectOnlyPage1LinkWithoutCurrentPageText(t *testing.T) {
	// Although there is a digital outlink to 1st page, there is no plain text "2" after it, so
	// there is no pagination.
	testURL := "http://www.google.com/forum/thread-20"
	outlink1 := "http://www.google.com/forum/thread-0"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_HandleOnlyHasNextPage(t *testing.T) {
	// The current doc is the 1st page and has only digital outlink which points to 2nd page.
	testURL := "http://www.google.com/forum?page=0"
	outlink1 := "http://www.google.com/forum?page=1"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink(outlink1, 2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, -1, paramInfo.Formula.Delta)
	assert.Equal(t, outlink1, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_RejectOnlyPage2LinkWithoutCurrentPageText(t *testing.T) {
	// Although there is a digital outlink to 2nd page, there is no plain text "1" before it, so
	// there is no pagination.
	testURL := "http://www.google.com/forum?page=0"
	outlink1 := "http://www.google.com/forum?page=1"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_HandleOnlyTwoPartialPages(t *testing.T) {
	testURL := "http://www.google.com/forum/1"
	outlink1 := "http://www.google.com/forum/0"
	outlink2 := "http://www.google.com/forum/2"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
		tu.PPCINumberInPlainText(2),
		tu.PPCINumericOutlink(outlink2, 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, -1, paramInfo.Formula.Delta)
	assert.Equal(t, outlink2, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_HandleCalendarPageDynamic(t *testing.T) {
	testURL := "http://www.google.com/forum/1"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://www.google.com/forum/1?m=20101201", 1),
		tu.PPCINumericOutlink("http://www.google.com/forum/1?m=20101202", 2),
		tu.PPCINumericOutlink("http://www.google.com/forum/1?m=20101204", 4),
		tu.PPCINumericOutlink("http://www.google.com/forum/1?m=20101206", 6),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_HandleCalendarPageStatic(t *testing.T) {
	testURL := "http://www.google.com/forum/foo"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://www.google.com/forum/foo/2010/12/01", 1),
		tu.PPCINumericOutlink("http://www.google.com/forum/foo/2010/12/02", 2),
		tu.PPCINumericOutlink("http://www.google.com/forum/foo/2010/12/04", 4),
		tu.PPCINumericOutlink("http://www.google.com/forum/foo/2010/12/06", 6),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_HandleOutlinksInCalendar(t *testing.T) {
	testURL := "http://www.google.com/forum/foo"
	allContentInfo := make([]*tu.PageParamContentInfo, 32)
	for i := 0; i < 32; i++ {
		outlink := fmt.Sprintf("http://www.google.com/forum/foo/%d", (i + 1))
		allContentInfo[i] = tu.PPCINumericOutlink(outlink, i+1)
	}

	// Outlinks "1" through "30", likely to be a calendar.
	paramInfo := detectPageParameterWithSize(testURL, allContentInfo, 30)
	expectEmptyPageParamInfo(t, paramInfo)

	// 27 outlinks, not a calendar.
	paramInfo = detectPageParameterWithSize(testURL, allContentInfo, 27)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 27)

	// 32 outlinks, not a calendar.
	paramInfo = detectPageParameterWithSize(testURL, allContentInfo, 32)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 32)
}

func Test_Pagination_Parser_Detector_FilterYearCalendar(t *testing.T) {
	testURL := "http://technet.microsoft.com/en-us/sharepoint/ff628785.aspx"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://technet.microsoft.com/en-us/sharepoint/fp123618", 2010),
		tu.PPCINumericOutlink("http://technet.microsoft.com/en-us/sharepoint/fp142366", 2013),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_HandleFirstPageWithoutParam1(t *testing.T) {
	// Some pages only have two partial pages while the first page doesn't specify the page
	// parameter, like www.slate.com/id/2278628.
	testURL := "http://www.google.com/id/2278628"
	outlink1 := "http://www.google.com/id/2278628/pagenum/2"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink(outlink1, 2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink1, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_HandleFirstPageWithoutParam2(t *testing.T) {
	// When there is only one digital link on the page, we insert the doc url into
	// page num vector as well.
	testURL := "http://www.google.com/id/2278628/pagenum/2"
	outlink1 := "http://www.google.com/id/2278628"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
		tu.PPCINumberInPlainText(2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, "", paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_HandleFirstPageWithoutParam3(t *testing.T) {
	testURL := "http://www.google.com/id/2278628/pagenum/2"
	outlink1 := "http://www.google.com/id/2278628"
	outlink2 := "http://www.google.com/id/2278628/pagenum/3"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
		tu.PPCINumberInPlainText(2),
		tu.PPCINumericOutlink(outlink2, 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink2, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_HandleFirstPageWithoutParam4(t *testing.T) {
	testURL := "http://www.google.com/id/2278628"
	outlink1 := "http://www.google.com/id/2278628?pagenum=2"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink(outlink1, 2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink1, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_BadPageParamName(t *testing.T) {
	testURL := "http://www.google.com/12345/tag/1"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink("http://www.google.com/12345/tag/2.html", 2),
		tu.PPCINumericOutlink("http://www.google.com/12345/tag/3.html", 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_CheckFirstPathComponent(t *testing.T) {
	// TODO(radhi): To be honest I'm not sure why original dom-distiller expect the param
	// info is empty, as far as I check the original code this should gives a valid param
	// info. Anyway, I've put a check in pattern.NewPathComponentPagePattern to make this
	// package works similar with the original dom-distiller.
	testURL := "http://www.google.com/1"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumberInPlainText(1),
		tu.PPCINumericOutlink("http://www.google.com/2", 2),
		tu.PPCINumericOutlink("http://www.google.com/3", 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_PatternValidation(t *testing.T) {
	// If the page param is or part of a path component, both the pattern and document url
	// must have similar path: "article" differs from "download/foo".
	{
		testURL := "http://www.google.com/download/foo"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumberInPlainText(1),
			tu.PPCINumericOutlink("http://www.google.com/article/2", 2),
			tu.PPCINumericOutlink("http://www.google.com/article/3", 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		expectEmptyPageParamInfo(t, paramInfo)
	}

	// If the page param is a query, both the pattern and document url must have same path
	// components: "article/foo" differs from "article/bar".
	{
		testURL := "http://www.google.com/article/foo"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumberInPlainText(1),
			tu.PPCINumericOutlink("http://www.google.com/article/bar?page=2", 2),
			tu.PPCINumericOutlink("http://www.google.com/article/bar?page=3", 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		expectEmptyPageParamInfo(t, paramInfo)
	}

	{
		testURL := "http://www.google.com/foo-bar-article"
		outlink1 := "http://www.google.com/foo-bar-article-2"
		outlink2 := "http://www.google.com/foo-bar-article-3"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumberInPlainText(1),
			tu.PPCINumericOutlink(outlink1, 2),
			tu.PPCINumericOutlink(outlink2, 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		assert.Equal(t, info.PageNumber, paramInfo.Type)
		assert.Len(t, paramInfo.AllPageInfo, 3)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 1, page.PageNumber)
		assert.Equal(t, testURL, page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, outlink1, page.URL)

		page = paramInfo.AllPageInfo[2]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, outlink2, page.URL)

		assert.NotNil(t, paramInfo.Formula)
		assert.Equal(t, 1, paramInfo.Formula.Coefficient)
		assert.Equal(t, 0, paramInfo.Formula.Delta)
		assert.Equal(t, outlink1, paramInfo.NextPagingURL)
	}

	// If the page param is or part of a path component, both the pattern and document url
	// must have similar path: "foo/bar/article" differs from "foo-bar-article".
	{
		testURL := "http://www.google.com/foo/bar/article"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumberInPlainText(1),
			tu.PPCINumericOutlink("http://www.google.com/foo-bar-article-2", 2),
			tu.PPCINumericOutlink("http://www.google.com/foo-bar-article-3", 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		expectEmptyPageParamInfo(t, paramInfo)
	}
}

func Test_Pagination_Parser_Detector_TooManyPagingDocuments(t *testing.T) {
	docURL := "http://www.google.com/thread/page/1"

	pages := &info.MonotonicPageInfoGroups{}
	pages.AddGroup()

	for i := 1; i <= parser.MaxPagingDocs; i++ {
		pages.AddNumber(i, fmt.Sprintf("http://www.google.com/thread/page/%d", i))
	}

	paramInfo := parser.DetectParamInfo(pages, docURL)
	assert.Len(t, paramInfo.AllPageInfo, parser.MaxPagingDocs)

	pages.AddNumber(parser.MaxPagingDocs+1,
		fmt.Sprintf("http://www.google.com/thread/page/%d", parser.MaxPagingDocs+1))
	paramInfo = parser.DetectParamInfo(pages, docURL)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_InsertFirstPage(t *testing.T) {
	// The current document is inserted as first page.
	{
		testURL := "http://www.google.com/article/bar"
		outlink1 := "http://www.google.com/article/bar?page=2"
		outlink2 := "http://www.google.com/article/bar?page=3"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumberInPlainText(1),
			tu.PPCINumericOutlink(outlink1, 2),
			tu.PPCINumericOutlink(outlink2, 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		assert.Equal(t, info.PageNumber, paramInfo.Type)
		assert.Len(t, paramInfo.AllPageInfo, 3)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 1, page.PageNumber)
		assert.Equal(t, testURL, page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, outlink1, page.URL)

		page = paramInfo.AllPageInfo[2]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, outlink2, page.URL)

		assert.NotNil(t, paramInfo.Formula)
		assert.Equal(t, 1, paramInfo.Formula.Coefficient)
		assert.Equal(t, 0, paramInfo.Formula.Delta)
		assert.Equal(t, outlink1, paramInfo.NextPagingURL)
	}

	// Current document url has same length as other paginated pages, so it shouldn't be
	// inserted as first page.
	{
		testURL := "http://www.google.com/article/bar?page=1"
		outlink1 := "http://www.google.com/article/bar?page=2"
		outlink2 := "http://www.google.com/article/bar?page=3"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumberInPlainText(1),
			tu.PPCINumericOutlink(outlink1, 2),
			tu.PPCINumericOutlink(outlink2, 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		assert.Equal(t, info.PageNumber, paramInfo.Type)
		assert.Len(t, paramInfo.AllPageInfo, 2)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, outlink1, page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, outlink2, page.URL)

		assert.NotNil(t, paramInfo.Formula)
		assert.Equal(t, 1, paramInfo.Formula.Coefficient)
		assert.Equal(t, 0, paramInfo.Formula.Delta)
		assert.Equal(t, outlink1, paramInfo.NextPagingURL)
	}

	// Current document url is the last page, so shouldn't be inserted as first page.
	{
		testURL := "http://www.google.com/article/bar?page=4"
		outlink1 := "http://www.google.com/article/bar"
		outlink2 := "http://www.google.com/article/bar?page=2&s=11"
		outlink3 := "http://www.google.com/article/bar?page=3&s=11"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumericOutlink(outlink1, 1),
			tu.PPCINumericOutlink(outlink2, 2),
			tu.PPCINumericOutlink(outlink3, 3),
			tu.PPCINumberInPlainText(4),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		assert.Equal(t, info.PageNumber, paramInfo.Type)
		assert.Len(t, paramInfo.AllPageInfo, 2)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, outlink2, page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, outlink3, page.URL)

		assert.NotNil(t, paramInfo.Formula)
		assert.Equal(t, 1, paramInfo.Formula.Coefficient)
		assert.Equal(t, 0, paramInfo.Formula.Delta)
		assert.Equal(t, "", paramInfo.NextPagingURL)
	}

	// Page one has an outlink to itself, should be inserted as first page.
	{
		testURL := "http://www.google.com/article/bar"
		outlink1 := "http://www.google.com/article/bar"
		outlink2 := "http://www.google.com/article/bar?page=2"
		outlink3 := "http://www.google.com/article/bar?page=3"
		allContentInfo := []*tu.PageParamContentInfo{
			tu.PPCINumericOutlink(outlink1, 1),
			tu.PPCINumericOutlink(outlink2, 2),
			tu.PPCINumericOutlink(outlink3, 3),
		}

		paramInfo := detectPageParameter(testURL, allContentInfo)
		assert.Equal(t, info.PageNumber, paramInfo.Type)
		assert.Len(t, paramInfo.AllPageInfo, 3)

		page := paramInfo.AllPageInfo[0]
		assert.Equal(t, 1, page.PageNumber)
		assert.Equal(t, outlink1, page.URL)

		page = paramInfo.AllPageInfo[1]
		assert.Equal(t, 2, page.PageNumber)
		assert.Equal(t, outlink2, page.URL)

		page = paramInfo.AllPageInfo[2]
		assert.Equal(t, 3, page.PageNumber)
		assert.Equal(t, outlink3, page.URL)

		assert.NotNil(t, paramInfo.Formula)
		assert.Equal(t, 1, paramInfo.Formula.Coefficient)
		assert.Equal(t, 0, paramInfo.Formula.Delta)
		assert.Equal(t, outlink2, paramInfo.NextPagingURL)
	}
}

func Test_Pagination_Parser_Detector_RejectNonAdjacentOutlinks(t *testing.T) {
	testURL := "http://forum.interia.pl/forum/praca,1094,2,2515,0,0"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://forum.interia.pl/forum/praca,1094,2,2515,0,0", 2),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://forum.interia.pl/forum/praca,1094,3,2515,0,0", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://forum.interia.pl/forum/praca,1094,4,2515,0,0", 4),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://forum.interia.pl/forum/praca,1094,5,2515,0,0", 5),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Parser_Detector_AcceptAdjacentOutlinks(t *testing.T) {
	testURL := "http://www.google.com/test/list.php?start=0"
	outlink1 := "http://www.google.com/test/foo/2"
	outlink2 := "http://www.google.com/test/list.php?start=0"
	outlink3 := "http://www.google.com/test/list.php?start=10"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink(outlink2, 1),
		tu.PPCINumericOutlink(outlink3, 2),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, -10, paramInfo.Formula.Delta)
	assert.Equal(t, outlink3, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_AcceptDuplicatePatterns(t *testing.T) {
	testURL := "http://forum.interia.pl/forum/praca,1094,2,2515,0,0"
	outlink1 := "http://forum.interia.pl/forum/praca,1094,2,2515,0,0"
	outlink2 := "http://forum.interia.pl/forum/praca,1094,3,2515,0,0"
	outlink3 := "http://forum.interia.pl/forum/praca,1094,4,2515,0,0"
	outlink4 := "http://forum.interia.pl/forum/praca,1094,5,2515,0,0"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCINumericOutlink(outlink2, 3),
		tu.PPCINumericOutlink(outlink3, 4),
		tu.PPCINumericOutlink(outlink4, 5),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCINumericOutlink(outlink2, 3),
		tu.PPCINumericOutlink(outlink3, 4),
		tu.PPCINumericOutlink(outlink4, 5),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 4)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 5, page.PageNumber)
	assert.Equal(t, outlink4, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 1, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, outlink2, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_PreferPageNumberThanPageSize(t *testing.T) {
	testURL := "http://www.google.com/test/list.php"
	outlink1 := "http://www.google.com/test/list.php"
	outlink2 := "http://www.google.com/test/list.php?start=10"
	outlink3 := "http://www.google.com/test/list.php?start=20"
	outlink4 := "http://www.google.com/test/list.php?start=30"
	outlink5 := "http://www.google.com/test/list.php?size=20"
	outlink6 := "http://www.google.com/test/list.php?size=50"
	outlink7 := "http://www.google.com/test/list.php?size=100"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
		tu.PPCINumericOutlink(outlink2, 2),
		tu.PPCINumericOutlink(outlink3, 3),
		tu.PPCINumericOutlink(outlink4, 4),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink(outlink5, 20),
		tu.PPCINumericOutlink(outlink6, 50),
		tu.PPCINumericOutlink(outlink7, 100),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 4)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, outlink4, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, -10, paramInfo.Formula.Delta)
	assert.Equal(t, outlink2, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_RejectMultiPageNumberPatterns(t *testing.T) {
	testURL := "http://www.google.com/test/list.php"
	outlink1 := "http://www.google.com/test/list.php?start=10"
	outlink2 := "http://www.google.com/test/list.php?start=20"
	outlink3 := "http://www.google.com/test/list.php?start=30"
	outlink4 := "http://www.google.com/test/list.php?offset=20"
	outlink5 := "http://www.google.com/test/list.php?offset=30"
	outlink6 := "http://www.google.com/test/list.php?offset=40"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 2),
		tu.PPCINumericOutlink(outlink2, 3),
		tu.PPCINumericOutlink(outlink3, 4),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink(outlink4, 2),
		tu.PPCINumericOutlink(outlink5, 3),
		tu.PPCINumericOutlink(outlink6, 4),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 4)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, testURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, -10, paramInfo.Formula.Delta)
	assert.Equal(t, outlink1, paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_PreferLinearFormulaPattern(t *testing.T) {
	testURL := "http://www.google.com/test/list.php"
	outlink1 := "http://www.google.com/test/list.php?start=10"
	outlink2 := "http://www.google.com/test/list.php?start=20"
	outlink3 := "http://www.google.com/test/list.php?start=30"
	outlink4 := "http://www.google.com/test/list.php?size=21324235"
	outlink5 := "http://www.google.com/test/list.php?size=21435252"
	outlink6 := "http://www.google.com/test/list.php?size=32523516"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink(outlink1, 1),
		tu.PPCINumericOutlink(outlink2, 2),
		tu.PPCINumericOutlink(outlink3, 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink(outlink4, 1),
		tu.PPCINumericOutlink(outlink5, 2),
		tu.PPCINumericOutlink(outlink6, 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, outlink1, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, outlink2, page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, outlink3, page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, "", paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_PreferLinearFormulaPattern2(t *testing.T) {
	testURL := "http://www.google.com/test/list.php"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=21324235", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=21435252", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=32523516", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=21324235", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=21435252", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=32523516", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=10", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=20", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=30", 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=10", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=20", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=30", page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, "", paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_PreferLinearFormulaPattern3(t *testing.T) {
	testURL := "http://www.google.com/test/list.php"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=21324235", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=21435252", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=32523516", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=10", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=20", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=30", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=21324235", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=21435252", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=32523516", 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=10", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=20", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=30", page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, "", paramInfo.NextPagingURL)
}

func Test_Pagination_Parser_Detector_PreferLinearFormulaPattern4(t *testing.T) {
	testURL := "http://www.google.com/test/list.php"
	allContentInfo := []*tu.PageParamContentInfo{
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=10", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=20", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?start=30", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=21324235", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=21435252", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?size=32523516", 3),
		tu.PPCIUnrelatedTerms(),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=21324235", 1),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=21435252", 2),
		tu.PPCINumericOutlink("http://www.google.com/test/list.php?foo=32523516", 3),
	}

	paramInfo := detectPageParameter(testURL, allContentInfo)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=10", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=20", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=30", page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, 0, paramInfo.Formula.Delta)
	assert.Equal(t, "", paramInfo.NextPagingURL)
}

func detectPageParameterWithSize(docURL string, allContentInfo []*tu.PageParamContentInfo, contentInfoSize int) *info.PageParamInfo {
	adjacentNumberGroups := &info.MonotonicPageInfoGroups{}
	adjacentNumberGroups.AddGroup()

	for i := 0; i < contentInfoSize; i++ {
		content := allContentInfo[i]

		switch content.Type {
		case tu.UnrelatedTerms:
			adjacentNumberGroups.AddGroup()

		case tu.NumberInPlainText, tu.NumericOutlink:
			url := ""
			if content.Type == tu.NumericOutlink {
				url = content.TargetURL
			}

			adjacentNumberGroups.AddNumber(content.Number, url)
		}
	}

	return parser.DetectParamInfo(adjacentNumberGroups, docURL)
}

func detectPageParameter(docURL string, allContentInfo []*tu.PageParamContentInfo) *info.PageParamInfo {
	return detectPageParameterWithSize(docURL, allContentInfo, len(allContentInfo))
}

func expectEmptyPageParamInfo(t *testing.T, paramInfo *info.PageParamInfo) {
	assert.Equal(t, info.Unset, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 0)
	assert.Nil(t, paramInfo.Formula)
	assert.True(t, paramInfo.NextPagingURL == "")
}
