// ORIGINAL: javatest/webdocument/TestWebTextBuilder.java

package testutil

import (
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

type TextBuilder struct {
	wordCounter stringutil.WordCounter
	textNodes   []*html.Node
}

func NewTextBuilder(wc stringutil.WordCounter) *TextBuilder {
	return &TextBuilder{wordCounter: wc}
}

func (tb *TextBuilder) CreateForText(str string) *webdoc.Text {
	return tb.create(str, false)
}

func (tb *TextBuilder) CreateForAnchorText(str string) *webdoc.Text {
	return tb.create(str, true)
}

func (tb *TextBuilder) CreateNestedText(str string, levels int) *webdoc.Text {
	div := dom.CreateElement("div")
	tmp := div

	for i := 0; i < levels-1; i++ {
		dom.AppendChild(tmp, dom.CreateElement("div"))
		tmp = dom.FirstElementChild(tmp)
	}

	dom.AppendChild(tmp, dom.CreateTextNode(str))
	tb.textNodes = append(tb.textNodes, tmp.FirstChild)

	idx := len(tb.textNodes) - 1
	numWords := tb.wordCounter.Count(str)

	return &webdoc.Text{
		Text:           str,
		TextNodes:      tb.textNodes,
		Start:          idx,
		End:            idx + 1,
		FirstWordNode:  idx,
		LastWordNode:   idx,
		NumWords:       numWords,
		NumLinkedWords: 0,
		TagLevel:       0,
		OffsetBlock:    idx,
	}
}

func (tb *TextBuilder) create(str string, isAnchor bool) *webdoc.Text {
	tb.textNodes = append(tb.textNodes, dom.CreateTextNode(str))

	idx := len(tb.textNodes) - 1
	numWords := tb.wordCounter.Count(str)
	numLinkedWords := numWords
	if !isAnchor {
		numLinkedWords = 0
	}

	return &webdoc.Text{
		Text:           str,
		TextNodes:      tb.textNodes,
		Start:          idx,
		End:            idx + 1,
		FirstWordNode:  idx,
		LastWordNode:   idx,
		NumWords:       numWords,
		NumLinkedWords: numLinkedWords,
		TagLevel:       0,
		OffsetBlock:    idx,
	}
}
