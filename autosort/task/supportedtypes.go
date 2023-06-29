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
	// populate state file with extension
	base := filepath.Base(state.OriginalFilename)
	base = strings.ToLower(base)
	ext := filepath.Ext(base)
	state.Ext = ext

	err := ShouldProcess(state.OriginalFilename, state)
	process := (err == nil)

	if !process {
		return core.StopWorkflow
	}
	return nil
}

var SupportedImageExt = map[string]bool{
	".heic": true,
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".mov":  true,
	".mp4":  true,
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

	supported, found := SupportedImageExt[state.Ext]
	if found == false {
		return fmt.Errorf("ext %s is not supported", state.Ext)
	}
	if supported == false {
		return fmt.Errorf("Not supported")
	}

	log.Tracef("extension %s supported: %s", state.Ext, path)

	return nil
}
