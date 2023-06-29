package autosort

import (
	"io/fs"
	"os"
	"path/filepath"
)

func (me *Autosort) ProcessFile(root, fullpath string, info fs.FileInfo, dstdir string) error {
	if ShouldProcess(fullpath, info) == false {
		log.Tracef("ProcessFile %s: skipped", fullpath)
		return nil
	}

	relpath, err := filepath.Rel(root, fullpath)
	if err != nil {
		return err
	}

	srcSize := info.Size()

	// determine destination path name
	date, err := GetDate(fullpath)
	if err != nil {
		return err
	}

	year := date.Format("2006")
	month := date.Format("01")
	basename := filepath.Base(fullpath)
	destpath := filepath.Join(dstdir, year, month, basename)

	dstExists := false
	dstSize := int64(0)
	dst, err := os.Stat(destpath)
	if err == nil {
		dstExists = true
		dstSize = dst.Size()
	} else {
		dstExists = false
	}
	sizeMismatch := srcSize != dstSize

	if !dstExists {
		log.Tracef(" destination exists: %v", dstExists)
	}

	log.Tracef("src:%s size:%d mismatch:%v destpath:%s", relpath, srcSize, sizeMismatch, destpath)

	// do we need to convert it?
	// heic ->  jpg

	// generate sha256

	// sha256 := ""
	// sha, err := Sha256File(fullpath)
	// if err != nil {
	// 	return err
	// }
	// sha256 = sha
	// log.Tracef("Sha256 %s: %s", fullpath, sha256)

	return nil
}
