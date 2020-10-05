// ORIGINAL: java/SchemaOrgParser.java

package schemaorg

import (
	"strings"

	"golang.org/x/net/html"
)

type ThingItem interface {
	addStringPropertyName(name string)
	addItemPropertyName(name string)
	getStringProperty(name string) string
	getItemProperty(name string) ThingItem
	getType() SchemaType
	isSupported() bool
	getElement() *html.Node

	// putStringValue stores `value` for property with `name`, unless the property
	// already has a non-empty value, in which case `value` will be ignored. This
	// means we only keep the first value.
	putStringValue(name string, value string)

	// putItemValue stores `value` for property with `name`, unless the property
	// already has a non-null value, in which case `value` will be ignored. This
	// means we only keep the first value.
	putItemValue(name string, value ThingItem)
}

type BaseThingItem struct {
	element          *html.Node
	schemaType       SchemaType
	stringProperties map[string]string
	itemProperties   map[string]ThingItem
}

func (ti *BaseThingItem) init(schemaType SchemaType, element *html.Node) {
	ti.element = element
	ti.schemaType = schemaType
	ti.stringProperties = make(map[string]string)
	ti.itemProperties = make(map[string]ThingItem)

	ti.addStringPropertyName(NameProp)
	ti.addStringPropertyName(URLProp)
	ti.addStringPropertyName(DescriptionProp)
	ti.addStringPropertyName(ImageProp)
}

func (ti *BaseThingItem) addStringPropertyName(name string) {
	ti.stringProperties[name] = ""
}

func (ti *BaseThingItem) addItemPropertyName(name string) {
	ti.itemProperties[name] = nil
}

func (ti *BaseThingItem) getStringProperty(name string) string {
	return ti.stringProperties[name]
}

func (ti *BaseThingItem) getItemProperty(name string) ThingItem {
	return ti.itemProperties[name]
}

func (ti *BaseThingItem) getType() SchemaType {
	return ti.schemaType
}

func (ti *BaseThingItem) isSupported() bool {
	return ti.schemaType != Unsupported
}

func (ti *BaseThingItem) putStringValue(name string, value string) {
	currentValue, exist := ti.stringProperties[name]
	if exist && currentValue == "" {
		ti.stringProperties[name] = strings.TrimSpace(value)
	}
}

func (ti *BaseThingItem) putItemValue(name string, value ThingItem) {
	currentValue, exist := ti.itemProperties[name]
	if exist && currentValue == nil {
		ti.itemProperties[name] = value
	}
}

func (ti *BaseThingItem) getElement() *html.Node {
	return ti.element
}
