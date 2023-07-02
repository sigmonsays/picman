package task

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sigmonsays/picman/core"
)

func NewCheckSupportedType(w *core.Workflow) *CheckSupportedType {
	ret := &CheckSupportedType{}
	ret.Workflow = w
	return ret
}

type CheckSupportedType struct {
	Workflow *core.Workflow
}

func (me *CheckSupportedType) Run(state *core.State) error {

	// if we've already processed this then just abort
	if state.DoNotProcess {
		log.Tracef("DoNotProcess is set, StopProcessing")
		return core.StopProcessing
	}

	// populate state file with extension
	base := filepath.Base(state.OriginalFilename)
	base = strings.ToLower(base)
	// strip leading dot for hidden files
	base = strings.TrimLeft(base, ".")
	ext := filepath.Ext(base)
	state.Ext = ext

	err := ShouldProcess(state.OriginalFilename, state)
	process := (err == nil)

	if !process {
		state.DoNotProcess = true
		state.Logf("%s", err)
		return core.StopProcessing
	}
	return nil
}

func ShouldProcess(path string, state *core.State) error {
	stat := state.Stat
	if stat.Size == 0 {
		return fmt.Errorf("skip zero byte file %s", path)
	}
	log.Tracef("size is non-zero %s", path)

	if state.Ext == "" {
		return fmt.Errorf("No extension provided")
	}

	err := core.IsFileExtSupported(state.Ext)
	if err != nil {
		return err
	}

	log.Tracef("extension %s supported: %s", state.Ext, path)

	return nil
}
