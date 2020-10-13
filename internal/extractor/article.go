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

// ExtractArticle extracts TextDocument. It is tuned towards news articles.
func ExtractArticle(doc *webdoc.TextDocument, wc stringutil.WordCounter, candidateTitles []string) bool {
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
	changed := false

	// Don't store changes for these two
	terminatingBlocksFinder.Process(doc)
	documentTitleMatch.Process(doc)

	changed = changed || numWordsRulesClassifier.Process(doc)
	changed = changed || labelNotContentToBoilerplate.Process(doc)
	changed = changed || similarSiblingContentExpansion1.Process(doc)
	changed = changed || similarSiblingContentExpansion2.Process(doc)
	changed = changed || headingFusion.Process(doc)
	changed = changed || blockProximityFusionPre.Process(doc)
	changed = changed || boilerplateBlockKeepTitle.Process(doc)
	changed = changed || blockProximityFusionPost.Process(doc)
	changed = changed || keepLargestBlockExpandToSibling.Process(doc)
	changed = changed || expandTitleToContent.Process(doc)
	changed = changed || largeBlockSameTagLevel.Process(doc)
	changed = changed || listAtEnd.Process(doc)
	return changed
}
