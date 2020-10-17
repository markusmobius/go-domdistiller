// ORIGINAL: java/PageParamInfo.java and
//           java/MonotonicPageInfosGroups.java

package info

import "fmt"

// PageInfo stores potential pagination info:
// - page number represented as original plain text in document URL
// - if the info is extracted from an anchor, its href.
type PageInfo struct {
	PageNumber int
	URL        string
}

func (pi *PageInfo) String() string {
	return fmt.Sprintf("pg%d: %s", pi.PageNumber, pi.URL)
}

type PageInfoGroup struct {
	List      []*PageInfo
	DeltaSign int
}
