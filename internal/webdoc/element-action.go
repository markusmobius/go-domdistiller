// ORIGINAL: java/webdocument/ElementAction.java

package webdoc

import (
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"golang.org/x/net/html"
)

const maxClassCount = 2

var rxComment = regexp.MustCompile(`(?i)\bcomments?\b`)

type ElementAction struct {
	Flush           bool
	IsAnchor        bool
	ChangesTagLevel bool
	Labels          []string
}

func GetActionForElement(element *html.Node) ElementAction {
	tagName := dom.TagName(element)

	// NEED-COMPUTE-CSS
	// In original dom-distiller, the `flush` and `changesTagLevel` values are decided depending
	// on element display syle. For example, inline element shouldn't change tag level. Unfortunately,
	// this is not possible since we can't compute stylesheet. As fallback, here we simply use the
	// default display for the tag name
	action := ElementAction{}
	display := getDisplayStyle(element)
	switch display {
	case "inline": // do nothing
	case "inline-block", "inline-flex":
		action.ChangesTagLevel = true
	default:
		action.Flush = true
		action.ChangesTagLevel = true
	}

	// Check if item is inside <li>
	if domutil.HasAncestor(element, "li", "summary") {
		action.Flush = false
		action.ChangesTagLevel = false
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
