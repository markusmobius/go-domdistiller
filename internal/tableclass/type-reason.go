// ORIGINAL: java/TableClassifier.java

package tableclass

type Reason uint

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

func (r Reason) String() string {
	switch r {
	case InsideEditableArea:
		return "InsideEditableArea"
	case RoleTable:
		return "RoleTable"
	case RoleDescendant:
		return "RoleDescendant"
	case Datatable0:
		return "Datatable0"
	case CaptionTheadTfootColgroupColTh:
		return "CaptionTheadTfootColgroupColTh"
	case AbbrHeadersScope:
		return "AbbrHeadersScope"
	case OnlyHasAbbr:
		return "OnlyHasAbbr"
	case More95PercentDocWidth:
		return "More95PercentDocWidth"
	case Summary:
		return "Summary"
	case NestedTable:
		return "NestedTable"
	case LessEq1Row:
		return "LessEq1Row"
	case LessEq1Col:
		return "LessEq1Col"
	case MoreEq5Cols:
		return "MoreEq5Cols"
	case CellsHaveBorder:
		return "CellsHaveBorder"
	case DifferentlyColoredRows:
		return "DifferentlyColoredRows"
	case MoreEq20Rows:
		return "MoreEq20Rows"
	case LessEq10Cells:
		return "LessEq10Cells"
	case EmbedObjectAppletIframe:
		return "EmbedObjectAppletIframe"
	case More90PercentDocHeight:
		return "More90PercentDocHeight"
	case Default:
		return "Default"
	}
	return "Unknown"
}
