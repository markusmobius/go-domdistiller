// ORIGINAL: java/TableClassifier.java

package tableclass

type Type uint

const (
	Data Type = iota
	Layout
)

func (t Type) String() string {
	switch t {
	case Data:
		return "Data"
	case Layout:
		return "Layout"
	}
	return ""
}
