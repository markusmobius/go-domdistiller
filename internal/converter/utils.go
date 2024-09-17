package converter

import (
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

var (
	unlikelyRoles = map[string]struct{}{
		"menu":          {},
		"menubar":       {},
		"complementary": {},
		"navigation":    {},
		"alert":         {},
		"alertdialog":   {},
		"dialog":        {},
	}
)

// isElementWithoutContent determines if node is empty
// or only filled with <br> and <hr>.
func isElementWithoutContent(node *html.Node) bool {
	brs := dom.GetElementsByTagName(node, "br")
	hrs := dom.GetElementsByTagName(node, "hr")
	childs := dom.Children(node)

	return node.Type == html.ElementNode &&
		strings.TrimSpace(dom.TextContent(node)) == "" &&
		(len(childs) == 0 || len(childs) == len(brs)+len(hrs))
}

func isByline(node *html.Node, matchString string) bool {
	rel := dom.GetAttribute(node, "rel")
	itemprop := dom.GetAttribute(node, "itemprop")
	nodeText := dom.TextContent(node)
	if (rel == "author" || strings.Contains(itemprop, "author") || re2go.IsByline(matchString)) &&
		isValidByline(nodeText) {
		return true
	}

	return false
}

// isValidByline checks whether the input string could be a byline.
// This verifies that the input is a string, and that the length
// is less than 100 chars.
func isValidByline(byline string) bool {
	byline = strings.TrimSpace(byline)
	nChar := stringutil.CharCount(byline)
	return nChar > 0 && nChar < 100
}
