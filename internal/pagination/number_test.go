// ORIGINAL: javatest/PageParameterParserTest.java

package pagination_test

import (
	nurl "net/url"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/pagination"
	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	BaseURL = "http://www.test.com/"
	TestURL = BaseURL + "foo/bar"
)

func Test_Pagination_Number_Basic(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`1<br>` +
		`<a href="/foo/bar/2">2</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	paramInfo = processDefaultDocumentWithPageNum(`1<br>` +
		`<a href="/foo/bar/2">2</a>` +
		`<a href="/foo/bar/3">3</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 3)
}

func Test_Pagination_Number_RejectOnlyPage2LinkWithoutCurrentPageText(t *testing.T) {
	// Although there is a digital outlink to 2nd page, there is no plain text "1"
	// before it, so there is no pagination.
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`If there were a '1', pagination should be detected. But there isn't.` +
		`<a href="/foo/bar/2">2</a>` +
		`Main content`)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Number_RejectNonAdjacentOutlinks(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`1<br>` +
		`Unrelated terms<br>` +
		`<a href="/foo/bar/2">2</a>` +
		`Unrelated terms<br>` +
		`<a href="/foo/bar/3">3</a>` +
		`<a href="/foo/bar/all">All</a>`)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Number_AcceptAdjacentOutlinks(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`Unrelated link: <a href="http://www.test.com/other/2">2</a>` +
		`<p>Main content</p>` +
		`1<br>` +
		`<a href="http://www.test.com/foo/bar/2">2</a>` +
		`<a href="http://www.test.com/foo/bar/3">3</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar/2", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar/3", page.URL)

	assert.Equal(t, BaseURL+"foo/bar/2", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_AcceptDuplicatePatterns(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`1<br>` +
		`<a href="http://www.test.com/foo/bar/2">2</a>` +
		`<a href="http://www.test.com/foo/bar/3">3</a>` +
		`<p>Main content</p>` +
		`1<br>` +
		`<a href="http://www.test.com/foo/bar/2">2</a>` +
		`<a href="http://www.test.com/foo/bar/3">3</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar/2", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar/3", page.URL)

	assert.Equal(t, BaseURL+"foo/bar/2", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_PreferPageNumber(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`<a href="http://www.test.com/foo/bar/size-25">25</a>` +
		`<a href="http://www.test.com/foo/bar/size-50">50</a>` +
		`<a href="http://www.test.com/foo/bar/size-100">100</a>` +
		`<p>Main content</p>` +
		`1<br>` +
		`<a href="http://www.test.com/foo/bar/2">2</a>` +
		`<a href="http://www.test.com/foo/bar/3">3</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar/2", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, BaseURL+"foo/bar/3", page.URL)

	assert.Equal(t, BaseURL+"foo/bar/2", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_RejectMultiplePageNumberPatterns(t *testing.T) {
	paramInfo := processDocumentWithPageNum("http://www.google.com/test/list.php", ""+
		`<a href="http://www.google.com/test/list.php?start=10">2</a>`+
		`<a href="http://www.google.com/test/list.php?start=20">3</a>`+
		`<a href="http://www.google.com/test/list.php?start=30">4</a>`+
		`<p>Main content</p>`+
		`<a href="http://www.google.com/test/list.php?offset=10">2</a>`+
		`<a href="http://www.google.com/test/list.php?offset=20">3</a>`+
		`<a href="http://www.google.com/test/list.php?offset=30">4</a>`+
		`<a href="http://www.google.com/test/list.php?offset=all">All</a>`)
	assert.Equal(t, info.PageNumber, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 4)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=10", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=20", page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, "http://www.google.com/test/list.php?start=30", page.URL)

	assert.NotNil(t, paramInfo.Formula)
	assert.Equal(t, 10, paramInfo.Formula.Coefficient)
	assert.Equal(t, -10, paramInfo.Formula.Delta)
	assert.Equal(t, "http://www.google.com/test/list.php?start=10", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_InvalidAndVoidLinks(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`1<br>` +
		`<a href="javascript:void(0)">2</a>`)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Number_DifferentHostLinks(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`1<br>` +
		`<a href="http://www.foo.com/foo/bar/2">2</a>`)
	expectEmptyPageParamInfo(t, paramInfo)
}

func Test_Pagination_Number_WhitespaceSibling(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`1<br>` +
		`       ` +
		`<a href="/foo/bar/2">2</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 2)
}

func Test_Pagination_Number_PunctuationSibling(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`<a href="/foo/bar/1">1</a>` +
		`,` +
		`<a href="/foo/bar/2">2</a>`)
	assert.Len(t, paramInfo.AllPageInfo, 2)
}

func Test_Pagination_Number_SeparatorSibling(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`<div>` +
		`1 | ` +
		`<a href="/foo/bar/2">2</a>` +
		` | ` +
		`<a href="/foo/bar/3">3</a>` +
		`</div>`)
	assert.Len(t, paramInfo.AllPageInfo, 3)
}

func Test_Pagination_Number_ParentSibling0(t *testing.T) {
	paramInfo := processDefaultDocumentWithPageNum(`` +
		`<div>begin` +
		`<strong>1</strong>` +
		`<div><a href="http://www.test.com/foo/bar/2">2</a></div>` +
		`<div><a href="http://www.test.com/foo/bar/3">3</a></div>` +
		`end</div>`)
	assert.Len(t, paramInfo.AllPageInfo, 3)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, TestURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, TestURL+"/2", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, TestURL+"/3", page.URL)

	assert.Equal(t, "http://www.test.com/foo/bar/2", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_ParentSibling1(t *testing.T) {
	paramInfo := processDocumentWithPageNum("http://www.test.com/foo/bar/2", ``+
		`<div>begin`+
		`<div><a href="http://www.test.com/foo/bar">1</a></div>`+
		`<strong>2</strong>`+
		`<div><a href="http://www.test.com/foo/bar/3">3</a></div>`+
		`end</div>`)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, TestURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 3, page.PageNumber)
	assert.Equal(t, TestURL+"/3", page.URL)

	assert.Equal(t, "http://www.test.com/foo/bar/3", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_ParentSibling2(t *testing.T) {
	paramInfo := processDocumentWithPageNum("http://www.test.com/foo/bar/3", ``+
		`<div>begin`+
		`<div><a href="http://www.test.com/foo/bar">1</a></div>`+
		`<div><a href="http://www.test.com/foo/bar/2">2</a></div>`+
		`<strong>3</strong>`+
		`end</div>`)
	assert.Len(t, paramInfo.AllPageInfo, 2)

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, TestURL, page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, TestURL+"/2", page.URL)

	assert.Equal(t, "", paramInfo.NextPagingURL)
}

func Test_Pagination_Number_NestedStructure(t *testing.T) {
	paramInfo := processDocumentWithPageNum("http://www.test.com/foo?page=3", ``+
		`<div>begin`+
		`<span><a href="http://www.test.com/foo?page=2">&lsaquo;&lsaquo; Prev</a></span>`+
		`<span><a href="http://www.test.com/foo?page=1">1</a></span>`+
		`<span><a href="http://www.test.com/foo?page=2">2</a></span>`+
		`<span>3</span>`+
		`<span><a href="http://www.test.com/foo?page=4">4</a></span>`+
		`<span><a href="http://www.test.com/foo?page=5">5</a></span>`+
		`<span>...</span>`+
		`<span><a href="http://www.test.com/foo?page=48">48</a></span>`+
		`<span><a href="http://www.test.com/foo?page=4">Next &rsaquo;&rsaquo;</a></span>`+
		`</div>`)
	assert.Len(t, paramInfo.AllPageInfo, 5)

	urlPrefix := "http://www.test.com/foo?page="

	page := paramInfo.AllPageInfo[0]
	assert.Equal(t, 1, page.PageNumber)
	assert.Equal(t, urlPrefix+"1", page.URL)

	page = paramInfo.AllPageInfo[1]
	assert.Equal(t, 2, page.PageNumber)
	assert.Equal(t, urlPrefix+"2", page.URL)

	page = paramInfo.AllPageInfo[2]
	assert.Equal(t, 4, page.PageNumber)
	assert.Equal(t, urlPrefix+"4", page.URL)

	page = paramInfo.AllPageInfo[3]
	assert.Equal(t, 5, page.PageNumber)
	assert.Equal(t, urlPrefix+"5", page.URL)

	page = paramInfo.AllPageInfo[4]
	assert.Equal(t, 48, page.PageNumber)
	assert.Equal(t, urlPrefix+"48", page.URL)

	assert.Equal(t, urlPrefix+"4", paramInfo.NextPagingURL)
}

func processDocumentWithPageNum(originalURL string, content string) *info.PageParamInfo {
	pageURL, _ := nurl.ParseRequestURI(originalURL)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, content)

	numberFinder := pagination.NewPageNumberFinder(stringutil.FastWordCounter{}, nil, nil)
	return numberFinder.FindOutlink(root, pageURL)
}

func processDefaultDocumentWithPageNum(content string) *info.PageParamInfo {
	return processDocumentWithPageNum(TestURL, content)
}

func expectEmptyPageParamInfo(t *testing.T, paramInfo *info.PageParamInfo) {
	assert.Equal(t, info.Unset, paramInfo.Type)
	assert.Len(t, paramInfo.AllPageInfo, 0)
	assert.Nil(t, paramInfo.Formula)
	assert.True(t, paramInfo.NextPagingURL == "")
}
