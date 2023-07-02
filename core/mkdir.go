package core

import (
	"os"
	"path/filepath"
)

func EnsureParentDirExists(path string) error {
	parentdir := filepath.Dir(path)
	_, err := os.Stat(parentdir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(parentdir, DirMask)
		if err != nil {
			return err
		}
	}
	return nil
}
