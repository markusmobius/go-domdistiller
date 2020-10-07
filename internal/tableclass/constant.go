// ORIGINAL: java/TableClassifier.java

package tableclass

type Type uint
type Reason uint

const (
	Data Type = iota
	Layout
)

const (
	Unknown Reason = iota
	InsideEditableArea
	RoleTable
	RoleDescendant
	Datatable0
	CaptionTheadTfootColgroupColTh
	AbbrHeadersScope
	OnlyHasAbbr
	More95PercentDocWidth
	Summary
	NestedTable
	LessEq1Row
	LessEq1Col
	MoreEq5Cols
	CellsHaveBorder
	DifferentlyColoredRows
	MoreEq20Rows
	LessEq10Cells
	EmbedObjectAppletIframe
	More90PercentDocHeight
	Default
)

var headerTags = map[string]bool{
	"colgroup": false,
	"col":      false,
	"th":       true,
}

var objectTags = map[string]bool{
	"embed":  false,
	"object": false,
	"applet": false,
	"iframe": false,
}

// ARIA roles for table, see http://www.w3.org/TR/wai-aria/roles#widget_roles_header.
var ariaTableRoles = map[string]struct{}{
	"grid":     {},
	"treegrid": {},
}

// ARIA roles for descendants of table, see :
// - http://www.w3.org/TR/wai-aria/roles#widget_roles_header.
// - http://www.w3.org/TR/wai-aria/roles#document_structure_roles_header.
var ariaTableDescendantRoles = map[string]struct{}{
	"gridcell":     {},
	"columnheader": {},
	"row":          {},
	"rowgroup":     {},
	"rowheader":    {},
}

// ARIA landmark roles, applicable to both table and its descendants
// - http://www.w3.org/TR/wai-aria/roles#landmark_roles_header.
var ariaRoles = map[string]struct{}{
	"application":   {},
	"banner":        {},
	"complementary": {},
	"contentinfo":   {},
	"form":          {},
	"main":          {},
	"navigation":    {},
	"search":        {},
}
