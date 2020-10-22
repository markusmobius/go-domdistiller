// ORIGINAL: java/LogUtil.java

package logutil

import (
	"time"

	"github.com/markusmobius/go-domdistiller/data"
)

func AddTimingInfo(timingInfo *data.TimingInfo, start time.Time, name string) {
	if timingInfo != nil {
		entry := data.TimingEntry{
			Name: name,
			Time: time.Now().Sub(start),
		}

		timingInfo.OtherTimes = append(timingInfo.OtherTimes, entry)
	}
}
