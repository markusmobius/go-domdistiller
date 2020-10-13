// ORIGINAL: java/filters/simple/BoilerplateBlockFilter.java

package simple

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// BoilerplateBlock removes TextBlocks which have explicitly been
// marked as "not content".
type BoilerplateBlock struct {
	labelToKeep string
}

func NewBoilerplateBlock(labelToKeep string) *BoilerplateBlock {
	return &BoilerplateBlock{labelToKeep: labelToKeep}
}

func (f *BoilerplateBlock) Process(doc *webdoc.TextDocument) bool {
	hasChanges := false
	textBlocks := doc.TextBlocks

	for i := 0; i < len(textBlocks); i++ {
		tb := textBlocks[i]
		if !tb.IsContent() && (f.labelToKeep != "" || !tb.HasLabel(label.Title)) {
			hasChanges = true

			// These lines is used to remove item from array.
			copy(textBlocks[i:], textBlocks[i+1:])
			textBlocks[len(textBlocks)-1] = nil
			textBlocks = textBlocks[:len(textBlocks)-1]
			i--
		}
	}

	return hasChanges
}
