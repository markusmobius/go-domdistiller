// ORIGINAL: java/MarkupParser.java

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

package markup

import "github.com/markusmobius/go-domdistiller/data"

// Accessor is the interface that all parsers must implement so that Parser
// can retrieve their properties.
type Accessor interface {
	// Title returns the markup title of the document, empty if none.
	Title() string

	// Type returns the markup type of the document, empty if none.
	Type() string

	// URL returns the markup url of the document, empty if none.
	URL() string

	// Images returns the properties of all markup images in the document.
	// The first image is the dominant (i.e. top or salient) one.
	Images() []data.MarkupImage

	// Description returns the markup description of the document, empty if none.
	Description() string

	// Publisher returns the markup publisher of the document, empty if none.
	Publisher() string

	// Copyright returns the markup copyright of the document, empty if none.
	Copyright() string

	// Author returns the full name of the markup author, empty if none.
	Author() string

	// Article returns the properties of the markup "article" object, null if none.
	Article() *data.MarkupArticle

	// OptOut returns true if page owner has opted out of distillation.
	OptOut() bool
}
