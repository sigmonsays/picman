package task

import (
	"os"

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
		log.Tracef("file already copied %s", state.DestinationFilename)
		return nil
	}

	core.EnsureParentDirExists(state.DestinationFilename)

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
	log.Infof("Copied %s to %s", state.OriginalFilename, state.DestinationFilename)

	return nil
}
