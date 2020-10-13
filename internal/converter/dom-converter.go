// ORIGINAL: java/webdocument/DomConverter.java

package converter

import (
	nurl "net/url"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/extractor/embed"
	"github.com/markusmobius/go-domdistiller/internal/tableclass"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// DomConverter converts a node and its children into a Document.
type DomConverter struct {
	builder         webdoc.DocumentBuilder
	embedExtractors []embed.EmbedExtractor
	embedTagNames   map[string]struct{}
	pageURL         *nurl.URL
}

func NewDomConverter(builder webdoc.DocumentBuilder, pageURL *nurl.URL) *DomConverter {
	extractors := []embed.EmbedExtractor{
		embed.NewImageExtractor(pageURL),
		embed.NewTwitterExtractor(pageURL),
		embed.NewVimeoExtractor(pageURL),
		embed.NewYouTubeExtractor(pageURL),
	}

	embedTagNames := make(map[string]struct{})
	for _, extractor := range extractors {
		for _, tagName := range extractor.RelevantTagNames() {
			embedTagNames[tagName] = struct{}{}
		}
	}

	return &DomConverter{
		builder:         builder,
		embedExtractors: extractors,
		embedTagNames:   embedTagNames,
		pageURL:         pageURL,
	}
}

func (dc *DomConverter) Convert(root *html.Node) {
	domutil.WalkNodes(root, dc.visitNodeHandler, dc.exitNodeHandler)
}

func (dc *DomConverter) visitNodeHandler(node *html.Node) bool {
	switch node.Type {
	case html.TextNode:
		dc.builder.AddTextNode(node)
		return false

	case html.ElementNode:
		return dc.visitElementNodeHandler(node)

	default:
		return false
	}
}

func (dc *DomConverter) exitNodeHandler(node *html.Node) {
	if node.Type == html.ElementNode {
		if tagName := dom.TagName(node); webdoc.CanBeNested(tagName) {
			dc.builder.AddTag(webdoc.NewTag(tagName, webdoc.TagEnd))
		}
	}

	dc.builder.EndNode()
}

func (dc *DomConverter) visitElementNodeHandler(node *html.Node) bool {
	// In original dom-distiller they skip invisible or uninteresting elements.
	// Unfortunately it's impossible here (NEED-COMPUTE-CSS), so we simply
	// assume everything is visible.

	// Node-type specific extractors check for elements they are interested in here.
	// Everything else will be filtered through the switch below.
	tagName := dom.TagName(node)
	if _, isEmbed := dc.embedTagNames[tagName]; isEmbed {
		// If the tag is marked as interesting, check the extractors.
		for _, extractor := range dc.embedExtractors {
			embed := extractor.Extract(node)
			if embed != nil {
				dc.builder.AddEmbed(embed)
				return false
			}
		}
	}

	// Skip social and sharing elements.
	// See crbug.com/692553, crbug.com/696556, and crbug.com/674557
	className := dom.GetAttribute(node, "class")
	component := dom.GetAttribute(node, "data-component")
	if className == "sharing" || className == "socialArea" || component == "share" {
		return false
	}

	// Create a placeholder for the elements we want to preserve.
	if webdoc.CanBeNested(tagName) {
		dc.builder.AddTag(webdoc.NewTag(tagName, webdoc.TagStart))
	}

	switch tagName {
	case "a":
		// The "section" parameter is to differentiate with "redlinks".
		// Ref: https://en.wikipedia.org/wiki/Wikipedia:Red_link
		href := dom.GetAttribute(node, "href")
		if strings.Contains(href, "action=edit&section=") {
			// Skip "edit section" on mediawiki.
			// See crbug.com/647667.
			return false
		}

	case "span":
		if className == "mw-editsection" {
			// Skip "[edit]" on mediawiki desktop version.
			// See crbug.com/647667.
			return false
		}

	case "br":
		dc.builder.AddLineBreak(node)
		return false

	case "table":
		tableType, _ := tableclass.Classify(node)
		if tableType == tableclass.Data {
			dc.builder.AddDataTable(node)
			return false
		}

	case "video":
		dc.builder.AddEmbed(webdoc.NewVideo(node, dc.pageURL, 0, 0))
		return false

	// These element types are all skipped (but may affect document construction).
	case "option", "object", "embed", "applet":
		dc.builder.SkipNode(node)
		return false

	// These types are skipped and don't affect document construction.
	case "head", "style", "script", "link", "noscript", "iframe", "svg":
		return false
	}

	dc.builder.StartNode(node)
	return true
}
