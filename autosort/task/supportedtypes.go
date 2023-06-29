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
	err := ShouldProcess(state.OriginalFilename, state.Stat)
	process := (err == nil)

	if !process {
		return fmt.Errorf("do not process: %s", err)
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

func ShouldProcess(path string, stat *core.Stat) error {
	if stat.Size == 0 {
		return fmt.Errorf("skip zero byte file %s", path)
	}

	base := filepath.Base(path)
	base = strings.ToLower(base)
	ext := filepath.Ext(base)

	supported, found := SupportedImageExt[ext]
	if found == false {
		return fmt.Errorf("ext %s is not supported", ext)
	}
	if supported == false {
		return fmt.Errorf("Not supported")
	}
	return nil
}
