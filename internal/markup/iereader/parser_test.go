// ORIGINAL: javatest/IEReadingViewParserTest.java

package iereader_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/markup/iereader"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_IEReader_Title(t *testing.T) {
	root := testutil.CreateHTML()
	head := dom.QuerySelector(root, "head")
	body := dom.QuerySelector(root, "body")

	expectedTitle := "Testing title"
	createMeta(root, "title", expectedTitle)
	dom.AppendChild(head, testutil.CreateTitle(expectedTitle))
	dom.AppendChild(body, testutil.CreateHeading(1, "start_h1: "+expectedTitle+" end_h1"))

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedTitle, parser.Title())
}

func Test_IEReader_DateInOneElemWithOneClass(t *testing.T) {
	expectedDate := "Monday January 1st 2011 01:01"
	div := testutil.CreateDiv(0)
	dom.SetAttribute(div, "class", "dateline")
	dom.SetInnerHTML(div, expectedDate)

	root := testutil.CreateHTML()
	dom.AppendChild(root, div)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedDate, parser.Article().PublishedTime)
}

func Test_IEReader_DateInSubstringClassName(t *testing.T) {
	div := testutil.CreateDiv(0)
	dom.SetAttribute(div, "class", "b4datelineaft")
	dom.SetInnerHTML(div, "Monday January 1st 2011 01:01")

	root := testutil.CreateHTML()
	dom.AppendChild(root, div)

	parser := iereader.NewParser(root)
	assert.Equal(t, "", parser.Article().PublishedTime)
}

func Test_IEReader_DateInOneElemWithMultiClasses(t *testing.T) {
	expectedDate := "Tuesday February 2nd 2012 02:02"
	div := testutil.CreateDiv(0)
	dom.SetAttribute(div, "class", "blah1 dateline blah2")
	dom.SetInnerHTML(div, expectedDate)

	root := testutil.CreateHTML()
	dom.AppendChild(root, div)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedDate, parser.Article().PublishedTime)
}

func Test_IEReader_DateInOneBranch(t *testing.T) {
	expectedDate := "Wednesday Mar 3rd 2013 03:03"
	div1 := testutil.CreateDiv(1)
	div2 := testutil.CreateDiv(2)
	div3 := testutil.CreateDiv(3)
	dom.SetAttribute(div1, "class", "blah11 blah12")
	dom.SetInnerHTML(div1, "blah11 blah12")
	dom.SetAttribute(div2, "class", "blah21")
	dom.SetInnerHTML(div2, "blah21 only")
	dom.SetAttribute(div3, "class", "blah31 dateline")
	dom.SetInnerHTML(div3, expectedDate)

	root := testutil.CreateHTML()
	dom.AppendChild(div2, div3)
	dom.AppendChild(div1, div2)
	dom.AppendChild(root, div1)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedDate, parser.Article().PublishedTime)
}

func Test_IEReader_DateInMultiBranches(t *testing.T) {
	expectedDate := "Thursday Apr 4th 2014 04:04"
	div1 := testutil.CreateDiv(1)
	div2 := testutil.CreateDiv(2)
	div3 := testutil.CreateDiv(3)
	dom.SetAttribute(div1, "class", "blah11 blah12")
	dom.SetInnerHTML(div1, "blah11 blah12")
	dom.SetAttribute(div2, "class", "blah12 dateline")
	dom.SetInnerHTML(div2, expectedDate)
	dom.SetAttribute(div3, "class", "blah31")
	dom.SetInnerHTML(div3, "blah31 only")

	root := testutil.CreateHTML()
	head := dom.QuerySelector(root, "head")
	body := dom.QuerySelector(root, "body")
	dom.AppendChild(head, div1)
	dom.AppendChild(body, div2)
	dom.AppendChild(root, div3)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedDate, parser.Article().PublishedTime)
}

func Test_IEReader_DateInMeta(t *testing.T) {
	root := testutil.CreateHTML()
	expectedDate := "Friday Apr 5th 2015 05:05"
	createMeta(root, "displaydate", expectedDate)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedDate, parser.Article().PublishedTime)
}

func Test_IEReader_AuthorInSiblings(t *testing.T) {
	expectedAuthor := "Jane Doe"
	div1 := testutil.CreateDiv(1)
	div2 := testutil.CreateDiv(2)
	dom.SetAttribute(div1, "class", "blah1")
	dom.SetInnerHTML(div1, "blah1 only")
	dom.SetAttribute(div2, "class", "byline-name")
	dom.SetInnerHTML(div2, expectedAuthor)

	root := testutil.CreateHTML()
	head := dom.QuerySelector(root, "head")
	dom.AppendChild(head, testutil.CreateTitle("testing author"))
	dom.AppendChild(head, div2)
	dom.AppendChild(head, div1)

	parser := iereader.NewParser(root)
	authors := parser.Article().Authors
	assert.Equal(t, 1, len(authors))
	assert.Equal(t, expectedAuthor, authors[0])
}

func Test_IEReader_AuthorInBranch(t *testing.T) {
	expectedAuthor := "Jane Doe"
	div2 := testutil.CreateDiv(2)
	dom.SetAttribute(div2, "class", "byline-name")
	dom.SetInnerHTML(div2, expectedAuthor)

	div1 := testutil.CreateDiv(1)
	dom.SetAttribute(div1, "class", "blah1")
	dom.SetInnerHTML(div1, "blah1 only")
	dom.AppendChild(div1, div2)

	root := testutil.CreateHTML()
	head := dom.QuerySelector(root, "head")
	dom.AppendChild(head, testutil.CreateTitle("testing author"))
	dom.AppendChild(head, div1)

	parser := iereader.NewParser(root)
	authors := parser.Article().Authors
	assert.Equal(t, 1, len(authors))
	assert.Equal(t, expectedAuthor, authors[0])
}

func Test_IEReader_PublisherFromPublisherAttr(t *testing.T) {
	expectedPublisher := "Publisher Attribute"
	div := testutil.CreateDiv(0)
	dom.SetAttribute(div, "publisher", expectedPublisher)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.AppendChild(body, div)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedPublisher, parser.Publisher())
}

func Test_IEReader_PublisherFromSourceOrganizationrAttr(t *testing.T) {
	expectedPublisher := "Source Organization Attribute"
	div := testutil.CreateDiv(0)
	dom.SetAttribute(div, "source_organization", expectedPublisher)

	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.AppendChild(body, div)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedPublisher, parser.Publisher())
}

func Test_IEReader_CopyrightInMeta(t *testing.T) {
	expectedCopyright := "Friday Apr 5th 2015 05:05"

	root := testutil.CreateHTML()
	createMeta(root, "copyright", expectedCopyright)

	parser := iereader.NewParser(root)
	assert.Equal(t, expectedCopyright, parser.Copyright())
}

func Test_IEReader_UncaptionedDominantImage(t *testing.T) {
	expectedUrl := "http://example.com/dominant_without_caption.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "600")
	dom.SetAttribute(img, "height", "400")

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, img)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 1, len(images))

	image := images[0]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, "", image.Caption)
	assert.Equal(t, 600, image.Width)
	assert.Equal(t, 400, image.Height)
}

func Test_IEReader_UncaptionedDominantImageWithInvalidSize(t *testing.T) {
	expectedUrl := "http://example.com/dominant_without_caption.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, img)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 0, len(images))
}

func Test_IEReader_CaptionedDominantImageWithSmallestAR(t *testing.T) {
	expectedUrl := "http://example.com/captioned_smallest_dominant.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "400")
	dom.SetAttribute(img, "height", "307")

	expectedCaption := "Captioned Dominant Image with Smallest AR"
	figure := createFigureWithCaption(expectedCaption)
	dom.AppendChild(figure, img)

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, figure)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 1, len(images))

	image := images[0]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, expectedCaption, image.Caption)
	assert.Equal(t, 400, image.Width)
	assert.Equal(t, 307, image.Height)
}

func Test_IEReader_CaptionedDominantImageWithBiggestAR(t *testing.T) {
	expectedUrl := "http://example.com/captioned_biggest_dominant.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "400")
	dom.SetAttribute(img, "height", "134")

	expectedCaption := "Captioned Dominant Image with Biggest AR"
	figure := createFigureWithCaption(expectedCaption)
	dom.AppendChild(figure, img)

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, figure)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 1, len(images))

	image := images[0]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, expectedCaption, image.Caption)
	assert.Equal(t, 400, image.Width)
	assert.Equal(t, 134, image.Height)
}

func Test_IEReader_CaptionedDominantImageWithInvalidSize(t *testing.T) {
	expectedUrl := "http://example.com/captioned_dominant_with_wrong_dimensions.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")

	expectedCaption := "Captioned Dominant Image with Invalid Size"
	figure := createFigureWithCaption(expectedCaption)
	dom.AppendChild(figure, img)

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, figure)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 1, len(images))

	image := images[0]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, expectedCaption, image.Caption)
	assert.Equal(t, 100, image.Width)
	assert.Equal(t, 100, image.Height)
}

func Test_IEReader_UncaptionedInlineImageWithSmallestAR(t *testing.T) {
	expectedUrl := "http://example.com/inline_without_caption.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "400")
	dom.SetAttribute(img, "height", "307")

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, createDefaultDominantFigure())
	dom.AppendChild(root, img)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 2, len(images))

	image := images[1]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, "", image.Caption)
	assert.Equal(t, 400, image.Width)
	assert.Equal(t, 307, image.Height)
}

func Test_IEReader_UncaptionedInlineImageWithBiggestAR(t *testing.T) {
	expectedUrl := "http://example.com/inline_without_caption.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "400")
	dom.SetAttribute(img, "height", "134")

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, createDefaultDominantFigure())
	dom.AppendChild(root, img)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 2, len(images))

	image := images[1]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, "", image.Caption)
	assert.Equal(t, 400, image.Width)
	assert.Equal(t, 134, image.Height)
}

func Test_IEReader_CaptionedInlineImageWithInvalidSize(t *testing.T) {
	expectedUrl := "http://example.com/captioned_smallest_inline.jpeg"
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", expectedUrl)
	dom.SetAttribute(img, "width", "400")
	dom.SetAttribute(img, "height", "400")

	expectedCaption := "Captioned Inline Image with Smallest AR"
	figure := createFigureWithCaption(expectedCaption)
	dom.AppendChild(figure, img)

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, createDefaultDominantFigure())
	dom.AppendChild(root, figure)

	parser := iereader.NewParser(root)
	images := parser.Images()
	assert.Equal(t, 2, len(images))

	image := images[1]
	assert.Equal(t, expectedUrl, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, "", image.Type)
	assert.Equal(t, expectedCaption, image.Caption)
	assert.Equal(t, 400, image.Width)
	assert.Equal(t, 400, image.Height)
}

func Test_IEReader_OptOut(t *testing.T) {
	root := testutil.CreateHTML()
	createMeta(root, "IE_RM_OFF", "true")

	parser := iereader.NewParser(root)
	assert.True(t, parser.OptOut())
}

func Test_IEReader_OptIn(t *testing.T) {
	root := testutil.CreateHTML()
	createMeta(root, "IE_RM_OFF", "false")

	parser := iereader.NewParser(root)
	assert.False(t, parser.OptOut())
}

func createMeta(root *html.Node, name, content string) {
	head := dom.QuerySelector(root, "head")
	if head != nil {
		meta := testutil.CreateMetaName(name, content)
		dom.AppendChild(head, meta)
	}
}

func createFigureWithCaption(caption string) *html.Node {
	figCaption := dom.CreateElement("figcaption")
	dom.SetInnerHTML(figCaption, caption)

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, figCaption)

	return figure
}

func createDefaultDominantFigure() *html.Node {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", "http://example.com/dominant.jpeg")
	dom.SetAttribute(img, "width", "600")
	dom.SetAttribute(img, "height", "400")

	figure := createFigureWithCaption("Default Dominant Image")
	dom.AppendChild(figure, img)

	return figure
}
