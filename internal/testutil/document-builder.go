// ORIGINAL: javatest/webdocument/FakeWebDocumentBuilder.java

package testutil

import (
	"bytes"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// WebDocumentBuilder is a simple builder that just creates an html-like string
// from the calls. Only used for test.
type WebDocumentBuilder struct {
	buffer bytes.Buffer
	nodes  []*html.Node
}

func NewWebDocumentBuilder() *WebDocumentBuilder {
	return &WebDocumentBuilder{}
}

func (db *WebDocumentBuilder) Build() string {
	return db.buffer.String()
}

func (db *WebDocumentBuilder) SkipNode(e *html.Node) {}

func (db *WebDocumentBuilder) StartNode(e *html.Node) {
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

func (db *WebDocumentBuilder) EndNode() {
	node := db.nodes[len(db.nodes)-1]
	db.nodes = db.nodes[:len(db.nodes)-1]
	db.buffer.WriteString("</" + dom.TagName(node) + ">")
}

func (db *WebDocumentBuilder) AddTextNode(textNode *html.Node) {
	db.buffer.WriteString(textNode.Data)
}

func (db *WebDocumentBuilder) AddLineBreak(node *html.Node) {
	db.buffer.WriteString("\n")
}

func (db *WebDocumentBuilder) AddDataTable(e *html.Node) {
	db.buffer.WriteString("<datatable/>")
}

func (db *WebDocumentBuilder) AddTag(tag *webdoc.Tag) {}

func (db *WebDocumentBuilder) AddEmbed(embed webdoc.Element) {}
