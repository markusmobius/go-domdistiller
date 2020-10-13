// ORIGINAL: java/filters/english/TerminatingBlocksFinder.java

package english

import (
	"regexp"
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

var rxTerminatingBlocks = regexp.MustCompile(`(?i)(` +
	`^(comments|© reuters|please rate this|post a comment|` +
	`\d+\s+(comments|users responded in)` +
	`)` +
	`|what you think\.\.\.` +
	`|add your comment` +
	`|add comment` +
	`|reader views` +
	`|have your say` +
	`|reader comments` +
	`|rätta artikeln` +
	`|^thanks for your comments - this feedback is now closed$` +
	`)`)

// TerminatingBlocksFinder finds blocks which are potentially indicating the end of
// an article text and marks them with label.StrictlyNotContent.
type TerminatingBlocksFinder struct{}

func NewTerminatingBlocksFinder() *TerminatingBlocksFinder {
	return &TerminatingBlocksFinder{}
}

func (f *TerminatingBlocksFinder) Process(doc *webdoc.TextDocument) bool {
	changes := false

	for _, block := range doc.TextBlocks {
		if f.isTerminating(block) {
			block.AddLabels(label.StrictlyNotContent)
			changes = true
		}
	}

	return changes
}

func (f *TerminatingBlocksFinder) isTerminating(tb *webdoc.TextBlock) bool {
	if tb.NumWords > 14 {
		return false
	}

	text := strings.TrimSpace(tb.Text)
	if stringutil.CharCount(text) >= 8 {
		return rxTerminatingBlocks.MatchString(text)
	} else if tb.LinkDensity == 1 {
		return text == "Comment"
	} else if text == "Shares" {
		// Skip social and sharing elements.
		// See crbug.com/692553
		return true
	}

	return false
}
