// ORIGINAL: java/filters/heuristics/KeepLargestBlockFilter.java

package heuristic

import (
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// KeepLargestBlock keeps the largest TextBlock only (by the number of words). In case of
// more than one block with the same number of words, the first block is chosen. All
// discarded blocks are marked "not content" and flagged as `label.MightBeContent`. Note
// that, by default, only TextBlocks marked as "content" are taken into consideration.
type KeepLargestBlock struct {
	expandToSiblings bool
}

func NewKeepLargestBlock(expandToSiblings bool) *KeepLargestBlock {
	return &KeepLargestBlock{
		expandToSiblings: expandToSiblings,
	}
}

func (f *KeepLargestBlock) Process(doc *webdoc.TextDocument) bool {
	textBlocks := doc.TextBlocks
	if len(textBlocks) < 2 {
		return false
	}

	maxNumWords := -1
	largestBlockIndex := -1
	var largestBlock *webdoc.TextBlock

	for i, tb := range textBlocks {
		if tb.IsContent() {
			if tb.NumWords > maxNumWords {
				largestBlock = tb
				maxNumWords = tb.NumWords
				largestBlockIndex = i
			}
		}
	}

	for _, tb := range textBlocks {
		if tb == largestBlock {
			tb.SetIsContent(true)
			tb.AddLabels(label.VeryLikelyContent)
		} else {
			tb.SetIsContent(false)
			tb.AddLabels(label.MightBeContent)
		}
	}

	if f.expandToSiblings && largestBlockIndex != -1 {
		f.maybeExpandContentToEarlierTextBlocks(textBlocks, largestBlock, largestBlockIndex)
		f.maybeExpandContentToLaterTextBlocks(textBlocks, largestBlock, largestBlockIndex)
	}

	return true
}

func (f *KeepLargestBlock) maybeExpandContentToEarlierTextBlocks(textBlocks []*webdoc.TextBlock, largestBlock *webdoc.TextBlock, largestBlockIndex int) {
	firstTextElement := domutil.GetParentElement(largestBlock.FirstNonWhitespaceTextNode())
	for i := largestBlockIndex - 1; i >= 0; i-- {
		candidate := textBlocks[i]
		candidateLastTextElement := domutil.GetParentElement(candidate.LastNonWhitespaceTextNode())
		if f.isSibling(firstTextElement, candidateLastTextElement) {
			candidate.SetIsContent(true)
			candidate.AddLabels(label.SiblingOfMainContent)
			firstTextElement = domutil.GetParentElement(candidate.FirstNonWhitespaceTextNode())
		}
	}
}

func (f *KeepLargestBlock) maybeExpandContentToLaterTextBlocks(textBlocks []*webdoc.TextBlock, largestBlock *webdoc.TextBlock, largestBlockIndex int) {
	lastTextElement := domutil.GetParentElement(largestBlock.LastNonWhitespaceTextNode())
	for i := largestBlockIndex + 1; i < len(textBlocks); i++ {
		candidate := textBlocks[i]
		candidateFirstTextElement := domutil.GetParentElement(candidate.FirstNonWhitespaceTextNode())
		if f.isSibling(lastTextElement, candidateFirstTextElement) {
			candidate.SetIsContent(true)
			candidate.AddLabels(label.SiblingOfMainContent)
			lastTextElement = domutil.GetParentElement(candidate.LastNonWhitespaceTextNode())
		}
	}
}

func (f *KeepLargestBlock) isSibling(e1, e2 *html.Node) bool {
	return domutil.GetParentElement(e1) == domutil.GetParentElement(e2)
}
