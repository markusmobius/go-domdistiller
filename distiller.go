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

// Result is the final output of the distiller
type Result struct {
	// Title is the title of the processed page.
	Title string

	// MarkupInfo is the metadata of the page. The metadata is extracted following three markup
	// specifications: OpenGraphProtocol, IEReadingView and SchemaOrg. For now, OpenGraph protocol
	// takes precedence because it uses specific meta tags and hence the fastest. The other
	// specifications is used as fallback in case some metadata not found.
	MarkupInfo data.MarkupInfo

	// TimingInfo is the record of the time it takes to do each step in the process of content extraction.
	TimingInfo data.TimingInfo

	// DebugInfo contains log of all process.
	DebugInfo data.DebugInfo

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

	// How much debug output to dump to log.
	// (0): Logs nothing
	// (1): Text Node data for each stage of processing
	// (2): (1) and some node visibility information
	// (3): (2) and extracted paging information
	DebugLevel uint

	// Original URL of the page, which is used in the heuristics in
	// detecting next/prev page links.
	OriginalURL *nurl.URL

	// Set to true to skip process for finding pagination.
	SkipPagination bool

	// Which algorithm to use for next page detection:
	// "next"    : detect anchors with "next" text
	// "pagenum" : detect anchors with numeric page numbers
	PaginationAlgo string
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
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	// Apply distiller to doc
	return Apply(doc, opts)
}

// Apply runs distiller for the specified parsed doc
func Apply(doc *html.Node, opts *Options) (*Result, error) {
	//check whether doc is valid
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

	// Start extractor
	ce := extractor.NewContentExtractor(doc, opts.OriginalURL)
	content, wordCount := ce.ExtractContent(opts.ExtractTextOnly)

	result := Result{}
	result.HTML = content
	result.WordCount = wordCount
	result.Title = ce.ExtractTitle()
	result.ContentImages = ce.ImageURLs
	result.MarkupInfo = ce.Parser.MarkupInfo()

	// Find pagination
	if !opts.SkipPagination && opts.OriginalURL != nil {
		if opts.PaginationAlgo == "pagenum" {
			finder := pagination.NewPageNumberFinder(ce.WordCounter, nil)
			result.PaginationInfo = finder.FindPagination(doc, opts.OriginalURL)
		} else {
			finder := pagination.NewPrevNextFinder()
			result.PaginationInfo = finder.FindPagination(doc, opts.OriginalURL)
		}
	}

	return &result, nil
}
