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
	"strings"
	"unicode"

	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"golang.org/x/text/unicode/rangetable"
)

var (
	// The following range includes broader alphabetical letters and Hangul Syllables.
	rtMatcher1 = createRangeTable([]rune("azAZ09__\u00C0\u1FFF\uAC00\uD7AF"))

	// The following range includes Hiragana, Katakana, and CJK Unified Ideographs.
	// Hangul Syllables are not included.
	rtMatcher2 = createRangeTable([]rune("\u3040\uA4CF"))

	// The following range includes broader alphabetical letters.
	rtMatcher3 = createRangeTable([]rune("azAZ09__\u00C0\u1FFF"))
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
	// Count alphabetical letters and Hangul Syllables.
	var nMatcher1 int
	for _, word := range strings.FieldsFunc(text, spaceRunes) {
		if containRune(word, rtMatcher1) {
			nMatcher1++
		}
	}

	// Count Hiragana, Katakana, and CJK Unified Ideographs.
	var nMatcher2 int
	for _, r := range text {
		if unicode.Is(rtMatcher2, r) {
			nMatcher2++
		}
	}

	count := nMatcher1 + int(math.Ceil(float64(nMatcher2)*0.55))
	return count
}

func (c LetterWordCounter) Count(text string) int {
	// Count alphabetical letters and Hangul Syllables.
	var count int
	for _, word := range strings.FieldsFunc(text, spaceRunes) {
		if containRune(word, rtMatcher1) {
			count++
		}
	}
	return count
}

func (c FastWordCounter) Count(text string) int {
	// Count broader alphabetical letters.
	var count int
	for _, word := range strings.FieldsFunc(text, spaceRunes) {
		if containRune(word, rtMatcher3) {
			count++
		}
	}
	return count
}

// SelectWordCounter picks the most suitable WordCounter depending on
// the specified text.
func SelectWordCounter(text string) WordCounter {
	switch {
	case re2go.UseFullWordCounter(text):
		return FullWordCounter{}
	case re2go.UseLetterWordCounter(text):
		return LetterWordCounter{}
	default:
		return FastWordCounter{}
	}
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, added to avoid using regex.
// =================================================================================

func createRangeTable(runePairs []rune) *unicode.RangeTable {
	// Make sure rune pairs are even
	nRunePairs := len(runePairs)
	if nRunePairs%2 != 0 {
		return nil
	}

	// Calculate total rune count
	var count int
	for i := 0; i < nRunePairs; i += 2 {
		start, end := runePairs[i], runePairs[i+1]
		if end < start {
			start, end = end, start
			runePairs[i], runePairs[i+1] = runePairs[i+1], runePairs[i]
		}
		count += int(end-start) + 1
	}

	// Convert rune pair to list of rune
	var cursor int
	runes := make([]rune, count)
	for i := 0; i < nRunePairs; i += 2 {
		for r := runePairs[i]; r <= runePairs[i+1]; r++ {
			runes[cursor] = r
			cursor++
		}
	}

	return rangetable.New(runes...)
}

func containRune(s string, rt *unicode.RangeTable) bool {
	for _, r := range s {
		if unicode.Is(rt, r) {
			return true
		}
	}
	return false
}

func spaceRunes(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	default:
		return false
	}
}
