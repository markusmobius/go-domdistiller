// ORIGINAL: java/extractors/embeds/YouTubeExtractor.java

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

// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

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

// YouTubeExtractor is used for extracting YouTube videos and relevant information.
type YouTubeExtractor struct {
	PageURL *nurl.URL
	logger  logutil.Logger
}

func NewYouTubeExtractor(pageURL *nurl.URL, logger logutil.Logger) *YouTubeExtractor {
	return &YouTubeExtractor{
		PageURL: pageURL,
		logger:  logger,
	}
}

func (ye *YouTubeExtractor) RelevantTagNames() []string {
	tagNames := []string{}
	for tagName := range relevantYouTubeTags {
		tagNames = append(tagNames, tagName)
	}
	return tagNames
}

func (ye *YouTubeExtractor) Extract(node *html.Node) webdoc.Element {
	if node == nil {
		return nil
	}

	nodeTagName := dom.TagName(node)
	if _, exist := relevantYouTubeTags[nodeTagName]; !exist {
		return nil
	}

	// Handle deprecated way to embed youtube.
	// Ref: https://www.w3.org/blog/2008/09/howto-insert-youtube-video/
	//      http://xahlee.info/js/html_embed_video.html
	src := dom.GetAttribute(node, "src")
	if nodeTagName == "object" {
		objType := dom.GetAttribute(node, "type")
		if objType == "application/x-shockwave-flash" {
			src = dom.GetAttribute(node, "data")
		} else {
			param := dom.QuerySelector(node, `param[name="movie"]`)
			if param != nil {
				src = dom.GetAttribute(param, "value")
			}
		}
	}

	// Wrong syntax like "http://www.youtube.com/v/<video-id>&param=value" has been
	// observed in the wild. Youtube seems to be resilient.
	if !strings.Contains(src, "?") {
		src = strings.Replace(src, "&", "?", 1)
	}

	src = stringutil.CreateAbsoluteURL(src, ye.PageURL)
	if !domutil.HasRootDomain(src, "youtube.com") && !domutil.HasRootDomain(src, "youtube-nocookie.com") {
		return nil
	}

	youtubeID, params := ye.getDataFromSrcURL(src)
	if youtubeID == "" {
		return nil
	}

	logMsg := fmt.Sprintf("YouTube embed extracted (ID: %s)", youtubeID)
	ye.printLog(logMsg)

	return &webdoc.Embed{
		Element: node,
		Type:    "youtube",
		ID:      youtubeID,
		Params:  params,
	}
}

func (ye *YouTubeExtractor) getDataFromSrcURL(srcURL string) (string, map[string]string) {
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
			if part != "embed" {
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

func (ye *YouTubeExtractor) printLog(args ...interface{}) {
	if ye.logger != nil && !ye.logger.InternallyNil() {
		ye.logger.PrintVisibilityInfo(args...)
	}
}
