// ORIGINAL: java/extractors/embeds/ImageExtractor.java

package extractor

import (
	nurl "net/url"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// ImageExtractor treats images as another type of embed and provides heuristics for
// lead image candidacy.
type ImageExtractor struct {
	PageURL *nurl.URL
}

func NewImageExtractor(pageURL *nurl.URL) *ImageExtractor {
	return &ImageExtractor{PageURL: pageURL}
}

func (ie *ImageExtractor) RelevantTagNames() []string {
	tagNames := []string{}
	for tagName := range relevantImageTags {
		tagNames = append(tagNames, tagName)
	}
	return tagNames
}

func (ie *ImageExtractor) Extract(node *html.Node) webdoc.Element {
	if node == nil {
		return nil
	}

	nodeTagName := dom.TagName(node)
	if _, exist := relevantImageTags[nodeTagName]; !exist {
		return nil
	}

	if nodeTagName == "figure" {
		img := domutil.GetFirstElementByTagNameInc(node, "picture")
		if img == nil {
			img = domutil.GetFirstElementByTagNameInc(node, "img")
			if img == nil {
				return nil
			}
		}

		figCaption := domutil.GetFirstElementByTagName(node, "figcaption")
		if figCaption == nil {
			figCaption = ie.createFigCaption(node)
		} else {
			links := dom.QuerySelectorAll(figCaption, "a[href]")
			if len(links) == 0 {
				// Here we look for links because some sites put non-caption elements into <figcaption>.
				// For example: image credit could contain a link. So we get the whole DOM structure
				// within <figcaption> only when it contains links, otherwise we get the innerText.
				figCaption = ie.createFigCaption(figCaption)
			}
		}

		imgSrc, width, height := ie.extractImageAttrs(img)

		return &webdoc.Figure{
			Image: webdoc.Image{
				Element:   img,
				SourceURL: imgSrc,
				Width:     width,
				Height:    height,
				PageURL:   ie.PageURL,
			},
			Caption: figCaption,
		}
	}

	if nodeTagName == "span" {
		className := dom.GetAttribute(node, "class")
		if !strings.Contains(className, "lazy-image-placeholder") {
			return nil
		}

		// Image lazy loading on Wikipedia
		img := dom.CreateElement("img")
		dom.SetAttribute(img, "srcset", dom.GetAttribute(node, "data-srcset"))
		imgSrc := dom.GetAttribute(node, "data-src")
		width, _ := strconv.Atoi(dom.GetAttribute(node, "data-width"))
		height, _ := strconv.Atoi(dom.GetAttribute(node, "data-height"))

		return &webdoc.Image{
			Element:   img,
			SourceURL: imgSrc,
			Width:     width,
			Height:    height,
			PageURL:   ie.PageURL,
		}
	}

	// At this point we assume that the node is image element
	imgSrc, width, height := ie.extractImageAttrs(node)
	return &webdoc.Image{
		Element:   node,
		SourceURL: imgSrc,
		Width:     width,
		Height:    height,
		PageURL:   ie.PageURL,
	}
}

// extractImageAttrs will fetch the image source. In original dom-distiller this function
// will also parse width and height of the image. Unfortunately it's not possible here,
// so we only check width and height in attribute.
// NEED-COMPUTE_CSS.
func (ie *ImageExtractor) extractImageAttrs(img *html.Node) (string, int, int) {
	// Try to get lazily-loaded images before falling back to get the src attribute.
	imgSrc := dom.GetAttribute(img, "src")
	for _, attrName := range lazyImageAttrs {
		if attrValue := dom.GetAttribute(img, attrName); attrValue != "" {
			imgSrc = attrValue
			break
		}
	}

	width, _ := strconv.Atoi(dom.GetAttribute(img, "width"))
	height, _ := strconv.Atoi(dom.GetAttribute(img, "height"))
	return imgSrc, width, height
}

func (ie *ImageExtractor) createFigCaption(base *html.Node) *html.Node {
	baseText := domutil.InnerText(base)
	figCaption := dom.CreateElement("figcaption")
	dom.SetTextContent(figCaption, strings.TrimSpace(baseText))
	return figCaption
}
