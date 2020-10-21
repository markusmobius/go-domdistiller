// ORIGINAL: javatest/PagingLinksFinderTest.java

package pagination_test

import (
	nurl "net/url"
	"strings"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/pagination"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

const ExampleURL = "http://example.com/path/toward/news.php"

// There are some tests in original dom-distiller that not reproduced here
// because the structure of our code a bit different :
// - Test_Pagination_PrevNext_BaseTag

func Test_Pagination_PrevNext_NoLink(t *testing.T) {
	doc := testutil.CreateHTML()
	assertDefaultDocumenOutlink(t, doc, nil, nil)
}

func Test_Pagination_PrevNext_1NextLink(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("next", "next page")
	dom.AppendChild(root, anchor)

	assertDefaultDocumenOutlink(t, doc, nil, nil)
}

func Test_Pagination_PrevNext_1NextLinkWithDifferentDomain(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("http://testing.com/page2", "next page")
	dom.AppendChild(root, anchor)

	assertDefaultDocumenOutlink(t, doc, nil, nil)
}

func Test_Pagination_PrevNext_1NextLinkWithOriginalDomain(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("http://testing.com/page2", "next page")
	dom.AppendChild(root, anchor)

	assertDocumentOutlink(t, "http://testing.com", doc, nil, anchor)
}

func Test_Pagination_PrevNext_CaseInsensitive(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("HTTP://testing.COM/page2", "next page")
	dom.AppendChild(root, anchor)

	assertDocumentOutlink(t, "http://testing.com", doc, nil, anchor)
}

func Test_Pagination_PrevNext_1PageNumberedLink(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("page2", "page 2")
	dom.AppendChild(root, anchor)

	// The word "page" in the link text increases its score confidently enough to
	// be considered as the previous paging link.
	assertDefaultDocumenOutlink(t, doc, anchor, anchor)
}

func Test_Pagination_PrevNext_3NumberedLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor1 := testutil.CreateAnchor("page1", "1")
	anchor2 := testutil.CreateAnchor("page2", "2")
	anchor3 := testutil.CreateAnchor("page3", "3")
	dom.AppendChild(root, anchor1)
	dom.AppendChild(root, anchor2)
	dom.AppendChild(root, anchor3)

	// Because link text contains only digits with no paging-related words, no link
	// has a score high enough to be confidently considered paging link.
	assertDefaultDocumenOutlink(t, doc, nil, nil)
}

func Test_Pagination_PrevNext_2NextLinksWithSameHref(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor1 := testutil.CreateAnchor("page2", "dummy link")
	anchor2 := testutil.CreateAnchor("page2", "next page")
	dom.AppendChild(root, anchor1)
	assertDefaultDocumenOutlink(t, doc, nil, nil)

	// anchor1 is not a confident next page link, but anchor2 is due to the link text.
	dom.AppendChild(root, anchor2)
	assertDefaultDocumenOutlink(t, doc, nil, anchor1)
}

func Test_Pagination_PrevNext_PagingParent(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	div := testutil.CreateDiv(1)
	dom.SetAttribute(div, "class", "page")
	dom.AppendChild(root, div)

	anchor := testutil.CreateAnchor("page1", "dummy link")
	dom.AppendChild(div, anchor)

	// While it may seem strange that both previous and next links are the same, this test
	// is testing that the anchor's parents will affect its paging score even if it has a
	// meaningless link text like "dummy link".
	assertDefaultDocumenOutlink(t, doc, anchor, anchor)
}

func Test_Pagination_PrevNext_1PrevLink(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("prev", "prev page")
	dom.AppendChild(root, anchor)

	assertDefaultDocumenOutlink(t, doc, anchor, nil)
}

func Test_Pagination_PrevNext_PrevAnd1NextLink(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	prevAnchor := testutil.CreateAnchor("prev", "prev page")
	nextAnchor := testutil.CreateAnchor("page2", "next page")
	dom.AppendChild(root, prevAnchor)
	dom.AppendChild(root, nextAnchor)

	assertDefaultDocumenOutlink(t, doc, prevAnchor, nextAnchor)
}

func Test_Pagination_PrevNext_PopularBadLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	nextAnchor := testutil.CreateAnchor("page2", "next page")
	dom.AppendChild(root, nextAnchor)

	// If the same bad URL can get scores accumulated across links,
	// it would wrongly get selected.
	bad1 := testutil.CreateAnchor("not-page1", "not")
	bad2 := testutil.CreateAnchor("not-page1", "not")
	bad3 := testutil.CreateAnchor("not-page1", "not")
	bad4 := testutil.CreateAnchor("not-page1", "not")
	bad5 := testutil.CreateAnchor("not-page1", "not")
	dom.AppendChild(root, bad1)
	dom.AppendChild(root, bad2)
	dom.AppendChild(root, bad3)
	dom.AppendChild(root, bad4)
	dom.AppendChild(root, bad5)

	assertDefaultDocumenOutlink(t, doc, nil, nextAnchor)
}

func Test_Pagination_PrevNext_HeldBackLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	nextAnchor := testutil.CreateAnchor("page2", "next page")
	dom.AppendChild(root, nextAnchor)

	// If "page2" gets bad scores from other links, it would be missed.
	bad := testutil.CreateAnchor("page2", "prev or next")
	dom.AppendChild(root, bad)

	assertDefaultDocumenOutlink(t, doc, nil, nextAnchor)
}

func Test_Pagination_PrevNext_FirstPageLinkAsFolderURL(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	// Some sites' first page links are the same as the folder URL,
	// previous page link needs to recognize this.
	href := ExampleURL[:strings.LastIndex(ExampleURL, "/")]
	anchor := testutil.CreateAnchor(href, "PREV")
	dom.AppendChild(root, anchor)

	assertDefaultDocumenOutlink(t, doc, anchor, nil)
}

func Test_Pagination_PrevNext_NonHttpOrHttpsLink(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("javascript:void(0)", "next")
	dom.AppendChild(root, anchor)
	assertDefaultDocumentNextLink(t, doc, nil)

	dom.SetAttribute(anchor, "href", "file://test.html")
	assertDefaultDocumentNextLink(t, doc, nil)
}

func Test_Pagination_PrevNext_NextArticleLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor1 := testutil.CreateAnchor("page2", "next article")
	dom.AppendChild(root, anchor1)
	assertDefaultDocumentNextLink(t, doc, nil)

	// The banned word "article" also affects anchor2 because it has the same href as anchor
	anchor2 := testutil.CreateAnchor("page2", "next page")
	dom.AppendChild(root, anchor2)
	assertDefaultDocumentNextLink(t, doc, nil)

	// Removing the banned word revives the link
	dom.SetInnerHTML(anchor1, "next thing")
	assertDefaultDocumenOutlink(t, doc, nil, anchor1)
}

func Test_Pagination_PrevNext_NextChineseArticleLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.SetAttribute(root, "class", "page")
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("page2", "下一篇")
	dom.AppendChild(root, anchor)
	assertDefaultDocumentNextLink(t, doc, nil)

	dom.SetInnerHTML(anchor, "下一頁")
	assertDefaultDocumenOutlink(t, doc, anchor, anchor)
}

func Test_Pagination_PrevNext_NextPostLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("page2", "next post")
	dom.AppendChild(root, anchor)
	assertDefaultDocumentNextLink(t, doc, nil)
}

func Test_Pagination_PrevNext_AsOneLinks(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor("page2", "view as one page")
	dom.AppendChild(root, anchor)
	assertDefaultDocumentNextLink(t, doc, nil)

	dom.SetInnerHTML(anchor, "next")
	assertDefaultDocumenOutlink(t, doc, nil, anchor)
}

func Test_Pagination_PrevNext_LinksWithLongText(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor(ExampleURL+"/page2", "page 2 with long text")
	dom.AppendChild(root, anchor)
	assertDefaultDocumentNextLink(t, doc, nil)
}

func Test_Pagination_PrevNext_NonTailPageInfo(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	root := testutil.CreateDiv(0)
	dom.AppendChild(body, root)

	anchor := testutil.CreateAnchor(ExampleURL+"/gap/12/somestuff", "page down")
	dom.AppendChild(root, anchor)
	assertDefaultDocumentNextLink(t, doc, nil)
}

func assertDefaultDocumenOutlink(t *testing.T, doc *html.Node, prevAnchor, nextAnchor *html.Node) {
	assertDocumentOutlink(t, ExampleURL, doc, prevAnchor, nextAnchor)
}

func assertDocumentOutlink(t *testing.T, pageURL string, doc *html.Node, prevAnchor, nextAnchor *html.Node) {
	url, err := nurl.ParseRequestURI(pageURL)
	assert.NoError(t, err)
	assert.NotNil(t, url)

	prevHref := pagination.NewPrevNextFinder().FindOutlink(doc, url, false)
	if prevAnchor == nil {
		assert.Equal(t, "", prevHref)
	} else {
		linkHref := dom.GetAttribute(prevAnchor, "href")
		linkHref = normalizeLinkHref(linkHref, url)
		assert.Equal(t, linkHref, prevHref)
	}

	nextHref := pagination.NewPrevNextFinder().FindOutlink(doc, url, true)
	if nextAnchor == nil {
		assert.Equal(t, "", nextHref)
	} else {
		linkHref := dom.GetAttribute(nextAnchor, "href")
		linkHref = normalizeLinkHref(linkHref, url)
		assert.Equal(t, linkHref, nextHref)
	}
}

func assertDefaultDocumentNextLink(t *testing.T, doc *html.Node, anchor *html.Node) {
	assertDocumentNextLink(t, ExampleURL, doc, anchor)
}

func assertDocumentNextLink(t *testing.T, pageURL string, doc *html.Node, anchor *html.Node) {
	url, err := nurl.ParseRequestURI(pageURL)
	assert.NoError(t, err)
	assert.NotNil(t, url)

	nextHref := pagination.NewPrevNextFinder().FindOutlink(doc, url, true)
	if anchor == nil {
		assert.Equal(t, "", nextHref)
	} else {
		linkHref := dom.GetAttribute(anchor, "href")
		linkHref = normalizeLinkHref(linkHref, url)
		assert.Equal(t, linkHref, nextHref)
	}
}

func normalizeLinkHref(linkHref string, pageURL *nurl.URL) string {
	// Try to convert relative URL in link href to absolute URL
	linkHref = stringutil.CreateAbsoluteURL(linkHref, pageURL)

	// Make sure the link href is absolute
	_, err := nurl.ParseRequestURI(linkHref)
	if err != nil {
		return linkHref
	}

	// Remove url anchor and then trailing '/' from link's href.
	tmp, _ := nurl.Parse(linkHref)
	tmp.RawQuery = ""
	tmp.Fragment = ""
	tmp.RawFragment = ""
	tmp.Path = strings.TrimSuffix(tmp.Path, "/")
	tmp.RawPath = tmp.Path
	return tmp.String()
}
