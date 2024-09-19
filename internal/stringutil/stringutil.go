// ORIGINAL: java/StringUtil.java

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

package stringutil

import (
	nurl "net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	rxAllDigits = regexp.MustCompile(`^\d+$`)
)

func IsStringAllWhitespace(str string) bool {
	for _, char := range str {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

func IsStringAllDigit(str string) bool {
	return rxAllDigits.MatchString(str)
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for handling string.
// =================================================================================

// EqualsIgnoreCase checks if two string is similar in case-insensitive mode.
func EqualsIgnoreCase(str1, str2 string) bool {
	return strings.EqualFold(str1, str2)
}

// HasPrefixIgnoreCase checks if str is stared with prefix in case-insensitive mode.
func HasPrefixIgnoreCase(str, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(str), strings.ToLower(prefix))
}

// CharCount returns number of char in str.
func CharCount(str string) int {
	return utf8.RuneCountInString(str)
}

// CreateAbsoluteURL convert url to absolute path based on base.
// However, if url is prefixed with hash (#), the url won't be changed.
func CreateAbsoluteURL(url string, base *nurl.URL) string {
	if url == "" || base == nil {
		return url
	}

	// If it is hash tag, return as it is
	if strings.HasPrefix(url, "#") {
		return url
	}

	// If it is data URI, return as it is
	if strings.HasPrefix(url, "data:") {
		return url
	}

	// If it is javascript URI, return as it is
	if strings.HasPrefix(url, "javascript:") {
		return url
	}

	// If it is already an absolute URL, return as it is
	tmp, err := nurl.ParseRequestURI(url)
	if err == nil && tmp.Scheme != "" && tmp.Hostname() != "" {
		return url
	}

	// Otherwise, resolve against base URI.
	tmp, err = nurl.Parse(url)
	if err != nil {
		return url
	}

	return base.ResolveReference(tmp).String()
}

func UnescapedString(u *nurl.URL) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Opaque != "" {
		buf.WriteString(u.Opaque)
	} else {
		if u.Scheme != "" || u.Host != "" || u.User != nil {
			if u.Host != "" || u.Path != "" || u.User != nil {
				buf.WriteString("//")
			}
			if ui := u.User; ui != nil {
				buf.WriteString(ui.String())
				buf.WriteByte('@')
			}
			if h := u.Host; h != "" {
				buf.WriteString(h)
			}
		}
		path := u.Path
		if path != "" && path[0] != '/' && u.Host != "" {
			buf.WriteByte('/')
		}
		if buf.Len() == 0 {
			// RFC 3986 §4.2
			// A path segment that contains a colon character (e.g., "this:that")
			// cannot be used as the first segment of a relative-path reference, as
			// it would be mistaken for a scheme name. Such a segment must be
			// preceded by a dot-segment (e.g., "./this:that") to make a relative-
			// path reference.
			if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
				buf.WriteString("./")
			}
		}
		buf.WriteString(path)
	}
	if u.ForceQuery || u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.EscapedFragment())
	}
	return buf.String()
}
