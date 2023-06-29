package task

import (
	"fmt"
	"time"

	"github.com/sigmonsays/picman/core"
)

// // determine destination path name
// date, err := GetDate(fullpath)
// if err != nil {
// 	return err
// }

func NewObtainDateTaken(w *core.Workflow) *ObtainDateTaken {
	ret := &ObtainDateTaken{}
	ret.Workflow = w
	return ret
}

type ObtainDateTaken struct {
	Workflow *core.Workflow
}

func (me *ObtainDateTaken) Run(state *core.State) error {
	log.Tracef("start %s", state.OriginalFilename)
	if state.ExifData.Values == nil {
		return state.StopProcessing("No exif data")
	}

	dateKey := state.ExifData.FindFirst("DateTimeOriginal", "CreateDate", "TrackCreateDate")
	if dateKey == "" {
		return fmt.Errorf("No date key found in exif data")
	}
	log.Tracef("DateKey is %s", dateKey)

	dateVal, err := state.ExifData.GetString(dateKey)
	if err != nil {
		return fmt.Errorf("%s key is not a string", dateKey)
	}
	switch dateKey {

	case "DateTimeOriginal":

		// 2023:03:02 18:13:31
		tm, err := time.Parse("2006:01:02 15:04:05", dateVal)
		if err != nil {
			return fmt.Errorf("Parse time: %s: %s", dateVal, err)
		}

		state.Date.Year = int(tm.Year())
		state.Date.Month = int(tm.Month())
		state.Date.Day = int(tm.Day())
		state.Date.Hour = int(tm.Hour())
		state.Date.Minute = int(tm.Minute())
		state.Date.Second = int(tm.Second())

	default:
		log.Warnf("Unknown date key %s", dateKey)

	}

	return nil
}
