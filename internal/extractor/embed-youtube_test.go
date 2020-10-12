// ORIGINAL: javatest/EmbedExtractorTest.java

package extractor_test

import (
	nurl "net/url"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/extractor"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_Extractor_YouTube_Extract(t *testing.T) {
	youtube := dom.CreateElement("iframe")
	dom.SetAttribute(youtube, "src", "//www.youtube.com/embed/M7lc1UVf-VE?autoplay=1&hl=zh_TW")

	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := extractor.NewYouTubeExtractor(pageURL)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "M7lc1UVf-VE", result.ID)
	assert.Equal(t, "1", result.Params["autoplay"])
	assert.Equal(t, "zh_TW", result.Params["hl"])

	// Begin negative test
	notYoutube := dom.CreateElement("iframe")
	dom.SetAttribute(notYoutube, "src", "http://www.notyoutube.com/embed/M7lc1UVf-VE?autoplay=1")

	result, _ = (extractor.Extract(notYoutube)).(*webdoc.Embed)
	assert.Nil(t, result)
}

func Test_Extractor_YouTube_ExtractID(t *testing.T) {
	youtube := dom.CreateElement("iframe")
	dom.SetAttribute(youtube, "src", "http://www.youtube.com/embed/M7lc1UVf-VE///?autoplay=1")

	extractor := extractor.NewYouTubeExtractor(nil)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "M7lc1UVf-VE", result.ID)

	// Begin negative test
	notYoutube := dom.CreateElement("iframe")
	dom.SetAttribute(notYoutube, "src", "http://www.youtube.com/embed")

	result, _ = (extractor.Extract(notYoutube)).(*webdoc.Embed)
	assert.Nil(t, result)
}

func Test_Extractor_YouTube_Object(t *testing.T) {
	html := `<object>` +
		`<param name="movie" ` +
		`value="//www.youtube.com/v/DML2WUhn2M0&hl=en_US&fs=1&hd=1">` +
		`</param>` +
		`<param name="allowFullScreen" value="true">` +
		`</param>` +
		`</object>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	youtube := dom.FirstElementChild(div)
	pageURL, _ := nurl.ParseRequestURI("http://example.com")
	extractor := extractor.NewYouTubeExtractor(pageURL)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "DML2WUhn2M0", result.ID)
	assert.Equal(t, "en_US", result.Params["hl"])
	assert.Equal(t, "1", result.Params["fs"])
	assert.Equal(t, "1", result.Params["hd"])
}

func Test_Extractor_YouTube_Object2(t *testing.T) {
	html := `<object type="application/x-shockwave-flash" ` +
		`data="http://www.youtube.com/v/ZuNNhOEzJGA&hl=fr&fs=1&rel=0&color1=0x006699&color2=0x54abd6&border=1">` +
		`</object>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	youtube := dom.FirstElementChild(div)
	extractor := extractor.NewYouTubeExtractor(nil)
	result, _ := (extractor.Extract(youtube)).(*webdoc.Embed)

	// Check YouTube specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "youtube", result.Type)
	assert.Equal(t, "ZuNNhOEzJGA", result.ID)
	assert.Equal(t, "fr", result.Params["hl"])
	assert.Equal(t, "1", result.Params["fs"])
	assert.Equal(t, "0", result.Params["rel"])
	assert.Equal(t, "0x006699", result.Params["color1"])
	assert.Equal(t, "0x54abd6", result.Params["color2"])
	assert.Equal(t, "1", result.Params["border"])
}
