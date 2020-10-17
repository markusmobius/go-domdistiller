// ORIGINAL: java/PageParamInfo.java

package pagination

const (
	minLinksToJustifyLinearMap = 2
)

// ParamType is types of page parameter values in paging URLs.
type ParamType uint

const (
	Unset      ParamType = iota // Initialized type to indicate empty PageParamInfo.
	PageNumber                  // Value is a page number.
	Unknown                     // None of the above.
)
