// ORIGINAL: java/LogUtil.java

// Copyright (c) 2020 Markus Mobius
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package distiller

import (
	"os"

	"github.com/rs/zerolog"
)

// LogFlag is enum to specify logging level.
type LogFlag uint

const (
	// LogNothing will disable the logger.
	LogNothing LogFlag = 0

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
	log   zerolog.Logger
	flags LogFlag
}

func newDistillerLogger(flags LogFlag) *distillerLogger {
	return &distillerLogger{
		log: zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04",
		}).With().Timestamp().Logger(),
		flags: flags,
	}
}

func (l *distillerLogger) InternallyNil() bool { return l == nil }

func (l *distillerLogger) IsLogExtraction() bool { return l.hasFlag(LogExtraction) }

func (l *distillerLogger) IsLogVisibility() bool { return l.hasFlag(LogVisibility) }

func (l *distillerLogger) IsLogPagination() bool { return l.hasFlag(LogPagination) }

func (l *distillerLogger) IsLogTiming() bool { return l.hasFlag(LogTiming) }

func (l *distillerLogger) PrintExtractionInfo(args ...interface{}) { l.print(LogExtraction, args...) }

func (l *distillerLogger) PrintVisibilityInfo(args ...interface{}) { l.print(LogVisibility, args...) }

func (l *distillerLogger) PrintPaginationInfo(args ...interface{}) { l.print(LogPagination, args...) }

func (l *distillerLogger) PrintTimingInfo(args ...interface{}) { l.print(LogTiming, args...) }

func (l *distillerLogger) hasFlag(flag LogFlag) bool {
	if l.InternallyNil() {
		return false
	}
	return l.flags&flag != 0
}

func (l *distillerLogger) print(flag LogFlag, args ...interface{}) {
	if !l.InternallyNil() && l.hasFlag(flag) {
		l.log.Println(args...)
	}
}
