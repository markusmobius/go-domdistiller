// ORIGINAL: java/webdocument/WebText.java

package webdoc

import (
	nurl "net/url"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"golang.org/x/net/html"
)

type Text struct {
	BaseElement

	Text           string
	NumWords       int
	NumLinkedWords int
	Labels         map[string]struct{}
	TagLevel       int
	OffsetBlock    int
	GroupNumber    int
	PageURL        *nurl.URL

	textNodes     []*html.Node
	start         int
	end           int
	firstWordNode int
	lastWordNode  int
}

func (t *Text) GenerateOutput(textOnly bool) string {
	if t.HasLabel(label.Title) {
		return ""
	}

	// TODO: Instead of doing this next part, in the future track font size weight
	// and etc. and wrap the nodes in a "p" tag.
	clonedRoot := domutil.TreeClone(t.TextNodes())

	// To keep formatting/structure, at least one parent element should be in the output.
	// This is necessary because many times a WebText is only a single text node.
	if clonedRoot.Type != html.ElementNode {
		parentClone := dom.Clone(t.TextNodes()[0].Parent, false)
		dom.AppendChild(parentClone, clonedRoot)
		clonedRoot = parentClone
	}

	// The body element should not be used in the output.
	if dom.TagName(clonedRoot) == "body" {
		div := dom.CreateElement("div")
		dom.SetInnerHTML(div, dom.InnerHTML(clonedRoot))
		clonedRoot = div
	}

	// Retain parent tags until the root is not an inline element, to make sure the
	// style is display:block.
	var srcRoot *html.Node
	for {
		if _, isInline := inlineTagNames[dom.TagName(clonedRoot)]; !isInline {
			break
		}

		if srcRoot == nil {
			srcRoot = domutil.GetNearestCommonAncestor(t.TextNodes()...)
			if srcRoot.Type != html.ElementNode {
				srcRoot = domutil.GetParentElement(srcRoot)
			}
		}

		srcRoot = domutil.GetParentElement(srcRoot)
		if dom.TagName(srcRoot) == "body" {
			break
		}

		parentClone := dom.Clone(srcRoot, false)
		dom.AppendChild(parentClone, clonedRoot)
		clonedRoot = parentClone
	}

	// Make sure links are absolute and IDs are gone.
	domutil.MakeAllLinksAbsolute(clonedRoot, t.PageURL)
	domutil.StripTargetAttributes(clonedRoot)
	domutil.StripIDs(clonedRoot)
	domutil.StripUnwantedClassNames(clonedRoot)
	domutil.StripFontColorAttributes(clonedRoot)
	domutil.StripStyleAttributes(clonedRoot)
	domutil.StripAllUnsafeAttributes(clonedRoot)
	// TODO: if we allow images in WebText later, add StripImageElements().

	// Since there are tag elements that are being wrapped by a pair of Tags,
	// we only need to get the innerHTML, otherwise these tags would be duplicated.
	if textOnly {
		return domutil.InnerText(clonedRoot)
	}

	if _, canBeNested := nestingTags[dom.TagName(clonedRoot)]; canBeNested {
		return dom.InnerHTML(clonedRoot)
	}

	return dom.OuterHTML(clonedRoot)
}

func (t *Text) AddLabel(s string) {
	if t.Labels == nil {
		t.Labels = make(map[string]struct{})
	}
	t.Labels[s] = struct{}{}
}

func (t *Text) HasLabel(s string) bool {
	_, exist := t.Labels[s]
	return exist
}

func (t *Text) TakeLabels() map[string]struct{} {
	res := t.Labels
	t.Labels = make(map[string]struct{})
	return res
}

func (t Text) FirstNonWhitespaceTextNode() *html.Node {
	return t.textNodes[t.firstWordNode]
}

func (t Text) LastNonWhitespaceTextNode() *html.Node {
	return t.textNodes[t.lastWordNode]
}

func (t Text) TextNodes() []*html.Node {
	return t.textNodes[t.start:t.end]
}
