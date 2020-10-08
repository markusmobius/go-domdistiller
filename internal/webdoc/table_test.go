// ORIGINAL: javatest/webdocument/WebTableTest.java

package webdoc_test

import (
	nurl "net/url"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

func Test_WebDoc_Table_GenerateOutput(t *testing.T) {
	html := `<table><tbody>` +
		`<tr>` +
		`<td>row1col1</td>` +
		`<td><img src="http://example.com/table.png"/></td>` +
		`<td><picture><img/></picture></td>` +
		`</tr>` +
		`</tbody></table>`

	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, html)

	table := dom.QuerySelector(div, "table")
	webTable := webdoc.Table{TableElement: table}

	// Output should be the same as the input in this case.
	got := webTable.GenerateOutput(false)
	assert.Equal(t, html, testutil.RemoveAllDirAttributes(got))

	// Test GetImageURLs as well.
	imgURLs := webTable.GetImageURLs()
	assert.Equal(t, 1, len(imgURLs))
	assert.Equal(t, "http://example.com/table.png", imgURLs[0])
}

func Test_WebDoc_Table_GetImageURLs(t *testing.T) {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, `
	<table>
	<tbody>
		<tr>
			<td>
				<img src="http://example.com/table.png" srcset="image100 100w, //example.org/image300 300w"/>
			</td>
			<td>
				<picture>
					<source srcset="image200 200w, //example.org/image400 400w"/>
					<img/>
				</picture>
			</td>
		</tr>
	</tbody>
	</table>`)

	table := dom.QuerySelector(div, "table")
	baseURL, _ := nurl.ParseRequestURI("http://example.com/")
	webTable := webdoc.Table{TableElement: table, PageURL: baseURL}

	urls := webTable.GetImageURLs()
	assert.Equal(t, 5, len(urls))
	assert.Equal(t, "http://example.com/table.png", urls[0])
	assert.Equal(t, "http://example.com/image100", urls[1])
	assert.Equal(t, "http://example.org/image300", urls[2])
	assert.Equal(t, "http://example.com/image200", urls[3])
	assert.Equal(t, "http://example.org/image400", urls[4])
}
