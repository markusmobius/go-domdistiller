// ORIGINAL: Protobuf model in proto/dom_distiller.proto

package data

import "time"

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

func (ti *TimingInfo) AddEntry(start time.Time, name string) {
	if ti == nil {
		return
	}

	ti.OtherTimes = append(ti.OtherTimes, TimingEntry{
		Name: name,
		Time: time.Now().Sub(start),
	})
}
