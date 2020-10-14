// ORIGINAL: javatest/webdocument/FakeWebDocumentBuilder.java

package testutil

import (
	"bytes"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// FakeWebDocumentBuilder is a simple builder that just creates an html-like string
// from the calls. Only used for dom-converter test.
type FakeWebDocumentBuilder struct {
	buffer bytes.Buffer
	nodes  []*html.Node
}

func NewFakeWebDocumentBuilder() *FakeWebDocumentBuilder {
	return &FakeWebDocumentBuilder{}
}

func (db *FakeWebDocumentBuilder) Build() string {
	return db.buffer.String()
}

func (db *FakeWebDocumentBuilder) SkipNode(e *html.Node) {}

func (db *FakeWebDocumentBuilder) StartNode(e *html.Node) {
	db.nodes = append(db.nodes, e)
	db.buffer.WriteString("<")
	db.buffer.WriteString(dom.TagName(e))
	for _, attr := range e.Attr {
		db.buffer.WriteString(" ")
		db.buffer.WriteString(attr.Key)
		db.buffer.WriteString(`="`)
		db.buffer.WriteString(attr.Val)
		db.buffer.WriteString(`"`)
	}
	db.buffer.WriteString(">")
}

func (db *FakeWebDocumentBuilder) EndNode() {
	node := db.nodes[len(db.nodes)-1]
	db.nodes = db.nodes[:len(db.nodes)-1]
	db.buffer.WriteString("</" + dom.TagName(node) + ">")
}

func (db *FakeWebDocumentBuilder) AddTextNode(textNode *html.Node) {
	db.buffer.WriteString(textNode.Data)
}

func (db *FakeWebDocumentBuilder) AddLineBreak(node *html.Node) {
	db.buffer.WriteString("\n")
}

func (db *FakeWebDocumentBuilder) AddDataTable(e *html.Node) {
	db.buffer.WriteString("<datatable/>")
}

func (db *FakeWebDocumentBuilder) AddTag(tag *webdoc.Tag) {}

func (db *FakeWebDocumentBuilder) AddEmbed(embed webdoc.Element) {}
