// ORIGINAL: javatest/TestTextBlockBuilder.java

package testutil

import (
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

type TextBlockBuilder struct {
	textBuilder *TextBuilder
}

func NewTextBlockBuilder(wc stringutil.WordCounter) *TextBlockBuilder {
	return &TextBlockBuilder{
		textBuilder: NewTextBuilder(wc),
	}
}

func (tbb *TextBlockBuilder) CreateForText(text string) *webdoc.TextBlock {
	wt := tbb.textBuilder.CreateForText(text)
	return webdoc.NewTextBlock(wt)
}

func (tbb *TextBlockBuilder) CreateForAnchorText(text string) *webdoc.TextBlock {
	wt := tbb.textBuilder.CreateForAnchorText(text)
	return webdoc.NewTextBlock(wt)
}
