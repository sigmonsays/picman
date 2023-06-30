package autosort

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sigmonsays/picman/core"
)

func (me *Autosort) ProcessFile(root, fullpath string, info fs.FileInfo, dstdir string, source string, opts *Options, stats *Stats) error {
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
	workflow.DestinationDir = dstdir
	workflow.NoCopy = opts.NoCopy

	state := core.NewState()

	err = RunWorkflow(workflow, state, opts, stats)

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
