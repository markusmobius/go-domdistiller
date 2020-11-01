// ORIGINAL: java/TableClassifier.java

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

package tableclass

type Reason uint

const (
	Unknown Reason = iota
	InsideEditableArea
	RoleTable
	RoleDescendant
	Datatable0
	CaptionTheadTfootColgroupColTh
	AbbrHeadersScope
	OnlyHasAbbr
	More95PercentDocWidth
	Summary
	NestedTable
	LessEq1Row
	LessEq1Col
	MoreEq5Cols
	CellsHaveBorder
	DifferentlyColoredRows
	MoreEq20Rows
	LessEq10Cells
	EmbedObjectAppletIframe
	More90PercentDocHeight
	Default
)

func (r Reason) String() string {
	switch r {
	case InsideEditableArea:
		return "InsideEditableArea"
	case RoleTable:
		return "RoleTable"
	case RoleDescendant:
		return "RoleDescendant"
	case Datatable0:
		return "Datatable0"
	case CaptionTheadTfootColgroupColTh:
		return "CaptionTheadTfootColgroupColTh"
	case AbbrHeadersScope:
		return "AbbrHeadersScope"
	case OnlyHasAbbr:
		return "OnlyHasAbbr"
	case More95PercentDocWidth:
		return "More95PercentDocWidth"
	case Summary:
		return "Summary"
	case NestedTable:
		return "NestedTable"
	case LessEq1Row:
		return "LessEq1Row"
	case LessEq1Col:
		return "LessEq1Col"
	case MoreEq5Cols:
		return "MoreEq5Cols"
	case CellsHaveBorder:
		return "CellsHaveBorder"
	case DifferentlyColoredRows:
		return "DifferentlyColoredRows"
	case MoreEq20Rows:
		return "MoreEq20Rows"
	case LessEq10Cells:
		return "LessEq10Cells"
	case EmbedObjectAppletIframe:
		return "EmbedObjectAppletIframe"
	case More90PercentDocHeight:
		return "More90PercentDocHeight"
	case Default:
		return "Default"
	}
	return "Unknown"
}
