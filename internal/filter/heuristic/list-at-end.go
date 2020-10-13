// ORIGINAL: java/filters/heuristics/ListAtEndFilter.java

package heuristic

import (
	"math"

	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// ListAtEnd marks nested list-item blocks after the end of the main content.
type ListAtEnd struct{}

func NewListAtEnd() *ListAtEnd {
	return &ListAtEnd{}
}

func (f *ListAtEnd) Process(doc *webdoc.TextDocument) bool {
	changes := false
	tagLevel := math.MaxInt16

	for _, tb := range doc.TextBlocks {
		if tb.IsContent() && tb.HasLabel(label.VeryLikelyContent) {
			tagLevel = tb.TagLevel
		} else {
			if tb.TagLevel > tagLevel && tb.HasLabel(label.MightBeContent) &&
				tb.HasLabel(label.Li) && tb.LinkDensity == 0 {
				tb.SetIsContent(true)
				changes = true
			} else {
				tagLevel = math.MaxInt16
			}
		}
	}

	return changes
}
