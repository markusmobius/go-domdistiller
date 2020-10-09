// ORIGINAL: java/document/TextBlock.java

package document

import (
	"fmt"
	"sort"
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// TextBlock describes a block of text. A block can be an "atomic" text node (i.e., a sequence
// of text that is not interrupted by any HTML markup) or a compound of such atomic elements.
type TextBlock struct {
	TextElements     []webdoc.Text
	TextIndexes      []int
	Text             string
	Labels           map[string]struct{}
	NumWords         int
	NumWordsInAnchor int
	LinkDensity      float64
	TagLevel         int

	isContent bool
}

func NewTextBlock(textElements []webdoc.Text, index int) *TextBlock {
	wt := textElements[index]
	tb := &TextBlock{
		TextElements:     textElements,
		TextIndexes:      []int{index},
		Text:             wt.Text,
		Labels:           wt.TakeLabels(),
		NumWords:         wt.NumWords,
		NumWordsInAnchor: wt.NumLinkedWords,
		TagLevel:         wt.TagLevel,
	}

	tb.LinkDensity = tb.calcLinkDensity()
	return tb
}

func (tb *TextBlock) IsContent() bool {
	return tb.isContent
}

// SetIsContent set the value of isContent.
// Returns true if isContent value changed.
func (tb *TextBlock) SetIsContent(isContent bool) bool {
	if isContent == tb.isContent {
		return false
	}

	tb.isContent = isContent
	return true
}

func (tb *TextBlock) MergeNext(other TextBlock) {
	tb.Text += "\n" + other.Text
	tb.NumWords += other.NumWords
	tb.NumWordsInAnchor += other.NumWordsInAnchor
	tb.LinkDensity = tb.calcLinkDensity()
	tb.isContent = tb.isContent || other.isContent
	tb.TextIndexes = append(tb.TextIndexes, other.TextIndexes...)

	for label := range other.Labels {
		tb.AddLabels(label)
	}

	if other.TagLevel < tb.TagLevel {
		tb.TagLevel = other.TagLevel
	}
}

func (tb *TextBlock) AddLabels(labels ...string) {
	for _, label := range labels {
		tb.Labels[label] = struct{}{}
	}
}

func (tb *TextBlock) RemoveLabels(labels ...string) {
	for _, label := range labels {
		delete(tb.Labels, label)
	}
}

func (tb *TextBlock) HasLabel(label string) bool {
	_, exist := tb.Labels[label]
	return exist
}

func (tb *TextBlock) OffsetBlocksStart() int {
	return tb.firstText().OffsetBlock
}

func (tb *TextBlock) OffsetBlocksEnd() int {
	return tb.lastText().OffsetBlock
}

func (tb *TextBlock) FirstNonWhitespaceTextNode() *html.Node {
	return tb.firstText().FirstNonWhitespaceTextNode()
}

func (tb *TextBlock) LastNonWhitespaceTextNode() *html.Node {
	return tb.firstText().LastNonWhitespaceTextNode()
}

func (tb *TextBlock) ApplyToModel() {
	if !tb.isContent {
		return
	}

	for _, idx := range tb.TextIndexes {
		wt := tb.TextElements[idx]
		wt.SetIsContent(true)
		if tb.HasLabel(label.Title) {
			wt.AddLabel(label.Title)
		}
	}
}

func (tb *TextBlock) String() string {
	str := "["
	str += fmt.Sprintf("%d-%d;", tb.OffsetBlocksStart(), tb.OffsetBlocksEnd())
	str += fmt.Sprintf("tl=%d;", tb.TagLevel)
	str += fmt.Sprintf("nw=%d;", tb.NumWords)
	str += fmt.Sprintf("ld=%.3f;", tb.LinkDensity)
	str += "]\t"

	if tb.isContent {
		str += "CONTENT,"
	} else {
		str += "boilerplate,"
	}

	str += tb.labelsDebugString() + "\n" + tb.Text
	return str
}

func (tb *TextBlock) calcLinkDensity() float64 {
	if tb.NumWords == 0 {
		return 0
	}

	return float64(tb.NumWordsInAnchor) / float64(tb.NumWords)
}

func (tb *TextBlock) firstText() webdoc.Text {
	return tb.TextElements[0]
}

func (tb *TextBlock) lastText() webdoc.Text {
	return tb.TextElements[len(tb.TextElements)-1]
}

func (tb *TextBlock) labelsDebugString() string {
	labels := []string{}
	for label := range tb.Labels {
		labels = append(labels, label)
	}

	sort.Strings(labels)
	return strings.Join(labels, ",")
}
