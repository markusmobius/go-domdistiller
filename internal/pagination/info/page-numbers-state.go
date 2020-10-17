// ORIGINAL: java/PageParamInfo.java

package info

// PageNumbersState is struct that returned by getPageNumbersState() after it
// has checked if the given list of PageLinkInfo's and PageInfo's are adjacent
// and consecutive, and if there's a gap in the list.
type PageNumbersState struct {
	IsAdjacent    bool
	IsConsecutive bool
	NextPagingURL string
}

func (pns *PageNumbersState) isPageNumberSequence(ascendingNumbers []*PageInfo) bool {
	if len(ascendingNumbers) <= 1 {
		return false
	}

	// The first one must have a URL unless it is the first page.
	firstPage := ascendingNumbers[0]
	if firstPage.PageNumber != 1 && firstPage.URL == "" {
		return false
	}

	// There's only one plain number without URL in ascending numbers group.
	hasPlainNum := false
	for _, page := range ascendingNumbers {
		if page.URL == "" {
			if hasPlainNum {
				return false
			}
			hasPlainNum = true
		} else if hasPlainNum && pns.NextPagingURL == "" {
			pns.NextPagingURL = page.URL
		}
	}

	// If there are only two pages, they must be siblings.
	if len(ascendingNumbers) == 2 {
		return firstPage.PageNumber+1 == ascendingNumbers[1].PageNumber
	}

	// Check if page numbers in ascendingNumbers are adjacent and consecutive.
	for i := 1; i < len(ascendingNumbers); i++ {
		// If two adjacent numbers are not consecutive, we accept them only when:
		// 1) one of them is head/tail, like [1],[n-i][n-i+1]..[n] or [1],[2], [3]...[i], [n].
		// 2) both of them have URLs.
		currentPage := ascendingNumbers[i]
		prevPage := ascendingNumbers[i-1]
		if currentPage.PageNumber-prevPage.PageNumber != 1 {
			if i != 1 && i != len(ascendingNumbers)-1 {
				return false
			}

			if currentPage.URL == "" || prevPage.URL == "" {
				return false
			}
		}
	}

	return true
}
