// ORIGINAL: javatest/HeadingFusionTest.java

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

func Test_Filter_Heuristic_HF_HeadingFused(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(headingText, label.Heading)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.True(t, changed)

	textBlocks := textDocument.TextBlocks
	assert.Len(t, textBlocks, 2)
	assert.False(t, textBlocks[0].HasLabel(label.Heading))
	assert.False(t, textBlocks[0].HasLabel(label.BoilerplateHeadingFused))

	docContent := testutil.GetContentFromTextDocument(textDocument)
	assert.True(t, strings.Contains(docContent, headingText))
	assert.True(t, strings.Contains(docContent, longText))
	assert.True(t, strings.Contains(docContent, shortText))
}

func Test_Filter_Heuristic_HF_BoilerplateHeadingFused(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddNonContentBlock(headingText, label.Heading)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.True(t, changed)

	textBlocks := textDocument.TextBlocks
	assert.Len(t, textBlocks, 2)
	assert.False(t, textBlocks[0].HasLabel(label.Heading))
	assert.True(t, textBlocks[0].HasLabel(label.BoilerplateHeadingFused))

	docContent := testutil.GetContentFromTextDocument(textDocument)
	assert.True(t, strings.Contains(docContent, headingText))
	assert.True(t, strings.Contains(docContent, longText))
	assert.True(t, strings.Contains(docContent, shortText))
}

func Test_Filter_Heuristic_HF_BeforeBoilerplate(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(headingText, label.Heading)
	builder.AddNonContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.True(t, changed)

	textBlocks := textDocument.TextBlocks
	assert.Len(t, textBlocks, 3)
	assert.False(t, textBlocks[0].IsContent())
}

func Test_Filter_Heuristic_HF_TitleNotFused(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(headingText, label.Heading, label.Title)
	builder.AddContentBlock(longText)
	builder.AddContentBlock(shortText)
	textDocument := builder.Build()

	changed := heuristic.NewHeadingFusion().Process(textDocument)
	assert.False(t, changed)
	assert.Len(t, textDocument.TextBlocks, 3)
}
