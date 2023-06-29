package task

import (
	"github.com/sigmonsays/picman/core"
)

func NewCheckExif(w *core.Workflow) *CheckExif {
	ret := &CheckExif{}
	ret.Workflow = w
	return ret
}

type CheckExif struct {
	Workflow *core.Workflow
}

func (me *CheckExif) Run(state *core.State) error {
	log.Tracef("start %s", state.OriginalFilename)
	if state.ExifData.Values == nil {
		return state.StopProcessing("No exif data")
	}

	// if FileType is TXT we dont wanna process it
	filetype, err := state.ExifData.GetString("FileType")
	if err == nil && filetype == "TXT" {
		return state.StopProcessing("Bad exif data FileType: %s", filetype)
	}

	return nil
}
