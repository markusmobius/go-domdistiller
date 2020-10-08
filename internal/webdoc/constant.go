// ORIGINAL: java/webdocument/WebText.java, java/webdocument/WebTag.java,
//           java/webdocument/WebImage.java

package webdoc

type TagType uint

const (
	TagStart TagType = iota
	TagEnd
)

var nestingTags = map[string]struct{}{
	"ul":         {},
	"ol":         {},
	"li":         {},
	"blockquote": {},
	"pre":        {},
}

// All inline elements except for impossible tags: br, object, and script.
// Please refer to DomConverter.visitElement() for skipped tags.
// Reference: https://developer.mozilla.org/en-US/docs/HTML/Inline_elements
var inlineTagNames = map[string]struct{}{
	"b":        {},
	"big":      {},
	"i":        {},
	"small":    {},
	"tt":       {},
	"abbr":     {},
	"acronym":  {},
	"cite":     {},
	"code":     {},
	"dfn":      {},
	"em":       {},
	"kbd":      {},
	"strong":   {},
	"samp":     {},
	"time":     {},
	"var":      {},
	"a":        {},
	"bdo":      {},
	"img":      {},
	"map":      {},
	"q":        {},
	"span":     {},
	"sub":      {},
	"sup":      {},
	"button":   {},
	"input":    {},
	"label":    {},
	"select":   {},
	"textarea": {},
}

var lazyImageAttrs = map[string]string{
	"data-srcset": "srcset",
}
