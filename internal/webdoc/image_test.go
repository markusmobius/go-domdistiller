// ORIGINAL: javatest/webdocument/WebImageTest.java

package webdoc_test

import (
	nurl "net/url"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_WebDoc_Image_GenerateOutput(t *testing.T) {
	html := `<picture>` +
		`<source srcset="image"/>` +
		`<img dirty-attributes/>` +
		`</picture>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	picture := dom.QuerySelector(div, "picture")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{Element: picture, PageURL: baseURL}

	expected := `<picture><source srcset="http://example.com/image"/><img/></picture>`
	assert.Equal(t, expected, webImage.GenerateOutput(false))
}

func Test_WebDoc_Image_GetSrcList(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", "image")
	dom.SetAttribute(img, "srcset", "image200 200w, image400 400w")

	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{
		Element:   img,
		PageURL:   baseURL,
		SourceURL: dom.GetAttribute(img, "src"),
	}

	urls := webImage.GetURLs()
	assert.Equal(t, 3, len(urls))
	assert.Equal(t, "http://example.com/image", urls[0])
	assert.Equal(t, "http://example.com/image200", urls[1])
	assert.Equal(t, "http://example.com/image400", urls[2])
}

func Test_WebDoc_Image_GetSrcListInPicture(t *testing.T) {
	html := `<picture>` +
		`<source data-srcset="image200 200w, //example.org/image400 400w"/>` +
		`<source srcset="image100 100w, //example.org/image300 300w"/>` +
		`<img/>` +
		`</picture>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	picture := dom.QuerySelector(div, "picture")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webImage := webdoc.Image{Element: picture, PageURL: baseURL}

	urls := webImage.GetURLs()
	assert.Equal(t, 4, len(urls))
	assert.Equal(t, "http://example.com/image200", urls[0])
	assert.Equal(t, "http://example.org/image400", urls[1])
	assert.Equal(t, "http://example.com/image100", urls[2])
	assert.Equal(t, "http://example.org/image300", urls[3])
}
