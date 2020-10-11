// ORIGINAL: java/webdocument/WebEmbed.java

package webdoc

import (
	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

// Embed is the base class for many site-specific embedded
// elements (Twitter, YouTube, etc.).
type Embed struct {
	BaseElement

	Element *html.Node
	ID      string
	Type    string
	Params  map[string]string
}

func (e *Embed) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	embed := dom.CreateElement("div")
	dom.SetAttribute(embed, "class", "embed-placeholder")
	dom.SetAttribute(embed, "data-type", e.Type)
	dom.SetAttribute(embed, "data-id", e.ID)

	// TODO:
	// I just realize the embed element never used in original dom-distiller. No wonder
	// Chromium doesn't render any embedded element. To be fair Readability.js doesn't
	// render it as well citing security concerns. In my opinion since dom-distiller
	// usually only used in page that we already visit, the embedded iframe should
	// automatically be trustworthy enough. Just to be save we should sanitize it.
	// For now I'll just allow blockquote for Twitter.
	if dom.TagName(e.Element) == "blockquote" {
		dom.AppendChild(embed, e.Element)
	}

	return dom.OuterHTML(embed)
}
