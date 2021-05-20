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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/extractor"
	"github.com/markusmobius/go-domdistiller/internal/pagination"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

	// HTML is the string which contains the distilled content in HTML format.
	HTML string

	// ContentImages is list of image URLs that used within the distilled content.
	ContentImages []string
}

// Options is configuration for the distiller.
type Options struct {
	// Whether to extract only the text (or to include the containing html).
	ExtractTextOnly bool

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
	doc, err := parseHTMLDocument(r)
	if err != nil {
		return nil, err
	}

	// Apply distiller to doc
	return Apply(doc, opts)
}

// Apply runs distiller for the specified parsed doc
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
	logger := newDistillerLogger(opts.LogFlags)

	// Start extractor
	ce := extractor.NewContentExtractor(doc, opts.OriginalURL, logger)
	content, wordCount := ce.ExtractContent(opts.ExtractTextOnly)

	result := Result{}
	result.HTML = content
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

	timingInfo.TotalTime = time.Now().Sub(distillerStart)
	result.TimingInfo = *timingInfo

	if logger.hasFlag(LogTiming) {
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

// ====================== INFORMATION ======================
// Methods below these point are used for making sure that
// a HTML document is parsed using UTF-8 encoder, so these
// are not exist in original DOM Distiller.
// =========================================================

var rxCharset = regexp.MustCompile(`(?i)charset\s*=\s*([^;\s"]+)`)

func isSoftHyphen(r rune) bool {
	return r == '\u00AD'
}

func containsErrorRune(bt []byte) bool {
	return bytes.ContainsRune(bt, utf8.RuneError)
}

// normalizeReaderEncoding convert text encoding from NFD to NFC.
// It also remove soft hyphen since apparently it's useless in web.
// See: https://web.archive.org/web/19990117011731/http://www.hut.fi/~jkorpela/shy.html
func normalizeReaderEncoding(r io.Reader) io.Reader {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.Predicate(isSoftHyphen)), norm.NFC)
	return transform.NewReader(r, transformer)
}

// parseHTMLDocument parses a reader and try to convert the character encoding into UTF-8.
func parseHTMLDocument(r io.Reader) (*html.Node, error) {
	// Check if there is invalid character.
	r, valid := validateRunes(r)

	// If already valid, we can parse and return it.
	if valid {
		r = normalizeReaderEncoding(r)
		return html.Parse(r)
	}

	// Find the decoder and parse HTML.
	r, charset := findHtmlCharset(r)
	r = transform.NewReader(r, charset.NewDecoder())
	r = normalizeReaderEncoding(r)
	return html.Parse(r)
}

// validateRunes check a reader to make sure all runes inside it is
// valid UTF-8 character. Since a reader only usable once, here
// we also return a mirror of the reader.
func validateRunes(r io.Reader) (io.Reader, bool) {
	buffer := bytes.NewBuffer(nil)
	tee := io.TeeReader(r, buffer)

	allValid := true
	scanner := bufio.NewScanner(tee)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !utf8.Valid(line) || containsErrorRune(line) {
			allValid = false
			break
		}
	}

	ioutil.ReadAll(tee)
	return buffer, allValid
}

// validateRunes check HTML page in the reader to find charset
// that used in the HTML page. If none found, we assume it's
// UTF-8. Since a reader only usable once, here we also return
// a mirror of the reader.
func findHtmlCharset(r io.Reader) (io.Reader, encoding.Encoding) {
	// Prepare mirror
	buffer := bytes.NewBuffer(nil)
	tee := io.TeeReader(r, buffer)

	// Parse HTML normally
	doc, err := html.Parse(tee)
	if err != nil {
		return buffer, unicode.UTF8
	}

	// Look for charset in <meta> elements
	var customCharset string
	for _, meta := range dom.GetElementsByTagName(doc, "meta") {
		// Look in charset
		charsetAttr := dom.GetAttribute(meta, "charset")
		if charsetAttr != "" {
			customCharset = strings.TrimSpace(charsetAttr)
			break
		}

		// Look in http-equiv="Content-Type"
		content := dom.GetAttribute(meta, "content")
		httpEquiv := dom.GetAttribute(meta, "http-equiv")
		if httpEquiv == "Content-Type" && content != "" {
			matches := rxCharset.FindStringSubmatch(content)
			if len(matches) > 0 {
				customCharset = matches[1]
				break
			}
		}
	}

	// If there are no custom charset specified, assume it's utf-8
	if customCharset == "" {
		return buffer, unicode.UTF8
	}

	// Find the encoder
	e, _ := charset.Lookup(customCharset)
	if e == nil {
		e = unicode.UTF8
	}

	return buffer, e
}
