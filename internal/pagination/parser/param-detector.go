// ORIGINAL: java/PageParameterDetector.java

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

// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package parser

import (
	nurl "net/url"

	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
)

// DetectParamInfo creates a PageParamInfo based on outlinks and numeric text around them.
// Always return PageParamInfo (never nil). If no page parameter is detected or
// determined to be best, its ParamType is Unset.
func DetectParamInfo(adjacentNumberGroups *info.MonotonicPageInfoGroups, docURL string, logger logutil.Logger) *info.PageParamInfo {
	// Make sure URL absolute and clean it
	parsedDocURL, err := nurl.ParseRequestURI(docURL)
	if err != nil || parsedDocURL.Scheme == "" || parsedDocURL.Hostname() == "" {
		return &info.PageParamInfo{}
	}
	parsedDocURL.User = nil

	// Start detection
	detectionState := &DetectionState{}
	for _, group := range adjacentNumberGroups.Groups {
		if len(group.List) < 2 {
			continue
		}

		strPattern := ""
		if !detectionState.isEmpty() {
			strPattern = detectionState.bestPageParamInfo.PagePattern
		}

		state := newDetectionStateFromMonotonicNumbers(
			group.List, group.DeltaSign < 0, parsedDocURL, strPattern)
		if state != nil {
			detectionState.compareAndUpdate(state)
		}
	}

	if detectionState.isEmpty() {
		return &info.PageParamInfo{}
	}

	// For now, if there're multiple page patterns, we take the first one.
	// If this doesn't work for most sites, we might have to return nothing.
	if detectionState.hasMultiPagePatterns && logger != nil && !logger.InternallyNil() {
		logger.PrintPaginationInfo("Detected multiple page pattern")
	}

	bestPageParamInfo := detectionState.bestPageParamInfo
	bestPageParamInfo.DetermineNextPagingURL(docURL)
	return bestPageParamInfo
}
