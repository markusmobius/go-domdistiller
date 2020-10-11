// ORIGINAL: java/extractors/embeds/EmbedExtractor.java

package extractor

import (
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// EmbedExtractor is interface for extracting embedded nodes int webdoc.Element.
type EmbedExtractor interface {
	// RelevantTagNames returns a set of HTML tag names that are relevant to this extractor.
	RelevantTagNames() []string
	// Extract detects if a node should be extracted as an embedded element; if not return nil.
	Extract(node *html.Node) webdoc.Element
}
