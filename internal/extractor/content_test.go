// ORIGINAL: javatest/ContentExtractorTest.java

package extractor_test

import (
	nurl "net/url"
	"strings"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/extractor"
	"github.com/markusmobius/go-domdistiller/internal/markup/opengraph"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

const (
	contentText = "Lorem Ipsum Lorem Ipsum Lorem Ipsum."
	titleText   = "I am the document title"
)

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_Extractor_Content_HiddenArticle

func Test_Extractor_Content_DoesNotExtractTitleInContent(t *testing.T) {
	titleDiv := testutil.CreateDiv(0)
	dom.AppendChild(titleDiv, dom.CreateTextNode(titleText))

	contentDiv1 := testutil.CreateDiv(1)
	dom.AppendChild(contentDiv1, dom.CreateTextNode(contentText))

	contentDiv2 := testutil.CreateDiv(2)
	dom.AppendChild(contentDiv2, dom.CreateTextNode(contentText))

	contentDiv3 := testutil.CreateDiv(3)
	dom.AppendChild(contentDiv3, dom.CreateTextNode(contentText))

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, titleDiv)
	dom.AppendChild(body, contentDiv1)
	dom.AppendChild(body, contentDiv2)
	dom.AppendChild(body, contentDiv3)

	// Title hasn't been set yet, everything should be content.
	ce := extractor.NewContentExtractor(doc, nil, nil)
	extractedContent := extractContent(ce)
	assert.True(t, strings.Contains(extractedContent, domutil.InnerText(contentDiv1)))
	assert.True(t, strings.Contains(extractedContent, domutil.InnerText(titleDiv)))

	// Now set the title and it should excluded from the content.
	head := dom.QuerySelector(doc, "head")
	dom.AppendChild(head, testutil.CreateTitle(titleText))

	ce = extractor.NewContentExtractor(doc, nil, nil)
	extractedContent = extractContent(ce)
	assert.True(t, strings.Contains(extractedContent, domutil.InnerText(contentDiv1)))
	assert.False(t, strings.Contains(extractedContent, domutil.InnerText(titleDiv)))
}

func Test_Extractor_Content_ExtractsEssentialWhitespace(t *testing.T) {
	div := testutil.CreateDiv(0)
	dom.AppendChild(div, testutil.CreateSpan(contentText))
	dom.AppendChild(div, dom.CreateTextNode(" "))
	dom.AppendChild(div, testutil.CreateSpan(contentText))
	dom.AppendChild(div, dom.CreateTextNode("\n"))
	dom.AppendChild(div, testutil.CreateSpan(contentText))
	dom.AppendChild(div, dom.CreateTextNode(" "))

	doc, body := createHTML()
	dom.AppendChild(body, div)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	extractedContent := extractContent(ce)
	assert.Equal(t, "<div><span>"+contentText+"</span> "+
		"<span>"+contentText+"</span>\n"+
		"<span>"+contentText+"</span> </div>",
		testutil.RemoveAllDirAttributes(extractedContent))
}

func Test_Extractor_Content_PrefersMarkupParserOverDocumentTitle(t *testing.T) {
	// Minimum fields for open-graph parser.
	markupParserTitle := "title from markup parser"
	doc := testutil.CreateHTML()
	createMeta(doc, "og:title", markupParserTitle)
	createMeta(doc, "og:type", "video.movie")
	createMeta(doc, "og:image", "http://test/image.jpeg")
	createMeta(doc, "og:url", "http://test/test.html")

	// Make sure open-graph parser works properly
	parser, _ := opengraph.NewParser(doc, nil)
	assert.NotNil(t, parser)
	assert.Equal(t, markupParserTitle, parser.Title())

	// Open-graph title should be picked over document title
	head := dom.QuerySelector(doc, "head")
	dom.AppendChild(head, testutil.CreateTitle(titleText))

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, markupParserTitle, ce.ExtractTitle())
}

func Test_Extractor_Content_Image(t *testing.T) {
	rawHTML := `<h1>` + contentText + `</h1>` +
		`<img id="a" style="typo" align="left" src="image" srcset="image200 200w, //example.org/image400 400w"/>` +
		`<figure><picture>` +
		`<source srcset="image200 200w, //example.org/image400 400w"/>` +
		`<source srcset="image100 100w, //example.org/image300 300w"/>` +
		`<img/>` +
		`</picture></figure>` +
		`<span class="lazy-image-placeholder" data-src="/image" ` +
		`data-srcset="/image2x 2x" data-width="20" data-height="10"></span>` +
		`<img id="b" style="align: left" alt="b" data-dummy="c" data-src="image2"/>` +
		`<table role="grid"><tbody><tr><td>` +
		`<img id="c" style="a" alt="b" src="/image" srcset="https://example.com/image2x 2x, /image4x 4x,"/>` +
		`<img id="d" style="a" align="left" src="/image2"/>` +
		`</td></tr></tbody></table>` +
		`<p>` + contentText + `</p>`

	expected := `<h1>` + contentText + `</h1>` +
		`<img src="http://example.com/path/image" ` +
		`srcset="http://example.com/path/image200 200w, http://example.org/image400 400w"/>` +
		`<figure><picture>` +
		`<source srcset="http://example.com/path/image200 200w, http://example.org/image400 400w"/>` +
		`<source srcset="http://example.com/path/image100 100w, http://example.org/image300 300w"/>` +
		`<img/>` +
		`</picture></figure>` +
		`<img src="http://example.com/image" srcset="http://example.com/image2x 2x" ` +
		`width="20" height="10"/>` +
		`<img alt="b" src="http://example.com/path/image2"/>` +
		`<table role="grid"><tbody><tr><td>` +
		`<img alt="b" src="http://example.com/image" ` +
		`srcset="https://example.com/image2x 2x, http://example.com/image4x 4x,"/>` +
		`<img src="http://example.com/image2"/>` +
		`</td></tr></tbody></table>` +
		`<p>` + contentText + `</p>`

	doc, body := createHTML()
	dom.SetInnerHTML(body, rawHTML)

	pageURL, _ := nurl.ParseRequestURI("http://example.com/path/")
	ce := extractor.NewContentExtractor(doc, pageURL, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func Test_Extractor_Content_RemoveFontColorAttributes(t *testing.T) {
	// The result of this test is different with the original dom-distiller because
	// we can't compute stylesheet that needed in . NEED-COMPUTE-CSS
	outerFontTag := dom.CreateElement("font")
	dom.SetAttribute(outerFontTag, "color", "blue")

	text := `<font color="red">` + contentText + `</font>`
	dom.AppendChild(outerFontTag, testutil.CreateSpan(text))
	dom.AppendChild(outerFontTag, dom.CreateTextNode(" "))
	dom.AppendChild(outerFontTag, testutil.CreateSpan(text))
	dom.AppendChild(outerFontTag, dom.CreateTextNode("\n"))
	dom.AppendChild(outerFontTag, testutil.CreateSpan(text))
	dom.AppendChild(outerFontTag, dom.CreateTextNode(" "))

	doc, body := createHTML()
	dom.AppendChild(body, outerFontTag)

	expected := "<div><span>" +
		"<span><span>" + contentText + "</span></span> " +
		"<span><span>" + contentText + "</span></span>\n" +
		"<span><span>" + contentText + "</span></span> " +
		"</span></div>"

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func Test_Extractor_Content_RemoveStyleAttributes(t *testing.T) {
	rawHTML := `<h1 style="font-weight: folder">` + contentText + `</h1>` +
		`<p style="">` + contentText + `</p>` +
		`<img style="align: left" data-src="/test.png"/>` +
		`<table style="position: absolute">` +
		`<tbody style="font-size: 2">` +
		`<tr style="z-index: 0">` +
		`<th style="top: 0px">` + contentText +
		`<img style="align: left" src="/test.png"/>` +
		`</th>` +
		`<th style="width: 20px">` + contentText + `</th>` +
		`</tr><tr style="left: 0">` +
		`<td style="display: block">` + contentText + `</td>` +
		`<td style="color: #123">` + contentText + `</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	expected := `<h1>` + contentText + `</h1>` +
		`<p>` + contentText + `</p>` +
		`<img src="http://example.com/test.png"/>` +
		`<table>` +
		`<tbody>` +
		`<tr>` +
		`<th>` + contentText +
		`<img src="http://example.com/test.png"/>` +
		`</th>` +
		`<th>` + contentText + `</th>` +
		`</tr><tr>` +
		`<td>` + contentText + `</td>` +
		`<td>` + contentText + `</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	doc, body := createHTML()
	dom.SetInnerHTML(body, rawHTML)

	pageURL, _ := nurl.ParseRequestURI("http://example.com/path/")
	ce := extractor.NewContentExtractor(doc, pageURL, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func Test_Extractor_Content_RemoveNonAllowlistedAttributes(t *testing.T) {
	rawHTML := `<h1 onclick="alert(0);">` +
		contentText +
		`</h1>` +
		`<p larry="console.error(0);">` +
		contentText +
		`</p>` +
		`<img sergey="alert(0);" data-src="/test.png">` +
		`<video onkeydown="window.location.href = 'foo';">` +
		`<source src="http://example.com/foo.ogg">` +
		`<track src="http://example.com/foo.vtt">` +
		`</video>` +
		`<table onscroll="new XMLHttpRequest();">` +
		`<tbody>` +
		`<tr larry="1">` +
		`<th>` + contentText +
		`<img src="/test.png">` +
		`</th>` +
		`<th sergey="2">` + contentText + `</th>` +
		`</tr><tr>` +
		`<td>` + contentText + `</td>` +
		`<td>` + contentText + `</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	expected := `<h1>` + contentText + `</h1>` +
		`<p>` + contentText + `</p>` +
		`<img src="http://example.com/test.png"/>` +
		`<video>` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>` +
		`<table>` +
		`<tbody>` +
		`<tr>` +
		`<th>` + contentText +
		`<img src="http://example.com/test.png"/>` +
		`</th>` +
		`<th>` + contentText + `</th>` +
		`</tr><tr>` +
		`<td>` + contentText + `</td>` +
		`<td>` + contentText + `</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	doc, body := createHTML()
	dom.SetInnerHTML(body, rawHTML)

	pageURL, _ := nurl.ParseRequestURI("http://example.com/path/")
	ce := extractor.NewContentExtractor(doc, pageURL, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func Test_Extractor_Content_KeepingWidthAndHeightAttributes(t *testing.T) {
	rawHTML := `<h1>` + contentText + `</h1>` +
		`<p>` + contentText + `</p>` +
		`<img style="align: left" src="/test.png" width="200" height="300"/>` +
		`<img style="align: left" src="/test.png" width="200"/>` +
		`<img style="align: left" src="/test.png"/>`

	expected := `<h1>` + contentText + `</h1>` +
		`<p>` + contentText + `</p>` +
		`<img src="http://example.com/test.png" width="200" height="300"/>` +
		`<img src="http://example.com/test.png" width="200"/>` +
		`<img src="http://example.com/test.png"/>`

	doc, body := createHTML()
	dom.SetInnerHTML(body, rawHTML)

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	ce := extractor.NewContentExtractor(doc, pageURL, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func Test_Extractor_Content_PreserveOrderedList(t *testing.T) {
	outerList := dom.CreateElement("ol")
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol>", extractContent(ce))
}

func Test_Extractor_Content_PreserveOrderedListWithSpan(t *testing.T) {
	li := dom.CreateElement("li")
	dom.AppendChild(li, testutil.CreateSpan(contentText))

	outerList := dom.CreateElement("ol")
	dom.AppendChild(outerList, li)
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol>"+
		"<li><span>"+contentText+"</span></li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol>", extractContent(ce))
}

func Test_Extractor_Content_PreserveNestedOrderedList(t *testing.T) {
	innerList := dom.CreateElement("ol")
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))

	outerListItem := dom.CreateElement("li")
	dom.AppendChild(outerListItem, innerList)

	outerList := dom.CreateElement("ol")
	dom.AppendChild(outerList, outerListItem)
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol>"+
		"<li>"+"<ol>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol>"+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol>", extractContent(ce))
}

func Test_Extractor_Content_PreserveNestedOrderedListWithOtherElementsInside(t *testing.T) {
	innerList := dom.CreateElement("ol")
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateParagraph(""))

	outerListItem := dom.CreateElement("li")
	dom.AppendChild(outerListItem, dom.CreateTextNode(contentText))
	dom.AppendChild(outerListItem, testutil.CreateParagraph(contentText))
	dom.AppendChild(outerListItem, innerList)

	outerList := dom.CreateElement("ol")
	dom.AppendChild(outerList, outerListItem)
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateParagraph(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol>"+
		"<li>"+contentText+
		"<p>"+contentText+"</p>"+
		"<ol>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol>"+
		"</li>"+
		"<li>"+contentText+"</li>"+
		"<p>"+contentText+"</p>"+
		"</ol>", extractContent(ce))
}

func Test_Extractor_Content_PreserveUnorderedList(t *testing.T) {
	outerList := dom.CreateElement("ul")
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ul>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ul>", extractContent(ce))
}

func Test_Extractor_Content_PreserveNestedUnorderedList(t *testing.T) {
	innerList := dom.CreateElement("ul")
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))

	outerListItem := dom.CreateElement("li")
	dom.AppendChild(outerListItem, innerList)

	outerList := dom.CreateElement("ul")
	dom.AppendChild(outerList, outerListItem)
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ul>"+
		"<li>"+"<ul>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ul>"+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ul>", extractContent(ce))
}

func Test_Extractor_Content_PreserveNestedUnorderedListWithOtherElementsInside(t *testing.T) {
	innerList := dom.CreateElement("ul")
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateListItem(contentText))
	dom.AppendChild(innerList, testutil.CreateParagraph(""))

	outerListItem := dom.CreateElement("li")
	dom.AppendChild(outerListItem, dom.CreateTextNode(contentText))
	dom.AppendChild(outerListItem, testutil.CreateParagraph(contentText))
	dom.AppendChild(outerListItem, innerList)

	outerList := dom.CreateElement("ul")
	dom.AppendChild(outerList, outerListItem)
	dom.AppendChild(outerList, testutil.CreateListItem(contentText))
	dom.AppendChild(outerList, testutil.CreateParagraph(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, outerList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ul>"+
		"<li>"+contentText+
		"<p>"+contentText+"</p>"+
		"<ul>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ul>"+
		"</li>"+
		"<li>"+contentText+"</li>"+
		"<p>"+contentText+"</p>"+
		"</ul>", extractContent(ce))
}

func Test_Extractor_Content_PreserveUnorderedListWithNestedOrderedList(t *testing.T) {
	orderedList := dom.CreateElement("ol")
	dom.AppendChild(orderedList, testutil.CreateListItem(contentText))
	dom.AppendChild(orderedList, testutil.CreateListItem(contentText))

	li := dom.CreateElement("li")
	dom.AppendChild(li, orderedList)

	unorderedList := dom.CreateElement("ul")
	dom.AppendChild(unorderedList, li)
	dom.AppendChild(unorderedList, testutil.CreateListItem(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, unorderedList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ul><li>"+
		"<ol>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol></li>"+
		"<li>"+contentText+"</li>"+
		"</ul>", extractContent(ce))
}

func Test_Extractor_Content_MalformedListStructureWithExtraLiTagEnd(t *testing.T) {
	unorderedList := dom.CreateElement("ul")
	dom.SetInnerHTML(unorderedList, "<li>"+contentText+"</li></li><li>"+contentText+"</li>")

	doc, body := createHTML()
	dom.AppendChild(body, unorderedList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ul>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ul>", extractContent(ce))
}

func Test_Extractor_Content_MalformedListStructureWithExtraLiTagStart(t *testing.T) {
	orderedList := dom.CreateElement("ol")
	dom.SetInnerHTML(orderedList, "<li><li>"+contentText+"</li><li>"+contentText+"</li>")

	doc, body := createHTML()
	dom.AppendChild(body, orderedList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol>", extractContent(ce))
}

func Test_Extractor_Content_MalformedListStructureWithExtraOlTagStart(t *testing.T) {
	orderedList := dom.CreateElement("ol")
	dom.SetInnerHTML(orderedList, "<ol><li>"+contentText+"</li><li>"+contentText+"</li>")

	doc, body := createHTML()
	dom.AppendChild(body, orderedList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol><ol>"+
		"<li>"+contentText+"</li>"+
		"<li>"+contentText+"</li>"+
		"</ol></ol>", extractContent(ce))
}

func Test_Extractor_Content_MalformedListStructureWithoutLiTag(t *testing.T) {
	orderedList := dom.CreateElement("ol")
	dom.SetInnerHTML(orderedList, ""+
		"<li>"+contentText+"</li>"+
		contentText+
		"<li>"+contentText+"</li>")

	doc, body := createHTML()
	dom.AppendChild(body, orderedList)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<ol>"+
		"<li>"+contentText+"</li>"+
		contentText+
		"<li>"+contentText+"</li>"+
		"</ol>", extractContent(ce))
}

func Test_Extractor_Content_PreserveChildElementWithinBlockquote(t *testing.T) {
	blockquote := dom.CreateElement("blockquote")
	dom.AppendChild(blockquote, testutil.CreateParagraph(contentText+
		contentText+contentText+contentText))

	doc, body := createHTML()
	dom.AppendChild(body, blockquote)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<blockquote>"+
		"<p>"+contentText+contentText+contentText+contentText+"</p>"+
		"</blockquote>", extractContent(ce))
}

func Test_Extractor_Content_PreserveChildrenElementsWithinBlockquote(t *testing.T) {
	blockquote := dom.CreateElement("blockquote")
	dom.AppendChild(blockquote, testutil.CreateParagraph(contentText))
	dom.AppendChild(blockquote, testutil.CreateParagraph(contentText))
	dom.AppendChild(blockquote, testutil.CreateParagraph(contentText))

	doc, body := createHTML()
	dom.AppendChild(body, blockquote)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, "<blockquote>"+
		"<p>"+contentText+"</p>"+
		"<p>"+contentText+"</p>"+
		"<p>"+contentText+"</p>"+
		"</blockquote>", extractContent(ce))
}

func Test_Extractor_Content_DiscardBlockquoteWithoutContent(t *testing.T) {
	assertContentExtractor(t, "", "<blockquote></blockquote>")
}

func Test_Extractor_Content_PreservePre(t *testing.T) {
	article := contentText + contentText + contentText
	rawHTML := "<h1>" + contentText + "</h1><pre><kbd>" + article + "</kbd></pre>"
	assertContentExtractor(t, rawHTML, rawHTML)
}

func Test_Extractor_Content_DropCap(t *testing.T) {
	html := "<h1>" + contentText + "</h1>" +
		`<p><strong><span style="float: left">T</span>est</strong>` + contentText + "</p>"

	expected := "<h1>" + contentText + "</h1>" +
		"<p><strong><span>T</span>est</strong>" + contentText + "</p>"

	doc, body := createHTML()
	dom.SetInnerHTML(body, html)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func Test_Extractor_Content_BlockyArticle(t *testing.T) {
	rawHTML := "<h1>" + contentText + "</h1>" +
		"<span>" + contentText + "</span>" +
		"<div><span>" + contentText + "</span></div>" +
		"<p><em>" + contentText + "</em></p>" +
		"<div><cite><span><span>" + contentText + "</span></span></cite></div>" +
		"<div><span>" + contentText + "</span><span>" + contentText + "</span></div>" +
		"<main><span><blockquote><cite>" +
		"<span><span>" + contentText + "</span></span><span>" + contentText + "</span>" +
		"</cite></blockquote></span></main>"

	expected := "<h1>" + contentText + "</h1>" +
		"<span>" + contentText + "</span>" +
		"<div><span>" + contentText + "</span></div>" +
		"<p><em>" + contentText + "</em></p>" +
		"<div><cite><span><span>" + contentText + "</span></span></cite></div>" +
		"<div><span>" + contentText + "</span><span>" + contentText + "</span></div>" +
		"<blockquote><cite>" +
		"<span><span>" + contentText + "</span></span><span>" + contentText + "</span>" +
		"</cite></blockquote>"

	assertContentExtractor(t, expected, rawHTML)
}

func Test_Extractor_Content_SandboxedIframe(t *testing.T) {
	assertContentExtractor(t, "", "<iframe sandbox></iframe>")
}

func Test_Extractor_Content_SpanArticle(t *testing.T) {
	rawHTML := "" +
		"<span>" + contentText + "</span>" +
		"<span>" + contentText + "</span>" +
		"<span>" + contentText + "</span>"

	expected := "<div>" + rawHTML + "</div>"
	assertContentExtractor(t, expected, rawHTML)
}

func Test_Extractor_Content_UnwantedIframe(t *testing.T) {
	rawHTML := "" +
		"<p>" + contentText + "</p>" +
		"<iframe>dummy</iframe>" +
		"<p>" + contentText + "</p>"

	expected := "" +
		"<p>" + contentText + "</p>" +
		"<p>" + contentText + "</p>"

	assertContentExtractor(t, expected, rawHTML)
}

func Test_Extractor_Content_StripUnwantedClassNames(t *testing.T) {
	rawHTML := "" +
		`<p class="test">` + contentText + `</p>` +
		`<p class="iscaption">` + contentText + `</p>`

	expected := "" +
		`<p>` + contentText + `</p>` +
		`<p class="caption">` + contentText + `</p>`

	assertContentExtractor(t, expected, rawHTML)
}

func createHTML() (doc, body *html.Node) {
	doc = testutil.CreateHTML()
	body = dom.QuerySelector(doc, "body")
	head := dom.QuerySelector(doc, "head")
	dom.AppendChild(head, testutil.CreateTitle(titleText))
	return
}

func createMeta(doc *html.Node, property, content string) {
	head := dom.QuerySelector(doc, "head")
	if head != nil {
		meta := testutil.CreateMetaProperty(property, content)
		dom.AppendChild(head, meta)
	}
}

func assertContentExtractor(t *testing.T, expected, rawHTML string) {
	doc, body := createHTML()
	dom.SetInnerHTML(body, rawHTML)

	ce := extractor.NewContentExtractor(doc, nil, nil)
	assert.Equal(t, expected, extractContent(ce))
}

func extractContent(ce *extractor.ContentExtractor) string {
	extractedContent, _ := ce.ExtractContent(false)
	return extractedContent
}
