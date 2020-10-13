// ORIGINAL: java/filters/heuristics/HeadingFusion.java

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// HeadingFusion fuses headings with the blocks after them. If the heading was
// marked as boilerplate, the fused block will be labeled to prevent
// BlockProximityFusion from merging through it.
type HeadingFusion struct{}

func NewHeadingFusion() *HeadingFusion {
	return &HeadingFusion{}
}

func (f *HeadingFusion) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	changes := false
	currentBlock := textBlocks[0]
	var prevBlock *webdoc.TextBlock

	for i := 1; i < len(textBlocks); i++ {
		prevBlock = currentBlock
		currentBlock = textBlocks[i]

		if !prevBlock.HasLabel(label.Heading) {
			continue
		}

		if prevBlock.HasLabel(label.StrictlyNotContent) || currentBlock.HasLabel(label.StrictlyNotContent) {
			continue
		}

		if prevBlock.HasLabel(label.Title) || currentBlock.HasLabel(label.Title) {
			continue
		}

		if currentBlock.IsContent() {
			changes = true

			headingWasContent := prevBlock.IsContent()
			prevBlock.MergeNext(currentBlock)
			currentBlock = prevBlock

			currentBlock.RemoveLabels(label.Heading)
			if !headingWasContent {
				currentBlock.AddLabels(label.BoilerplateHeadingFused)
			}

			// These lines is used to remove item from array.
			copy(textBlocks[i:], textBlocks[i+1:])
			textBlocks[len(textBlocks)-1] = nil
			textBlocks = textBlocks[:len(textBlocks)-1]
			i--
		} else if prevBlock.IsContent() {
			changes = true
			prevBlock.SetIsContent(false)
		}
	}

	doc.TextBlocks = textBlocks
	return changes
}
