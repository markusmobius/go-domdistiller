// ORIGINAL: java/SchemaOrgParser.java

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
