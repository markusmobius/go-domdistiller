// ORIGINAL: javatest/BlockProximityFusionTest.java

package heuristic_test

import (
	"strings"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_Heuristic_BPF_MergeShortLeadingContent(t *testing.T) {
	bpfTestMergeShortLeadingContent(t, heuristic.NewBlockProximityFusion(false))
	bpfTestMergeShortLeadingContent(t, heuristic.NewBlockProximityFusion(true))
}

func Test_Filter_Heuristic_BPF_DoNotMergeShortLeadingLiNonContent(t *testing.T) {
	bpfTestDoNotMergeShortLeadingLiNonContent(t, heuristic.NewBlockProximityFusion(false))
	bpfTestDoNotMergeShortLeadingLiNonContent(t, heuristic.NewBlockProximityFusion(true))
}

func Test_Filter_Heuristic_BPF_DoNotMergeShortLeadingNonContent(t *testing.T) {
	bpfTestDoNotMergeShortLeadingNonContent(t, heuristic.NewBlockProximityFusion(false))
	bpfTestDoNotMergeShortLeadingNonContent(t, heuristic.NewBlockProximityFusion(true))
}

func Test_Filter_Heuristic_BPF_MergeLotsOfContent(t *testing.T) {
	bpfTestMergeLotsOfContent(t, heuristic.NewBlockProximityFusion(false))
	bpfTestMergeLotsOfContent(t, heuristic.NewBlockProximityFusion(true))
}

func Test_Filter_Heuristic_BPF_SkipNonContentInBody(t *testing.T) {
	bpfTestSkipNonContentInBody(t, heuristic.NewBlockProximityFusion(false))
	bpfTestSkipNonContentInBody(t, heuristic.NewBlockProximityFusion(true))
}

func Test_Filter_Heuristic_BPF_PreFilteringSkipNonContentListInBody(t *testing.T) {
	// If "content" flag is ignored, a single non-content Li in the body is not merged.
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(longLeadingText)
	builder.AddContentBlock(longLeadingText)
	builder.AddNonContentBlock(shortText, label.Li)
	builder.AddContentBlock(longText)
	textDocument := builder.Build()
	docContent := testutil.GetContentFromTextDocument(textDocument)

	classifier := heuristic.NewBlockProximityFusion(false)
	classifier.Process(textDocument)

	assert.Equal(t, 3, len(textDocument.TextBlocks))
	assert.True(t, strings.Contains(docContent, longLeadingText))
	assert.False(t, strings.Contains(docContent, shortText))
	assert.True(t, strings.Contains(docContent, longText))
}

func bpfTestMergeShortLeadingContent(t *testing.T, classifier *heuristic.BlockProximityFusion) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(shortText)
	builder.AddContentBlock(longText)
	textDocument := builder.Build()
	docContent := testutil.GetContentFromTextDocument(textDocument)

	classifier.Process(textDocument)
	assert.Equal(t, 1, len(textDocument.TextBlocks))
	assert.True(t, strings.Contains(docContent, shortText))
	assert.True(t, strings.Contains(docContent, longText))
}

func bpfTestDoNotMergeShortLeadingLiNonContent(t *testing.T, classifier *heuristic.BlockProximityFusion) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddNonContentBlock(shortText, label.Li)
	builder.AddContentBlock(longText)
	textDocument := builder.Build()
	docContent := testutil.GetContentFromTextDocument(textDocument)

	classifier.Process(textDocument)
	assert.Equal(t, 2, len(textDocument.TextBlocks))
	assert.False(t, strings.Contains(docContent, shortText))
	assert.True(t, strings.Contains(docContent, longText))
}

func bpfTestDoNotMergeShortLeadingNonContent(t *testing.T, classifier *heuristic.BlockProximityFusion) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddNonContentBlock(shortText)
	builder.AddContentBlock(longText)
	textDocument := builder.Build()
	docContent := testutil.GetContentFromTextDocument(textDocument)

	classifier.Process(textDocument)
	assert.Equal(t, 2, len(textDocument.TextBlocks))
	assert.False(t, strings.Contains(docContent, shortText))
	assert.True(t, strings.Contains(docContent, longText))
}

func bpfTestMergeLotsOfContent(t *testing.T, classifier *heuristic.BlockProximityFusion) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(longLeadingText)
	builder.AddContentBlock(longLeadingText)
	builder.AddContentBlock(shortText)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()
	docContent := testutil.GetContentFromTextDocument(textDocument)

	classifier.Process(textDocument)
	assert.Equal(t, 1, len(textDocument.TextBlocks))
	assert.True(t, strings.Contains(docContent, longLeadingText))
	assert.True(t, strings.Contains(docContent, shortText))
	assert.True(t, strings.Contains(docContent, longText))
}

func bpfTestSkipNonContentInBody(t *testing.T, classifier *heuristic.BlockProximityFusion) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(longLeadingText)
	builder.AddContentBlock(longLeadingText)
	builder.AddNonContentBlock(shortText)
	builder.AddContentBlock(longText)
	textDocument := builder.Build()
	docContent := testutil.GetContentFromTextDocument(textDocument)

	classifier.Process(textDocument)
	assert.Equal(t, 3, len(textDocument.TextBlocks))
	assert.True(t, strings.Contains(docContent, longLeadingText))
	assert.False(t, strings.Contains(docContent, shortText))
	assert.True(t, strings.Contains(docContent, longText))
}
