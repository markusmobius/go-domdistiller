// ORIGINAL: java/filters/heuristics/BlockProximityFusion.java

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// BlockProximityFusion fuses adjacent blocks if their distance (in blocks) does not
// exceed a certain limit. This probably makes sense only in cases where an upstream
// filter already has removed some blocks.
type BlockProximityFusion struct {
	postFiltering bool
}

func NewBlockProximityFusion(postFiltering bool) *BlockProximityFusion {
	return &BlockProximityFusion{
		postFiltering: postFiltering,
	}
}

func (f *BlockProximityFusion) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	changes := false
	prevBlock := textBlocks[0]

	for i := 1; i < len(textBlocks); i++ {
		block := textBlocks[i]
		if !block.IsContent() || !prevBlock.IsContent() {
			prevBlock = block
			continue
		}

		diffBlocks := block.OffsetBlocksStart() - prevBlock.OffsetBlocksEnd() - 1
		if diffBlocks <= 1 {
			ok := true
			if f.postFiltering {
				if prevBlock.TagLevel != block.TagLevel {
					ok = false
				}
			} else {
				if block.HasLabel(label.BoilerplateHeadingFused) {
					ok = false
				}
			}

			if prevBlock.HasLabel(label.StrictlyNotContent) != block.HasLabel(label.StrictlyNotContent) {
				ok = false
			}

			if prevBlock.HasLabel(label.Title) != block.HasLabel(label.Title) {
				ok = false
			}

			if (!prevBlock.IsContent() && prevBlock.HasLabel(label.Li)) && !block.HasLabel(label.Li) {
				ok = false
			}

			if ok {
				changes = true
				prevBlock.MergeNext(block)

				// These lines is used to remove item from array.
				copy(textBlocks[i:], textBlocks[i+1:])
				textBlocks[len(textBlocks)-1] = nil
				textBlocks = textBlocks[:len(textBlocks)-1]
				i--
			}
		} else {
			prevBlock = block
		}
	}

	doc.TextBlocks = textBlocks
	return changes
}
