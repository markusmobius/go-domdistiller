// ORIGINAL: java/webdocument/WebDocumentBuilder.java and
//           java/webdocument/WebDocumentBuilderInterface.java

package webdoc

import (
	nurl "net/url"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type DocumentBuilder interface {
	SkipNode(e *html.Node)
	StartNode(e *html.Node)
	EndNode()
	AddTextNode(textNode *html.Node)
	AddLineBreak(node *html.Node)
	AddDataTable(e *html.Node)
	AddTag(tag *Tag)
	AddEmbed(embed Element)
}

type WebDocumentBuilder struct {
	tagLevel      int
	nextTextIndex int
	groupNumber   int
	flush         bool
	document      *Document
	textBuilder   *TextBuilder
	actionStack   []ElementAction
	pageURL       *nurl.URL
}

func NewWebDocumentBuilder(wc stringutil.WordCounter, pageURL *nurl.URL) *WebDocumentBuilder {
	return &WebDocumentBuilder{
		document:    &Document{},
		textBuilder: NewTextBuilder(wc),
		pageURL:     pageURL,
	}
}

func (db *WebDocumentBuilder) SkipNode(e *html.Node) {
	db.flush = true
}

func (db *WebDocumentBuilder) StartNode(e *html.Node) {
	action := GetActionForElement(e)
	db.actionStack = append(db.actionStack, action)

	if action.ChangesTagLevel {
		db.tagLevel++
	}

	if action.IsAnchor {
		db.textBuilder.EnterAnchor()
	}

	db.flush = db.flush || action.Flush
}

func (db *WebDocumentBuilder) EndNode() {
	nActions := len(db.actionStack)
	if nActions == 0 {
		return
	}

	lastAction := db.actionStack[nActions-1]

	if lastAction.ChangesTagLevel {
		db.tagLevel--
	}

	if db.flush || lastAction.Flush {
		db.flushBlock(db.groupNumber)
		db.groupNumber++
	}

	if lastAction.IsAnchor {
		db.textBuilder.ExitAnchor()
	}

	db.actionStack = db.actionStack[:nActions-1]
}

func (db *WebDocumentBuilder) AddTextNode(textNode *html.Node) {
	if db.flush {
		db.flushBlock(db.groupNumber)
		db.groupNumber++
		db.flush = false
	}

	db.textBuilder.AddTextNode(textNode, db.tagLevel)
}

func (db *WebDocumentBuilder) AddLineBreak(br *html.Node) {
	if db.flush {
		db.flushBlock(db.groupNumber)
		db.groupNumber++
		db.flush = false
	}

	db.textBuilder.AddLineBreak(br)
}

func (db *WebDocumentBuilder) AddDataTable(table *html.Node) {
	db.flushBlock(db.groupNumber)
	db.document.AddElements(&Table{
		Element: table,
		PageURL: db.pageURL,
	})
}

func (db *WebDocumentBuilder) AddTag(tag *Tag) {
	db.flushBlock(db.groupNumber)
	db.document.AddElements(tag)
}

func (db *WebDocumentBuilder) AddEmbed(embed Element) {
	db.flushBlock(db.groupNumber)
	db.document.AddElements(embed)
}

func (db *WebDocumentBuilder) Build() *Document {
	db.flushBlock(db.groupNumber)
	return db.document
}

func (db *WebDocumentBuilder) flushBlock(group int) {
	if text := db.textBuilder.Build(db.nextTextIndex); text != nil {
		text.GroupNumber = group
		db.nextTextIndex++
		db.addText(*text)
	}
}

func (db *WebDocumentBuilder) addText(text Text) {
	for _, action := range db.actionStack {
		for _, label := range action.Labels {
			text.AddLabel(label)
		}
	}

	db.document.AddElements(&text)
}
