// ORIGINAL: java/SchemaOrgParser.java

// Copyright (c) 2020 Markus Mobius
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

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
