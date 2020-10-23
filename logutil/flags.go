// ORIGINAL: java/LogUtil.java

package logutil

// Flag is enum to specify logging level.
type Flag uint

const (
	// If DistillPhases is set logger will print changes of each process when extracting article.
	DistillPhases Flag = 1 << iota

	// If VisibilityInfo is set logger will print info on why an element is visible.
	VisibilityInfo

	// If PaginationInfo is set logger will print info of pagination process.
	PaginationInfo
)
