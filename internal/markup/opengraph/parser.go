// ORIGINAL: java/OpenGraphProtocolParser.java

package opengraph

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/model"
	"golang.org/x/net/html"
)

var (
	rxOgpNsPrefix         = regexp.MustCompile(`(?i)((\w+):\s+(http:\/\/ogp.me\/ns(\/\w+)*#))\s*`)
	rxOgpNsNonPrefixName  = regexp.MustCompile(`(?i)^xmlns:(\w+)`)
	rxOgpNsNonPrefixValue = regexp.MustCompile(`(?i)^http:\/\/ogp.me\/ns(\/\w+)*#`)
)

type Parser struct {
	prefixes      PrefixNameList
	propertyTable map[string]string
	imageParser   ImagePropParser
	profileParser ProfilePropParser
	articleParser ArticlePropParser
}

func NewParser(root *html.Node, timingInfo *model.TimingInfo) (*Parser, error) {
	// Initiate parser
	ps := &Parser{}
	ps.prefixes = make(PrefixNameList)
	ps.propertyTable = make(map[string]string)

	start := time.Now()
	ps.findPrefixes(root)
	logutil.AddTimingInfo(timingInfo, start, "OpenGraphProtocolParser.findPrefixes")

	start = time.Now()
	ps.parseMetaTags(root)
	logutil.AddTimingInfo(timingInfo, start, "OpenGraphProtocolParser.parseMetaTags")

	start = time.Now()
	ps.imageParser.Verify()
	logutil.AddTimingInfo(timingInfo, start, "OpenGraphProtocolParser.imageParser.verify")

	prefix := ps.prefixes[OG] + ":"
	switch {
	case ps.propertyTable[TitleProp] == "":
		return nil, fmt.Errorf("required \"%s:title\" property is missing", prefix)

	case ps.propertyTable[TypeProp] == "":
		return nil, fmt.Errorf("required \"%s:type\" property is missing", prefix)

	case ps.propertyTable[URLProp] == "":
		return nil, fmt.Errorf("required \"%s:url\" property is missing", prefix)

	case len(ps.imageParser.ImageList) == 0:
		return nil, fmt.Errorf("required \"%s:image\" property is missing", prefix)
	}

	return ps, nil
}

func (ps *Parser) Images() []model.MarkupImage {
	return ps.imageParser.ImageList
}

func (ps *Parser) FullName() string {
	return ps.profileParser.GetFullName(ps.propertyTable)
}

func (ps *Parser) Authors() []string {
	return ps.articleParser.Authors
}

func (ps *Parser) findPrefixes(root *html.Node) {
	strPrefixes := ""

	// See if HTML tag has "prefix" attribute.
	htmlNode := dom.QuerySelector(root, "html")
	if htmlNode != nil {
		strPrefixes = dom.GetAttribute(htmlNode, "prefix")
	}

	// Otherwise, see if HEAD tag has "prefix" attribute.
	if strPrefixes == "" {
		head := dom.QuerySelector(root, "head")
		if head != nil {
			strPrefixes = dom.GetAttribute(head, "prefix")
		}
	}

	// If there's "prefix" attribute, its value is something like
	// "og: http://ogp.me/ns# profile: http://ogp.me/ns/profile# article: http://ogp.me/ns/article#".
	if strPrefixes != "" {
		matches := rxOgpNsPrefix.FindAllStringSubmatch(strPrefixes, -1)
		for _, groups := range matches {
			ps.prefixes.addObjectType(groups[2], groups[4])
		}
	} else {
		// Still no "prefix" attribute, see if HTMl tag has "xmlns" attributes e.g.:
		// - "xmlns:og="http://ogp.me/ns#"
		// - "xmlns:profile="http://ogp.me/ns/profile#"
		// - "xmlns:article="http://ogp.me/ns/article#".
		for _, attr := range htmlNode.Attr {
			attrName := strings.ToLower(attr.Key)
			nameMatch := rxOgpNsNonPrefixName.FindStringSubmatch(attrName)
			if nameMatch == nil {
				continue
			}

			valueMatch := rxOgpNsNonPrefixValue.FindStringSubmatch(attr.Val)
			if valueMatch != nil {
				ps.prefixes.addObjectType(nameMatch[1], valueMatch[1])
			}
		}
	}

	ps.prefixes.setDefault()
}

func (ps *Parser) parseMetaTags(root *html.Node) {
	// Fetch meta nodes
	var metaNodes []*html.Node
	if doPrefixFiltering {
		// Attribute selectors with prefix
		// https://developer.mozilla.org/en-US/docs/Web/CSS/Attribute_selectors
		query := ""
		for _, prefix := range ps.prefixes {
			query += `meta[property^=` + prefix + `],`
		}

		query = strings.TrimSuffix(query, ",")
		metaNodes = dom.QuerySelectorAll(root, query)
	} else {
		metaNodes = dom.QuerySelectorAll(root, "meta[property]")
	}

	// Parse property
	for _, meta := range metaNodes {
		content := dom.GetAttribute(meta, "content")
		property := dom.GetAttribute(meta, "property")
		property = strings.ToLower(property)

		// Only store properties that we care about for distillation.
		for _, importantProperty := range importantProperties {
			prefixWithColon := ps.prefixes[importantProperty.Prefix] + ":"

			// Note that `==` won't work here because importantProperties uses "image:"
			// (ImageStructPropPfx) for all image structured properties, so as to prevent
			// repetitive property name comparison - here and then again in ImageParser.
			if !strings.HasPrefix(property, prefixWithColon+importantProperty.Name) {
				continue
			}

			addProperty := true
			property = strings.TrimPrefix(property, prefixWithColon)
			switch importantProperty.Type {
			case "image":
				addProperty = ps.imageParser.Parse(property, content, ps.propertyTable)
			case "profile":
				addProperty = ps.profileParser.Parse(property, content, ps.propertyTable)
			case "article":
				addProperty = ps.articleParser.Parse(property, content, ps.propertyTable)
			}

			if addProperty {
				ps.propertyTable[importantProperty.Name] = content
			}

			break
		}
	}
}
