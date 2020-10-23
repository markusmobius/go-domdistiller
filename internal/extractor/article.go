// ORIGINAL: java/extractors/ArticleExtractor.java

package extractor

import (
	"github.com/markusmobius/go-domdistiller/internal/filter/english"
	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/filter/simple"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/markusmobius/go-domdistiller/logutil"
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

	printArticleLog(doc, true, "Start")

	// Run filters
	// Intentionally don't print changes from these two
	terminatingBlocksFinder.Process(doc)
	documentTitleMatch.Process(doc)

	changed := numWordsRulesClassifier.Process(doc)
	printArticleLog(doc, changed, "Classification complete")

	changed = labelNotContentToBoilerplate.Process(doc)
	printArticleLog(doc, changed, "Ignore strictly not content blocks")

	changed = similarSiblingContentExpansion1.Process(doc)
	printArticleLog(doc, changed, "Cross headings SimilarSiblingContentExpansion")

	changed = similarSiblingContentExpansion2.Process(doc)
	printArticleLog(doc, changed, "Mixed tags SimilarSiblingContentExpansion")

	changed = headingFusion.Process(doc)
	printArticleLog(doc, changed, "HeadingFusion")

	changed = blockProximityFusionPre.Process(doc)
	printArticleLog(doc, changed, "BlockProximityFusion for distance=1")

	changed = boilerplateBlockKeepTitle.Process(doc)
	printArticleLog(doc, changed, "BlockFilter keep title")

	changed = blockProximityFusionPost.Process(doc)
	printArticleLog(doc, changed, "BlockProximityFusion for same level content-only")

	changed = keepLargestBlockExpandToSibling.Process(doc)
	printArticleLog(doc, changed, "Keep largest block")

	changed = expandTitleToContent.Process(doc)
	printArticleLog(doc, changed, "Expand title to content")

	changed = largeBlockSameTagLevel.Process(doc)
	printArticleLog(doc, changed, "Largest block with same tag level to content")

	changed = listAtEnd.Process(doc)
	printArticleLog(doc, changed, "List at end filter")

	return true
}

func printArticleLog(doc *webdoc.TextDocument, changed bool, header string) {
	logMsg := ""
	if !changed {
		logMsg = header + ": NO CHANGES"
	} else {
		logMsg = header + ":\n" + doc.DebugString()
	}

	logutil.PrintDistillPhaseInfo(logMsg)
}
