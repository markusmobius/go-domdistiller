// ORIGINAL: java/MarkupParser.java

package markup

import "github.com/markusmobius/go-domdistiller/internal/model"

// Accessor is the interface that all parsers must implement so that Parser
// can retrieve their properties.
type Accessor interface {
	// Title returns the markup title of the document, empty if none.
	Title() string

	// Type returns the markup type of the document, empty if none.
	Type() string

	// URL returns the markup url of the document, empty if none.
	URL() string

	// Images returns the properties of all markup images in the document.
	// The first image is the dominant (i.e. top or salient) one.
	Images() []model.MarkupImage

	// Description returns the markup description of the document, empty if none.
	Description() string

	// Publisher returns the markup publisher of the document, empty if none.
	Publisher() string

	// Copyright returns the markup copyright of the document, empty if none.
	Copyright() string

	// Author returns the full name of the markup author, empty if none.
	Author() string

	// Article returns the properties of the markup "article" object, null if none.
	Article() *model.MarkupArticle

	// OptOut returns true if page owner has opted out of distillation.
	OptOut() bool
}
