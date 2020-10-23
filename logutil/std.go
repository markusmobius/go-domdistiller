package logutil

import "github.com/sirupsen/logrus"

var std = &Logger{
	Logger: logrus.New(),
}

func SetFlags(flags Flag) {
	std.SetFlags(flags)
}

func HasFlag(flag Flag) bool {
	return std.HasFlag(flag)
}

func PrintDistillPhaseInfo(args ...interface{}) {
	std.PrintDistillPhaseInfo(args...)
}

func PrintVisibilityInfo(args ...interface{}) {
	std.PrintVisibilityInfo(args...)
}

func PrintPaginationInfo(args ...interface{}) {
	std.PrintPaginationInfo(args...)
}
