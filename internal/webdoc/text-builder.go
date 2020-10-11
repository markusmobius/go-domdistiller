// ORIGINAL: java/webdocument/WebText.java

package webdoc

import (
	"bytes"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type TextBuilder struct {
	textBuffer             bytes.Buffer
	numWords               int
	numAnchorWords         int
	blockTagLevel          int
	inAnchor               bool
	textNodes              []*html.Node
	firstNode              int
	firstNonWhitespaceNode int
	lastNonWhitespaceNode  int
	wordCounter            stringutil.WordCounter
}

func NewTextBuilder(wc stringutil.WordCounter) *TextBuilder {
	return &TextBuilder{
		blockTagLevel:          -1,
		firstNonWhitespaceNode: -1,
		wordCounter:            wc,
	}
}

func (tb *TextBuilder) AddTextNode(textNode *html.Node, tagLevel int) {
	if textNode.Type != html.TextNode {
		return
	}

	text := textNode.Data
	if text == "" {
		return
	}

	tb.textBuffer.WriteString(text)
	tb.textNodes = append(tb.textNodes, textNode)

	if stringutil.IsStringAllWhitespace(text) {
		return
	}

	wordCount := tb.wordCounter.Count(text)
	tb.numWords += wordCount
	if tb.inAnchor {
		tb.numAnchorWords += wordCount
	}

	tb.lastNonWhitespaceNode = len(tb.textNodes) - 1
	if tb.firstNonWhitespaceNode < tb.firstNode {
		tb.firstNonWhitespaceNode = tb.lastNonWhitespaceNode
	}

	if tb.blockTagLevel == -1 {
		tb.blockTagLevel = tagLevel
	}
}

func (tb *TextBuilder) AddLineBreak(node *html.Node) {
	tb.textBuffer.WriteString("\n")
	tb.textNodes = append(tb.textNodes, node)
}

func (tb *TextBuilder) Reset() {
	tb.textBuffer.Reset()
	tb.numWords = 0
	tb.numAnchorWords = 0
	tb.firstNode = len(tb.textNodes)
	tb.blockTagLevel = -1
}

func (tb *TextBuilder) EnterAnchor() {
	tb.inAnchor = true
	tb.textBuffer.WriteString(" ")
}

func (tb *TextBuilder) ExitAnchor() {
	tb.inAnchor = false
	tb.textBuffer.WriteString(" ")
}

func (tb *TextBuilder) Build(offsetBlock int) *Text {
	if tb.firstNode == len(tb.textNodes) {
		return nil
	}

	if tb.firstNonWhitespaceNode < tb.firstNode {
		tb.Reset()
		return nil
	}

	text := Text{
		Text:           tb.textBuffer.String(),
		TextNodes:      tb.textNodes,
		Start:          tb.firstNode,
		End:            len(tb.textNodes),
		FirstWordNode:  tb.firstNonWhitespaceNode,
		LastWordNode:   tb.lastNonWhitespaceNode,
		NumWords:       tb.numWords,
		NumLinkedWords: tb.numAnchorWords,
		TagLevel:       tb.blockTagLevel,
		OffsetBlock:    offsetBlock,
	}

	tb.Reset()
	return &text
}
