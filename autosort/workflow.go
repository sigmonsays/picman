package autosort

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sigmonsays/picman/autosort/task"
	"github.com/sigmonsays/picman/core"
)

var StateSubDir = ".picman/state"
var ErrorSubDir = ".picman/error"

// run the import workflow for a single file
func RunWorkflow(workflow *core.Workflow, state *core.State, opts *Options, stats *Stats) error {
	log.Tracef("")
	log.Tracef("start %s", workflow.Fullpath)
	stats.Processed++

	steps := []struct {
		Name string
		Task core.ImageProcessor
	}{
		{
			Name: "StartWorkflow",
			Task: task.NewStartWorkflow(workflow),
		},
		{
			Name: "CheckSupportedType",
			Task: task.NewCheckSupportedType(workflow),
		},
		{
			Name: "PopulateExif",
			Task: task.NewPopulateExif(workflow),
		},
		{
			Name: "CheckExif",
			Task: task.NewCheckExif(workflow),
		},
		{
			Name: "ObtainDateTaken",
			Task: task.NewObtainDateTaken(workflow),
		},
		{
			Name: "ChecksumFile",
			Task: task.NewChecksumFile(workflow),
		},
		{
			Name: "GenerateFinalName",
			Task: task.NewGenerateFinalName(workflow),
		},
		{
			Name: "CopyFile",
			Task: task.NewCopyFile(workflow),
		},
	}

	// determine the state path
	cs := sha256.New()
	fmt.Fprintf(cs, workflow.Fullpath)
	sha := cs.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	sha6 := shaStr[:6]
	sha2 := shaStr[:2]
	basename := filepath.Base(workflow.Fullpath)
	statebasename := basename + "-" + sha6 + ".json"
	statefile := filepath.Join(workflow.Root, StateSubDir, sha2, statebasename)
	errorfile := filepath.Join(workflow.Root, ErrorSubDir, sha2, statebasename)
	log.Tracef("state file %s", statefile)

	if opts.Force {
		log.Tracef("Force used, removing state file %s", statefile)
		os.Remove(statefile)
	}

	if st, err := os.Stat(statefile); err == nil && st.IsDir() == false {
		state.Load(statefile)
	}

	var taskErr error

	fileCopied := state.FileCopied

	for _, step := range steps {
		log.Tracef("run task %s for %s", step.Name, workflow.Fullpath)
		taskErr = step.Task.Run(state)

		err := state.Save(statefile)
		if err != nil {
			log.Warnf("save state file %s failed: %s", statefile, err)
			break
		}

		if taskErr != nil {
			// tasks have one way to stop processing
			if taskErr == core.StopProcessing {
				return core.StopProcessing

			} else if taskErr == core.SkipStep {
				continue
			}

			log.Warnf("step:%s %s failed: %s", step.Name, workflow.Fullpath, taskErr)
			break
		}

	}

	// check if the file has just been copied this run
	if fileCopied == false && state.FileCopied == true {
		stats.Copied += 1
	}

	if taskErr != nil {
		os.Symlink(statefile, errorfile)
	}

	return nil
}
