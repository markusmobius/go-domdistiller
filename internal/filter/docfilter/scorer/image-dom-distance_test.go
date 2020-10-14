// ORIGINAL: javatest/ImageHeuristicsTest.java

package scorer_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter/scorer"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_Scorer_ImageDomDistanceScorer(t *testing.T) {
	root := testutil.CreateDiv(0)
	content := testutil.CreateDiv(1)
	image := dom.CreateElement("img")
	dom.SetAttribute(image, "style", "width: 100px; height: 100px; display: block;")

	dom.AppendChild(content, image)
	dom.AppendChild(root, content)

	// Build long chain of divs to separate image from content.
	currentDiv := testutil.CreateDiv(3)
	dom.AppendChild(root, currentDiv)
	for i := 0; i < 7; i++ {
		child := testutil.CreateDiv(i + 4)
		dom.AppendChild(currentDiv, child)
		currentDiv = child
	}

	normalScorer := scorer.NewImageDomDistanceScorer(50, content)
	farContentScorer := scorer.NewImageDomDistanceScorer(50, currentDiv)

	assert.True(t, normalScorer.GetImageScore(image) > 0)
	assert.Equal(t, 0, farContentScorer.GetImageScore(image))
	assert.Equal(t, 0, normalScorer.GetImageScore(nil))
}
