// ORIGINAL: javatest/EmbedExtractorTest.java

package embed_test

import (
	nurl "net/url"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/extractor/embed"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/stretchr/testify/assert"
)

const imageBase64 = "data:image/png;base64,iVBORw0KGgo" +
	"AAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/" +
	"w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_Embed_Image_WithSettingDimension
// - Test_Embed_Image_WithHeightCSS
// - Test_Embed_Image_WithWidthHeightCSSPx
// - Test_Embed_Image_WithWidthAttributeHeightCSSPx
// - Test_Embed_Image_WithWidthAttributeHeightCSS
// - Test_Embed_Image_WithAttributeCSS
// - Test_Embed_Image_WithAttributesCSSHeightCMAndWidthAttrb
// - Test_Embed_Image_WithAttributesCSSHeightCM

func Test_Embed_Image_HasWidthHeightAttributes(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", imageBase64)
	dom.SetAttribute(img, "width", "32")
	dom.SetAttribute(img, "height", "32")

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(img)).(*webdoc.Image)

	assert.NotNil(t, result)
	assert.Equal(t, 32, result.Width)
	assert.Equal(t, 32, result.Height)
}

func Test_Embed_Image_HasNoAttributes(t *testing.T) {
	img := dom.CreateElement("img")

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(img)).(*webdoc.Image)

	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Width)
	assert.Equal(t, 0, result.Height)
}

func Test_Embed_Image_HasWidthAttribute(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "src", imageBase64)
	dom.SetAttribute(img, "width", "32")

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(img)).(*webdoc.Image)

	assert.NotNil(t, result)
	assert.Equal(t, 32, result.Width)
	assert.Equal(t, 0, result.Height)
}

func Test_Embed_Image_LazyLoadedImage(t *testing.T) {
	extractLazyLoadedImage(t, "data-src")
	extractLazyLoadedImage(t, "datasrc")
	extractLazyLoadedImage(t, "data-original")
	extractLazyLoadedImage(t, "data-url")

	extractLazyLoadedFigure(t, "data-src")
	extractLazyLoadedFigure(t, "datasrc")
	extractLazyLoadedFigure(t, "data-original")
	extractLazyLoadedFigure(t, "data-url")
}

func Test_Embed_Image_FigureWithoutCaptionWithNoscript(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	noscript := dom.CreateElement("noscript")
	dom.SetInnerHTML(noscript, "<span>text</span>")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, noscript)

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	// In original dom-distiller the text is not included because by default
	// noscript is hidden element so inner text wouldn't capture it. However
	// we works in server side so we will catch it as well.
	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>text</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, 100, result.Width)
	assert.Equal(t, 100, result.Height)
	assert.Equal(t, expected, result.GenerateOutput(false))
}

func Test_Embed_Image_FigureWithoutImageAndCaption(t *testing.T) {
	figure := dom.CreateElement("figure")

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	assert.Nil(t, result)
}

func Test_Embed_Image_FigureWithPictureWithoutImg(t *testing.T) {
	source := dom.CreateElement("source")
	dom.SetAttribute(source, "srcset", "http://www.example.com/image-240-200.jpg")
	dom.SetAttribute(source, "media", "(min-width: 800px)")

	picture := dom.CreateElement("picture")
	dom.AppendChild(picture, source)

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, picture)

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	expected := `<figure><picture>` +
		`<source srcset="http://www.example.com/image-240-200.jpg" media="(min-width: 800px)"/>` +
		`<img srcset="http://www.example.com/image-240-200.jpg"/>` +
		`</picture></figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
}

func Test_Embed_Image_FigureWithPictureWithoutSourceAndImg(t *testing.T) {
	picture := dom.CreateElement("picture")
	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, picture)

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)
	expected := `<figure><picture></picture></figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
}

func Test_Embed_Image_FigureCaptionTextOnly(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	figcaption := dom.CreateElement("figcaption")
	dom.SetTextContent(figcaption, "This is a caption")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, figcaption)

	extractor := embed.NewImageExtractor(nil)
	result := extractor.Extract(figure)

	assert.NotNil(t, result)
	assert.Equal(t, "This is a caption", result.GenerateOutput(true))
}

func Test_Embed_Image_FigureCaptionWithAnchor(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	anchor := dom.CreateElement("a")
	dom.SetAttribute(anchor, "href", "link_page.html")
	dom.SetInnerHTML(anchor, "caption<br>link")

	figcaption := dom.CreateElement("figcaption")
	dom.AppendChild(figcaption, dom.CreateTextNode("This is a "))
	dom.AppendChild(figcaption, anchor)

	figure := dom.CreateElement("figure")
	dom.SetAttribute(figure, "attribute", "value")
	dom.SetAttribute(figure, "class", "test")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, figcaption)

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL)
	result := extractor.Extract(figure)

	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>` +
		`This is a <a href="http://example.com/link_page.html">caption<br/>link</a>` +
		`</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, "This is a caption\nlink", result.GenerateOutput(true))
}

func Test_Embed_Image_FigureCaptionWithoutAnchor(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	figcaption := dom.CreateElement("figcaption")
	dom.SetInnerHTML(figcaption, "<div><span>This is a caption</span><a></a></div>")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, figcaption)

	extractor := embed.NewImageExtractor(nil)
	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)

	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>This is a caption</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, 100, result.Width)
	assert.Equal(t, 100, result.Height)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, "This is a caption", result.GenerateOutput(true))
}

func Test_Embed_Image_FigureDivCaption(t *testing.T) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, "width", "100")
	dom.SetAttribute(img, "height", "100")
	dom.SetAttribute(img, "src", "http://wwww.example.com/image.jpeg")

	div := dom.CreateElement("div")
	dom.SetTextContent(div, "This is a caption")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)
	dom.AppendChild(figure, div)

	extractor := embed.NewImageExtractor(nil)
	result := extractor.Extract(figure)

	expected := `<figure>` +
		`<img width="100" height="100" src="http://wwww.example.com/image.jpeg"/>` +
		`<figcaption>This is a caption</figcaption>` +
		`</figure>`

	assert.NotNil(t, result)
	assert.Equal(t, expected, result.GenerateOutput(false))
	assert.Equal(t, "This is a caption", result.GenerateOutput(true))
}

func extractLazyLoadedImage(t *testing.T, attr string) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, attr, "image.png")

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL)

	result, _ := (extractor.Extract(img)).(*webdoc.Image)
	assert.NotNil(t, result)
	assert.Equal(t, `<img src="http://example.com/image.png"/>`, result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image.png"}, result.GetURLs())
}

func extractLazyLoadedFigure(t *testing.T, attr string) {
	img := dom.CreateElement("img")
	dom.SetAttribute(img, attr, "image.png")

	figure := dom.CreateElement("figure")
	dom.AppendChild(figure, img)

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := embed.NewImageExtractor(pageURL)

	result, _ := (extractor.Extract(figure)).(*webdoc.Figure)
	assert.NotNil(t, result)
	assert.Equal(t, `<figure><img src="http://example.com/image.png"/></figure>`, result.GenerateOutput(false))
	assert.Equal(t, []string{"http://example.com/image.png"}, result.GetURLs())
}
