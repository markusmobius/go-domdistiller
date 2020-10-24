// ORIGINAL: java/LogUtil.java

package logutil

// Logger is the base interface for logging process of distiller.
type Logger interface {
	IsLogExtraction() bool
	IsLogVisibility() bool
	IsLogPagination() bool
	IsLogTiming() bool

	PrintExtractionInfo(args ...interface{})
	PrintVisibilityInfo(args ...interface{})
	PrintPaginationInfo(args ...interface{})
	PrintTimingInfo(args ...interface{})
}
