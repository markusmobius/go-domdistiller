// ORIGINAL: javatest/TreeCloneBuilderTest.java

package domutil_test

import (
	"regexp"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_DomUtil_TreeClone_FullTreeBuilder(t *testing.T) {
	expectedHTML := `
		<div id="0">
			<div id="1">
				<div id="2">
					<div id="3"></div>
					<div id="4"></div>
				</div>
				<div id="5"></div>
			</div>
			<div id="8">
				<div id="12">
					<div id="14"></div>
				</div>
			</div>
		</div>`

	expectedHTML = regexp.MustCompile(`(?mi)^\s+`).ReplaceAllString(expectedHTML, "")
	expectedHTML = regexp.MustCompile(`(?i)\n`).ReplaceAllString(expectedHTML, "")

	divs := testutil.CreateDivTree()
	leafNodes := []*html.Node{divs[3], divs[4], divs[5], divs[14]}

	root := domutil.TreeClone(leafNodes)
	assert.Equal(t, expectedHTML, dom.OuterHTML(root))
}

func Test_DomUtil_TreeClone_SingleLeafNode(t *testing.T) {
	leafNodes := []*html.Node{dom.CreateTextNode("some content")}

	root := domutil.TreeClone(leafNodes)
	assert.Equal(t, dom.TextContent(leafNodes[0]), dom.TextContent(root))
}

func Test_DomUtil_TreeClone_NoCommonAncestors(t *testing.T) {
	divs := testutil.CreateDivTree()
	leafNodes := []*html.Node{divs[3], dom.CreateElement("div")}

	root := domutil.TreeClone(leafNodes)
	assert.Nil(t, root)
}
