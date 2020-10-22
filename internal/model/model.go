// ORIGINAL: Protobuf model in proto/dom_distiller.proto

package model

import "time"

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

type TimingEntry struct {
	Name string
	Time time.Duration
}

type TimingInfo struct {
	MarkupParsingTime        time.Duration
	DocumentConstructionTime time.Duration
	ArticleProcessingTime    time.Duration
	FormattingTime           time.Duration
	TotalTime                time.Duration

	// A place to hold arbitrary breakdowns of time. The perf scoring/server
	// should display these entries with appropriate names.
	OtherTimes []TimingEntry
}

type DebugInfo struct {
	Log string
}
