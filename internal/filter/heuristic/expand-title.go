// ORIGINAL: java/filters/heuristics/ExpandTitleToContentFilter.java

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// ExpandTitleToContent marks all TextBlocks "content" which are between the headline and the
// part that has already been marked content, if they are marked with label.MightBeContent.
// This filter is quite specific to the news domain.
type ExpandTitleToContent struct{}

func NewExpandTitleToContent() *ExpandTitleToContent {
	return &ExpandTitleToContent{}
}

func (f *ExpandTitleToContent) Process(doc *webdoc.TextDocument) bool {
	title := -1
	contentStart := -1
	for i, tb := range doc.TextBlocks {
		if contentStart == -1 && tb.HasLabel(label.Title) {
			title = i
			contentStart = -1
		}

		if contentStart == -1 && tb.IsContent() {
			contentStart = i
		}
	}

	if contentStart <= title || title == -1 {
		return false
	}

	changes := false
	for _, tb := range doc.TextBlocks[title:contentStart] {
		if tb.HasLabel(label.MightBeContent) {
			changed := tb.SetIsContent(true)
			changes = changes || changed
		}
	}

	return changes
}
