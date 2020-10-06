// ORIGINAL: javatest/TestUtilTest.java

package testutil_test

import (
	"regexp"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
)

var (
	rxCleanWhitespaces = regexp.MustCompile(`(?mi)^\s+`)
	rxNewlines         = regexp.MustCompile(`(?i)\n`)
)

func Test_CreateDivTree(t *testing.T) {
	expectedHTML := `
		<div id="0">
			<div id="1">
				<div id="2">
					<div id="3"></div>
					<div id="4"></div>
				</div>
				<div id="5">
					<div id="6"></div>
					<div id="7"></div>
				</div>
			</div>
			<div id="8">
				<div id="9">
					<div id="10"></div>
					<div id="11"></div>
				</div>
				<div id="12">
					<div id="13"></div>
					<div id="14"></div>
				</div>
			</div>
		</div>`

	expectedHTML = rxCleanWhitespaces.ReplaceAllString(expectedHTML, "")
	expectedHTML = rxNewlines.ReplaceAllString(expectedHTML, "")

	divs := testutil.CreateDivTree()
	assert.Equal(t, expectedHTML, dom.OuterHTML(divs[0]))
}
