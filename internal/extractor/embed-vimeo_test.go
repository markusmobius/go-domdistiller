// ORIGINAL: javatest/EmbedExtractorTest.java

package extractor_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/extractor"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_Extractor_Vimeo_Extract(t *testing.T) {
	vimeo := dom.CreateElement("iframe")
	dom.SetAttribute(vimeo, "src", "http://player.vimeo.com/video/12345?portrait=0")

	extractor := extractor.NewVimeoExtractor()
	result, _ := (extractor.Extract(vimeo)).(*webdoc.Embed)

	// Check Vimeo specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "vimeo", result.Type)
	assert.Equal(t, "12345", result.ID)
	assert.Equal(t, "0", result.Params["portrait"])

	// Begin negative test
	wrongDomain := dom.CreateElement("iframe")
	dom.SetAttribute(wrongDomain, "src", "http://vimeo.com/video/09876?portrait=1")

	result, _ = (extractor.Extract(wrongDomain)).(*webdoc.Embed)
	assert.Nil(t, result)
}

func Test_Extractor_Vimeo_ExtractID(t *testing.T) {
	vimeo := dom.CreateElement("iframe")
	dom.SetAttribute(vimeo, "src", "http://player.vimeo.com/video/12345?portrait=0")

	extractor := extractor.NewVimeoExtractor()
	result, _ := (extractor.Extract(vimeo)).(*webdoc.Embed)

	// Check Vimeo specific attributes
	assert.NotNil(t, result)
	assert.Equal(t, "vimeo", result.Type)
	assert.Equal(t, "12345", result.ID)

	// Begin negative test
	wrongDomain := dom.CreateElement("iframe")
	dom.SetAttribute(wrongDomain, "src", "http://player.vimeo.com/video")

	result, _ = (extractor.Extract(wrongDomain)).(*webdoc.Embed)
	assert.Nil(t, result)
}
