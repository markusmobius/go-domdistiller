// ORIGINAL: java/webdocument/WebDocument.java

package webdoc

import (
	"bytes"
)

// Document is a simplified view of the underlying webpage. It contains the
// logical elements (blocks of text, image + caption, video, etc).
type Document struct {
	Elements []Element
}

func (doc *Document) AddElements(elements ...Element) {
	doc.Elements = append(doc.Elements, elements...)
}

func (doc *Document) GenerateOutput(textOnly bool) string {
	buffer := bytes.NewBuffer(nil)
	for _, e := range doc.Elements {
		if !e.IsContent() {
			continue
		}

		buffer.WriteString(e.GenerateOutput(textOnly))
		if textOnly {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}

// CreateTextDocument generates a web document to be processed by distiller.
// Text groups have been introduced to help retain element order when adding
// images and embeds.
func (doc *Document) CreateTextDocument() *TextDocument {
	textBlocks := []*TextBlock{}
	firstTextIdx := doc.getNextTextIndex(0)
	if firstTextIdx == len(doc.Elements) {
		return NewTextDocument(nil)
	}

	currentText := doc.Elements[firstTextIdx].(*Text)
	currentGroup := currentText.GroupNumber
	previousGroup := currentGroup

	currentBlockTexts := []*Text{}
	for _, element := range doc.Elements {
		text, isText := element.(*Text)
		if !isText {
			continue
		}

		currentGroup = text.GroupNumber
		if currentGroup == previousGroup {
			currentBlockTexts = append(currentBlockTexts, text)
		} else {
			textBlocks = append(textBlocks, NewTextBlock(currentBlockTexts...))
			previousGroup = currentGroup
			currentBlockTexts = []*Text{text}
		}
	}

	textBlocks = append(textBlocks, NewTextBlock(currentBlockTexts...))
	return NewTextDocument(textBlocks)
}

// GetImageURLs returns list of source URLs of all image inside the document.
func (doc *Document) GetImageURLs() []string {
	imageURLs := []string{}
	for _, e := range doc.Elements {
		if !e.IsContent() {
			continue
		}

		// TODO: if we allow images in Text later, handle it here.
		switch element := e.(type) {
		case *Image:
			imageURLs = append(imageURLs, element.GetURLs()...)
		case *Table:
			imageURLs = append(imageURLs, element.GetImageURLs()...)
		}
	}

	return imageURLs
}

func (doc *Document) getNextTextIndex(startIndex int) int {
	for i := startIndex; i < len(doc.Elements); i++ {
		if _, isText := doc.Elements[i].(*Text); isText {
			return i
		}
	}

	return len(doc.Elements)
}
