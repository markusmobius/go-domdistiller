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
	"github.com/markusmobius/go-domdistiller/internal/extractor"
	"github.com/markusmobius/go-domdistiller/internal/model"
	"golang.org/x/net/html"
)

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

	// Which algorithm to use for next page detection:
	// "next" : detect anchors with "next" text
	// "pagenum" : detect anchors with numeric page numbers
	PaginationAlgo string
}

// ApplyForURL runs distiller for the specified URL.
func ApplyForURL(url string, timeout time.Duration, opts *Options) (*model.DistillerResult, error) {
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
	return Apply(resp.Body, opts)
}

// ApplyForFile runs distiller for the specified file.
func ApplyForFile(path string, opts *Options) (*model.DistillerResult, error) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	// Apply distiller to file
	return Apply(f, opts)
}

// Apply runs distiller for the specified io.Reader.
func Apply(r io.Reader, opts *Options) (*model.DistillerResult, error) {
	// Parse input
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

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
	result := model.DistillerResult{}
	ce := extractor.NewContentExtractor(doc, opts.OriginalURL)

	result.HTML = ce.ExtractContent(opts.ExtractTextOnly)
	result.Title = ce.ExtractTitle()
	result.ContentImages = ce.ImageURLs
	result.TextDirection = ce.TextDirection()
	result.MarkupInfo = ce.Parser.MarkupInfo()
	result.StatisticsInfo = ce.StatisticInfo

	return &result, nil
}
