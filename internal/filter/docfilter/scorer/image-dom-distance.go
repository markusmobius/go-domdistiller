// ORIGINAL: java/webdocument/filters/images/DomDistanceScorer.java

package scorer

import (
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

// ImageDomDistanceScorer uses DOM distance as its heuristic.
type ImageDomDistanceScorer struct {
	maxScore         int
	firstContentNode *html.Node
}

// NewImageDomDistanceScorer returns and initiates the ImageDomDistanceScorer.
func NewImageDomDistanceScorer(maxScore int, firstContent *html.Node) *ImageDomDistanceScorer {
	return &ImageDomDistanceScorer{
		maxScore:         maxScore,
		firstContentNode: firstContent,
	}
}

func (s *ImageDomDistanceScorer) GetImageScore(node *html.Node) int {
	var score int
	if node != nil {
		score = s.compute(node)
	}

	if score < s.maxScore {
		return score
	}

	return s.maxScore
}

func (s *ImageDomDistanceScorer) GetMaxScore() int {
	return s.maxScore
}

func (s *ImageDomDistanceScorer) compute(node *html.Node) int {
	if s.firstContentNode == nil {
		return 0
	}

	depthDiff := domutil.GetNodeDepth(s.firstContentNode) -
		domutil.GetNodeDepth(domutil.GetNearestCommonAncestor(s.firstContentNode, node))

	var multiplier float64
	if depthDiff < 4 {
		multiplier = 1
	} else if depthDiff < 6 {
		multiplier = 0.6
	} else if depthDiff < 8 {
		multiplier = 0.2
	}

	return int(float64(s.maxScore) * multiplier)
}
