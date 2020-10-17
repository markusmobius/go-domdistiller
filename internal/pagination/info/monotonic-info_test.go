// ORIGINAL: javatest/MonotonicPageInfosGroupsTest.java

package info_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/info"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_MPIG_BasicAscending(t *testing.T) {
	allNums := []int{1, 2, 3}
	groups := createMonotonicPageInfoGroups(allNums)

	assert.Len(t, groups.Groups, 1)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, allNums)
}

func Test_Pagination_MPIG_BasicDescending(t *testing.T) {
	allNums := []int{3, 2, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	assert.Len(t, groups.Groups, 1)
	assertMonotonicPageInfoGroups(t, groups, 0, -1, allNums)
}

func Test_Pagination_MPIG_One(t *testing.T) {
	allNums := []int{1}
	groups := createMonotonicPageInfoGroups(allNums)

	assert.Len(t, groups.Groups, 1)
	assertMonotonicPageInfoGroups(t, groups, 0, 0, allNums)
}

func Test_Pagination_MPIG_TwoSame(t *testing.T) {
	allNums := []int{1, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums := []int{1}
	assert.Len(t, groups.Groups, 1)
	assertMonotonicPageInfoGroups(t, groups, 0, 0, expectedNums)
}

func Test_Pagination_MPIG_AscendingAndDescending1(t *testing.T) {
	allNums := []int{1, 2, 3, 3, 2, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{1, 2, 3}
	expectedNums1 := []int{3, 2, 1}
	assert.Len(t, groups.Groups, 2)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, -1, expectedNums1)
}

func Test_Pagination_MPIG_AscendingAndDescending2(t *testing.T) {
	allNums := []int{1, 2, 3, 2, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{1, 2, 3}
	expectedNums1 := []int{3, 2, 1}
	assert.Len(t, groups.Groups, 2)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, -1, expectedNums1)
}

func Test_Pagination_MPIG_AscendingAndDescending3(t *testing.T) {
	allNums := []int{1, 3, 5, 4, 2, 1, 10, 999, 0}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{1, 3, 5}
	expectedNums1 := []int{5, 4, 2, 1}
	expectedNums2 := []int{1, 10, 999}
	expectedNums3 := []int{999, 0}
	assert.Len(t, groups.Groups, 4)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, -1, expectedNums1)
	assertMonotonicPageInfoGroups(t, groups, 2, 1, expectedNums2)
	assertMonotonicPageInfoGroups(t, groups, 3, -1, expectedNums3)
}

func Test_Pagination_MPIG_DuplicateAscending1(t *testing.T) {
	allNums := []int{1, 1, 2, 3}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums := []int{1, 2, 3}
	assert.Len(t, groups.Groups, 1)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, expectedNums)
}

func Test_Pagination_MPIG_DuplicateAscending2(t *testing.T) {
	allNums := []int{1, 2, 2, 3}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{1, 2}
	expectedNums1 := []int{2, 3}
	assert.Len(t, groups.Groups, 2)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, 1, expectedNums1)
}

func Test_Pagination_MPIG_DuplicateAscending3(t *testing.T) {
	allNums := []int{1, 2, 3, 3}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{1, 2, 3}
	expectedNums1 := []int{3}
	assert.Len(t, groups.Groups, 2)
	assertMonotonicPageInfoGroups(t, groups, 0, 1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, 0, expectedNums1)
}

func Test_Pagination_MPIG_DuplicateDescending1(t *testing.T) {
	allNums := []int{3, 2, 1, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{3, 2, 1}
	expectedNums1 := []int{1}
	assert.Len(t, groups.Groups, 2)
	assertMonotonicPageInfoGroups(t, groups, 0, -1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, 0, expectedNums1)
}

func Test_Pagination_MPIG_DuplicateDescending2(t *testing.T) {
	allNums := []int{3, 2, 2, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{3, 2}
	expectedNums1 := []int{2, 1}
	assert.Len(t, groups.Groups, 2)
	assertMonotonicPageInfoGroups(t, groups, 0, -1, expectedNums0)
	assertMonotonicPageInfoGroups(t, groups, 1, -1, expectedNums1)
}

func Test_Pagination_MPIG_DuplicateDescending3(t *testing.T) {
	allNums := []int{3, 3, 2, 1}
	groups := createMonotonicPageInfoGroups(allNums)

	expectedNums0 := []int{3, 2, 1}
	assert.Len(t, groups.Groups, 1)
	assertMonotonicPageInfoGroups(t, groups, 0, -1, expectedNums0)
}

func createMonotonicPageInfoGroups(numbers []int) *info.MonotonicPageInfoGroups {
	groups := &info.MonotonicPageInfoGroups{}
	groups.AddGroup()

	for _, number := range numbers {
		groups.AddNumber(number, "")
	}

	return groups
}

func assertMonotonicPageInfoGroups(t *testing.T, groups *info.MonotonicPageInfoGroups, index int, expectedDeltaSign int, expectedNums []int) {
	group := groups.Groups[index]
	assert.Equal(t, expectedDeltaSign, group.DeltaSign)
	assert.Equal(t, len(expectedNums), len(group.List))
	for i, num := range expectedNums {
		assert.Equal(t, num, group.List[i].PageNumber)
	}
}
