// ORIGINAL: java/ContentExtractor.java

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

package extractor

import (
	nurl "net/url"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/converter"
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/markup"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

const (
	documentCharThreshold = 500
)

type ContentExtractor struct {
	Parser      *markup.Parser
	TimingInfo  *data.TimingInfo
	ImageURLs   []string
	WordCounter stringutil.WordCounter

	pageURL         *nurl.URL
	documentElement *html.Node
	candidateTitles []string
	logger          logutil.Logger
}

func NewContentExtractor(root *html.Node, pageURL *nurl.URL, logger logutil.Logger) *ContentExtractor {
	timingInfo := &data.TimingInfo{}

	document := dom.QuerySelector(root, "html")
	if document == nil {
		document = root
	}
	start := time.Now()
	parser := markup.NewParser(document, timingInfo)
	timingInfo.MarkupParsingTime = time.Since(start)

	textContent := dom.TextContent(document)
	wordCounter := stringutil.SelectWordCounter(textContent)

	return &ContentExtractor{
		Parser:      parser,
		TimingInfo:  timingInfo,
		WordCounter: wordCounter,

		documentElement: document,
		pageURL:         pageURL,
		logger:          logger,
	}
}

func (ce *ContentExtractor) ExtractTitle() string {
	ce.ensureTitleInitialized()
	if len(ce.candidateTitles) > 0 {
		return ce.candidateTitles[0]
	}
	return ""
}

func (ce *ContentExtractor) ExtractContent() (*webdoc.Document, int) {
	start := time.Now()
	webDocument := ce.createWebDocumentInfoFromPage(converter.SkipUnlikelies)
	wordCount := ce.processDocument(webDocument)

	if wordCount < documentCharThreshold {
		webDocument = ce.createWebDocumentInfoFromPage(converter.Default)
		wordCount = ce.processDocument(webDocument)
	}

	ce.TimingInfo.DocumentConstructionTime = time.Since(start)

	start = time.Now()
	docfilter.NewRelevantElements().Process(webDocument)
	docfilter.NewLeadImageFinder(ce.logger).Process(webDocument)
	docfilter.NewNestedElementRetainer().Process(webDocument)
	ce.TimingInfo.ArticleProcessingTime = time.Since(start)

	ce.ImageURLs = webDocument.GetImageURLs()
	return webDocument, wordCount
}

// ensureTitleInitialized populates list of candidate titles in
// descending priority order:
// 1) meta-information
// 2) The document's title element, modified based on some readability heuristics
// 3) The document's title element, if it's a string
func (ce *ContentExtractor) ensureTitleInitialized() {
	if len(ce.candidateTitles) > 0 {
		return
	}

	title := ce.Parser.Title()
	if title != "" {
		ce.candidateTitles = append(ce.candidateTitles, title)
	}

	documentTitle := getDocumentTitle(ce.documentElement, ce.WordCounter)
	ce.candidateTitles = append(ce.candidateTitles, documentTitle)
}

// createWebDocumentInfoFromPage converts the original HTML page into a webdoc.Document for analysis.
func (ce *ContentExtractor) createWebDocumentInfoFromPage(flags converter.ConverterFlag) *webdoc.Document {
	docBuilder := webdoc.NewWebDocumentBuilder(ce.WordCounter, ce.pageURL)
	converter.NewDomConverter(flags, docBuilder, ce.pageURL, ce.logger).Convert(ce.documentElement)
	webDocument := docBuilder.Build()
	ce.ensureTitleInitialized()
	return webDocument
}

// processDocument do the actual analysis of the page content,
// identifying the core elements of the page. Returns word count
// inside document.
func (ce *ContentExtractor) processDocument(doc *webdoc.Document) int {
	textDocument := doc.CreateTextDocument()

	NewArticleExtractor(ce.logger).Extract(textDocument, ce.WordCounter, ce.candidateTitles)
	wordCount := textDocument.CountWordsInContent()

	textDocument.ApplyToModel()
	return wordCount
}
