// ORIGINAL: java/OpenGraphProtocolParser.java

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

package opengraph

import "strings"

type Prefix uint

const (
	OG Prefix = iota
	Profile
	Article
)

type PrefixNameList map[Prefix]string

func (prefixes PrefixNameList) addObjectType(prefix, objType string) {
	if objType == "" {
		prefixes[OG] = prefix
		return
	}

	objType = strings.TrimPrefix(objType, "/")
	if objType == ProfileObjtype {
		prefixes[Profile] = prefix
		return
	}

	if objType == ArticleObjtype {
		prefixes[Article] = prefix
	}
}

func (prefixes PrefixNameList) setDefault() {
	// For any unspecified prefix, use common ones:
	// - "og": http://ogp.me/ns#
	// - "profile": http://ogp.me/ns/profile#
	// - "article": http://ogp.me/ns/article#.
	if _, exist := prefixes[OG]; !exist {
		prefixes[OG] = "og"
	}

	if _, exist := prefixes[Profile]; !exist {
		prefixes[Profile] = ProfileObjtype
	}

	if _, exist := prefixes[Article]; !exist {
		prefixes[Article] = ArticleObjtype
	}
}
