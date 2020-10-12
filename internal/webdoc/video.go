// ORIGINAL: java/webdocument/WebVideo.java

package webdoc

import (
	nurl "net/url"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type Video struct {
	BaseElement

	// TODO: Handle multiple nested "source" and "track" tags.
	Element *html.Node
	Width   int
	Height  int
	PageURL *nurl.URL
}

func NewVideo(node *html.Node, pageURL *nurl.URL, width, height int) *Video {
	return &Video{
		Element: node,
		PageURL: pageURL,
		Width:   width,
		Height:  height,
	}
}

func (v *Video) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	vNode := dom.Clone(v.Element, false)
	for _, child := range dom.Children(v.Element) {
		childTag := dom.TagName(child)
		if childTag == "source" || childTag == "track" {
			dom.AppendChild(vNode, dom.Clone(child, false))
		}
	}

	if poster := dom.GetAttribute(vNode, "poster"); poster != "" {
		poster = stringutil.CreateAbsoluteURL(poster, v.PageURL)
		dom.SetAttribute(vNode, "poster", poster)
	}

	domutil.StripIDs(vNode)
	domutil.StripAllUnsafeAttributes(vNode)
	domutil.MakeAllSrcAttributesAbsolute(vNode, v.PageURL)

	return dom.OuterHTML(vNode)
}
