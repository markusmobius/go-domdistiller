// ORIGINAL: java/webdocument/ElementAction.java

package webdoc

import (
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"golang.org/x/net/html"
)

const maxClassCount = 2

var (
	rxComment       = regexp.MustCompile(`(?i)\bcomments?\b`)
	rxDisplay       = regexp.MustCompile(`(?i)display:\s*`)
	rxInlineDisplay = regexp.MustCompile(`(?i)display:\s*inline`)
)

type ElementAction struct {
	Flush           bool
	IsAnchor        bool
	ChangesTagLevel bool
	Labels          []string
}

func GetActionForElement(element *html.Node) ElementAction {
	tagName := dom.TagName(element)
	styleAttr := dom.GetAttribute(element, "style")

	// In original dom-distiller, the `flush` and `changesTagLevel` values are decided depending
	// on element display syle. For example, inline element shouldn't change tag level. Unfortunately,
	// this is not possible since we can't compute stylesheet. As fallback, here we simply check if
	// tag is inline by default (see https://developer.mozilla.org/en-US/docs/Web/HTML/Inline_elements).
	// We also check if element's inline style has `display: inline`. NEED-COMPUTE-CSS.
	action := ElementAction{Labels: make([]string, 0)}
	_, isInlineElement := inlineTagNames[tagName]
	hasDisplayStyleAttr := rxDisplay.MatchString(styleAttr)
	hasInlineDisplayStyleAttr := rxInlineDisplay.MatchString(styleAttr)

	switch {
	case isInlineElement && !hasDisplayStyleAttr,
		hasInlineDisplayStyleAttr:

	default:
		action.Flush = true
		action.ChangesTagLevel = true
	}

	if tagName != "html" && tagName != "body" && tagName != "article" {
		id := dom.GetAttribute(element, "id")
		className := dom.GetAttribute(element, "class")
		classCount := len(strings.Fields(className))
		if (rxComment.MatchString(id) || rxComment.MatchString(className)) && classCount <= maxClassCount {
			action.Labels = append(action.Labels, label.StrictlyNotContent)
		}

		switch tagName {
		case "aside", "nav":
			action.Labels = append(action.Labels, label.StrictlyNotContent)
		case "li":
			action.Labels = append(action.Labels, label.Li)
		case "h1":
			action.Labels = append(action.Labels, label.H1, label.Heading)
		case "h2":
			action.Labels = append(action.Labels, label.H2, label.Heading)
		case "h3":
			action.Labels = append(action.Labels, label.H3, label.Heading)
		case "h4", "h5", "h6":
			action.Labels = append(action.Labels, label.Heading)
		case "a":
			// TODO: Anchors probably shouldn't unconditionally change the tag level.
			action.ChangesTagLevel = true
			action.IsAnchor = dom.HasAttribute(element, "href")
		}
	}

	return action
}
