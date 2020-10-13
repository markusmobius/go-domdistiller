// ORIGINAL: javatest/document/TextDocumentTestUtil.java

package testutil

import (
	"bytes"

	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func GetContentFromTextDocument(doc *webdoc.TextDocument) string {
	buffer := bytes.NewBuffer(nil)
	for _, tb := range doc.TextBlocks {
		if tb.IsContent() {
			buffer.WriteString(tb.Text + "\n")
		}
	}
	return buffer.String()
}
