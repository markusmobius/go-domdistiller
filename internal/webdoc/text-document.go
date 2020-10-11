// ORIGINAL: java/document/TextDocument.java and
//           java/document/TextDocumentStatistics.java

package webdoc

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

func (td *TextDocument) CountWordsInContent() int {
	numWords := 0
	for _, tb := range td.TextBlocks {
		if tb.IsContent() {
			numWords += tb.NumWords
		}
	}
	return numWords
}
