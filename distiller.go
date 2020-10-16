package distiller

import (
	"errors"
	"io"
	nurl "net/url"

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
	OriginalURL string

	// Which algorithm to use for next page detection:
	// "next" : detect anchors with "next" text
	// "pagenum" : detect anchors with numeric page numbers
	PaginationAlgo string
}

// Apply runs distiller to specified HTML page.
func Apply(r io.Reader, opts Options) (*model.DistillerResult, error) {
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

	// Validate page URL
	pageURL, err := nurl.ParseRequestURI(opts.OriginalURL)
	if err != nil {
		return nil, err
	}

	// Start extractor
	result := model.DistillerResult{}
	ce := extractor.NewContentExtractor(doc, pageURL)

	result.HTML = ce.ExtractContent(opts.ExtractTextOnly)
	result.Title = ce.ExtractTitle()
	result.ContentImages = ce.ImageURLs
	result.TextDirection = ce.TextDirection()
	result.MarkupInfo = ce.Parser.MarkupInfo()
	result.StatisticsInfo = ce.StatisticInfo

	return &result, nil
}
