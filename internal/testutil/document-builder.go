// ORIGINAL: javatest/webdocument/TestWebDocumentBuilder.java

package testutil

import (
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// WebDocumentBuilder is a simple builder for testing.
type WebDocumentBuilder struct {
	document    *webdoc.Document
	textBuilder *TextBuilder
}

func NewWebDocumentBuilder() *WebDocumentBuilder {
	return &WebDocumentBuilder{
		document:    webdoc.NewDocument(),
		textBuilder: NewTextBuilder(stringutil.FastWordCounter{}),
	}
}

func (db *WebDocumentBuilder) AddText(text string) *webdoc.Text {
	wt := db.textBuilder.CreateForText(text)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddNestedText(text string) *webdoc.Text {
	wt := db.textBuilder.CreateNestedText(text, 5)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddAnchorText(text string) *webdoc.Text {
	wt := db.textBuilder.CreateForAnchorText(text)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddTable(innerHTML string) *webdoc.Table {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<table>"+innerHTML+"</table>")

	table := dom.QuerySelector(div, "table")
	wt := &webdoc.Table{Element: table}
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddImage() *webdoc.Image {
	image := dom.CreateElement("img")
	wi := &webdoc.Image{
		Element:   image,
		Width:     100,
		Height:    100,
		SourceURL: "http://www.example.com/foo.jpg",
	}

	db.document.AddElements(wi)
	return wi
}

func (db *WebDocumentBuilder) AddLeadImage() *webdoc.Image {
	image := dom.CreateElement("img")
	dom.SetAttribute(image, "width", "600")
	dom.SetAttribute(image, "height", "400")
	wi := &webdoc.Image{
		Element:   image,
		Width:     100,
		Height:    100,
		SourceURL: "http://www.example.com/lead.bmp",
	}

	db.document.AddElements(wi)
	return wi
}

func (db *WebDocumentBuilder) AddTagStart(tagName string) *webdoc.Tag {
	wt := webdoc.NewTag(tagName, webdoc.TagStart)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) AddTagEnd(tagName string) *webdoc.Tag {
	wt := webdoc.NewTag(tagName, webdoc.TagEnd)
	db.document.AddElements(wt)
	return wt
}

func (db *WebDocumentBuilder) Build() *webdoc.Document {
	return db.document
}
