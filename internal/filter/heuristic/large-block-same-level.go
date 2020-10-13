// ORIGINAL: java/filters/heuristics/LargeBlockSameTagLevelToContentFilter.java

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// LargeBlockSameTagLevelToContent marks all blocks as content that:
// - are on the same tag-level as very likely main content (usually the level of the largest block)
// - have a significant number of words, currently: at least 100
type LargeBlockSameTagLevelToContent struct{}

func NewLargeBlockSameTagLevelToContent() *LargeBlockSameTagLevelToContent {
	return &LargeBlockSameTagLevelToContent{}
}

func (f *LargeBlockSameTagLevelToContent) Process(doc *webdoc.TextDocument) bool {
	tagLevel := -1
	for _, tb := range doc.TextBlocks {
		if tb.IsContent() && tb.HasLabel(label.VeryLikelyContent) {
			tagLevel = tb.TagLevel
			break
		}
	}

	if tagLevel == -1 {
		return false
	}

	changes := false
	for _, tb := range doc.TextBlocks {
		if !tb.IsContent() {
			if tb.NumWords >= 100 && tb.TagLevel == tagLevel {
				tb.SetIsContent(true)
				changes = true
			}
		}
	}

	return changes
}
