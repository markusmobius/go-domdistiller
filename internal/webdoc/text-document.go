// ORIGINAL: java/document/TextDocument.java and
//           java/document/TextDocumentStatistics.java

package webdoc

import "bytes"

// TextDocument is a text document, consisting of one or more TextBlock.
type TextDocument struct {
	TextBlocks []*TextBlock
}

func NewTextDocument(textBlocks []*TextBlock) *TextDocument {
	return &TextDocument{textBlocks}
}

func (td *TextDocument) ApplyToModel() {
	for _, tb := range td.TextBlocks {
		tb.ApplyToModel()
	}
}

// CountWordsInContent returns the sum of number of words in content blocks.
func (td *TextDocument) CountWordsInContent() int {
	numWords := 0
	for _, tb := range td.TextBlocks {
		if tb.IsContent() {
			numWords += tb.NumWords
		}
	}
	return numWords
}

// DebugString returns detailed debugging information about the contained TextBlocks.
func (td *TextDocument) DebugString() string {
	buffer := bytes.NewBuffer(nil)
	for _, tb := range td.TextBlocks {
		buffer.WriteString(tb.String() + "\n")
	}
	return buffer.String()
}
