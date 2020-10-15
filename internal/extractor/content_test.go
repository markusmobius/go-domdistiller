// ORIGINAL: javatest/ContentExtractorTest.java

package extractor_test

import (
	"fmt"
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
	ce := extractor.NewContentExtractor(doc, nil)
	extractedContent := ce.ExtractContent(false)
	assert.True(t, strings.Contains(extractedContent, domutil.InnerText(contentDiv1)))
	assert.True(t, strings.Contains(extractedContent, domutil.InnerText(titleDiv)))

	// Now set the title and it should excluded from the content.
	head := dom.QuerySelector(doc, "head")
	dom.AppendChild(head, testutil.CreateTitle(titleText))

	ce = extractor.NewContentExtractor(doc, nil)
	extractedContent = ce.ExtractContent(false)
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

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, div)

	ce := extractor.NewContentExtractor(doc, nil)
	extractedContent := ce.ExtractContent(false)
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

	ce := extractor.NewContentExtractor(doc, nil)
	assert.Equal(t, markupParserTitle, ce.ExtractTitle())
}

func Test_Extractor_Content_Image(t *testing.T) {
	// Test the absolute and different kinds of relative URLs for image sources,
	// and also add an extra comma (,) as malformed srcset syntax for robustness.
	// Also test images in WebImage and WebTable.
	//
	// The result of this test is different with the original dom-distiller because
	// we can't compute stylesheet so our LeadImageFinder gives different result.
	// NEED-COMPUTE-CSS
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

	expected := `<img src="http://example.com/path/image" ` +
		`srcset="http://example.com/path/image200 200w, ` +
		`http://example.org/image400 400w"/>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, rawHTML)

	pageURL, _ := nurl.ParseRequestURI("http://example.com/path/")
	ce := extractor.NewContentExtractor(doc, pageURL)
	assert.Equal(t, expected, ce.ExtractContent(false))
}

func Test_Extractor_Content_RemoveFontColorAttributes(t *testing.T) {
	// The result of this test is different with the original dom-distiller because
	// we can't compute stylesheet so our LeadImageFinder gives different result.
	// NEED-COMPUTE-CSS
	outerFontTag := dom.CreateElement("font")
	dom.SetAttribute(outerFontTag, "color", "blue")

	text := `<font color="red">` + contentText + `</font>`
	dom.AppendChild(outerFontTag, testutil.CreateSpan(text))
	dom.AppendChild(outerFontTag, dom.CreateTextNode(" "))
	dom.AppendChild(outerFontTag, testutil.CreateSpan(text))
	dom.AppendChild(outerFontTag, dom.CreateTextNode("\n"))
	dom.AppendChild(outerFontTag, testutil.CreateSpan(text))
	dom.AppendChild(outerFontTag, dom.CreateTextNode(" "))

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, outerFontTag)

	expected := `<font>Lorem Ipsum Lorem Ipsum Lorem Ipsum.</font>` +
		`<font>Lorem Ipsum Lorem Ipsum Lorem Ipsum.</font>` +
		`<font>Lorem Ipsum Lorem Ipsum Lorem Ipsum.</font>`

	ce := extractor.NewContentExtractor(doc, nil)
	assert.Equal(t, expected, ce.ExtractContent(false))
}

func Test_Extractor_Content_RemoveStyleAttributes(t *testing.T) {
	rawHTML := `<h1 style="font-weight: folder">` +
		contentText +
		`</h1>` +
		`<p style="">` +
		contentText +
		`</p>` +
		`<img style="align: left" data-src="/test.png">` +
		`<table style="position: absolute">` +
		`<tbody style="font-size: 2">` +
		`<tr style="z-index: 0">` +
		`<th style="top: 0px">` + contentText +
		`<img style="align: left" src="/test.png">` +
		`</th>` +
		`<th style="width: 20px">` + contentText + `</th>` +
		`</tr><tr style="left: 0">` +
		`<td style="display: block">` + contentText + `</td>` +
		`<td style="color: #123">` + contentText + `</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, rawHTML)

	pageURL, _ := nurl.ParseRequestURI("http://example.com/path/")
	ce := extractor.NewContentExtractor(doc, pageURL)
	fmt.Println(ce.ExtractContent(false))
	// assert.Equal(t, expected, ce.ExtractContent(false))
}

func createMeta(doc *html.Node, property, content string) {
	head := dom.QuerySelector(doc, "head")
	if head != nil {
		meta := testutil.CreateMetaProperty(property, content)
		dom.AppendChild(head, meta)
	}
}
