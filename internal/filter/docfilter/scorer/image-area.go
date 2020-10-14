// ORIGINAL: java/webdocument/filters/images/AreaScorer.java

package scorer

import "golang.org/x/net/html"

// ImageAreaScorer uses image area (length*width) as its heuristic.
// Unfortunately to do that we need to compute CSS which is impossible
// in Go, so this scorer do nothing. NEED-COMPUTE-CSS.
type ImageAreaScorer struct {
	maxScore int
	minArea  int
	maxArea  int
}

// NewImageAreaScorer returns and initiates the ImageAreaScorer.
func NewImageAreaScorer(maxScore, minArea, maxArea int) *ImageAreaScorer {
	return &ImageAreaScorer{
		maxScore: maxScore,
		minArea:  minArea,
		maxArea:  maxArea,
	}
}

func (s *ImageAreaScorer) GetImageScore(_ *html.Node) int {
	return 0
}

func (s *ImageAreaScorer) GetMaxScore() int {
	return s.maxScore
}
