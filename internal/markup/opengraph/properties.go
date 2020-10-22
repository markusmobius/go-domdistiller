// ORIGINAL: java/OpenGraphProtocolParser.java

package opengraph

import (
	"strconv"
	"strings"

	"github.com/markusmobius/go-domdistiller/data"
)

type PropParser interface {
	Parse(property, content string, propertyTable map[string]string) bool
}

type ImagePropParser struct {
	ImageList []data.MarkupImage
}

func (pp *ImagePropParser) Parse(property, content string, propertyTable map[string]string) bool {
	// Root property means end of current structure.
	if property == ImageProp {
		image := data.MarkupImage{Root: content}
		pp.ImageList = append(pp.ImageList, image)
		return false
	}

	// Non ImageProp property means it's for current structure.
	currentIdx := len(pp.ImageList) - 1
	image := data.MarkupImage{}
	if currentIdx >= 0 {
		image = pp.ImageList[currentIdx]
	}

	imageChanged := true
	switch property {
	case ImageURLProp:
		image.URL = content
	case ImageSecureURLProp:
		image.SecureURL = content
	case ImageTypeProp:
		image.Type = content
	case ImageWidthProp:
		image.Width, _ = strconv.Atoi(content)
	case ImageHeightProp:
		image.Height, _ = strconv.Atoi(content)
	default:
		imageChanged = false
	}

	if imageChanged {
		if currentIdx < 0 {
			pp.ImageList = append(pp.ImageList, image)
		} else {
			pp.ImageList[currentIdx] = image
		}
	}

	return false
}

func (pp *ImagePropParser) Verify() {
	validImages := []data.MarkupImage{}

	for _, img := range pp.ImageList {
		if img.Root == "" {
			continue
		}

		if img.URL == "" {
			img.URL = img.Root
		}

		img.Root = ""
		validImages = append(validImages, img)
	}

	pp.ImageList = validImages
}

type ProfilePropParser struct {
	typeChecked   bool
	isProfileType bool
}

func (pp *ProfilePropParser) Parse(property, content string, propertyTable map[string]string) bool {
	// Check that "type" property exists and has "profile" value.
	if !pp.typeChecked {
		requiredType := propertyTable[TypeProp]
		pp.isProfileType = strings.ToLower(requiredType) == ProfileObjtype
		pp.typeChecked = true
	}

	// If it's profile object, insert into property table.
	return pp.isProfileType
}

func (pp *ProfilePropParser) GetFullName(propertyTable map[string]string) string {
	if !pp.isProfileType {
		return ""
	}

	fullName := propertyTable[ProfileFirstnameProp]
	lastName := propertyTable[ProfileLastnameProp]
	if fullName != "" && lastName != "" {
		fullName += " " + lastName
	}

	return fullName
}

type ArticlePropParser struct {
	isArticleType bool
	Authors       []string
}

func (pp *ArticlePropParser) Parse(property, content string, propertyTable map[string]string) bool {
	// Check that "type" property exists and has "article" value.
	if !pp.isArticleType {
		requiredType := propertyTable[TypeProp]
		pp.isArticleType = strings.ToLower(requiredType) == ArticleObjtype
	}

	if !pp.isArticleType {
		return false
	}

	// "author" property is an array of URLs, so we special-handle it here.
	if property == ArticleAuthorProp {
		pp.Authors = append(pp.Authors, content)
		return false
	}

	// Other properties are stateless, so inserting into property table is good enough.
	return true
}
