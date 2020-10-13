// ORIGINAL: java/BoilerpipeFilter.java

package filter

import "github.com/markusmobius/go-domdistiller/internal/webdoc"

// Filter is interface for filter that process a TextDocument.
type Filter interface {
	// Process processes the given document.
	// Returns true if changes have been made to the document.
	Process(doc *webdoc.TextDocument) bool
}
