// ORIGINAL: java/MarkupParser.java

package markup

import (
	"github.com/markusmobius/go-domdistiller/internal/model"
	"golang.org/x/net/html"
)

// Parser loads the different parsers that are based on different markup specifications, and
// allows retrieval of different distillation-related markup properties from a document. It retrieves
// the requested properties from one or more parsers.  If necessary, it may merge the information
// from multiple parsers.
//
// Currently, three markup format are supported: OpenGraphProtocol, IEReadingView and SchemaOrg.
// For now, OpenGraphProtocolParser takes precedence because it uses specific meta tags and hence
// extracts information the fastest; it also demands conformance to rules. If the rules are broken
// or the properties retrieved are null or empty, we try with SchemaOrg then IEReadingView.
//
// The properties that matter to distilled content are:
// - individual properties: title, page type, page url, description, publisher, author, copyright
// - dominant and inline images and their properties: url, secure_url, type, caption, width, height
// - article and its properties: section name, published time, modified time, expiration time,
//   authors.
//
// TODO: for some properties, e.g. dominant and inline images, we might want to retrieve from
// multiple parsers; IEReadingViewParser provides more information as it scans all images in the
// document.  If we do so, we would need to merge the multiple versions in a meaningful way.
type Parser struct {
	accessors  []Accessor
	timingInfo *model.TimingInfo
}

func NewParser(root *html.Node, timingInfo *model.TimingInfo) *Parser {
	return &Parser{}
}
