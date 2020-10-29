// ORIGINAL: java/filters/heuristics/LargeBlockSameTagLevelToContentFilter.java

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// LargeBlockAroundTagLevelToContent marks all blocks as content that:
// - are on the same or adjacent tag-level as very likely main content (usually the level of the largest block)
// - have a significant number of words, currently: at least 100
type LargeBlockAroundTagLevelToContent struct{}

func NewLargeBlockAroundTagLevelToContent() *LargeBlockAroundTagLevelToContent {
	return &LargeBlockAroundTagLevelToContent{}
}

func (f *LargeBlockAroundTagLevelToContent) Process(doc *webdoc.TextDocument) bool {
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
		if tb.IsContent() || tb.NumWords < 100 {
			continue
		}

		switch tb.TagLevel {
		case tagLevel, tagLevel - 1, tagLevel + 1:
			tb.SetIsContent(true)
			changes = true
		}
	}

	return changes
}
