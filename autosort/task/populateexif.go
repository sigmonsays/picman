package task

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/sigmonsays/picman/core"
)

func NewPopulateExif(w *core.Workflow) *PopulateExif {
	ret := &PopulateExif{}
	ret.Workflow = w
	return ret
}

type PopulateExif struct {
	Workflow *core.Workflow
}

func (me *PopulateExif) Run(state *core.State) error {
	log.Tracef("start %s", state.OriginalFilename)
	if state.ExifData.Values != nil {
		log.Tracef("Already have exif data, skipping")
		return nil
	}

	stdout := bytes.NewBuffer(nil)

	args := []string{
		"exiftool",
		"-j",
		state.OriginalFilename,
	}

	c := exec.Command(args[0], args[1:]...)
	c.Stdout = stdout
	err := c.Run()
	if err != nil {
		state.StopProcessing("exiftool error: %s", err)
		return err
	}

	wrap := []map[string]interface{}{}

	err = json.Unmarshal(stdout.Bytes(), &wrap)
	if err != nil {
		state.StopProcessing("unmarshal error: %s", err)
		return err
	}
	exifdata := wrap[0]
	log.Tracef("loaded %d values from exiftool", len(exifdata))
	state.ExifData.Values = exifdata

	return nil
}
