// ORIGINAL: java/SchemaOrgParser.java

package schemaorg

import (
	"github.com/markusmobius/go-domdistiller/internal/model"
	"golang.org/x/net/html"
)

type ArticleItem struct {
	BaseThingItem
}

func NewArticleItem(element *html.Node) *ArticleItem {
	item := &ArticleItem{}
	item.init(Article, element)

	item.addStringPropertyName(HeadlineProp)
	item.addStringPropertyName(PublisherProp)
	item.addStringPropertyName(CopyrightHolderProp)
	item.addStringPropertyName(CopyrightYearProp)
	item.addStringPropertyName(DateModifiedProp)
	item.addStringPropertyName(DatePublishedProp)
	item.addStringPropertyName(AuthorProp)
	item.addStringPropertyName(CreatorProp)
	item.addStringPropertyName(SectionProp)

	item.addItemPropertyName(PublisherProp)
	item.addItemPropertyName(CopyrightHolderProp)
	item.addItemPropertyName(AuthorProp)
	item.addItemPropertyName(CreatorProp)
	item.addItemPropertyName(AssociatedMediaProp)
	item.addItemPropertyName(EncodingProp)

	return item
}

func (ai *ArticleItem) getArticle() *model.MarkupArticle {
	author := ai.getPersonOrOrganizationName(AuthorProp)
	if author == "" {
		author = ai.getPersonOrOrganizationName(CreatorProp)
	}

	var authors []string
	if author != "" {
		authors = []string{author}
	}

	return &model.MarkupArticle{
		PublishedTime: ai.getStringProperty(DatePublishedProp),
		ModifiedTime:  ai.getStringProperty(DateModifiedProp),
		Section:       ai.getStringProperty(SectionProp),
		Authors:       authors,
	}
}

func (ai *ArticleItem) getCopyright() string {
	copyright := ai.getStringProperty(CopyrightYearProp)
	copyrightHolder := ai.getPersonOrOrganizationName(CopyrightHolderProp)
	if copyright != "" && copyrightHolder != "" {
		copyright += " "
	}
	copyright += copyrightHolder

	if copyright != "" {
		return "Copyright " + copyright
	}
	return ""
}

func (ai *ArticleItem) getPersonOrOrganizationName(propertyName string) string {
	// Returns either the string value of `propertyName` or the value
	// returned by getName() of PersonItem or OrganizationItem.
	value := ai.getStringProperty(propertyName)
	if value != "" {
		return value
	}

	valueItem := ai.getItemProperty(propertyName)
	if valueItem != nil {
		switch valueItem.getType() {
		case Person:
			if personItem, ok := valueItem.(*PersonItem); ok {
				return personItem.getName()
			}

		case Organization:
			if orgItem, ok := valueItem.(*OrganizationItem); ok {
				return orgItem.getName()
			}
		}
	}

	return ""
}

func (ai *ArticleItem) getRepresentativeImageItem() *ImageItem {
	// Returns the corresponding ImageItem for "associatedMedia" or "encoding" property.
	item := ai.getItemProperty(AssociatedMediaProp)
	if item == nil {
		item = ai.getItemProperty(EncodingProp)
	}

	if item != nil && item.getType() == Image {
		if imageItem, ok := item.(*ImageItem); ok {
			return imageItem
		}
	}

	return nil
}

func (ai *ArticleItem) getImage() *model.MarkupImage {
	// Use value of "image" property to create a MarkupParser.Image.
	imageURL := ai.getStringProperty(ImageProp)
	if imageURL == "" {
		return nil
	}

	return &model.MarkupImage{URL: imageURL}
}
