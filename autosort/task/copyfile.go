package task

import (
	"os"
	"path/filepath"

	"github.com/sigmonsays/picman/core"
)

func NewCopyFile(w *core.Workflow) *CopyFile {
	ret := &CopyFile{}
	ret.Workflow = w
	return ret
}

type CopyFile struct {
	Workflow *core.Workflow
}

func (me *CopyFile) Run(state *core.State) error {
	if me.Workflow.NoCopy {
		log.Tracef("Not copying file due to no-copy flag")
		return nil
	}

	log.Tracef("start %s", state.OriginalFilename)

	if state.DestinationFilename == "" {
		log.Tracef("DestinationFilename required")
		return nil
	}

	// prevent accidental invocation where source and destination are the same file
	if state.OriginalFilename == state.DestinationFilename {
		log.Warnf("Cowardly refusing to copy with same source and destination path: %s", state.OriginalFilename)
		return nil
	}

	if state.FileCopied {
		log.Tracef("file already copied")
		return nil
	}

	// check if destination dir exists
	destdir := filepath.Dir(state.DestinationFilename)
	_, err := os.Stat(destdir)
	if os.IsNotExist(err) {
		os.MkdirAll(destdir, core.DirMask)
		state.Logf("Created destdir %s", destdir)
	}

	// check if the destination exists already
	destExists := false
	st, err := os.Stat(state.DestinationFilename)
	if err == nil && st.IsDir() == false {
		destExists = true
	}

	if destExists && st.Size() == 0 {
		log.Warnf("Zero byte destination %s", state.DestinationFilename)
	}

	written, err := core.CopyFile(state.OriginalFilename, state.DestinationFilename)
	if err != nil {
		return err
	}
	log.Tracef("file copy: %d bytes to destination %s",
		written, state.DestinationFilename)

	state.FileCopied = true

	return nil
}
