// ORIGINAL: java/webdocument/filters/images/ImageScorer.java

package scorer

import "golang.org/x/net/html"

// ImageScorer is used to represent a single heuristic used in image extraction.
// The provided image will be given a score based on the heuristic and a max score.
type ImageScorer interface {
	// GetImageScore returns a particular image a score based on the heuristic
	// implemented in this ImageScorer and what the max score is set to.
	GetImageScore(e *html.Node) int

	// GetMaxScore returns the maximum possible score that this ImageScorer can return.
	GetMaxScore() int
}
