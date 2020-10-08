// ORIGINAL: javatest/DocumentTitleGetterTest.java

package internal

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"golang.org/x/net/html"
)

// Since implementation for getDocumentTitle between this package and the
// original dom-distilled code, the unit tests over here are different as
// well. However most of the test has same scenario as the original.

func Test_DocTitle_NoRoot(t *testing.T) {
	title := getDocumentTitle(nil, nil)
	assert.Equal(t, "", title)
}

func Test_DocTitle_TitlelessRoot(t *testing.T) {
	root := testutil.CreateDiv(0)
	title := getDocumentTitle(root, nil)
	assert.Equal(t, "", title)
}

func Test_DocTitle_TitledRoot(t *testing.T) {
	originalTitle := "testing non-string document.title with a titled root"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, originalTitle, title)
}

func Test_DocTitle_MultiTitledRoot(t *testing.T) {
	titleString1 := "first testing non-string document.title with a titled root"
	titleString2 := "second testing non-string document.title with a titled root"
	wc := stringutil.SelectWordCounter(titleString1)

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, testutil.CreateTitle(titleString1))
	dom.AppendChild(root, testutil.CreateTitle(titleString2))

	title := getDocumentTitle(root, wc)
	assert.Equal(t, titleString1, title)
}

func Test_DocTitle_1Dash2ShortParts(t *testing.T) {
	originalTitle := "before dash - after dash"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "before dash - after dash", title)
}

func Test_DocTitle_1Dash2LongParts(t *testing.T) {
	originalTitle := "part with 6 words before dash - part with 6 words after dash"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "part with 6 words before dash", title)
}

func Test_DocTitle_1Dash2LongPartsChinese(t *testing.T) {
	originalTitle := "比較長一點的句子 - 這是不要的部分"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "比較長一點的句子", title)
}

func Test_DocTitle_1DashLongAndShortParts(t *testing.T) {
	originalTitle := "part with 6 words before dash - after dash"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "part with 6 words before dash", title)
}

func Test_DocTitle_1DashShortAndLongParts(t *testing.T) {
	originalTitle := "before dash - part with 6 words after dash"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "part with 6 words after dash", title)
}

func Test_DocTitle_1DashShortAndLongPartsChinese(t *testing.T) {
	originalTitle := "短語 - 比較長一點的句子"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "比較長一點的句子", title)
}

func Test_DocTitle_2DashesShortParts(t *testing.T) {
	originalTitle := "before dash - between dash0 and dash1 - after dash1"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "before dash - between dash0 and dash1", title)
}

func Test_DocTitle_2DashesShortAndLongParts(t *testing.T) {
	// TODO(kuan): if using RegExp.split, this fails with "ant test.prod".
	originalTitle := "before - - part with 6 words after dash"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "- part with 6 words after dash", title)
}

func Test_DocTitle_1Bar2ShortParts(t *testing.T) {
	originalTitle := "before bar | after bar"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "before bar | after bar", title)
}

func Test_DocTitle_2ColonsShortParts(t *testing.T) {
	originalTitle := "start : midway : end"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "start : midway : end", title)
}

func Test_DocTitle_2ColonsShortPartsChinese(t *testing.T) {
	originalTitle := "開始 : 中間 : 最後"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "開始 : 中間 : 最後", title)
}

func Test_DocTitle_2ColonsShortAndLongParts(t *testing.T) {
	originalTitle := "start : midway : part with 6 words at end"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "part with 6 words at end", title)
}

func Test_DocTitle_2ColonsShortAndLongPartsChinese(t *testing.T) {
	originalTitle := "開始 : 中間 : 最後比較長的部分"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "最後比較長的部分", title)
}

func Test_DocTitle_2ColonsShortAndLongAndShortParts(t *testing.T) {
	originalTitle := "start : part with 6 words at midway : end"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "part with 6 words at midway : end", title)
}

func Test_DocTitle_2ColonsShortAndLongAndShortPartsChinese(t *testing.T) {
	originalTitle := "開始 : 中間要的部分 : 最後"
	wc := stringutil.SelectWordCounter(originalTitle)

	root := createTitledRoot(originalTitle)
	title := getDocumentTitle(root, wc)
	assert.Equal(t, "中間要的部分 : 最後", title)
}

func Test_DocTitle_H1AsTitle(t *testing.T) {
	headingText := "long heading with 5 words"
	wc := stringutil.SelectWordCounter(headingText)

	root := testutil.CreateDiv(0)
	h1 := testutil.CreateHeading(1, headingText)
	dom.AppendChild(root, h1)

	title := getDocumentTitle(root, wc)
	assert.Equal(t, "long heading with 5 words", title)
}

func Test_DocTitle_MultiHeadingsWithLongText(t *testing.T) {
	heading1Text := "long heading1 with 5 words"
	heading2Text := "long heading2 with 5 words"
	wc := stringutil.SelectWordCounter(heading1Text)

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, testutil.CreateHeading(1, heading1Text))
	dom.AppendChild(root, testutil.CreateHeading(2, heading2Text))

	title := getDocumentTitle(root, wc)
	assert.Equal(t, "long heading1 with 5 words", title)
}

func Test_DocTitle_H1WithLongHTML(t *testing.T) {
	headingHTML := `<a href="http://longheading.com"><b>long heading</b></a> with <br>5 words`
	h1 := testutil.CreateHeading(1, headingHTML)
	wc := stringutil.SelectWordCounter(domutil.InnerText(h1))

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, h1)

	title := getDocumentTitle(root, wc)
	assert.Equal(t, "long heading with 5 words", title)
}

func Test_DocTitle_H1WithLongHTMLWithNbsp(t *testing.T) {
	headingHTML := `<a href="http://longheading.com"><b> &nbsp;long heading</b></a> with <br>5 words &nbsp; `
	h1 := testutil.CreateHeading(1, headingHTML)
	wc := stringutil.SelectWordCounter(domutil.InnerText(h1))

	root := testutil.CreateDiv(0)
	dom.AppendChild(root, h1)

	title := getDocumentTitle(root, wc)
	assert.Equal(t, "long heading with 5 words", title)
}

func createTitledRoot(title string) *html.Node {
	root := testutil.CreateDiv(0)
	dom.AppendChild(root, testutil.CreateTitle(title))
	return root
}
