// ORIGINAL: javatest/webdocument/WebTagTest.java

package webdoc_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_WebDoc_Tag_OLGenerateOutput(t *testing.T) {
	olStartTag := webdoc.Tag{Name: "ol", Type: webdoc.TagStart}
	olEndTag := webdoc.Tag{Name: "ol", Type: webdoc.TagEnd}
	startResult := olStartTag.GenerateOutput(false)
	endResult := olEndTag.GenerateOutput(false)
	assert.Equal(t, "<ol>", startResult)
	assert.Equal(t, "</ol>", endResult)
}

func Test_WebDoc_Tag_GenerateOutput(t *testing.T) {
	startTag := webdoc.Tag{Name: "anytext", Type: webdoc.TagStart}
	endTag := webdoc.Tag{Name: "anytext", Type: webdoc.TagEnd}
	startResult := startTag.GenerateOutput(false)
	endResult := endTag.GenerateOutput(false)
	assert.Equal(t, "<anytext>", startResult)
	assert.Equal(t, "</anytext>", endResult)
}
