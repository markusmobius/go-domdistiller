// ORIGINAL: java/webdocument/WebTable.java

package webdoc

import (
	nurl "net/url"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"golang.org/x/net/html"
)

type Table struct {
	BaseElement

	Element *html.Node
	PageURL *nurl.URL

	cloned *html.Node
}

func (t *Table) ElementType() string {
	return "table"
}

func (t *Table) GenerateOutput(textOnly bool) string {
	if t.cloned == nil {
		t.cloned = domutil.CloneAndProcessTree(t.Element, t.PageURL)
	}

	if textOnly {
		return domutil.InnerText(t.cloned)
	}

	return dom.OuterHTML(t.cloned)
}

// GetImageURLs returns list of source URLs of all image inside the table.
func (t *Table) GetImageURLs() []string {
	if t.cloned == nil {
		t.cloned = domutil.CloneAndProcessTree(t.Element, t.PageURL)
	}

	imgURLs := []string{}
	for _, img := range dom.QuerySelectorAll(t.cloned, "img,source") {
		src := dom.GetAttribute(img, "src")
		if src != "" {
			imgURLs = append(imgURLs, src)
		}

		imgURLs = append(imgURLs, domutil.GetAllSrcSetURLs(img)...)
	}

	return imgURLs
}
