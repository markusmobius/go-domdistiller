// ORIGINAL: java/OpenGraphProtocolParserAccessor.java

package opengraph

import (
	"strings"

	"github.com/markusmobius/go-domdistiller/internal/model"
	"golang.org/x/net/html"
)

type Accessor struct {
	parser *Parser
}

func NewAccessor(root *html.Node, timingInfo *model.TimingInfo) (*Accessor, error) {
	parser, err := NewParser(root, timingInfo)
	if err != nil {
		return nil, err
	}

	return &Accessor{
		parser: parser,
	}, nil
}

// Title returns the required "title" of the document.
func (a *Accessor) Title() string {
	if a.parser == nil {
		return ""
	}

	return a.parser.propertyTable[TitleProp]
}

// Type returns the required "type" of the document if it's an
// article, empty string otherwise.
func (a *Accessor) Type() string {
	if a.parser == nil {
		return ""
	}

	objType := a.parser.propertyTable[TypeProp]
	if strings.ToLower(objType) == ArticleObjtype {
		return "Article"
	}

	return ""
}

// URL returns the required "url" of the document.
func (a *Accessor) URL() string {
	if a.parser == nil {
		return ""
	}

	return a.parser.propertyTable[URLProp]
}

// Images returns the structured properties of all "image"
// structures. Each "image" structure consists of image, image:url,
// image:secure_url, image:type, image:width, and image:height.
func (a *Accessor) Images() []model.MarkupImage {
	if a.parser == nil {
		return nil
	}

	return a.parser.Images()
}

// Description returns the optional "description" of the document.
func (a *Accessor) Description() string {
	if a.parser == nil {
		return ""
	}

	return a.parser.propertyTable[DescriptionProp]
}

// Publisher returns the optional "site_name" of the document.
func (a *Accessor) Publisher() string {
	if a.parser == nil {
		return ""
	}

	return a.parser.propertyTable[SiteNameProp]
}

// Copyright returns empty since OpenGraph not support it.
func (a *Accessor) Copyright() string {
	return ""
}

// Author returns the concatenated first_name and last_name
// (delimited by a whitespace) of the "profile" object when
// value of "og:type" is "profile".
func (a *Accessor) Author() string {
	if a.parser == nil {
		return ""
	}

	return a.parser.FullName()
}

// Article returns the properties of the "article" object when
// value of "og:type" is "article". The properties are published_time,
// modified_time and expiration_time, section, and a list of URLs
// to each author's profile.
func (a *Accessor) Article() *model.MarkupArticle {
	if a.parser == nil {
		return nil
	}

	article := model.MarkupArticle{
		PublishedTime:  a.parser.propertyTable[ArticlePublishedTimeProp],
		ModifiedTime:   a.parser.propertyTable[ArticleModifiedTimeProp],
		ExpirationTime: a.parser.propertyTable[ArticleExpirationTimeProp],
		Section:        a.parser.propertyTable[ArticleSectionProp],
		Authors:        a.parser.Authors(),
	}

	if article.Section == "" &&
		article.PublishedTime == "" &&
		article.ModifiedTime == "" &&
		article.ExpirationTime == "" &&
		len(article.Authors) == 0 {
		return nil
	}

	return &article
}

// OptOut returns false since OpenGraph not support it. While
// this is not directly supported, the page owner can simply
// omit the required tags and init() will return a null
// OpenGraphProtocolParser.
func (a *Accessor) OptOut() bool {
	return false
}
