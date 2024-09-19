// ORIGINAL: java/extractors/embeds/ImageExtractor.java

// Copyright (c) 2020 Markus Mobius
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package embed

import (
	nurl "net/url"
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

	switch nodeTagName {
	case "figure":
		// Find the real image inside the figure. Some sites put their real image
		// (instead of the lazy ones) inside <noscript> tags, so that take precedence.
		image := ie.findRealFigureImage(node)
		if image == nil {
			return nil
		}

		// Sometimes there are sites that use <picture> without any <img> inside it.
		// For these cases, we use one of the <source> as <img>.
		if dom.TagName(image) == "picture" {
			ie.processPicture(image)
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
		return &webdoc.Figure{
			Image: webdoc.Image{
				Element: image,
				PageURL: ie.PageURL,
			},
			Caption: figCaption,
		}

	case "span": // This is for some Wikipedia lazy-loaded images
		className := dom.GetAttribute(node, "class")
		if !strings.Contains(className, "lazy-image-placeholder") {
			return nil
		}

		// Image lazy loading on Wikipedia
		img := dom.CreateElement("img")
		dom.SetAttribute(img, "src", dom.GetAttribute(node, "data-src"))
		dom.SetAttribute(img, "srcset", dom.GetAttribute(node, "data-srcset"))

		return &webdoc.Image{
			Element: img,
			PageURL: ie.PageURL,
		}

	case "picture":
		ie.processPicture(node)
		fallthrough

	default: // At this point we assume that the node is image element
		ie.replaceLazyAttr(node)
		return &webdoc.Image{
			Element: node,
			PageURL: ie.PageURL,
		}
	}
}

func (ie *ImageExtractor) findRealFigureImage(figure *html.Node) *html.Node {
	var image *html.Node
	noscript := dom.QuerySelector(figure, "noscript")

	// Picture is more prioritized than img
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
				dom.PrependChild(figure, image)
			}
		}

		// If image is not found inside noscript, we check the one that directly inside figure.
		if image == nil {
			image = dom.QuerySelector(figure, imageTagName)
		}

		if image != nil {
			return image
		}
	}

	return nil
}

func (ie *ImageExtractor) processPicture(picture *html.Node) {
	// Picture should only contains sources and images
	for _, node := range dom.GetElementsByTagName(picture, "*") {
		tagName := dom.TagName(node)
		if tagName != "img" && tagName != "source" {
			node.Parent.RemoveChild(node)
		}
	}

	// Sometimes there are sites that use <picture> without any <img> inside it.
	// For these cases, we use one of the <source> as <img>.
	imgs := dom.GetElementsByTagName(picture, "img")
	sources := dom.GetElementsByTagName(picture, "source")
	if len(imgs) == 0 && len(sources) > 0 {
		sources[0].Data = "img"
	}
}

func (ie *ImageExtractor) replaceLazyAttr(base *html.Node) {
	nodes := dom.QuerySelectorAll(base, "img,source")
	nodes = append(nodes, base)

	for _, node := range nodes {
		ie.replaceLazySrcAttr(node)
		if dom.GetAttribute(node, "src") == "" {
			ie.replaceLazySrcsetAttr(node)
		}
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

	return b64length >= 133
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
	if ie.logger != nil && !ie.logger.InternallyNil() {
		ie.logger.PrintVisibilityInfo(args...)
	}
}
