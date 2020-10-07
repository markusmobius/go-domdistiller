// ORIGINAL: javatest/TableClassifierTest.java

package tableclass_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/tableclass"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"golang.org/x/net/html"
)

// NEED-COMPUTE-CSS
// There are some unit tests in original dom-distiller that can't be
// implemented because they require to compute the stylesheets :
// - Test_DocumentWidth
// - Test_WideTable
// - Test_BorderAroundCells
// - Test_NoBorderAroundCells
// - Test_DifferentlyColoredRows
// - Test_TallTable

func Test_InputElement(t *testing.T) {
	input := dom.CreateElement("input")
	dom.SetAttribute(input, "type", "text")

	table := createDefaultTableWithTH()
	dom.AppendChild(input, table)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.InsideEditableArea, reason)
}

func Test_ContentEditableAttribute(t *testing.T) {
	div := testutil.CreateDiv(0)
	dom.SetAttribute(div, "contenteditable", "true")

	table := createDefaultTableWithTH()
	dom.AppendChild(div, table)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.InsideEditableArea, reason)
}

func Test_RolePresentation(t *testing.T) {
	table := createDefaultTableWithTH()
	dom.SetAttribute(table, "role", "presentation")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.RoleTable, reason)
}

func Test_RoleGrid(t *testing.T) {
	table := createDefaultTableWithNoTH()
	dom.SetAttribute(table, "role", "grid")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleTable, reason)
}

func Test_RoleGridNested(t *testing.T) {
	table := createDefaultTableWithNoTH()
	nestedTable := createDefaultNestedTableWithNoTH(table)
	dom.SetAttribute(nestedTable, "role", "grid")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.NestedTable, reason)

	tableType, reason = tableclass.Classify(nestedTable)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleTable, reason)
}

func Test_RoleTreeGrid(t *testing.T) {
	table := createDefaultTableWithNoTH()
	dom.SetAttribute(table, "role", "treegrid")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleTable, reason)
}

func Test_RoleGridCell(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setRoleForFirstElement(table, "td", "gridcell")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_RoleGridCellNested(t *testing.T) {
	table := createDefaultTableWithNoTH()
	nestedTable := createDefaultNestedTableWithNoTH(table)
	setRoleForFirstElement(nestedTable, "td", "gridcell")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.NestedTable, reason)

	tableType, reason = tableclass.Classify(nestedTable)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_RoleColumnHeader(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setRoleForFirstElement(table, "td", "columnheader")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_RoleRow(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setRoleForFirstElement(table, "tr", "row")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_RoleRowGroup(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setRoleForFirstElement(table, "tbody", "rowgroup")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_RoleRowHeader(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setRoleForFirstElement(table, "tr", "rowheader")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_RoleLandmark(t *testing.T) {
	// Test landmark role in <table> element.
	table := createDefaultTableWithNoTH()
	dom.SetAttribute(table, "role", "application")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleTable, reason)

	// Test landmark role in table's descendant.
	dom.RemoveAttribute(table, "role")
	tableType, reason = tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq10Cells, reason)

	setRoleForFirstElement(table, "tr", "navigation")
	tableType, reason = tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.RoleDescendant, reason)
}

func Test_DatatableAttribute(t *testing.T) {
	table := createDefaultTableWithTH()
	dom.SetAttribute(table, "datatable", "0")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.Datatable0, reason)
}

func Test_CaptionTag(t *testing.T) {
	table := createDefaultTableWithNoTH()
	caption := dom.CreateElement("caption")
	dom.SetInnerHTML(caption, "Testing Caption")
	dom.PrependChild(table, caption)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_EmptyCaptionTag(t *testing.T) {
	table := createDefaultTableWithNoTH()
	caption := dom.CreateElement("caption")
	dom.PrependChild(table, caption)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq10Cells, reason)
}

func Test_AllWhitespacedCaptionTag(t *testing.T) {
	table := createDefaultTableWithNoTH()
	caption := dom.CreateElement("caption")
	dom.SetInnerHTML(caption, "&nbsp;  &nbsp;")
	dom.PrependChild(table, caption)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq10Cells, reason)
}

func Test_THeadTag(t *testing.T) {
	th1 := dom.CreateElement("th")
	th2 := dom.CreateElement("th")
	dom.SetInnerHTML(th1, "heading 1")
	dom.SetInnerHTML(th2, "heading 2")

	tr := dom.CreateElement("tr")
	dom.AppendChild(tr, th1)
	dom.AppendChild(tr, th2)

	thead := dom.CreateElement("thead")
	dom.AppendChild(thead, tr)

	table := createDefaultTableWithNoTH()
	dom.PrependChild(table, thead)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_TFootTag(t *testing.T) {
	td1 := dom.CreateElement("td")
	td2 := dom.CreateElement("td")
	dom.SetInnerHTML(td1, "total 1")
	dom.SetInnerHTML(td2, "total 2")

	tr := dom.CreateElement("tr")
	dom.AppendChild(tr, td1)
	dom.AppendChild(tr, td2)

	tfoot := dom.CreateElement("tfoot")
	dom.AppendChild(tfoot, tr)

	table := createDefaultTableWithNoTH()
	dom.PrependChild(table, tfoot)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_ColGroupTag(t *testing.T) {
	col1 := dom.CreateElement("col")
	col2 := dom.CreateElement("col")
	dom.SetAttribute(col1, "span", "2")
	dom.SetAttribute(col2, "align", "left")

	colgroup := dom.CreateElement("colgroup")
	dom.AppendChild(colgroup, col1)
	dom.AppendChild(colgroup, col2)

	table := createDefaultTableWithNoTH()
	dom.PrependChild(table, colgroup)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_ColTag(t *testing.T) {
	col := dom.CreateElement("col")
	dom.SetAttribute(col, "span", "2")

	table := createDefaultTableWithNoTH()
	dom.PrependChild(table, col)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_THTag(t *testing.T) {
	table := createDefaultTableWithTH()
	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_THTagNested(t *testing.T) {
	table := createDefaultTableWithNoTH()
	nestedTable := createDefaultNestedTableWithTH(table)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.NestedTable, reason)

	tableType, reason = tableclass.Classify(nestedTable)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.CaptionTheadTfootColgroupColTh, reason)
}

func Test_EmptyTHTag(t *testing.T) {
	table := createTable(`
	<tbody>
		<tr>
			<th>&nbsp;&nbsp;</th>
			<th>  </th>
		</tr>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
		</tr>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq10Cells, reason)
}

func Test_AllWhitespacedTHTag(t *testing.T) {
	table := createTable(`
	<tbody>
		<tr>
			<th>&nbsp;&nbsp;</th>
			<th>  </th>
		</tr>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
		</tr>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq10Cells, reason)
}

func Test_AbbrAttribute(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setAttributeForFirstElement(table, "td", "abbr", "HTML")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.AbbrHeadersScope, reason)
}

func Test_HeadersAttribute(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setAttributeForFirstElement(table, "td", "headers", "heading1")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.AbbrHeadersScope, reason)
}

func Test_ScopeAttribute(t *testing.T) {
	table := createDefaultTableWithNoTH()
	setAttributeForFirstElement(table, "td", "scope", "colgroup")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.AbbrHeadersScope, reason)
}

func Test_SingleAbbrTag(t *testing.T) {
	abbr := dom.CreateElement("abbr")
	dom.SetInnerHTML(abbr, "html")

	td := dom.CreateElement("td")
	dom.AppendChild(td, abbr)

	table := createDefaultTableWithNoTH()
	tr := getFirstElement(table, "tr")
	dom.AppendChild(tr, td)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.OnlyHasAbbr, reason)
}

func Test_SummaryAttribute(t *testing.T) {
	table := createDefaultTableWithNoTH()
	dom.SetAttribute(table, "summary", "Testing summary attribute")

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.Summary, reason)
}

func Test_EmptyTable(t *testing.T) {
	table := createTable(`
	<tbody>
		<p>empty table</p>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq1Row, reason)
}

func Test_1Row(t *testing.T) {
	table := createTable(`
	<tbody>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
		</tr>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq1Row, reason)
}

func Test_1ColInSameCols(t *testing.T) {
	table := createTable(`
	<tbody>
		<tr>
			<td>row1col1</td>
		</tr>
		<tr>
			<td>row2col1</td>
		</tr>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq1Col, reason)
}

func Test_1ColInDifferentCols(t *testing.T) {
	table := createTable(`
	<tbody>
		<tr>
			<td>row1col1</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
		</tr>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.LessEq10Cells, reason)
}

func Test_5Cols(t *testing.T) {
	table := createTable(`
	<tbody>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
			<td>row1col3</td>
			<td>row1col4</td>
			<td>row1col5</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
			<td>row2col3</td>
			<td>row2col4</td>
			<td>row2col5</td>
		</tr>
	</tbody>`)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.MoreEq5Cols, reason)
}

func Test_20Rows(t *testing.T) {
	table := createDefaultTableWithNoTH()
	tbody := dom.QuerySelector(table, "tbody")

	for i := 2; i < 20; i++ {
		td := dom.CreateElement("td")
		dom.SetTextContent(td, fmt.Sprintf("row %d, col%d", i, i))

		tr := dom.CreateElement("tr")
		dom.AppendChild(tr, td)

		dom.AppendChild(tbody, tr)
	}

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Data, tableType)
	assert.Equal(t, tableclass.MoreEq20Rows, reason)
}

func Test_EmbedElement(t *testing.T) {
	table := createBigDefaultTableWithNoTH()
	embed := dom.CreateElement("embed")
	dom.AppendChild(getFirstElement(table, "td"), embed)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.EmbedObjectAppletIframe, reason)
}

func Test_ObjectElement(t *testing.T) {
	table := createBigDefaultTableWithNoTH()
	embed := dom.CreateElement("object")
	dom.AppendChild(getFirstElement(table, "td"), embed)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.EmbedObjectAppletIframe, reason)
}

func Test_AppletElement(t *testing.T) {
	table := createBigDefaultTableWithNoTH()
	embed := dom.CreateElement("applet")
	dom.AppendChild(getFirstElement(table, "td"), embed)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.EmbedObjectAppletIframe, reason)
}

func Test_IframeElement(t *testing.T) {
	table := createBigDefaultTableWithNoTH()
	embed := dom.CreateElement("iframe")
	dom.AppendChild(getFirstElement(table, "td"), embed)

	tableType, reason := tableclass.Classify(table)
	assert.Equal(t, tableclass.Layout, tableType)
	assert.Equal(t, tableclass.EmbedObjectAppletIframe, reason)
}

func createTable(rawHTML string) *html.Node {
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, "<table>"+rawHTML+"</div>")
	return dom.QuerySelector(div, "table")
}

func createDefaultTableWithTH() *html.Node {
	return createTable(`
	<tbody>
		<tr>
			<th>heading1</th>
			<th>heading2</th>
		</tr>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
		</tr>
	</tbody>`)
}

func createDefaultTableWithNoTH() *html.Node {
	return createTable(`
	<tbody>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
		</tr>
	</tbody>`)
}

func createBigDefaultTableWithNoTH() *html.Node {
	return createTable(`
	<tbody>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
			<td>row1col3</td>
			<td>row1col4</td>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
			<td>row2col3</td>
			<td>row2col4</td>
		</tr>
		<tr>
			<td>row3col1</td>
			<td>row3col2</td>
			<td>row3col3</td>
			<td>row3col4</td>
		</tr>
	</tbody>`)
}

func getFirstElement(table *html.Node, tagName string) *html.Node {
	elements := dom.GetElementsByTagName(table, tagName)
	if len(elements) > 0 {
		return elements[0]
	}
	return nil
}

func setAttributeForFirstElement(table *html.Node, tagName, attrName, attrValue string) {
	element := getFirstElement(table, tagName)
	if element != nil {
		dom.SetAttribute(element, attrName, attrValue)
	}
}

func setRoleForFirstElement(table *html.Node, tagName, role string) {
	setAttributeForFirstElement(table, tagName, "role", role)
}

func createNestedTable(parentTable *html.Node, nestedTableHTML string) *html.Node {
	nestedTable := createTable(nestedTableHTML)

	// Insert nested table into 1st row of `parentTable`.
	rows := dom.GetElementsByTagName(parentTable, "tr")
	if len(rows) > 0 {
		dom.AppendChild(rows[0], nestedTable)
	}

	return nestedTable
}

func createDefaultNestedTableWithTH(parentTable *html.Node) *html.Node {
	return createNestedTable(parentTable, `
	<tbody>
		<tr>
			<th>row1col1</th>
			<th>row1col2</th>
		</tr>
		<tr>
			<td>row2col1</td>
			<td>row2col2</td>
		</tr>
	</tbody>`)
}

func createDefaultNestedTableWithNoTH(parentTable *html.Node) *html.Node {
	return createNestedTable(parentTable, `
	<tbody>
		<tr>
			<td>row1col1</td>
			<td>row1col2</td>
		</tr>
	</tbody>`)
}
