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

var headerTags = map[string]bool{
	"colgroup": false,
	"col":      false,
	"th":       true,
}

var objectTags = map[string]bool{
	"embed":  false,
	"object": false,
	"applet": false,
	"iframe": false,
}

// ARIA roles for table, see http://www.w3.org/TR/wai-aria/roles#widget_roles_header.
var ariaTableRoles = map[string]struct{}{
	"grid":     {},
	"treegrid": {},
}

// ARIA roles for descendants of table, see :
// - http://www.w3.org/TR/wai-aria/roles#widget_roles_header.
// - http://www.w3.org/TR/wai-aria/roles#document_structure_roles_header.
var ariaTableDescendantRoles = map[string]struct{}{
	"gridcell":     {},
	"columnheader": {},
	"row":          {},
	"rowgroup":     {},
	"rowheader":    {},
}

// ARIA landmark roles, applicable to both table and its descendants
// - http://www.w3.org/TR/wai-aria/roles#landmark_roles_header.
var ariaRoles = map[string]struct{}{
	"application":   {},
	"banner":        {},
	"complementary": {},
	"contentinfo":   {},
	"form":          {},
	"main":          {},
	"navigation":    {},
	"search":        {},
}
