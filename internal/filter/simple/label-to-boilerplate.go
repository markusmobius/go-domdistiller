// ORIGINAL: java/filters/simple/LabelToBoilerplateFilter.java

package simple

import "github.com/markusmobius/go-domdistiller/internal/webdoc"

// LabelToBoilerplate marks all blocks that contain a given label as "boilerplate".
type LabelToBoilerplate struct {
	labels []string
}

func NewLabelToBoilerplate(labels ...string) *LabelToBoilerplate {
	return &LabelToBoilerplate{labels: labels}
}

func (f *LabelToBoilerplate) Process(doc *webdoc.TextDocument) bool {
	changes := false

blockLoop:
	for _, tb := range doc.TextBlocks {
		if tb.IsContent() {
			for _, label := range f.labels {
				if tb.HasLabel(label) {
					tb.SetIsContent(false)
					changes = true
					continue blockLoop
				}
			}
		}
	}

	return changes
}
