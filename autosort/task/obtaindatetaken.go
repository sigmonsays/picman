package task

import (
	"fmt"
	"strings"
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

	// order of keys we should try
	dateKeys := []string{
		"DateTimeOriginal",
		"CreateDate",
		"TrackCreateDate",
		"FileModifyDate",
	}
	dateKey := ""
	dateVal := ""
	var err error
	for _, dk := range dateKeys {
		if state.ExifData.KeyExists(dk) == false {
			continue
		}
		dateVal, err = state.ExifData.GetString(dk)
		if err != nil {
			continue
		}

		if strings.HasPrefix(dateVal, "0000") {
			continue
		}
		dateKey = dk
		break
	}

	if dateKey == "" {
		return state.StopProcessing("No date key found in exif")
	}

	log.Tracef("using dateKey %s from exif", dateKey)

	if dateKey == "FileModifyDate" {
		// 2023:06:29 18:41:02+00:00
		idx := strings.Index(dateVal, "+")
		if idx > 0 {
			dateVal = dateVal[:idx]
			log.Tracef("Dropping weird timezone, new date %s", dateVal)
		}
	}

	switch dateKey {

	case "DateTimeOriginal", "CreateDate", "TrackCreateDate", "FileModifyDate":

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
