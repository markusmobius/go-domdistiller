// ORIGINAL: javatest/PageParamInfoTest.java

package info_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Info_PLI_GetPageNumbersState(t *testing.T) {
	allNums := []int{1, 2}
	selectedNums := []int{1, 2}
	state := getPageNumbersState(selectedNums, allNums)
	assert.True(t, state.IsAdjacent)
	assert.True(t, state.IsConsecutive)

	allNums = []int{1, 2, 3}
	selectedNums = []int{2, 3}
	state = getPageNumbersState(selectedNums, allNums)
	assert.True(t, state.IsAdjacent)
	assert.True(t, state.IsConsecutive)

	allNums = []int{1, 5, 6, 7, 10}
	selectedNums = []int{1, 5, 7, 10}
	state = getPageNumbersState(selectedNums, allNums)
	assert.True(t, state.IsAdjacent)
	assert.True(t, state.IsConsecutive)

	// No consecutive pairs.
	// TODO: Consider to mark it as consecutive because some news site separate their
	// page number with consistent multiplier, eg [2, 4, 6]
	allNums = []int{10, 25, 50}
	selectedNums = []int{10, 25, 50}
	state = getPageNumbersState(selectedNums, allNums)
	assert.True(t, state.IsAdjacent)
	assert.False(t, state.IsConsecutive)

	// This list doesn't satisfy consecutive rule. There should be "22" on the left of "23",
	// or "25" on the right of "24", or "29" on the left of "30".
	allNums = []int{23, 24, 30}
	selectedNums = []int{23, 24, 30}
	state = getPageNumbersState(selectedNums, allNums)
	assert.True(t, state.IsAdjacent)
	assert.False(t, state.IsConsecutive)

	// Has two gaps
	allNums = []int{1, 2, 3, 4, 5}
	selectedNums = []int{1, 3, 5}
	state = getPageNumbersState(selectedNums, allNums)
	assert.False(t, state.IsAdjacent)
	assert.False(t, state.IsConsecutive)

	// Has a gap of two numbers
	allNums = []int{2, 3, 4, 5}
	selectedNums = []int{2, 5}
	state = getPageNumbersState(selectedNums, allNums)
	assert.False(t, state.IsAdjacent)
	assert.False(t, state.IsConsecutive)
}

func getPageNumbersState(selectedNums, allNums []int) *info.PageNumbersState {
	ascendingNumbers := []*info.PageInfo{}
	numberToPos := make(map[int]int)

	for i := 0; i < len(allNums); i++ {
		number := allNums[i]
		numberToPos[number] = i
		ascendingNumbers = append(ascendingNumbers, &info.PageInfo{
			PageNumber: number,
		})
	}

	allLinkInfo := info.ListLinkInfo{}
	for _, number := range selectedNums {
		allLinkInfo = append(allLinkInfo, &info.PageLinkInfo{
			PageNumber:         number,
			PageParamValue:     number,
			PosInAscendingList: numberToPos[number],
		})
	}

	return allLinkInfo.PageNumbersState(ascendingNumbers)
}
