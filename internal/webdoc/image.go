// ORIGINAL: java/webdocument/WebImage.java

package webdoc

import (
	nurl "net/url"
	"strconv"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type Image struct {
	BaseElement

	Element   *html.Node // node for the image
	Width     int        // width of image in pixel
	Height    int        // height of image in pixel
	SourceURL string     // source url of the image
	PageURL   *nurl.URL  // url of page where image is placed

	cloned *html.Node
}

func (i *Image) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	if i.cloned == nil {
		i.cloned = i.cloneAndProcessNode()
	}

	return dom.OuterHTML(i.cloned)
}

// GetURLs returns the list of source URLs of this image.
func (i *Image) GetURLs() []string {
	if i.cloned == nil {
		i.cloned = i.cloneAndProcessNode()
	}

	urls := []string{}
	if i.SourceURL != "" {
		urls = append(urls, i.SourceURL)
	}

	urls = append(urls, domutil.GetAllSrcSetURLs(i.cloned)...)
	return urls
}

func (i *Image) cloneAndProcessNode() *html.Node {
	cloned := dom.Clone(i.Element, true)
	img := domutil.GetFirstElementByTagNameInc(cloned, "img")
	if i.SourceURL != "" {
		i.SourceURL = stringutil.CreateAbsoluteURL(i.SourceURL, i.PageURL)
		dom.SetAttribute(img, "src", i.SourceURL)
	}

	if i.Width > 0 && i.Height > 0 {
		dom.SetAttribute(img, "width", strconv.Itoa(i.Width))
		dom.SetAttribute(img, "height", strconv.Itoa(i.Height))
	}
	domutil.StripImageElement(img)

	for _, source := range dom.GetElementsByTagName(cloned, "source") {
		for lazyAttrName, realAttrName := range lazyImageAttrs {
			lazyAttrValue := dom.GetAttribute(source, lazyAttrName)
			if lazyAttrValue != "" {
				dom.SetAttribute(source, realAttrName, lazyAttrValue)
				dom.RemoveAttribute(source, lazyAttrName)
			}
		}
	}

	domutil.MakeAllSrcAttributesAbsolute(cloned, i.PageURL)
	domutil.MakeAllSrcSetAbsolute(cloned, i.PageURL)
	return cloned
}
