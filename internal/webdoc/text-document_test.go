// ORIGINAL: javatest/TextDocumentStatisticsTest.java

package webdoc_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

const ThreeWords = "I love statistics"

func Test_WebDoc_TextDocument_OnlyContent(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(ThreeWords)
	builder.AddContentBlock(ThreeWords)
	builder.AddContentBlock(ThreeWords)

	doc := builder.Build()
	assert.Equal(t, 9, doc.CountWordsInContent())
}

func Test_WebDoc_TextDocument_OnlyNonContent(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddNonContentBlock(ThreeWords)
	builder.AddNonContentBlock(ThreeWords)
	builder.AddNonContentBlock(ThreeWords)

	doc := builder.Build()
	assert.Equal(t, 0, doc.CountWordsInContent())
}

func Test_WebDoc_TextDocument_MixedContent(t *testing.T) {
	builder := testutil.NewTextDocumentBuilder(stringutil.FastWordCounter{})
	builder.AddContentBlock(ThreeWords)
	builder.AddNonContentBlock(ThreeWords)
	builder.AddContentBlock(ThreeWords)
	builder.AddNonContentBlock(ThreeWords)

	doc := builder.Build()
	assert.Equal(t, 6, doc.CountWordsInContent())
}
