// ORIGINAL: java/webdocument/WebImage.java

package webdoc

import (
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

type Figure struct {
	Image
	Caption *html.Node
}

func (f *Figure) GenerateOutput(textOnly bool) string {
	figCaption := domutil.CloneAndProcessTree(f.Caption, f.PageURL)
	if textOnly {
		return domutil.InnerText(figCaption)
	}

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, f.getProcessedNode())
	if dom.InnerHTML(f.Caption) != "" {
		dom.AppendChild(figure, figCaption)
	}

	return dom.OuterHTML(figure)
}
