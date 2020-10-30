// ORIGINAL: javatest/DomUtilTest.java

package domutil_test

import (
	nurl "net/url"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_DomUtil_GetOutputNodesWithHiddenChildren

func Test_DomUtil_HasRootDomain(t *testing.T) {
	// Positive tests.
	assert.True(t, domutil.HasRootDomain("//www.foo.bar/foo/bar.html", "foo.bar"))
	assert.True(t, domutil.HasRootDomain("http://www.foo.bar/foo/bar.html", "foo.bar"))
	assert.True(t, domutil.HasRootDomain("https://www.m.foo.bar/foo/bar.html", "foo.bar"))
	assert.True(t, domutil.HasRootDomain("https://www.m.foo.bar/foo/bar.html", "www.m.foo.bar"))
	assert.True(t, domutil.HasRootDomain("http://localhost/foo/bar.html", "localhost"))
	assert.True(t, domutil.HasRootDomain("https://www.m.foo.bar.baz", "foo.bar.baz"))
	// Negative tests.
	assert.False(t, domutil.HasRootDomain("https://www.m.foo.bar.baz", "x.foo.bar.baz"))
	assert.False(t, domutil.HasRootDomain("https://www.foo.bar.baz", "foo.bar"))
	assert.False(t, domutil.HasRootDomain("http://foo", "m.foo"))
	assert.False(t, domutil.HasRootDomain("https://www.badfoobar.baz", "foobar.baz"))
	assert.False(t, domutil.HasRootDomain("", "foo"))
	assert.False(t, domutil.HasRootDomain("http://foo.bar", ""))
}

func Test_DomUtil_NearestCommonAncestor(t *testing.T) {
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

func Test_DomUtil_NearestCommonAncestorIsRoot(t *testing.T) {
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

func Test_DomUtil_MakeAllLinksAbsolute(t *testing.T) {
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

func Test_DomUtil_GetSrcSetURLs(t *testing.T) {
	html := `<img src="http://example.com/image" ` +
		`srcset="http://example.com/image200 200w, http://example.com/image400 400w"/>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	img := dom.QuerySelector(div, "img")
	srcsetURLs := domutil.GetSrcSetURLs(img)
	assert.Equal(t, 2, len(srcsetURLs))
	assert.Equal(t, "http://example.com/image200", srcsetURLs[0])
	assert.Equal(t, "http://example.com/image400", srcsetURLs[1])
}

func Test_DomUtil_GetAllSrcSetURLs(t *testing.T) {
	html := `<picture>` +
		`<source srcset="image200 200w, //example.org/image400 400w"/>` +
		`<source srcset="image100 100w, //example.org/image300 300w"/>` +
		`<img/>` +
		`</picture>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	srcsetURLs := domutil.GetAllSrcSetURLs(div)
	assert.Equal(t, 4, len(srcsetURLs))
	assert.Equal(t, "image200", srcsetURLs[0])
	assert.Equal(t, "//example.org/image400", srcsetURLs[1])
	assert.Equal(t, "image100", srcsetURLs[2])
	assert.Equal(t, "//example.org/image300", srcsetURLs[3])
}

func Test_DomUtil_StripImageElements(t *testing.T) {
	html := `<img id="a" alt="alt" dir="rtl" title="t" style="typo" align="left"` +
		`src="image" class="a" srcset="image200 200w" data-dummy="a"/>` +
		`<img mulformed="nothing" data-empty data-dup="1" data-dup="2" src="image" src="second"/>`

	expected := `<img alt="alt" dir="rtl" title="t" src="image" srcset="image200 200w"/>` +
		`<img src="image" src="second"/>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, html)

	// Test if the root element is handled properly.
	for _, child := range dom.Children(body) {
		domutil.StripAttributes(child)
	}
	assert.Equal(t, expected, dom.InnerHTML(body))

	dom.SetInnerHTML(body, html)
	domutil.StripAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_DomUtil_StripTableBackgroundColorAttributes(t *testing.T) {
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

	domutil.StripAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_DomUtil_StripStyleAttributes(t *testing.T) {
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
		domutil.StripAttributes(child)
	}
	assert.Equal(t, expected, dom.InnerHTML(body))

	dom.SetInnerHTML(body, html)
	domutil.StripAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_DomUtil_StripUnwantedClassNames(t *testing.T) {
	html := `<br class="iscaptiontxt"/>` +
		`<br id="a"/>` +
		`<br class=""/>` +
		`<div class="tion cap">` +
		`<br class="hascaption"/>` +
		`<br class="not_me"/>` +
		`</div>`

	expected := `<br/>` +
		`<br/>` +
		`<br/>` +
		`<div>` +
		`<br/>` +
		`<br/>` +
		`</div>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, html)

	// Test if the root element is handled properly.
	for _, child := range dom.Children(body) {
		domutil.StripAttributes(child)
	}
	assert.Equal(t, expected, dom.InnerHTML(body))

	dom.SetInnerHTML(body, html)
	domutil.StripAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_DomUtil_StripAllUnsafeAttributes(t *testing.T) {
	unsafeHTML := `<h1 class='foo' onclick='alert(123);'>Foo</h1>` +
		`<img alt='bar' invalidattr='alert("Stop");'/>` +
		`<div tabIndex=0 onScroll='alert("Unsafe");'>Baz</div>`

	expected := `<div><h1>Foo</h1>` +
		`<img alt="bar"/>` +
		`<div tabindex="0">Baz</div></div>`

	div := dom.CreateElement("div")
	dom.SetAttribute(div, "onfocus", "window.location.href = 'bad.com';")
	dom.SetInnerHTML(div, unsafeHTML)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, div)

	domutil.StripAttributes(body)
	assert.Equal(t, expected, dom.InnerHTML(body))
}

func Test_DomUtil_GetOutputNodes(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, `<p>`+
		`<span>Some content</span>`+
		`<img src="./image.png"/>`+
		`</p>`)

	// Expected nodes: <div><p><span>#text<img>.
	contentNodes := domutil.GetOutputNodes(div)
	assert.Len(t, contentNodes, 5)

	node := contentNodes[0]
	assert.Equal(t, html.ElementNode, node.Type)
	assert.Equal(t, "div", dom.TagName(node))

	node = contentNodes[1]
	assert.Equal(t, html.ElementNode, node.Type)
	assert.Equal(t, "p", dom.TagName(node))

	node = contentNodes[2]
	assert.Equal(t, html.ElementNode, node.Type)
	assert.Equal(t, "span", dom.TagName(node))

	node = contentNodes[3]
	assert.Equal(t, html.TextNode, node.Type)

	node = contentNodes[4]
	assert.Equal(t, html.ElementNode, node.Type)
	assert.Equal(t, "img", dom.TagName(node))
}

func Test_DomUtil_GetOutputNodesNestedTable(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, `<table>`+
		`<tbody><tr>`+
		`<td><table><tbody><tr><td>nested</td></tr></tbody></table></td>`+
		`<td>outer</td>`+
		`</tr></tbody>`+
		`</table>`)

	table := dom.QuerySelector(div, "table")
	contentNodes := domutil.GetOutputNodes(table)
	assert.Len(t, contentNodes, 11)
}

func Test_DomUtil_NodeDepth(t *testing.T) {
	div := testutil.CreateDiv(1)
	div2 := testutil.CreateDiv(2)
	div3 := testutil.CreateDiv(3)

	dom.AppendChild(div, div2)
	dom.AppendChild(div2, div3)

	assert.Equal(t, 2, domutil.GetNodeDepth(div3))
}

func Test_DomUtil_ZeroOrNoNodeDepth(t *testing.T) {
	div := testutil.CreateDiv(0)
	assert.Equal(t, 0, domutil.GetNodeDepth(div))
	assert.Equal(t, -1, domutil.GetNodeDepth(nil))
}
