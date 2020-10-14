// ORIGINAL: java/webdocument/filters/NestedElementRetainer.java

package docfilter

import (
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

type NestedElementRetainer struct{}

func NewNestedElementRetainer() *NestedElementRetainer {
	return &NestedElementRetainer{}
}

func (f *NestedElementRetainer) Process(doc *webdoc.Document) bool {
	isContent := false
	stackMark := -1
	stack := []*webdoc.Tag{}

	for _, e := range doc.Elements {
		if webTag, isTag := e.(*webdoc.Tag); !isTag {
			if !isContent {
				isContent = e.IsContent()
			}
		} else {
			if webTag.Type == webdoc.TagStart {
				webTag.SetIsContent(isContent)
				stack = append(stack, webTag)
				isContent = false
			} else {
				startWebTag := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				isContent = isContent || stackMark >= len(stack)
				if isContent {
					stackMark = len(stack) - 1
				}

				wasContent := startWebTag.IsContent()
				startWebTag.SetIsContent(isContent)
				webTag.SetIsContent(isContent)
				isContent = wasContent
			}
		}
	}

	return true
}
