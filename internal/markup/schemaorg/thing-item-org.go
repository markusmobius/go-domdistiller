// ORIGINAL: java/SchemaOrgParser.java

package schemaorg

import (
	"golang.org/x/net/html"
)

type OrganizationItem struct {
	BaseThingItem
}

func NewOrganizationItem(element *html.Node) *OrganizationItem {
	item := &OrganizationItem{}
	item.init(Organization, element)
	item.addStringPropertyName(LegalNameProp)
	return item
}

func (oi *OrganizationItem) getName() string {
	// Returns either the value of NameProp, or LegalNameProp.
	if name := oi.getStringProperty(NameProp); name != "" {
		return name
	}

	return oi.getStringProperty(LegalNameProp)
}
