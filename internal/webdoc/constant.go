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

func CanBeNested(tagName string) bool {
	_, canBeNested := nestingTags[tagName]
	return canBeNested
}

// All inline elements except for impossible tags: br, object, and script.
// Please refer to DomConverter.visitElement() for skipped tags.
// Reference: https://developer.mozilla.org/en-US/docs/HTML/Inline_elements
var inlineTagNames = map[string]struct{}{
	"a":        {},
	"abbr":     {},
	"acronym":  {},
	"b":        {},
	"bdi":      {},
	"bdo":      {},
	"big":      {},
	"button":   {},
	"cite":     {},
	"code":     {},
	"dfn":      {},
	"em":       {},
	"i":        {},
	"img":      {},
	"input":    {},
	"kbd":      {},
	"label":    {},
	"map":      {},
	"q":        {},
	"s":        {},
	"samp":     {},
	"select":   {},
	"small":    {},
	"span":     {},
	"strong":   {},
	"sub":      {},
	"sup":      {},
	"textarea": {},
	"time":     {},
	"tt":       {},
	"u":        {},
	"var":      {},
}

var lazyImageAttrs = map[string]string{
	"data-srcset": "srcset",
}
