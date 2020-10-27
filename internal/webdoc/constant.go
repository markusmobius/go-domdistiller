// ORIGINAL: java/webdocument/WebText.java, java/webdocument/WebTag.java,
//           java/webdocument/WebImage.java

package webdoc

type TagType uint

const (
	TagStart TagType = iota
	TagEnd
)

var lazyImageAttrs = map[string]string{
	"data-srcset": "srcset",
}

func CanBeNested(tagName string) bool {
	switch tagName {
	case "ul", "ol", "li", "blockquote", "pre":
		return true

	default:
		return false
	}
}

// All inline elements except for impossible tags: br, object, and script.
// Please refer to DomConverter.visitElement() for skipped tags.
// Reference: https://developer.mozilla.org/en-US/docs/HTML/Inline_elements
var inlineTagNames = map[string]struct{}{}
