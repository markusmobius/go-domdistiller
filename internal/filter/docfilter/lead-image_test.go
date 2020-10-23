// ORIGINAL: javatest/webdocument/filters/LeadImageFinderTest.java

package docfilter_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_LIF_EmptyDocument(t *testing.T) {
	document := webdoc.NewDocument()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.Len(t, document.Elements, 0)
}

func Test_Filter_DocFilter_LIF_RelevantLeadImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddLeadImage()
	builder.AddText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.True(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, wi.IsContent())
}

func Test_Filter_DocFilter_LIF_NoContent(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddLeadImage()
	builder.AddText("text 1").SetIsContent(false)

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.False(t, wi.IsContent())
}

func Test_Filter_DocFilter_LIF_MultipleLeadImageCandidates(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	priority := builder.AddLeadImage()
	ignored := builder.AddLeadImage()
	builder.AddText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.True(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, priority.IsContent())
	assert.False(t, ignored.IsContent())
}

func Test_Filter_DocFilter_LIF_IrrelevantLeadImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	builder.AddText("text 2").SetIsContent(true)
	priority := builder.AddLeadImage()

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.False(t, priority.IsContent())
}

func Test_Filter_DocFilter_LIF_MultipleLeadImageCandidatesWithinTexts(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	builder.AddText("text 1").SetIsContent(true)
	priority := builder.AddLeadImage()
	builder.AddText("text 2").SetIsContent(true)
	ignored := builder.AddLeadImage()

	document := builder.Build()
	assert.True(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, priority.IsContent())
	assert.False(t, ignored.IsContent())
}

func Test_Filter_DocFilter_LIF_IrrelevantLeadImageWithContentImage(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	smallImage := builder.AddImage()
	smallImage.SetIsContent(true)
	largeImage := builder.AddLeadImage()
	builder.AddNestedText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.True(t, smallImage.IsContent())
	assert.False(t, largeImage.IsContent())
}

func Test_Filter_DocFilter_LIF_IrrelevantLeadImageAfterSingleText(t *testing.T) {
	builder := testutil.NewWebDocumentBuilder()
	wi := builder.AddImage()
	builder.AddNestedText("text 1").SetIsContent(true)

	document := builder.Build()
	assert.False(t, docfilter.NewLeadImageFinder(nil).Process(document))
	assert.False(t, wi.IsContent())
}
