// ORIGINAL: java/webdocument/WebElement.java

package webdoc

// Element is some logical part of a web document (text block, image, video, table, etc.)
type Element interface {
	// GenerateOutput generates HTML output for this Element.
	GenerateOutput(textOnly bool) string
}

// BaseElement is base of any other element.
type BaseElement struct {
	IsContent bool
}
