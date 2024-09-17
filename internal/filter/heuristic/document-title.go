// ORIGINAL: java/filters/heuristics/DocumentTitleMatchClassifier.java

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

// boilerpipe
//
// Copyright (c) 2009 Christian Kohlschütter
//
// The author licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package heuristic

import (
	"regexp"
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/re2go"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

var (
	rxDtmLongestPartPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)[ ]*[\|»|-][ ]*`),
		regexp.MustCompile(`(?i)[ ]*[\|»|:][ ]*`),
		regexp.MustCompile(`(?i)[ ]*[\|»|:\(\)][ ]*`),
		regexp.MustCompile(`(?i)[ ]*[\|»|:\(\)\-][ ]*`),
		regexp.MustCompile(`(?i)[ ]*[\|»|,|:\(\)\-][ ]*`),
		regexp.MustCompile(`(?i)[ ]*[\|»|,|:\(\)\-\x{00a0}][ ]*`),
	}

	rxDtmPotentialTitlePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)[ ]+[\|][ ]+`),
		regexp.MustCompile(`(?i)[ ]+[\-][ ]+`),
	}

	rxDtmPotentialTitleReplacePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i) - [^\-]+$`),
		regexp.MustCompile(`(?i)^[^\-]+ - `),
	}
)

// DocumentTitleMatch marks TextBlocks which contain parts of the HTML
// `title` tag, using some heuristics which are quite specific to the news domain.
type DocumentTitleMatch struct {
	wordCounter     stringutil.WordCounter
	potentialTitles map[string]struct{}
}

func NewDocumentTitleMatch(wc stringutil.WordCounter, titles ...string) *DocumentTitleMatch {
	dtm := DocumentTitleMatch{
		wordCounter:     wc,
		potentialTitles: make(map[string]struct{}),
	}

	for _, title := range titles {
		dtm.processPotentialTitle(title)
	}

	return &dtm
}

func (f *DocumentTitleMatch) Process(doc *webdoc.TextDocument) bool {
	if len(f.potentialTitles) == 0 {
		return false
	}

	changes := false
	for _, tb := range doc.TextBlocks {
		text := tb.Text
		text = strings.ReplaceAll(text, string('\u00a0'), " ")
		text = strings.ReplaceAll(text, "'", "")
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)
		if _, exist := f.potentialTitles[text]; exist {
			tb.AddLabels(label.Title)
			changes = true
			continue
		}

		text = re2go.RemoveDtmCharacters(text)
		text = strings.TrimSpace(text)
		if _, exist := f.potentialTitles[text]; exist {
			tb.AddLabels(label.Title)
			changes = true
		}
	}

	return changes
}

func (f *DocumentTitleMatch) processPotentialTitle(title string) {
	title = strings.ReplaceAll(title, string('\u00a0'), " ")
	title = strings.ReplaceAll(title, "'", "")
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)
	if title == "" {
		return
	}

	if _, exist := f.potentialTitles[title]; exist {
		return
	}

	for _, rx := range rxDtmLongestPartPatterns {
		if p := f.getLongestPart(title, rx); p != "" {
			f.potentialTitles[p] = struct{}{}
		}
	}

	for _, rx := range rxDtmPotentialTitlePatterns {
		f.addPotentialTitles(title, rx, 4)
	}

	for _, rx := range rxDtmPotentialTitleReplacePatterns {
		f.potentialTitles[rx.ReplaceAllString(title, "")] = struct{}{}
	}
}

func (f *DocumentTitleMatch) addPotentialTitles(title string, rx *regexp.Regexp, minWords int) {
	parts := rx.Split(title, -1)
	if len(parts) == 1 {
		return
	}

	for _, p := range parts {
		if strings.Contains(p, ".com") {
			continue
		}

		numWords := f.wordCounter.Count(p)
		if numWords >= minWords {
			f.potentialTitles[p] = struct{}{}
		}
	}
}

func (f *DocumentTitleMatch) getLongestPart(title string, rx *regexp.Regexp) string {
	parts := rx.Split(title, -1)
	if len(parts) == 1 {
		return ""
	}

	longestPart := ""
	longestNumWords := 0
	longestPartLength := 0
	for _, p := range parts {
		if strings.Contains(p, ".com") {
			continue
		}

		numWords := f.wordCounter.Count(p)
		partLength := stringutil.CharCount(p)
		if numWords > longestNumWords || partLength > longestPartLength {
			longestPart = p
			longestNumWords = numWords
			longestPartLength = partLength
		}
	}

	if longestPartLength > 0 {
		return strings.TrimSpace(longestPart)
	}

	return ""
}
