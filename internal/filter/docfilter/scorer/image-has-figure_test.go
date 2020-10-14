// ORIGINAL: javatest/ImageHeuristicsTest.java

package scorer_test

import (
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter/scorer"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_DocFilter_Scorer_ImageHasFigureScorer(t *testing.T) {
	root := testutil.CreateDiv(0)
	fig := dom.CreateElement("figure")

	goodImage := dom.CreateElement("img")
	dom.SetAttribute(goodImage, "style", "width: 100px; height: 100px; display: block;")

	badImage := dom.CreateElement("img")
	dom.SetAttribute(badImage, "style", "width: 100px; height: 100px; display: block;")

	dom.AppendChild(fig, goodImage)
	dom.AppendChild(root, fig)
	dom.AppendChild(root, badImage)

	imgScorer := scorer.NewImageHasFigureScorer(50)

	assert.True(t, imgScorer.GetImageScore(goodImage) > 0)
	assert.Equal(t, 0, imgScorer.GetImageScore(badImage))
	assert.Equal(t, 0, imgScorer.GetImageScore(nil))
}
