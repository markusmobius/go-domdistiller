// ORIGINAL: java/extractors/embeds/VimeoExtractor.java

package embed

import (
	"fmt"
	nurl "net/url"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/logutil"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/webdoc"
	"golang.org/x/net/html"
)

// VimeoExtractor is used for extracting Vimeo videos and relevant information.
type VimeoExtractor struct {
	PageURL *nurl.URL
	logger  logutil.Logger
}

func NewVimeoExtractor(pageURL *nurl.URL, logger logutil.Logger) *VimeoExtractor {
	return &VimeoExtractor{
		PageURL: pageURL,
		logger:  logger,
	}
}

func (ve *VimeoExtractor) RelevantTagNames() []string {
	tagNames := []string{}
	for tagName := range relevantVimeoTags {
		tagNames = append(tagNames, tagName)
	}
	return tagNames
}

func (ve *VimeoExtractor) Extract(node *html.Node) webdoc.Element {
	if node == nil {
		return nil
	}

	nodeTagName := dom.TagName(node)
	if _, exist := relevantVimeoTags[nodeTagName]; !exist {
		return nil
	}

	src := dom.GetAttribute(node, "src")
	src = stringutil.CreateAbsoluteURL(src, ve.PageURL)
	if !domutil.HasRootDomain(src, "player.vimeo.com") {
		return nil
	}

	vimeoID, params := ve.getDataFromSrcURL(src)
	if vimeoID == "" {
		return nil
	}

	logMsg := fmt.Sprintf("Vimeo embed extracted (ID: %s)", vimeoID)
	ve.printLog(logMsg)

	return &webdoc.Embed{
		Element: node,
		Type:    "vimeo",
		ID:      vimeoID,
		Params:  params,
	}
}

func (ve *VimeoExtractor) getDataFromSrcURL(srcURL string) (string, map[string]string) {
	// Parse src url
	if strings.HasPrefix(srcURL, "//") {
		srcURL = "http:" + srcURL
	}

	parsedURL, err := nurl.ParseRequestURI(srcURL)
	if err != nil {
		return "", nil
	}

	// Get video ID which will be the last part of the path, account
	// for possible tail slash/empty path sections.
	var videoID string
	pathParts := strings.Split(parsedURL.Path, "/")
	for i := len(pathParts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(pathParts[i])
		if part != "" {
			if part != "video" {
				videoID = part
			}
			break
		}
	}

	// Get parameters from URL. In case of queries that specified several times,
	// only use the last value.
	params := make(map[string]string)
	for key, values := range parsedURL.Query() {
		if nValue := len(values); nValue > 0 {
			params[key] = values[nValue-1]
		}
	}

	return videoID, params
}

func (ve *VimeoExtractor) printLog(args ...interface{}) {
	if ve.logger != nil {
		ve.logger.PrintVisibilityInfo(args...)
	}
}
