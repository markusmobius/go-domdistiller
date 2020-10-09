// ORIGINAL: javatest/webdocument/WebTextBuilderTest.java

package webdoc_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_WebDoc_TextBuilder_SimpleBlocks(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := webdoc.NewTextBuilder(wc)

	block := builder.Build(0)
	assert.Nil(t, block)

	tbAddText(builder, "Two words.", 0)
	block = builder.Build(0)
	assert.Equal(t, 2, block.NumWords)
	assert.Equal(t, 0, block.NumLinkedWords)
	assert.Equal(t, "Two words.", block.Text)
	assert.Equal(t, 1, len(block.TextNodes()))
	assert.Equal(t, 0, len(block.Labels))
	assert.Equal(t, 0, block.OffsetBlock)

	tbAddText(builder, "More", 0)
	tbAddText(builder, " than", 0)
	tbAddText(builder, " two", 0)
	tbAddText(builder, " words.", 0)
	block = builder.Build(1)
	assert.Equal(t, 4, block.NumWords)
	assert.Equal(t, 0, block.NumLinkedWords)
	assert.Equal(t, "More than two words.", block.Text)
	assert.Equal(t, 4, len(block.TextNodes()))
	assert.Equal(t, 1, block.OffsetBlock)

	assert.Nil(t, builder.Build(0))
	assert.Nil(t, builder.Build(0))
}

func Test_WebDoc_TextBuilder_BlockWithAnchors(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := webdoc.NewTextBuilder(wc)

	tbAddText(builder, "one", 0)
	builder.EnterAnchor()
	tbAddText(builder, "two", 0)
	tbAddText(builder, " three", 0)
	builder.ExitAnchor()

	block := builder.Build(0)
	assert.Equal(t, 3, block.NumWords)
	assert.Equal(t, 2, block.NumLinkedWords)
	assert.Equal(t, "one two three ", block.Text)

	builder.EnterAnchor()
	tbAddText(builder, "one", 0)
	block = builder.Build(0)
	assert.Equal(t, 1, block.NumWords)
	assert.Equal(t, 1, block.NumLinkedWords)
	assert.Equal(t, " one", block.Text)

	// Should still be in the previous anchor.
	tbAddText(builder, "one", 0)
	block = builder.Build(0)
	assert.Equal(t, 1, block.NumWords)
	assert.Equal(t, 1, block.NumLinkedWords)
	assert.Equal(t, "one", block.Text)
}

func Test_WebDoc_TextBuilder_ComplicatedText(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := webdoc.NewTextBuilder(wc)

	tbAddText(builder, "JULIE'S CALAMARI", 0)
	assert.Equal(t, 2, builder.Build(0).NumWords)
	tbAddText(builder, "all-purpose flour", 0)
	assert.Equal(t, 2, builder.Build(0).NumWords)
	tbAddText(builder, "1/2 cups flour", 0)
	assert.Equal(t, 3, builder.Build(0).NumWords)
	tbAddText(builder, "email foo@bar.com", 0)
	assert.Equal(t, 2, builder.Build(0).NumWords)
	tbAddText(builder, "$2.38 million", 0)
	assert.Equal(t, 2, builder.Build(0).NumWords)
	tbAddText(builder, "goto website.com", 0)
	assert.Equal(t, 2, builder.Build(0).NumWords)
	tbAddText(builder, "Deal expires:7d:04h:23m", 0)
	assert.Equal(t, 2, builder.Build(0).NumWords)
}

func Test_WebDoc_TextBuilder_WhitespaceAroundAnchor(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := webdoc.NewTextBuilder(wc)

	tbAddText(builder, "The ", 0)
	builder.EnterAnchor()
	tbAddText(builder, "Overview", 0)
	builder.ExitAnchor()
	tbAddText(builder, " is", 0)
	tb := builder.Build(0)
	assert.Equal(t, "The  Overview  is", tb.Text)
}

func Test_WebDoc_TextBuilder_WhitespaceNodes(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := webdoc.NewTextBuilder(wc)

	tbAddText(builder, "one", 0)
	tb := builder.Build(0)
	assert.Equal(t, tb.FirstNonWhitespaceTextNode(), tb.LastNonWhitespaceTextNode())

	tbAddText(builder, " ", 0)
	tbAddText(builder, "one", 0)
	tbAddText(builder, " ", 0)
	tb = builder.Build(0)
	assert.Equal(t, tb.FirstNonWhitespaceTextNode(), tb.LastNonWhitespaceTextNode())
	assert.Equal(t, 3, len(tb.TextNodes()))
	assert.False(t, tb.TextNodes()[0] == tb.FirstNonWhitespaceTextNode())
	assert.False(t, tb.TextNodes()[2] == tb.FirstNonWhitespaceTextNode())

	tbAddText(builder, "one", 0)
	tbAddText(builder, "two", 0)
	tb = builder.Build(0)
	assert.False(t, tb.FirstNonWhitespaceTextNode() == tb.LastNonWhitespaceTextNode())
}

func Test_WebDoc_TextBuilder_BrElement(t *testing.T) {
	wc := stringutil.FastWordCounter{}
	builder := webdoc.NewTextBuilder(wc)

	tbAddText(builder, "words", 0)
	builder.AddLineBreak(dom.CreateElement("br"))
	tbAddText(builder, "split", 0)
	builder.AddLineBreak(dom.CreateElement("br"))
	tbAddText(builder, "with", 0)
	builder.AddLineBreak(dom.CreateElement("br"))
	tbAddText(builder, "lines", 0)

	webText := builder.Build(0)
	assert.Equal(t, 7, len(webText.TextNodes()))
	assert.Equal(t, "words\nsplit\nwith\nlines", webText.Text)
}

func tbAddText(builder *webdoc.TextBuilder, text string, tagLevel int) {
	builder.AddTextNode(dom.CreateTextNode(text), tagLevel)
}
