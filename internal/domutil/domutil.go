// ORIGINAL: java/DomUtil.java

package domutil

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

var rxPunctuation = regexp.MustCompile(`\s+([.?!,:;])(\S+)`)

// =================================================================================
// Functions below these point are functions that exist in original Dom-Distiller
// code but that can't be perfectly replicated by this package. This is because
// in original Dom-Distiller they uses GWT which able to compute stylesheet.
// Unfortunately, Go can't do this unless we are using some kind of headless
// browser, so here we only do some kind of workaround (which works but obviously
// not as good as GWT) or simply ignore it.
// =================================================================================

// InnerText in JS and GWT is used to capture text from an element while excluding
// text from hidden children. A child is hidden if it's computed width is 0, whether
// because its CSS (e.g `display: none`, `visibility: hidden`, etc), or if the child
// has `hidden` attribute. Since we can't compute stylesheet, we only look at `hidden`
// attribute here.
//
// Besides excluding text from hidden children, difference between this function and
// `dom.TextContent` is the latter will skip <br> tag while this function will preserve
// <br> as whitespace. NEED-COMPUTE-CSS
func InnerText(node *html.Node) string {
	var buffer bytes.Buffer
	var finder func(*html.Node)

	finder = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			buffer.WriteString(" " + n.Data + " ")

		case html.ElementNode:
			if n.Data == "br" {
				buffer.WriteString("\n")
			} else if dom.HasAttribute(n, "hidden") {
				return
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}

	finder(node)
	text := buffer.String()
	text = strings.Join(strings.Fields(text), " ")
	text = rxPunctuation.ReplaceAllString(text, "$1 $2")
	return text
}

// GetArea in original code returns area of a node by multiplying
// offsetWidth and offsetHeight. Since it's not possible in Go, we
// simply return 0. NEED-COMPUTE-CSS
func GetArea(node *html.Node) int {
	return 0
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for dom management.
// =================================================================================

// SomeNode iterates over a NodeList, return true if any of the
// provided iterate function calls returns true, false otherwise.
func SomeNode(nodeList []*html.Node, fn func(*html.Node) bool) bool {
	for i := 0; i < len(nodeList); i++ {
		if fn(nodeList[i]) {
			return true
		}
	}
	return false
}
