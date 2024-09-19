// ORIGINAL: java/extractors/ArticleExtractor.java

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

package extractor

import (
	"github.com/markusmobius/go-domdistiller/internal/filter/english"
	"github.com/markusmobius/go-domdistiller/internal/filter/heuristic"
	"github.com/markusmobius/go-domdistiller/internal/filter/simple"
	"github.com/markusmobius/go-domdistiller/internal/label"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
)

type ArticleExtractor struct {
	logger logutil.Logger
}

func NewArticleExtractor(logger logutil.Logger) *ArticleExtractor {
	return &ArticleExtractor{logger: logger}
}

// Extract extracts TextDocument. It is tuned towards news articles.
func (ae *ArticleExtractor) Extract(doc *webdoc.TextDocument, wc stringutil.WordCounter, candidateTitles []string) bool {
	// Prepare filters
	terminatingBlocksFinder := english.NewTerminatingBlocksFinder()
	documentTitleMatch := heuristic.NewDocumentTitleMatch(wc, candidateTitles...)
	numWordsRulesClassifier := english.NewNumWordsRulesClassifier()
	labelNotContentToBoilerplate := simple.NewLabelToBoilerplate(label.StrictlyNotContent)

	similarSiblingContentExpansion1 := heuristic.NewSimilarSiblingContentExpansion()
	similarSiblingContentExpansion1.AllowCrossHeadings = true
	similarSiblingContentExpansion1.MaxLinkDensity = 0.5
	similarSiblingContentExpansion1.MaxBlockDistance = 10

	similarSiblingContentExpansion2 := heuristic.NewSimilarSiblingContentExpansion()
	similarSiblingContentExpansion2.AllowCrossHeadings = true
	similarSiblingContentExpansion2.AllowMixedTags = true
	similarSiblingContentExpansion2.MaxBlockDistance = 10

	headingFusion := heuristic.NewHeadingFusion()
	blockProximityFusionPre := heuristic.NewBlockProximityFusion(false)
	boilerplateBlockKeepTitle := simple.NewBoilerplateBlock(label.Title)
	blockProximityFusionPost := heuristic.NewBlockProximityFusion(true)
	keepLargestBlockExpandToSibling := heuristic.NewKeepLargestBlock(true)
	expandTitleToContent := heuristic.NewExpandTitleToContent()
	largeBlockAroundTagLevel := heuristic.NewLargeBlockAroundTagLevelToContent()
	listAtEnd := heuristic.NewListAtEnd()

	ae.printArticleLog(doc, true, "Start")

	// Run filters
	// Intentionally don't print changes from these two
	terminatingBlocksFinder.Process(doc)
	documentTitleMatch.Process(doc)

	changed := numWordsRulesClassifier.Process(doc)
	ae.printArticleLog(doc, changed, "Classification complete")

	changed = labelNotContentToBoilerplate.Process(doc)
	ae.printArticleLog(doc, changed, "Ignore strictly not content blocks")

	changed = similarSiblingContentExpansion1.Process(doc)
	ae.printArticleLog(doc, changed, "Cross headings SimilarSiblingContentExpansion")

	changed = similarSiblingContentExpansion2.Process(doc)
	ae.printArticleLog(doc, changed, "Mixed tags SimilarSiblingContentExpansion")

	changed = headingFusion.Process(doc)
	ae.printArticleLog(doc, changed, "HeadingFusion")

	changed = blockProximityFusionPre.Process(doc)
	ae.printArticleLog(doc, changed, "BlockProximityFusion for distance=1")

	changed = boilerplateBlockKeepTitle.Process(doc)
	ae.printArticleLog(doc, changed, "BlockFilter keep title")

	changed = blockProximityFusionPost.Process(doc)
	ae.printArticleLog(doc, changed, "BlockProximityFusion for same level content-only")

	changed = keepLargestBlockExpandToSibling.Process(doc)
	ae.printArticleLog(doc, changed, "Keep largest block")

	changed = expandTitleToContent.Process(doc)
	ae.printArticleLog(doc, changed, "Expand title to content")

	changed = largeBlockAroundTagLevel.Process(doc)
	ae.printArticleLog(doc, changed, "Largest block with same tag level to content")

	changed = listAtEnd.Process(doc)
	ae.printArticleLog(doc, changed, "List at end filter")

	return true
}

func (ae *ArticleExtractor) printArticleLog(doc *webdoc.TextDocument, changed bool, header string) {
	if ae.logger == nil || ae.logger.InternallyNil() {
		return
	}

	logMsg := ""
	if !changed {
		logMsg = header + ": NO CHANGES"
	} else {
		logMsg = header + ":\n" + doc.DebugString()
	}

	ae.logger.PrintExtractionInfo(logMsg)
}
