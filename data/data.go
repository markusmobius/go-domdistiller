// ORIGINAL: Protobuf model in proto/dom_distiller.proto

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

package data

type PaginationInfo struct {
	NextPage string
	PrevPage string
}

// MarkupArticle is object to contains the properties of an article document.
type MarkupArticle struct {
	PublishedTime  string
	ModifiedTime   string
	ExpirationTime string
	Section        string
	Authors        []string
}

// MarkupImage is used to contains the properties of an image in the document.
type MarkupImage struct {
	Root      string
	URL       string
	SecureURL string
	Type      string
	Caption   string
	Width     int
	Height    int
}

type MarkupInfo struct {
	Title       string
	Type        string
	URL         string
	Description string
	Publisher   string
	Copyright   string
	Author      string
	Article     MarkupArticle
	Images      []MarkupImage
}
