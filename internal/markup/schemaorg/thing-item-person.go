// ORIGINAL: java/SchemaOrgParser.java

package schemaorg

import (
	"golang.org/x/net/html"
)

type PersonItem struct {
	BaseThingItem
}

func NewPersonItem(element *html.Node) *PersonItem {
	item := &PersonItem{}
	item.init(Person, element)
	item.addStringPropertyName(FamilyNameProp)
	item.addStringPropertyName(GivenNameProp)
	return item
}

func (pi *PersonItem) getName() string {
	// Returns either the value of NameProp, or concatenated values
	// of GivenNameProp and FamilyNameProp delimited by a whitespace.
	if name := pi.getStringProperty(NameProp); name != "" {
		return name
	}

	givenName := pi.getStringProperty(GivenNameProp)
	familyName := pi.getStringProperty(FamilyNameProp)
	if givenName != "" && familyName != "" {
		givenName += " "
	}

	return givenName + familyName
}
