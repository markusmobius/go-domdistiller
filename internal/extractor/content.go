// ORIGINAL: java/ContentExtractor.java

package extractor

import (
	nurl "net/url"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/markusmobius/go-domdistiller/internal/converter"
	"github.com/markusmobius/go-domdistiller/internal/filter/docfilter"
	"github.com/markusmobius/go-domdistiller/internal/markup"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"github.com/markusmobius/go-domdistiller/logutil"
	"golang.org/x/net/html"
)

type ContentExtractor struct {
	Parser      *markup.Parser
	TimingInfo  *data.TimingInfo
	ImageURLs   []string
	WordCounter stringutil.WordCounter

	pageURL         *nurl.URL
	documentElement *html.Node
	candidateTitles []string
	logger          *logutil.Logger
}

func NewContentExtractor(root *html.Node, pageURL *nurl.URL, logger *logutil.Logger) *ContentExtractor {
	timingInfo := &data.TimingInfo{}

	document := dom.QuerySelector(root, "html")
	if document == nil {
		document = root
	}
	start := time.Now()
	parser := markup.NewParser(document, timingInfo)
	timingInfo.MarkupParsingTime = time.Now().Sub(start)

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

func (ce *ContentExtractor) ExtractContent(textOnly bool) (string, int) {
	start := time.Now()
	webDocument := ce.createWebDocumentInfoFromPage()
	ce.TimingInfo.DocumentConstructionTime = time.Now().Sub(start)

	start = time.Now()
	wordCount := ce.processDocument(webDocument)
	docfilter.NewRelevantElements().Process(webDocument)
	docfilter.NewLeadImageFinder(ce.logger).Process(webDocument)
	docfilter.NewNestedElementRetainer().Process(webDocument)
	ce.TimingInfo.ArticleProcessingTime = time.Now().Sub(start)

	start = time.Now()
	strHTML := webDocument.GenerateOutput(textOnly)
	ce.TimingInfo.FormattingTime = time.Now().Sub(start)

	ce.ImageURLs = webDocument.GetImageURLs()
	return strHTML, wordCount
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

// createWebDocumentInfoFromPage converts the original HTML page into a
// webdoc.Document for analysis.
func (ce *ContentExtractor) createWebDocumentInfoFromPage() *webdoc.Document {
	docBuilder := webdoc.NewWebDocumentBuilder(ce.WordCounter, ce.pageURL)
	converter.NewDomConverter(docBuilder, ce.pageURL, ce.logger).Convert(ce.documentElement)

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
