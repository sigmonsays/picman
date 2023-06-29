package autosort

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sigmonsays/picman/core"
)

func (me *Autosort) ProcessFile(root, fullpath string, info fs.FileInfo, dstdir string, source string) error {
	relpath, err := filepath.Rel(root, fullpath)
	if err != nil {
		return err
	}

	if strings.HasPrefix(relpath, StateSubDir) {
		return nil
	}

	// log.Tracef("process file %s relpath:%s", fullpath, relpath)

	workflow := &core.Workflow{}
	workflow.Source = source
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
