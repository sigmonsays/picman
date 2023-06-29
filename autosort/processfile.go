package autosort

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sigmonsays/picman/core"
)

func (me *Autosort) ProcessFile(root, fullpath string, info fs.FileInfo, dstdir string, source string, opts *Options) error {
	relpath, err := filepath.Rel(root, fullpath)
	if err != nil {
		return err
	}
	log.Tracef("process file %s relpath:%s", fullpath, relpath)

	if strings.HasPrefix(relpath, StateSubDir) {
		return nil
	}

	workflow := &core.Workflow{}
	workflow.Source = source
	workflow.Root = root
	workflow.Fullpath = fullpath
	workflow.Info = info
	workflow.RelPath = relpath

	state := core.NewState()

	err = RunWorkflow(workflow, state, opts)

	// if a test file is set, add extra info
	if opts.OneFile != "" {
		buf, _ := json.MarshalIndent(state, "", "  ")
		if err != nil {
			fmt.Printf("error:%s\n", err)
		}
		fmt.Printf("state file:\n%s\n", buf)
	}

	if err != nil {
		return err
	}

	return nil
}

// year := date.Format("2006")
// month := date.Format("01")
// basename := filepath.Base(fullpath)
// destpath := filepath.Join(dstdir, year, month, basename)

// dstExists := false
// dstSize := int64(0)
// dst, err := os.Stat(destpath)
// if err == nil {
// 	dstExists = true
// 	dstSize = dst.Size()
// } else {
// 	dstExists = false
// }
// sizeMismatch := srcSize != dstSize

// if !dstExists {
// 	log.Tracef(" destination exists: %v", dstExists)
// }

// log.Tracef("src:%s size:%d mismatch:%v destpath:%s", relpath, srcSize, sizeMismatch, destpath)
