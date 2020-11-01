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
	"strconv"
	"strings"

	"github.com/markusmobius/go-domdistiller/data"
	"golang.org/x/net/html"
)

type ImageItem struct {
	BaseThingItem
}

func NewImageItem(element *html.Node) *ImageItem {
	item := &ImageItem{}
	item.init(Image, element)
	item.addStringPropertyName(ContentURLProp)
	item.addStringPropertyName(EncodingFormatProp)
	item.addStringPropertyName(CaptionProp)
	item.addStringPropertyName(RepresentativeProp)
	item.addStringPropertyName(WidthProp)
	item.addStringPropertyName(HeightProp)
	return item
}

func (ii *ImageItem) isRepresentativeOfPage() bool {
	propValue := ii.getStringProperty(RepresentativeProp)
	return strings.ToLower(propValue) == "true"
}

func (ii *ImageItem) getImage() *data.MarkupImage {
	width, _ := strconv.Atoi(ii.getStringProperty(WidthProp))
	height, _ := strconv.Atoi(ii.getStringProperty(HeightProp))
	imageURL := ii.getStringProperty(ContentURLProp)
	if imageURL == "" {
		imageURL = ii.getStringProperty(URLProp)
	}

	return &data.MarkupImage{
		URL:     imageURL,
		Type:    ii.getStringProperty(EncodingFormatProp),
		Caption: ii.getStringProperty(CaptionProp),
		Width:   width,
		Height:  height,
	}
}
