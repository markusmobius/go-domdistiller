// ORIGINAL: java/PageParameterDetector.java

package parser

import (
	nurl "net/url"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/markusmobius/go-domdistiller/logutil"
)

// DetectParamInfo creates a PageParamInfo based on outlinks and numeric text around them.
// Always return PageParamInfo (never nil). If no page parameter is detected or
// determined to be best, its ParamType is Unset.
func DetectParamInfo(adjacentNumberGroups *info.MonotonicPageInfoGroups, docURL string) *info.PageParamInfo {
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
	if detectionState.hasMultiPagePatterns {
		logutil.PrintPaginationInfo("Detected multiple page pattern")
	}

	bestPageParamInfo := detectionState.bestPageParamInfo
	bestPageParamInfo.DetermineNextPagingURL(docURL)
	return bestPageParamInfo
}
