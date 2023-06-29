package task

import (
	"fmt"
	"path/filepath"

	"github.com/sigmonsays/picman/core"
)

func NewGenerateFinalName(w *core.Workflow) *GenerateFinalName {
	ret := &GenerateFinalName{}
	ret.Workflow = w
	return ret
}

type GenerateFinalName struct {
	Workflow *core.Workflow
}

func (me *GenerateFinalName) Run(state *core.State) error {
	log.Tracef("start %s", state.OriginalFilename)
	if state.DestinationFilename != "" {
		log.Tracef("DestinationFilename already generated, skipping")
		return nil
	}
	log.Tracef("Generating final name path under destination-dir %s", me.Workflow.DestinationDir)

	ym := fmt.Sprintf("%04d/%02d", state.Date.Year, state.Date.Month)
	cs6 := state.Checksum.Sha256[:6]

	// preserve basename
	// origbasename := filepath.Base(state.OriginalFilename)
	// name := origbasename[:len(origbasename)-len(state.Ext)]
	d := state.Date
	name := fmt.Sprintf("%04d%02d%02d-%02d%02d%02d",
		d.Year, d.Month, d.Day, d.Hour, d.Minute, d.Second)

	basename := fmt.Sprintf("%s-%s", name, cs6) + state.Ext
	fullpath := filepath.Join(me.Workflow.DestinationDir, ym, basename)
	state.DestinationFilename = fullpath

	return nil
}
