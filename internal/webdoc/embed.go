// ORIGINAL: java/webdocument/WebEmbed.java

package webdoc

import (
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
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

func (e *Embed) ElementType() string {
	return "embed"
}

func (e *Embed) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	embed := dom.CreateElement("div")
	dom.SetAttribute(embed, "class", "embed-placeholder")
	dom.SetAttribute(embed, "data-type", e.Type)
	dom.SetAttribute(embed, "data-id", e.ID)

	// Radhi:
	// I just realize the embed element never used in original dom-distiller. No wonder
	// Chromium doesn't render any embedded element. To be fair Readability.js doesn't
	// render some embed  as well citing security concerns. In my opinion since dom-
	// distiller usually only used in page that we already visit, the embedded iframe
	// should automatically be trustworthy enough.
	// TODO: Maybe just to be save we should sanitize it.
	tagName := dom.TagName(e.Element)
	if tagName == "blockquote" || tagName == "iframe" {
		domutil.StripAttributes(e.Element)
		dom.AppendChild(embed, e.Element)
	}

	return dom.OuterHTML(embed)
}
