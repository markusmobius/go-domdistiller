package converter

import (
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

var (
	rxUnlikelyCandidates   = regexp.MustCompile(`(?i)-ad-|ai2html|banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|footer|gdpr|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote`)
	rxOkMaybeItsACandidate = regexp.MustCompile(`(?i)and|article|body|column|content|main|shadow`)
	rxByline               = regexp.MustCompile(`(?i)byline|author|dateline|writtenby|p-author`)
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
	if (rel == "author" || strings.Contains(itemprop, "author") || rxByline.MatchString(matchString)) &&
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
