// ORIGINAL: java/LogUtil.java

package distiller

import "github.com/sirupsen/logrus"

// LogFlag is enum to specify logging level.
type LogFlag uint

const (
	// If LogEverything is set DistillerLogger will enable all logs.
	LogEverything LogFlag = LogExtraction | LogVisibility | LogPagination | LogTiming

	// If LogExtraction is set DistillerLogger will print info of each process when extracting article.
	LogExtraction LogFlag = 1 << iota

	// If LogVisibility is set DistillerLogger will print info on why an element is visible.
	LogVisibility

	// If LogPagination is set DistillerLogger will print info of pagination process.
	LogPagination

	// If LogTiming is set DistillerLogger will print info of duration of each process when extracting article.
	LogTiming
)

// distillerLogger is the main logger for dom-distiller
type distillerLogger struct {
	*logrus.Logger
	flags LogFlag
}

func newDistillerLogger(flags LogFlag) *distillerLogger {
	return &distillerLogger{
		Logger: logrus.New(),
		flags:  flags,
	}
}

func (l *distillerLogger) IsLogExtraction() bool { return l.hasFlag(LogExtraction) }

func (l *distillerLogger) IsLogVisibility() bool { return l.hasFlag(LogVisibility) }

func (l *distillerLogger) IsLogPagination() bool { return l.hasFlag(LogPagination) }

func (l *distillerLogger) IsLogTiming() bool { return l.hasFlag(LogTiming) }

func (l *distillerLogger) PrintExtractionInfo(args ...interface{}) { l.print(LogExtraction, args...) }

func (l *distillerLogger) PrintVisibilityInfo(args ...interface{}) { l.print(LogVisibility, args...) }

func (l *distillerLogger) PrintPaginationInfo(args ...interface{}) { l.print(LogPagination, args...) }

func (l *distillerLogger) PrintTimingInfo(args ...interface{}) { l.print(LogTiming, args...) }

func (l *distillerLogger) hasFlag(flag LogFlag) bool {
	return l.flags&flag != 0
}

func (l *distillerLogger) print(flag LogFlag, args ...interface{}) {
	if l.hasFlag(flag) {
		l.Println(args...)
	}
}
