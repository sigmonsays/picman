package autosort

import (
	"io/fs"
	"path/filepath"

	"github.com/sigmonsays/picman/core"
)

func (me *Autosort) ProcessFile(root, fullpath string, info fs.FileInfo, dstdir string) error {
	relpath, err := filepath.Rel(root, fullpath)
	if err != nil {
		return err
	}

	workflow := &core.Workflow{}
	workflow.Root = root
	workflow.Fullpath = fullpath
	workflow.Info = info
	workflow.RelPath = relpath

	err = RunWorkflow(workflow)
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
