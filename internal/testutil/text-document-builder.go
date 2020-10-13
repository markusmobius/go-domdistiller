// ORIGINAL: javatest/TestTextDocumentBuilder.java

package testutil

import (
	"net/url"

	"github.com/markusmobius/go-domdistiller/internal/converter"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

type TextDocumentBuilder struct {
	textBlocks  []*webdoc.TextBlock
	textBuilder *TextBuilder
}

func NewTextDocumentBuilder(wc stringutil.WordCounter) *TextDocumentBuilder {
	return &TextDocumentBuilder{
		textBuilder: NewTextBuilder(wc),
	}
}

func (tdb *TextDocumentBuilder) AddContentBlock(str string, labels ...string) *webdoc.TextBlock {
	tb := tdb.addBlock(str, labels...)
	tb.SetIsContent(true)
	return tb
}

func (tdb *TextDocumentBuilder) AddNonContentBlock(str string, labels ...string) *webdoc.TextBlock {
	tb := tdb.addBlock(str, labels...)
	tb.SetIsContent(false)
	return tb
}

func (tdb *TextDocumentBuilder) Build() *webdoc.TextDocument {
	return webdoc.NewTextDocument(tdb.textBlocks)
}

func (tdb *TextDocumentBuilder) addBlock(str string, labels ...string) *webdoc.TextBlock {
	wt := tdb.textBuilder.CreateForText(str)
	for _, label := range labels {
		wt.AddLabel(label)
	}

	tdb.textBlocks = append(tdb.textBlocks, webdoc.NewTextBlock(wt))
	return tdb.textBlocks[len(tdb.textBlocks)-1]
}

func NewTextDocumentFromPage(doc *html.Node, wc stringutil.WordCounter, pageURL *url.URL) *webdoc.TextDocument {
	builder := webdoc.NewWebDocumentBuilder(wc, pageURL)
	converter.NewDomConverter(builder, pageURL).Convert(doc)
	return builder.Build().CreateTextDocument()
}
