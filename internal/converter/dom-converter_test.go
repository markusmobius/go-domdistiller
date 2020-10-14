// ORIGINAL: javatest/webdocument/DomConverterTest.java

package converter_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/converter"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_Converter_DisplayNone
// - Test_Converter_VisibilityHidden
// - Test_Converter_InvisibleInVisible
// - Test_Converter_VisibleInInvisible
// - Test_Converter_VisibleInInvisible2
// - Test_Converter_InvisibleInInvisible
// - Test_Converter_DifferentChildrenInVisible
// - Test_Converter_DifferentChildrenInInvisible
// - Test_Converter_KeepHidden
// - Test_Converter_KeepHiddenNested
// - Test_Converter_KeepContinue
// - Test_Converter_KeepContinueNested
// - Test_Converter_WikipediaFoldedSections

func Test_Converter_VisibleText(t *testing.T) {
	html := "visible text"
	runTest(t, html, html)
}

func Test_Converter_VisibleElement(t *testing.T) {
	html := "<div>visible element</div>"
	runTest(t, html, html)
}

func Test_Converter_VisibleInVisible(t *testing.T) {
	html := "<div>visible parent" +
		"<div>visible child</div>" +
		"</div>"
	runTest(t, html, html)
}

func Test_Converter_DataTable(t *testing.T) {
	html := `<table align="left" role="grid">` + // role=grid make this a data table.
		`<tbody align="left">` +
		`<tr>` +
		`<td>row1col1</td>` +
		`<td>row1col2</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	runTest(t, html, "<datatable/>")
}

func Test_Converter_NonDataTable(t *testing.T) {
	html := `<table align="left">` +
		`<tbody align="left">` +
		`<tr>` +
		`<td>row1col1</td>` +
		`<td>row1col2</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	runTest(t, html, html)
}

func Test_Converter_IgnorableElements(t *testing.T) {
	runTest(t, "<head></head>", "")
	runTest(t, "<style></style>", "")
	runTest(t, "<script></script>", "")
	runTest(t, "<link></link>", "")
	runTest(t, "<noscript></noscript>", "")
	runTest(t, "<iframe></iframe>", "")
	runTest(t, "<svg></svg>", "")
	runTest(t, "<option></option>", "")
	runTest(t, "<object></object>", "")
	runTest(t, "<embed></embed>", "")
	runTest(t, "<applet></applet>", "")
}

func Test_Converter_SvgTagNameCase(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<SVG></SVG>")
	assert.Equal(t, "svg", dom.TagName(dom.FirstElementChild(div)))
}

func Test_Converter_ElementOrder(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, `Text content <img src="http://example.com/1.jpg"> more content`)

	builder := webdoc.NewWebDocumentBuilder(stringutil.FastWordCounter{}, nil)
	converter.NewDomConverter(builder, nil).Convert(div)

	doc := builder.Build()
	elements := doc.Elements

	assert.Equal(t, 3, len(elements))
	assert.IsType(t, &webdoc.Text{}, elements[0])
	assert.IsType(t, &webdoc.Image{}, elements[1])
	assert.IsType(t, &webdoc.Text{}, elements[2])
}

func Test_Converter_LineBreak(t *testing.T) {
	html := "text<br>split<br/>with<br/>lines"
	runTest(t, html, "text\nsplit\nwith\nlines")
}

func Test_Converter_List(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<ol><li>some text1</li><li>some text2</li></ol>")

	builder := webdoc.NewWebDocumentBuilder(stringutil.FastWordCounter{}, nil)
	converter.NewDomConverter(builder, nil).Convert(div)

	doc := builder.Build()
	elements := doc.Elements

	assert.Equal(t, 8, len(elements))
	assert.IsType(t, &webdoc.Tag{}, elements[0])
	assert.Equal(t, webdoc.TagStart, (elements[0].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Tag{}, elements[1])
	assert.Equal(t, webdoc.TagStart, (elements[1].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Text{}, elements[2])

	assert.IsType(t, &webdoc.Tag{}, elements[3])
	assert.Equal(t, webdoc.TagEnd, (elements[3].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Tag{}, elements[4])
	assert.Equal(t, webdoc.TagStart, (elements[4].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Text{}, elements[5])

	assert.IsType(t, &webdoc.Tag{}, elements[6])
	assert.Equal(t, webdoc.TagEnd, (elements[6].(*webdoc.Tag).Type))

	assert.IsType(t, &webdoc.Tag{}, elements[7])
	assert.Equal(t, webdoc.TagEnd, (elements[7].(*webdoc.Tag).Type))
}

func Test_Converter_SocialElements(t *testing.T) {
	runTest(t, `<div></div>`, `<div></div>`)
	runTest(t, `<div data-component="share"></div>`, ``)
	runTest(t, `<div class="socialArea"></div>`, ``)
	runTest(t, `<li></li>`, `<li></li>`)
	runTest(t, `<li class="sharing"></li>`, ``)
}

func Test_Converter_WikiEditLinks(t *testing.T) {
	html := `<a href="index.php?action=edit&redlink=1"></a>`
	runTest(t, html, html)

	html = `<a href="index.php?action=edit&section=3" class="mw-ui-icon"></a>`
	runTest(t, html, "")
}

func Test_Converter_WikiEditSection(t *testing.T) {
	html := `<span class="mw-headline"></span>`
	runTest(t, html, html)

	html = `<span class="mw-editsection"></span>`
	runTest(t, html, "")
}

func runTest(t *testing.T, innerHTML, expectedHTML string) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, innerHTML)

	builder := testutil.NewFakeWebDocumentBuilder()
	converter.NewDomConverter(builder, nil).Convert(div)

	expected := "<div>" + expectedHTML + "</div>"
	assert.Equal(t, expected, builder.Build())
}
