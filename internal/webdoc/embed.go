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

	EmbedNodes []*html.Node
	ID         string
	Type       string
	Params     map[string]string
}

func (e *Embed) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	embed := dom.CreateElement("div")
	dom.SetAttribute(embed, "class", "embed-placeholder")
	dom.SetAttribute(embed, "data-type", e.Type)
	dom.SetAttribute(embed, "data-id", e.ID)
	return dom.OuterHTML(embed)
}
