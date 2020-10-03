// ORIGINAL: javatest/TestUtil.java

package testutil

import (
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

func CreateMetaProperty(property string, content string) *html.Node {
	meta := dom.CreateElement("meta")
	dom.SetAttribute(meta, "property", property)
	dom.SetAttribute(meta, "content", content)
	return meta
}
