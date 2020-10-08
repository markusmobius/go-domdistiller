// ORIGINAL: javatest/DomUtilTest.java

package domutil_test

import (
	nurl "net/url"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
)

func Test_NearestCommonAncestor(t *testing.T) {
	// The tree graph is
	// 1 - 2 - 3
	//		 \ 4 - 5
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	div1 := testutil.CreateDiv(1)
	dom.AppendChild(body, div1)

	div2 := testutil.CreateDiv(2)
	dom.AppendChild(div1, div2)

	currDiv := testutil.CreateDiv(3)
	dom.AppendChild(div2, currDiv)
	finalDiv1 := currDiv

	currDiv = testutil.CreateDiv(4)
	dom.AppendChild(div2, currDiv)
	dom.AppendChild(currDiv, testutil.CreateDiv(5))

	assert.Equal(t, div2, domutil.GetNearestCommonAncestor(finalDiv1, currDiv.FirstChild))
	nodeList := dom.QuerySelectorAll(doc, `[id="3"],[id="5"]`)
	assert.Equal(t, div2, domutil.GetNearestCommonAncestor(nodeList...))
}

func Test_NearestCommonAncestorIsRoot(t *testing.T) {
	// The tree graph is
	// 1 - 2 - 3
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")

	div1 := testutil.CreateDiv(1)
	dom.AppendChild(body, div1)

	div2 := testutil.CreateDiv(2)
	dom.AppendChild(div1, div2)

	div3 := testutil.CreateDiv(3)
	dom.AppendChild(div2, div3)

	assert.Equal(t, div1, domutil.GetNearestCommonAncestor(div1, div3))
	nodeList := dom.QuerySelectorAll(doc, `[id="1"],[id="3"]`)
	assert.Equal(t, div1, domutil.GetNearestCommonAncestor(nodeList...))
}

func Test_MakeAllLinksAbsolute(t *testing.T) {
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")

	html := `<a href="link"></a>` +
		`<img src="image" srcset="image200 200w, image400 400w"/>` +
		`<img src="image2"/>` +
		`<video src="video" poster="poster">` +
		`<source src="source" srcset="s2, s3"/>` +
		`<track src="track"/>` +
		`</video>`

	expected := `<a href="http://example.com/link"></a>` +
		`<img src="http://example.com/image" srcset="http://example.com/image200 200w, http://example.com/image400 400w"/>` +
		`<img src="http://example.com/image2"/>` +
		`<video src="http://example.com/video" poster="http://example.com/poster">` +
		`<source src="http://example.com/source" srcset="http://example.com/s2, http://example.com/s3"/>` +
		`<track src="http://example.com/track"/>` +
		`</video>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, html)

	domutil.MakeAllLinksAbsolute(doc, baseURL)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_StripTableBackgroundColorAttributes(t *testing.T) {
	html := `<table bgcolor="red">` +
		`<tbody>` +
		`<tr bgcolor="red">` +
		`<th bgcolor="red">text</th>` +
		`<th bgcolor="red">text</th>` +
		`</tr><tr bgcolor="red">` +
		`<td bgcolor="red">text</td>` +
		`<td bgcolor="red">text</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	expected := `<table>` +
		`<tbody>` +
		`<tr>` +
		`<th>text</th>` +
		`<th>text</th>` +
		`</tr><tr>` +
		`<td>text</td>` +
		`<td>text</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, html)

	domutil.StripTableBackgroundColorAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_StripStyleAttributes(t *testing.T) {
	html := `<div style="font-weight: folder">text</div>` +
		`<table style="position: absolute">` +
		`<tbody style="font-size: 2">` +
		`<tr style="z-index: 0">` +
		`<th style="top: 0px">text</th>` +
		`<th style="width: 20px">text</th>` +
		`</tr><tr style="left: 0">` +
		`<td style="display: block">text</td>` +
		`<td style="color: #123">text</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	expected := `<div>text</div>` +
		`<table>` +
		`<tbody>` +
		`<tr>` +
		`<th>text</th>` +
		`<th>text</th>` +
		`</tr><tr>` +
		`<td>text</td>` +
		`<td>text</td>` +
		`</tr>` +
		`</tbody>` +
		`</table>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, html)

	// Test if the root element is handled properly.
	for _, child := range dom.Children(body) {
		domutil.StripStyleAttributes(child)
	}
	assert.Equal(t, expected, dom.InnerHTML(body))

	dom.SetInnerHTML(body, html)
	domutil.StripStyleAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_StripUnwantedClassNames(t *testing.T) {
	html := `<br class="iscaptiontxt"/>` +
		`<br id="a"/>` +
		`<br class=""/>` +
		`<div class="tion cap">` +
		`<br class="hascaption"/>` +
		`<br class="not_me"/>` +
		`</div>`

	expected := `<br class="caption"/>` +
		`<br id="a"/>` +
		`<br/>` +
		`<div>` +
		`<br class="caption"/>` +
		`<br/>` +
		`</div>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, html)

	// Test if the root element is handled properly.
	for _, child := range dom.Children(body) {
		domutil.StripUnwantedClassNames(child)
	}
	assert.Equal(t, expected, dom.InnerHTML(body))

	dom.SetInnerHTML(body, html)
	domutil.StripUnwantedClassNames(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_StripAllUnsafeAttributes(t *testing.T) {
	unsafeHTML := `<h1 class='foo' onclick='alert(123);'>Foo</h1>` +
		`<img alt='bar' invalidattr='alert("Stop");'/>` +
		`<div tabIndex=0 onScroll='alert("Unsafe");'>Baz</div>`

	expected := `<div><h1 class="foo">Foo</h1>` +
		`<img alt="bar"/>` +
		`<div tabindex="0">Baz</div></div>`

	div := dom.CreateElement("div")
	dom.SetAttribute(div, "onfocus", "window.location.href = 'bad.com';")
	dom.SetInnerHTML(div, unsafeHTML)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, div)

	domutil.StripAllUnsafeAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}
