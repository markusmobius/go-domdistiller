// ORIGINAL: javatest/webdocument/filters/RelevantElementsTest.java

package docfilter_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_RE_EmptyDocument(t *testing.T) {
	document := webdoc.NewDocument()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.True(t, len(document.Elements) == 0)
}

func Test_Filter_DocFilter_RE_NoContent(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1")
	builder.AddText("text 2")
	builder.AddTable("<tbody><tr><td>t1</td></tr></tbody>")

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	for _, e := range document.Elements {
		assert.False(t, e.IsContent())
	}
}

func Test_Filter_DocFilter_RE_RelevantTable(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	wt := builder.AddTable("<tbody><tr><td>t1</td></tr></tbody>")

	document := builder.Build()
	assert.True(t, docfilter.NewRelevantElements().Process(document))
	assert.True(t, wt.IsContent())
}

func Test_Filter_DocFilter_RE_NonRelevantTable(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	builder.AddText("text 2")
	wt := builder.AddTable("<tbody><tr><td>t1</td></tr></tbody>")

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.False(t, wt.IsContent())
}

func Test_Filter_DocFilter_RE_RelevantImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	wi := builder.AddImage()

	document := builder.Build()
	assert.True(t, docfilter.NewRelevantElements().Process(document))
	assert.True(t, wi.IsContent())
}

func Test_Filter_DocFilter_RE_NonRelevantImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddImage()
	builder.AddText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.False(t, wi.IsContent())
}

func Test_Filter_DocFilter_RE_ImageAfterNonContent(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	builder.AddText("text 2").SetIsContent(false)
	wi := builder.AddImage()

	document := builder.Build()
	assert.False(t, docfilter.NewRelevantElements().Process(document))
	assert.False(t, wi.IsContent())
}
