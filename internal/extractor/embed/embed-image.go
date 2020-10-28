// ORIGINAL: java/extractors/embeds/ImageExtractor.java

package embed

import (
	nurl "net/url"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// ImageExtractor treats images as another type of embed and provides heuristics for
// lead image candidacy.
type ImageExtractor struct {
	PageURL *nurl.URL
	logger  logutil.Logger
}

func NewImageExtractor(pageURL *nurl.URL, logger logutil.Logger) *ImageExtractor {
	return &ImageExtractor{
		PageURL: pageURL,
		logger:  logger,
	}
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
		// Find the real image inside the figure. Some sites put their real image
		// (instead of the lazy ones) inside <noscript> tags, so that take precedence.
		var image *html.Node
		noscript := dom.QuerySelector(node, "noscript")

		for _, imageTagName := range []string{"picture", "img"} {
			if noscript != nil {
				// Noscript is a bit weird in Go.
				// Sometimes its content is treated as *html.Node, so QuerySelector will works.
				image = dom.QuerySelector(noscript, "img")

				// Other times, it's treated as plain text content so we need to parse it first.
				if image == nil {
					tmp := dom.CreateElement("div")
					dom.SetInnerHTML(tmp, dom.TextContent(noscript))
					image = dom.QuerySelector(tmp, imageTagName)
				}

				// If there is image inside noscript, put it outside into the figure
				if image != nil {
					dom.PrependChild(node, image)
				}
			}

			if image == nil {
				image = dom.QuerySelector(node, imageTagName)
			}

			if image != nil {
				break
			}
		}

		if image == nil {
			return nil
		}

		// Sometimes there are sites that use <picture> without any <img> inside it.
		// For these cases, we use one of the <source> as <img>.
		if dom.TagName(image) == "picture" {
			sources := dom.GetElementsByTagName(image, "source")
			imgElements := dom.GetElementsByTagName(image, "img")
			if len(imgElements) == 0 && len(sources) > 0 {
				srcset := dom.GetAttribute(sources[0], "srcset")

				img := dom.CreateElement("img")
				dom.SetAttribute(img, "srcset", srcset)
				dom.AppendChild(image, img)
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

		ie.replaceLazyAttr(image)
		width, height := ie.extractImageAttrs(image)

		return &webdoc.Figure{
			Image: webdoc.Image{
				Element: image,
				Width:   width,
				Height:  height,
				PageURL: ie.PageURL,
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
		dom.SetAttribute(img, "src", dom.GetAttribute(node, "data-src"))
		dom.SetAttribute(img, "srcset", dom.GetAttribute(node, "data-srcset"))

		width, _ := strconv.Atoi(dom.GetAttribute(node, "data-width"))
		height, _ := strconv.Atoi(dom.GetAttribute(node, "data-height"))

		return &webdoc.Image{
			Element: img,
			Width:   width,
			Height:  height,
			PageURL: ie.PageURL,
		}
	}

	// At this point we assume that the node is image element
	ie.replaceLazyAttr(node)
	width, height := ie.extractImageAttrs(node)
	return &webdoc.Image{
		Element: node,
		Width:   width,
		Height:  height,
		PageURL: ie.PageURL,
	}
}

// extractImageAttrs will fetch the image source. In original dom-distiller this function
// will also parse width and height of the image. Unfortunately it's not possible here,
// so we only check width and height in attribute.
// NEED-COMPUTE_CSS.
func (ie *ImageExtractor) extractImageAttrs(img *html.Node) (int, int) {
	width, _ := strconv.Atoi(dom.GetAttribute(img, "width"))
	height, _ := strconv.Atoi(dom.GetAttribute(img, "height"))
	return width, height
}

func (ie *ImageExtractor) replaceLazyAttr(img *html.Node) {
	ie.replaceLazySrcAttr(img)
	if dom.GetAttribute(img, "src") == "" {
		ie.replaceLazySrcsetAttr(img)
	}
}

func (ie *ImageExtractor) replaceLazySrcAttr(img *html.Node) {
	// In some sites (e.g. Kotaku), they put 1px square image as data uri in the src attribute.
	// So, here we check if the data uri is too short, just might as well remove it.
	imgSrc := dom.GetAttribute(img, "src")
	if imgSrc != "" && !ie.imageSrcIsValid(imgSrc) {
		dom.RemoveAttribute(img, "src")
		imgSrc = ""
	}

	// Try to get the common lazily-loaded src attrs first.
	for _, attrName := range lazyImageSrcAttrs {
		if attrValue := dom.GetAttribute(img, attrName); attrValue != "" {
			imgSrc = attrValue
			break
		}
	}

	// If the image source still not found, it's possible that they don't put image source in
	// common lazy attributes, so we look at all attributes to find attribute value that looks
	// like an image source.
	if imgSrc == "" {
		for _, attr := range img.Attr {
			if rxLazyImageSrc.MatchString(attr.Val) {
				imgSrc = attr.Val
				break
			}
		}
	}

	if imgSrc != "" {
		ie.printLog("Extracted image src:", imgSrc)
		dom.SetAttribute(img, "src", imgSrc)
	}
}

func (ie *ImageExtractor) replaceLazySrcsetAttr(img *html.Node) {
	// Try to get the common lazily-loaded srcset attrs first.
	imgSrcset := dom.GetAttribute(img, "srcset")
	for _, attrName := range lazyImageSrcsetAttrs {
		if attrValue := dom.GetAttribute(img, attrName); attrValue != "" {
			imgSrcset = attrValue
			break
		}
	}

	// If the srcset still not found, it's possible that they don't put it in common
	// lazy attributes, so we look at all attributes to find attribute value that
	// looks like an image srcset.
	if imgSrcset == "" {
		for _, attr := range img.Attr {
			if rxLazyImageSrcset.MatchString(attr.Val) {
				imgSrcset = attr.Val
				break
			}
		}
	}

	if imgSrcset != "" {
		ie.printLog("Extracted image srcset:", imgSrcset)
		dom.SetAttribute(img, "srcset", imgSrcset)
	}
}

// dataURLIsValid checks if the image src doesn't contains small data URL
// image which often used as placeholder.
func (ie *ImageExtractor) imageSrcIsValid(src string) bool {
	// Check if it's base64 encoded image.
	// If it's not, we assume it's valid image.
	parts := rxB64DataURL.FindStringSubmatch(src)
	if len(parts) == 0 {
		return true
	}

	// If it's SVG, we assume it's valid because SVG can have a meaningful
	// image in under 133 bytes.
	if parts[1] == "image/svg+xml" {
		return true
	}

	// If image is less than 100 bytes (or 133B after encoded to base64),
	// it will be too small therefore it's not a valid image.
	b64starts := strings.Index(src, "base64") + 7
	b64length := len(src) - b64starts
	if b64length < 133 {
		return false
	}

	return true
}

func (ie *ImageExtractor) createFigCaption(base *html.Node) *html.Node {
	// In some sites noscript is put inside figure caption (eg Medium).
	// So, before fetching the inner text we need to parse it first.
	tmp := dom.CreateElement("div")
	dom.SetInnerHTML(tmp, domutil.InnerText(base))
	baseText := domutil.InnerText(tmp)

	figCaption := dom.CreateElement("figcaption")
	dom.SetTextContent(figCaption, strings.TrimSpace(baseText))
	return figCaption
}

func (ie *ImageExtractor) printLog(args ...interface{}) {
	if ie.logger != nil {
		ie.logger.PrintVisibilityInfo(args...)
	}
}
