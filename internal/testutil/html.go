// ORIGINAL: javatest/TestUtil.java

package testutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

// CreateHTML returns an <html> that consist of empty <head> and <body>.
// This is an additional method and doesn't exist in original Java code.
func CreateHTML() *html.Node {
	rawHTML := `
		<!DOCTYPE html>
		<html lang="en">
			<head></head>
			<body></body>
		</html>`

	root, _ := html.Parse(strings.NewReader(rawHTML))
	return root
}

// CreateDiv creates a div with the integer id as its id.
func CreateDiv(id int) *html.Node {
	div := dom.CreateElement("div")
	dom.SetAttribute(div, "id", strconv.Itoa(id))
	return div
}

func CreateTitle(value string) *html.Node {
	title := dom.CreateElement("title")
	dom.SetInnerHTML(title, value)
	return title
}

func CreateHeading(n int, value string) *html.Node {
	h := dom.CreateElement(fmt.Sprintf("h%d", n))
	dom.SetInnerHTML(h, value)
	return h
}

func CreateMetaProperty(property string, content string) *html.Node {
	meta := dom.CreateElement("meta")
	dom.SetAttribute(meta, "property", property)
	dom.SetAttribute(meta, "content", content)
	return meta
}

func CreateMetaName(name string, content string) *html.Node {
	meta := dom.CreateElement("meta")
	dom.SetAttribute(meta, "name", name)
	dom.SetAttribute(meta, "content", content)
	return meta
}
