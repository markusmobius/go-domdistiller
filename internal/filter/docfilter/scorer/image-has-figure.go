// ORIGINAL: java/webdocument/filters/images/HasFigureScorer.java

package scorer

import (
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

// ImageHasFigureScorer scores based on if the image has a "figure" node as an ancestor.
type ImageHasFigureScorer struct {
	maxScore int
}

// NewImageHasFigureScorer returns and initiates the ImageHasFigureScorer.
func NewImageHasFigureScorer(maxScore int) *ImageHasFigureScorer {
	return &ImageHasFigureScorer{
		maxScore: maxScore,
	}
}

func (s *ImageHasFigureScorer) GetImageScore(node *html.Node) int {
	var score int
	if node != nil {
		score = s.compute(node)
	}

	if score < s.maxScore {
		return score
	}

	return s.maxScore
}

func (s *ImageHasFigureScorer) GetMaxScore() int {
	return s.maxScore
}

func (s *ImageHasFigureScorer) compute(node *html.Node) int {
	parents := domutil.GetParentNodes(node)
	for _, n := range parents {
		if n.Type == html.ElementNode && n.Data == "figure" {
			return s.maxScore
		}
	}
	return 0
}
