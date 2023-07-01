package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dstdir := filepath.Dir(dst)
	dstbase := filepath.Base(dst)
	dsttmp := filepath.Join(dstdir, ".tmp"+dstbase+".tmp")

	destination, err := os.Create(dsttmp)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return nBytes, err
	}

	err = os.Rename(dsttmp, dst)
	if err != nil {
		return 0, err
	}

	return nBytes, err
}
