// ORIGINAL: java/webdocument/filters/images/DimensionsRatioScorer.java

package scorer

import "golang.org/x/net/html"

// ImageRatioScorer uses image ratio (length/width) as its heuristic.
// Unfortunately to do that we need to compute CSS which is impossible
// in Go, so this scorer do nothing. NEED-COMPUTE-CSS.
type ImageRatioScorer struct {
	maxScore int
}

// NewImageRatioScorer returns and initiates the ImageRatioScorer.
func NewImageRatioScorer(maxScore int) *ImageRatioScorer {
	return &ImageRatioScorer{
		maxScore: maxScore,
	}
}

func (s *ImageRatioScorer) GetImageScore(_ *html.Node) int {
	return 0
}

func (s *ImageRatioScorer) GetMaxScore() int {
	return s.maxScore
}
