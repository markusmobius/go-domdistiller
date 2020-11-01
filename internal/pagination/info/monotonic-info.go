// ORIGINAL: java/MonotonicPageInfosGroups.java

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

package info

// MonotonicPageInfoGroups stores all numeric content (both outlinks and plain text pieces)
// parsed from the document, grouped by their monotonicity and adjacency. Basically, it's a
// list of groups of monotonic and adjacent PageInfo's (with or without links) in the document.
type MonotonicPageInfoGroups struct {
	Groups       []*PageInfoGroup
	prevPageInfo *PageInfo
}

// AddGroup adds a new group because a non-plain-number has been encountered in the document
// being parsed i.e. there's a break in the adjacency of plain numbers.
func (mig *MonotonicPageInfoGroups) AddGroup() {
	if len(mig.Groups) == 0 || len(mig.lastGroup().List) > 0 {
		mig.Groups = append(mig.Groups, &PageInfoGroup{})
		mig.prevPageInfo = nil
	}
}

// AddPageInfo adds the given PageInfo, ensuring the group stays monotonic:
// - add in the current group if the page number is strictly increasing or decreasing
// - otherwise, start a new group.
func (mig *MonotonicPageInfoGroups) AddPageInfo(pageInfo *PageInfo) {
	group := mig.lastGroup()
	if group == nil {
		return
	}

	if len(group.List) == 0 {
		group.List = append(group.List, pageInfo)
		mig.prevPageInfo = pageInfo
		return
	}

	deltaSign := 0
	delta := pageInfo.PageNumber - mig.prevPageInfo.PageNumber

	if delta > 0 {
		deltaSign = 1
	} else if delta < 0 {
		deltaSign = -1
	}

	if deltaSign != group.DeltaSign {
		// group.mDeltaSign = 0 means the group only has 1 entry, and hence no ordering yet;
		// the new deltaSign would determine the ordering.
		// Otherwise, the group has been strictly ascending/descending until this number, in
		// which case, start a new group:
		// - with this number if it's the same as previous (deltaSign = 0)
		// - with the previous, followed by this, if both are different numbers.
		if group.DeltaSign != 0 {
			group = &PageInfoGroup{}
			if deltaSign != 0 {
				group.List = []*PageInfo{mig.prevPageInfo}
			}

			mig.Groups = append(mig.Groups, group)
		}
	} else if deltaSign == 0 {
		// This number is same as previous (i.e. the only entry in the group), and will be added
		// (below), so clear the group to avoid duplication.
		group.List = []*PageInfo{}
	}

	group.List = append(group.List, pageInfo)
	group.DeltaSign = deltaSign
	mig.prevPageInfo = pageInfo
}

// AddNumber adds a PageInfo for the given page number and URL, ensuring the group stays monotonic.
func (mig *MonotonicPageInfoGroups) AddNumber(number int, url string) {
	mig.AddPageInfo(&PageInfo{
		URL:        url,
		PageNumber: number,
	})
}

// CleanUp removes last empty adjacent number group.
func (mig *MonotonicPageInfoGroups) CleanUp() {
	if len(mig.Groups) != 0 && len(mig.lastGroup().List) == 0 {
		lastIdx := len(mig.Groups) - 1
		mig.Groups[lastIdx] = nil
		mig.Groups = mig.Groups[:lastIdx]
	}
}

func (mig *MonotonicPageInfoGroups) lastGroup() *PageInfoGroup {
	nGroup := len(mig.Groups)
	if nGroup == 0 {
		return nil
	}

	return mig.Groups[nGroup-1]
}
