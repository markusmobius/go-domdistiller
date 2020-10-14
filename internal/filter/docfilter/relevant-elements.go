// ORIGINAL: java/webdocument/filters/RelevantElements.java

package docfilter

import "github.com/markusmobius/go-domdistiller/internal/webdoc"

type RelevantElements struct{}

func NewRelevantElements() *RelevantElements {
	return &RelevantElements{}
}

func (f *RelevantElements) Process(doc *webdoc.Document) bool {
	changes := false
	inContent := false

	for _, e := range doc.Elements {
		if e.IsContent() {
			inContent = true
		} else if _, isText := e.(*webdoc.Text); isText {
			inContent = false
		} else {
			if inContent {
				e.SetIsContent(true)
				changes = true
			}
		}
	}

	return changes
}
