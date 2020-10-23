// ORIGINAL: java/LogUtil.java

package logutil

import "github.com/sirupsen/logrus"

type Logger struct {
	*logrus.Logger
	flags Flag
}

func (l *Logger) SetFlags(flags Flag) {
	l.flags = flags
}

func (l *Logger) HasFlag(flag Flag) bool {
	return l.flags&flag != 0
}

func (l *Logger) PrintDistillPhaseInfo(args ...interface{}) {
	if l.HasFlag(DistillPhases) {
		l.Println(args...)
	}
}

func (l *Logger) PrintVisibilityInfo(args ...interface{}) {
	if l.HasFlag(VisibilityInfo) {
		l.Println(args...)
	}
}

func (l *Logger) PrintPaginationInfo(args ...interface{}) {
	if l.HasFlag(PaginationInfo) {
		l.Println(args...)
	}
}
