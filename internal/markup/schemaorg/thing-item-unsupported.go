// ORIGINAL: java/SchemaOrgParser.java

package schemaorg

import (
	"golang.org/x/net/html"
)

type UnsupportedItem struct {
	BaseThingItem
}

func NewUnsupportedItem(element *html.Node) *UnsupportedItem {
	item := &UnsupportedItem{}
	item.init(Unsupported, element)
	return item
}
