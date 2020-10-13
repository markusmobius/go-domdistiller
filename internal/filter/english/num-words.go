// ORIGINAL: java/filters/english/NumWordsRulesClassifier.java

package english

import "github.com/markusmobius/go-domdistiller/internal/webdoc"

// NumWordsRulesClassifier classifies several TextBlock as content or not-content through
// rules that have been determined using the C4.8 machine learning algorithm, as described
// in the paper "Boilerplate Detection using Shallow Text Features" (WSDM 2010), particularly
// using number of words per block and link density per block.
type NumWordsRulesClassifier struct{}

func NewNumWordsRulesClassifier() *NumWordsRulesClassifier {
	return &NumWordsRulesClassifier{}
}

func (f *NumWordsRulesClassifier) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) == 0 {
		return false
	}

	hasChanges := false
	for i, block := range textBlocks {
		var prevBlock, nextBlock *webdoc.TextBlock
		if i > 0 {
			prevBlock = textBlocks[i-1]
		}
		if i < len(textBlocks)-1 {
			nextBlock = textBlocks[i+1]
		}

		hasChanges = hasChanges || f.classify(prevBlock, block, nextBlock)
	}

	return hasChanges
}

func (f *NumWordsRulesClassifier) classify(prev, current, next *webdoc.TextBlock) bool {
	isContent := false

	if current.LinkDensity < 0.333333 {
		if prev == nil || prev.LinkDensity <= 0.555556 {
			if current.NumWords <= 16 {
				if next == nil || next.NumWords <= 15 {
					if prev == nil || prev.NumWords <= 4 {
						isContent = false
					} else {
						isContent = true
					}
				} else {
					isContent = true
				}
			} else {
				isContent = true
			}
		} else {
			if current.NumWords <= 40 {
				if next == nil || next.NumWords <= 17 {
					isContent = false
				} else {
					isContent = true
				}
			} else {
				isContent = true
			}
		}
	} else {
		isContent = false
	}

	return current.SetIsContent(isContent)
}
