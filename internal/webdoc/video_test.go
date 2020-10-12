// ORIGINAL: javatest/webdocument/WebVideoTest.java

package webdoc_test

import (
	nurl "net/url"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_WebDoc_Video_GenerateOutput(t *testing.T) {
	video := dom.CreateElement("video")
	dom.SetAttribute(video, "onfocus", "new XMLHttpRequest();") // should be stripped

	child := dom.CreateElement("source")
	dom.SetAttribute(child, "src", "http://example.com/foo.ogg")
	dom.AppendChild(video, child)

	child = dom.CreateElement("track")
	dom.SetAttribute(child, "src", "http://example.com/foo.vtt")
	dom.SetAttribute(child, "onclick", "alert(1)") // should be stripped
	dom.AppendChild(video, child)

	expected := `<video>` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>`

	webVideo := webdoc.NewVideo(video, nil, 0, 0)
	assert.Equal(t, expected, webVideo.GenerateOutput(false))
}

func Test_WebDoc_Video_GenerateOutputInvalidChildren(t *testing.T) {
	video := dom.CreateElement("video")
	dom.SetAttribute(video, "onfocus", "new XMLHttpRequest();") // should be stripped

	child := dom.CreateElement("source")
	dom.SetAttribute(child, "src", "http://example.com/foo.ogg")
	dom.AppendChild(video, child)

	child = dom.CreateElement("track")
	dom.SetAttribute(child, "src", "http://example.com/foo.vtt")
	dom.SetAttribute(child, "onclick", "alert(1)") // should be stripped
	dom.AppendChild(video, child)

	child = dom.CreateElement("div")
	dom.SetTextContent(child, "We do not use custom error messages!")
	dom.AppendChild(video, child)

	// Output should ignore anything other than "track" and "source" tags.
	expected := `<video>` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>`

	webVideo := webdoc.NewVideo(video, nil, 0, 0)
	assert.Equal(t, expected, webVideo.GenerateOutput(false))
}

func Test_WebDoc_Video_GenerateOutputRelativeURL(t *testing.T) {
	video := dom.CreateElement("video")
	dom.SetAttribute(video, "poster", "bar.png")                // should be stripped
	dom.SetAttribute(video, "onfocus", "new XMLHttpRequest();") // should be stripped

	child := dom.CreateElement("source")
	dom.SetAttribute(child, "src", "foo.ogg")
	dom.AppendChild(video, child)

	child = dom.CreateElement("track")
	dom.SetAttribute(child, "src", "foo.vtt")
	dom.SetAttribute(child, "onclick", "alert(1)") // should be stripped
	dom.AppendChild(video, child)

	expected := `<video poster="http://example.com/bar.png">` +
		`<source src="http://example.com/foo.ogg"/>` +
		`<track src="http://example.com/foo.vtt"/>` +
		`</video>`

	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webVideo := webdoc.NewVideo(video, baseURL, 0, 0)
	assert.Equal(t, expected, webVideo.GenerateOutput(false))
}

func Test_WebDoc_Video_PosterEmpty(t *testing.T) {
	video := dom.CreateElement("video")
	webVideo := webdoc.NewVideo(video, nil, 400, 300)
	assert.Equal(t, "<video></video>", webVideo.GenerateOutput(false))
}
