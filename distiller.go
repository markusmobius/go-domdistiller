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

package distiller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"os"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/extractor"
	"github.com/markusmobius/go-domdistiller/internal/pagination"
	"golang.org/x/net/html"
)

// PaginationAlgo is the algorithm to find the pagination links.
type PaginationAlgo uint

const (
	// PrevNext is the algorithm to find pagination links that work by scoring  each anchor
	// in documents using various heuristics on its href, text, class name and ID. It's quite
	// accurate and used as default algorithm. Unfortunately it uses a lot of regular expressions,
	// so it's a bit slow.
	PrevNext PaginationAlgo = iota

	// PageNumber is algorithm to find pagination links that work by collecting groups of adjacent plain
	// text numbers and outlinks with digital anchor text. A lot faster than PrevNext, but also less
	// accurate.
	PageNumber
)

// Result is the final output of the distiller
type Result struct {
	// URL is the URL of the processed page.
	URL string

	// Title is the title of the processed page.
	Title string

	// MarkupInfo is the metadata of the page. The metadata is extracted following three markup
	// specifications: OpenGraphProtocol, IEReadingView and SchemaOrg. For now, OpenGraph protocol
	// takes precedence because it uses specific meta tags and hence the fastest. The other
	// specifications is used as fallback in case some metadata not found.
	MarkupInfo data.MarkupInfo

	// TimingInfo is the record of the time it takes to do each step in the process of content extraction.
	TimingInfo data.TimingInfo

	// PaginationInfo contains link to previous and next partial page. This is useful for long article or
	// that may be partitioned into several partial pages by its webmaster.
	PaginationInfo data.PaginationInfo

	// WordCount is the count of words within document.
	WordCount int

	// Node is the *html.Node which contain the distilled content.
	Node *html.Node

	// Text is the string which contains the distilled content in text format.
	Text string

	// ContentImages is list of image URLs that used within the distilled content.
	ContentImages []string
}

// Options is configuration for the distiller.
type Options struct {
	// Flags to specify which info to dump to log.
	LogFlags LogFlag

	// Original URL of the page, which is used in the heuristics in detecting
	// next/prev page links. Will be ignored if Option is used in ApplyForURL.
	OriginalURL *nurl.URL

	// Set to true to skip process for finding pagination.
	SkipPagination bool

	// Algorithm to use for next page detection.
	PaginationAlgo PaginationAlgo
}

// ApplyForURL runs distiller for the specified URL.
func ApplyForURL(url string, timeout time.Duration, opts *Options) (*Result, error) {
	// Make sure URL absolute
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil {
		return nil, err
	}

	// Fetch page from URL
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	// Make sure content type is HTML
	cp := resp.Header.Get("Content-Type")
	if !strings.Contains(cp, "text/html") {
		return nil, fmt.Errorf("URL is not a HTML document")
	}

	// Apply distiller to response body
	if opts == nil {
		opts = &Options{}
	}

	opts.OriginalURL = parsedURL
	return ApplyForReader(resp.Body, opts)
}

// ApplyForFile runs distiller for the specified file.
func ApplyForFile(path string, opts *Options) (*Result, error) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	// Apply distiller to file
	return ApplyForReader(f, opts)
}

// Apply runs distiller for the specified io.Reader.
func ApplyForReader(r io.Reader, opts *Options) (*Result, error) {
	// Parse input
	doc, err := dom.Parse(r)
	if err != nil {
		return nil, err
	}

	// Apply distiller to doc
	return Apply(doc, opts)
}

// Apply runs distiller for the specified parsed document.
func Apply(doc *html.Node, opts *Options) (*Result, error) {
	// Mark the start time
	distillerStart := time.Now()

	// Check whether doc is valid
	if doc.Type != html.ElementNode {
		doc = dom.QuerySelector(doc, "*")
		if doc == nil {
			return nil, errors.New("input doesn't have a valid element")
		}
	}

	// Create default options
	if opts == nil {
		opts = &Options{}
	}

	// Prepare logger
	var logger *distillerLogger
	if opts.LogFlags != LogNothing {
		logger = newDistillerLogger(opts.LogFlags)
	}

	// Start extractor
	ce := extractor.NewContentExtractor(doc, opts.OriginalURL, logger)
	extractedDocument, wordCount := ce.ExtractContent()

	// Generate output
	start := time.Now()
	extractedText := extractedDocument.GenerateOutput(true)
	extractedHTML := extractedDocument.GenerateOutput(false)
	ce.TimingInfo.FormattingTime = time.Since(start)

	// Convert generated html string into node
	container := dom.CreateElement("div")
	dom.SetInnerHTML(container, extractedHTML)

	// Prepare result
	result := Result{}
	result.Node = container
	result.Text = extractedText
	result.WordCount = wordCount
	result.Title = ce.ExtractTitle()
	result.ContentImages = ce.ImageURLs
	result.MarkupInfo = ce.Parser.MarkupInfo()

	if opts.OriginalURL != nil {
		result.URL = opts.OriginalURL.String()
	}

	// Find pagination
	timingInfo := ce.TimingInfo
	if !opts.SkipPagination && opts.OriginalURL != nil {
		paginationStart := time.Now()

		if opts.PaginationAlgo == PageNumber {
			finder := pagination.NewPageNumberFinder(ce.WordCounter, nil, logger)
			result.PaginationInfo = finder.FindPagination(doc, opts.OriginalURL)
			logger.PrintPaginationInfo("Paging by PageNum, prev: " + result.PaginationInfo.PrevPage)
			logger.PrintPaginationInfo("Paging by PageNum, next: " + result.PaginationInfo.NextPage)
		} else {
			finder := pagination.NewPrevNextFinder(logger)
			result.PaginationInfo = finder.FindPagination(doc, opts.OriginalURL)
			logger.PrintPaginationInfo("Paging by PrevNext, prev: " + result.PaginationInfo.PrevPage)
			logger.PrintPaginationInfo("Paging by PrevNext, next: " + result.PaginationInfo.NextPage)
		}

		timingInfo.AddEntry(paginationStart, "Pagination")
	}

	timingInfo.TotalTime = time.Since(distillerStart)
	result.TimingInfo = *timingInfo

	if logger != nil && !logger.InternallyNil() && logger.hasFlag(LogTiming) {
		for _, entry := range ce.TimingInfo.OtherTimes {
			logger.PrintTimingInfo("Timing:", entry.Name, "=", entry.Time)
		}

		logger.PrintTimingInfo("TimingMarkupParsingTime =", ce.TimingInfo.MarkupParsingTime)
		logger.PrintTimingInfo("TimingDocumentConstructionTime =", ce.TimingInfo.DocumentConstructionTime)
		logger.PrintTimingInfo("TimingArticleProcessingTime =", ce.TimingInfo.ArticleProcessingTime)
		logger.PrintTimingInfo("TimingFormattingTime =", ce.TimingInfo.FormattingTime)
		logger.PrintTimingInfo("TimingTotalTime =", ce.TimingInfo.TotalTime)
	}

	return &result, nil
}
