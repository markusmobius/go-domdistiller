// ORIGINAL: Part of java/StringUtil.java

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
	"math"
	"regexp"
)

var (
	rxFullWordCounter   = regexp.MustCompile(`[\x{3040}-\x{A4CF}]`)
	rxLetterWordCounter = regexp.MustCompile(`[\x{AC00}-\x{D7AF}]`)

	rxWordMatcher1 = regexp.MustCompile(`(\S*[\w\x{00C0}-\x{1FFF}\x{AC00}-\x{D7AF}]\S*)`)
	rxWordMatcher2 = regexp.MustCompile(`([\x{3040}-\x{A4CF}])`)
	rxWordMatcher3 = regexp.MustCompile(`(\S*[\w\x{00C0}-\x{1FFF}]\S*)`)
)

// WordCounter is object for counting the number of words. For some languages,
// doing this relies on non-trivial word segmentation algorithms, or even huge
// look-up tables. However, for our purpose this function needs to be reasonably
// fast, so the word count for some languages would only be an approximation.
// Read https://crbug.com/484750 for more info.
type WordCounter interface {
	Count(string) int
}

type FullWordCounter struct{}
type LetterWordCounter struct{}
type FastWordCounter struct{}

func (c FullWordCounter) Count(text string) int {
	// The following range includes broader alphabetical letters and Hangul Syllables.
	matches := rxWordMatcher1.FindAllString(text, -1)
	count := len(matches)

	// The following range includes Hiragana, Katakana, and CJK Unified Ideographs.
	// Hangul Syllables are not included.
	matches = rxWordMatcher2.FindAllString(text, -1)
	count += int(math.Ceil(float64(len(matches)) * 0.55))
	return count
}

func (c LetterWordCounter) Count(text string) int {
	// The following range includes broader alphabetical letters and Hangul Syllables.
	matches := rxWordMatcher1.FindAllString(text, -1)
	return len(matches)
}

func (c FastWordCounter) Count(text string) int {
	// The following range includes broader alphabetical letters.
	matches := rxWordMatcher3.FindAllString(text, -1)
	return len(matches)
}

// SelectWordCounter picks the most suitable WordCounter depending on
// the specified text.
func SelectWordCounter(text string) WordCounter {
	switch {
	case rxFullWordCounter.MatchString(text):
		return FullWordCounter{}
	case rxLetterWordCounter.MatchString(text):
		return LetterWordCounter{}
	default:
		return FastWordCounter{}
	}
}
