// ORIGINAL: java/webdocument/filters/LeadImageFinder.java

package docfilter

import (
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter/scorer"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// In original dom-distiller they specified this value as 26. However since
// we only use two heuristics instead of four I decided to halved it.
const imageMinimumAcceptedScore = 13

// LeadImageFinder is used to identify a lead image for an article (sometimes known as a mast).
// Each candidate image is scored based on several heuristics including:
// - If an ancestor node of the image is the "figure" tag.
// - If the image is close with content.
//
// In original dom-distiller they have two more heuristics which cannot be implemented here
// because it would require us to compute the stylesheet (NEED-COMPUTE-CSS):
// - The ratio of width/height.
// - The area of the image (width * height) relative to its container.
type LeadImageFinder struct{}

func NewLeadImageFinder() *LeadImageFinder {
	return &LeadImageFinder{}
}

func (f *LeadImageFinder) Process(doc *webdoc.Document) bool {
	candidates := []*webdoc.Image{}
	var firstContent, lastContent *webdoc.Text

	// TODO: WebDocument should have a separate list that point to all images
	// in the document.
	for _, e := range doc.Elements {
		wt, isText := e.(*webdoc.Text)
		if !isText || !e.IsContent() {
			continue
		}

		if firstContent == nil {
			firstContent = wt
		}

		lastContent = wt
	}

	if lastContent == nil {
		return false
	}

	for _, e := range doc.Elements {
		// If the element is an image and not already considered content.
		webImage, isWebImage := e.(*webdoc.Image)
		if (isWebImage && e.IsContent()) || e == lastContent {
			// If we hit the last content or a image that is
			// content, then we are no longer searching for a
			// "lead" image.
			break
		} else if isWebImage {
			candidates = append(candidates, webImage)
		}
	}

	return f.findLeadImage(candidates, firstContent)
}

func (f *LeadImageFinder) findLeadImage(candidates []*webdoc.Image, firstContent *webdoc.Text) bool {
	if len(candidates) == 0 {
		return false
	}

	var contentElement *html.Node
	if firstContent != nil {
		contentElement = firstContent.FirstNonWhitespaceTextNode()
	}

	bestScore := 0
	heuristics := f.getLeadHeuristics(contentElement)
	var bestImage *webdoc.Image

	for _, wi := range candidates {
		currentScore := f.getImageScore(wi, heuristics)
		if currentScore > imageMinimumAcceptedScore {
			if bestImage == nil || bestScore < currentScore {
				bestImage = wi
				bestScore = currentScore
			}
		}
	}

	if bestImage == nil {
		return false
	}

	bestImage.SetIsContent(true)
	return true
}

func (f *LeadImageFinder) getImageScore(wi *webdoc.Image, heuristics []scorer.ImageScorer) int {
	if wi == nil || len(heuristics) == 0 {
		return 0
	}

	score := 0
	imgNode := wi.Element
	for _, ir := range heuristics {
		score += ir.GetImageScore(imgNode)
	}

	return score
}

func (f *LeadImageFinder) getLeadHeuristics(firstContent *html.Node) []scorer.ImageScorer {
	return []scorer.ImageScorer{
		// scorer.NewImageAreaScorer(25, 75_000, 200_000),
		// scorer.NewImageRatioScorer(25),
		scorer.NewImageDomDistanceScorer(25, firstContent),
		scorer.NewImageHasFigureScorer(15),
	}
}
