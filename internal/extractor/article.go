// ORIGINAL: java/extractors/ArticleExtractor.java

package extractor

import (
	"github.com/markusmobius/go-domdistiller/internal/filter/english"
	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/filter/simple"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

// extractArticle extracts TextDocument. It is tuned towards news articles.
func extractArticle(doc *webdoc.TextDocument, wc stringutil.WordCounter, candidateTitles []string) bool {
	// Prepare filters
	terminatingBlocksFinder := english.NewTerminatingBlocksFinder()
	documentTitleMatch := heuristic.NewDocumentTitleMatch(wc, candidateTitles...)
	numWordsRulesClassifier := english.NewNumWordsRulesClassifier()
	labelNotContentToBoilerplate := simple.NewLabelToBoilerplate(label.StrictlyNotContent)

	similarSiblingContentExpansion1 := heuristic.NewSimilarSiblingContentExpansion()
	similarSiblingContentExpansion1.AllowCrossHeadings = true
	similarSiblingContentExpansion1.MaxLinkDensity = 0.5
	similarSiblingContentExpansion1.MaxBlockDistance = 10

	similarSiblingContentExpansion2 := heuristic.NewSimilarSiblingContentExpansion()
	similarSiblingContentExpansion2.AllowCrossHeadings = true
	similarSiblingContentExpansion2.AllowMixedTags = true
	similarSiblingContentExpansion2.MaxBlockDistance = 10

	headingFusion := heuristic.NewHeadingFusion()
	blockProximityFusionPre := heuristic.NewBlockProximityFusion(false)
	boilerplateBlockKeepTitle := simple.NewBoilerplateBlock(label.Title)
	blockProximityFusionPost := heuristic.NewBlockProximityFusion(true)
	keepLargestBlockExpandToSibling := heuristic.NewKeepLargestBlock(true)
	expandTitleToContent := heuristic.NewExpandTitleToContent()
	largeBlockSameTagLevel := heuristic.NewLargeBlockSameTagLevelToContent()
	listAtEnd := heuristic.NewListAtEnd()

	// Run filters
	terminatingBlocksFinder.Process(doc)
	documentTitleMatch.Process(doc)
	numWordsRulesClassifier.Process(doc)
	labelNotContentToBoilerplate.Process(doc)
	similarSiblingContentExpansion1.Process(doc)
	similarSiblingContentExpansion2.Process(doc)
	headingFusion.Process(doc)
	blockProximityFusionPre.Process(doc)
	boilerplateBlockKeepTitle.Process(doc)
	blockProximityFusionPost.Process(doc)
	keepLargestBlockExpandToSibling.Process(doc)
	expandTitleToContent.Process(doc)
	largeBlockSameTagLevel.Process(doc)
	listAtEnd.Process(doc)

	return true
}
