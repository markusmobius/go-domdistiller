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

import (
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

// Classifier is object that classifies the table
type Classifier struct {
	logger logutil.Logger
}

func NewClassifier(logger logutil.Logger) *Classifier {
	return &Classifier{logger: logger}
}

// Classify classifies a <table> element as layout or data type, based on the set of heuristics at
// http://asurkov.blogspot.com/2011/10/data-vs-layout-table.html, with some modifications to suit
// our distillation needs.
func (c *Classifier) Classify(t *html.Node) (Type, Reason) {
	// Different from url above, table created by CSS display style is layout table, because
	// we only handle actual <table> elements.

	// 1) Table inside editable area is layout table, different from said url because we ignore
	//    editable areas during distillation.
	parent := t.Parent
	for parent != nil {
		parentTagName := dom.TagName(parent)
		parentEditable := dom.GetAttribute(parent, "contenteditable")
		if parentTagName == "input" || strings.ToLower(parentEditable) == "true" {
			return c.logAndReturn(Layout, InsideEditableArea)
		}
		parent = parent.Parent
	}

	// 2) Table having role="presentation" is layout table.
	tableRole := strings.ToLower(dom.GetAttribute(t, "role"))
	if tableRole == "presentation" {
		return c.logAndReturn(Layout, RoleTable)
	}

	// 3) Table having ARIA table-related roles is data table.
	_, ariaRoleExist := ariaRoles[tableRole]
	_, ariaTableRoleExist := ariaTableRoles[tableRole]
	if ariaRoleExist || ariaTableRoleExist {
		return c.logAndReturn(Data, RoleTable)
	}

	// 4) Table having ARIA table-related roles in its descendants is data table.
	// This may have deviated from said url if it only checks for <table> element but not its
	// descendants.
	directDescendants := c.getDirectDescendants(t)
	for _, element := range directDescendants {
		role := strings.ToLower(dom.GetAttribute(element, "role"))
		_, ariaRoleExist := ariaRoles[role]
		_, ariaDescendantExist := ariaTableDescendantRoles[role]
		if ariaRoleExist || ariaDescendantExist {
			return c.logAndReturn(Data, RoleDescendant)
		}
	}

	// 5) Table having datatable="0" attribute is layout table.
	if dom.GetAttribute(t, "datatable") == "0" {
		return c.logAndReturn(Layout, Datatable0)
	}

	// 6) Table having nested table(s) is layout table.
	// The order here and #7 (table having <=1 row/col is layout table) is different from said
	// url: the latter has these heuristics after #10 (table having "summary" attribute is
	// data table), but our eval sets indicate the need to bump these way up to here, because
	// many (old) pages have layout tables that are nested or with <TH>/<CAPTION>s but only 1
	// row or col.
	if c.hasNestedTables(t) {
		return c.logAndReturn(Layout, NestedTable)
	}

	// 7) Table having only one row or column is layout table.
	// See comments for #6 about deviation from said url.
	rowCount, columnCount := c.getRowAndColumnCount(t)
	if rowCount <= 1 {
		return c.logAndReturn(Layout, LessEq1Row)
	}
	if columnCount <= 1 {
		return c.logAndReturn(Layout, LessEq1Col)
	}

	// 8) Table having legitimate data table structures is data table :
	// a. table has <caption>, <thead>, <tfoot>, <colgroup>, <col>, or <th> elements
	caption := dom.QuerySelector(t, "caption")
	tHead := dom.QuerySelector(t, "thead")
	tFoot := dom.QuerySelector(t, "tfoot")
	if (caption != nil && c.hasValidText(caption)) || tHead != nil || tFoot != nil ||
		c.hasOneOfElements(directDescendants, headerTags) {
		return c.logAndReturn(Data, CaptionTheadTfootColgroupColTh)
	}

	// Extract all <td> elements from direct descendants, for easier/faster multiple access.
	directTDs := []*html.Node{}
	for _, element := range directDescendants {
		if dom.TagName(element) == "td" {
			directTDs = append(directTDs, element)
		}
	}

	for _, td := range directTDs {
		// b) table cell has abbr, headers, or scope attributes
		if dom.HasAttribute(td, "abbr") || dom.HasAttribute(td, "headers") || dom.HasAttribute(td, "scope") {
			return c.logAndReturn(Data, AbbrHeadersScope)
		}

		// c) table cell has <abbr> element as a single child element.
		tdChildren := dom.GetElementsByTagName(td, "*")
		if len(tdChildren) == 1 && dom.TagName(tdChildren[0]) == "abbr" {
			return c.logAndReturn(Data, OnlyHasAbbr)
		}
	}

	// 9) Table occupying > 95% of document width without viewport meta is layout table;
	// viewport condition is not in said url, added here for typical mobile-optimized sites.
	// The order here is different from said url: the latter has it after #14 (>=20 rows is
	// data table), but our eval sets indicate the need to bump this way up to here, because
	// many (old) pages have layout tables with the "summary" attribute (#10).
	//
	// Unfortunately, to do this we need to compute the stylesheets which is not possible
	// right now in Go. So we will skip it. NEED-COMPUTE-CSS.

	// 10) Table having summary attribute is data table.
	// This is different from said url: the latter lumps "summary" attribute with #8, but we
	// split it so as to insert #9 in between. Many (old) pages have tables that are clearly
	// layout: their "summary" attributes say they're for layout. They also occupy > 95% of
	// document width, so #9 coming before #10 will correctly classify them as layout.
	if dom.HasAttribute(t, "summary") {
		return c.logAndReturn(Data, Summary)
	}

	// 11) Table having >=5 columns is data table.
	if columnCount >= 5 {
		return c.logAndReturn(Data, MoreEq5Cols)
	}

	// 12) Table having borders around cells is data table.
	// Again, this is impossible to do right now. NEED-COMPUTE-CSS.

	// 13) Table having differently-colored rows is data table.
	// Like before, impossible to do right now. NEED-COMPUTE-CSS.

	// 14) Table having >=20 rows is data table.
	if rowCount >= 20 {
		return c.logAndReturn(Data, MoreEq20Rows)
	}

	// 15) Table having <=10 cells is layout table.
	if len(directTDs) <= 10 {
		return c.logAndReturn(Layout, LessEq10Cells)
	}

	// 16) Table containing <embed>, <object>, <applet> or <iframe> elements (typical
	//     advertisement elements) is layout table.
	if c.hasOneOfElements(directDescendants, objectTags) {
		return c.logAndReturn(Layout, EmbedObjectAppletIframe)
	}

	// 17) Table occupying > 90% of document height is layout table.
	// This is not in said url, added here because many (old) pages have tables that
	// don't fall into any of the above heuristics but are for layout, and hence
	// shouldn't default to data by #18.
	//
	// And, unfortunately it's impossible to implement here. NEED-COMPUTE-CSS.

	// 18) Otherwise, it's data table.
	return c.logAndReturn(Data, Default)
}

func (c *Classifier) hasNestedTables(t *html.Node) bool {
	return len(dom.GetElementsByTagName(t, "table")) > 0
}

func (c *Classifier) getDirectDescendants(t *html.Node) []*html.Node {
	// Get all elements inside table
	allDescendants := dom.GetElementsByTagName(t, "*")

	// If there are no nested tables, all descendants is direct descendants
	if !c.hasNestedTables(t) {
		return allDescendants
	}

	directDescendants := []*html.Node{}
	for _, descendant := range allDescendants {
		// Check if the current element is a direct descendant of the `t`
		// table element in question, as opposed to being a descendant of
		// a nested table in `t`.
		parent := descendant.Parent
		for parent != nil {
			if dom.TagName(parent) == "table" {
				if parent == t {
					directDescendants = append(directDescendants, descendant)
				}
				break
			}
			parent = parent.Parent
		}
	}

	return directDescendants
}

func (c *Classifier) hasValidText(e *html.Node) bool {
	text := domutil.InnerText(e)
	return text != "" && !stringutil.IsStringAllWhitespace(text)
}

func (c *Classifier) hasOneOfElements(elements []*html.Node, tags map[string]bool) bool {
	for _, element := range elements {
		tagName := dom.TagName(element)
		if value, exist := tags[tagName]; exist {
			return !value || c.hasValidText(element)
		}
	}
	return false
}

// getRowAndColumnCount returns how many rows and columns this table has.
// This method doesn't exist in original dom-distiller because GWT already provides
// method for calculating count of rows and columns. As workaround, we use method
// from go-readability.
func (c *Classifier) getRowAndColumnCount(t *html.Node) (int, int) {
	rows := 0
	columns := 0
	trs := dom.GetElementsByTagName(t, "tr")
	for i := 0; i < len(trs); i++ {
		strRowSpan := dom.GetAttribute(trs[i], "rowspan")
		rowSpan, _ := strconv.Atoi(strRowSpan)
		if rowSpan == 0 {
			rowSpan = 1
		}
		rows += rowSpan

		// Now look for column-related info
		columnsInThisRow := 0
		cells := dom.GetElementsByTagName(trs[i], "td")
		for j := 0; j < len(cells); j++ {
			strColSpan := dom.GetAttribute(cells[j], "colspan")
			colSpan, _ := strconv.Atoi(strColSpan)
			if colSpan == 0 {
				colSpan = 1
			}
			columnsInThisRow += colSpan
		}

		if columnsInThisRow > columns {
			columns = columnsInThisRow
		}
	}

	return rows, columns
}

func (c *Classifier) logAndReturn(tableType Type, reason Reason) (Type, Reason) {
	if c.logger != nil && !c.logger.InternallyNil() {
		c.logger.PrintVisibilityInfo(reason, "=>", tableType)
	}
	return tableType, reason
}
