// ORIGINAL: java/LogUtil.java

package logutil

import (
	"time"

	"github.com/markusmobius/go-domdistiller/internal/model"
)

func AddTimingInfo(timingInfo *model.TimingInfo, start time.Time, name string) {
	if timingInfo != nil {
		entry := model.TimingEntry{
			Name: name,
			Time: time.Now().Sub(start),
		}

		timingInfo.OtherTimes = append(timingInfo.OtherTimes, entry)
	}
}
