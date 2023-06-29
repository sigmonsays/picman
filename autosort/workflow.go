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

func RunWorkflow(workflow *core.Workflow, state *core.State, opts *Options) error {
	log.Tracef("")
	log.Tracef("start %s", workflow.Fullpath)

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
	}

	// determine the state path
	cs := sha256.New()
	fmt.Fprintf(cs, workflow.Fullpath)
	sha := cs.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	basename := filepath.Base(workflow.Fullpath)
	statebasename := basename + "-" + shaStr[:6] + ".json"
	statefile := filepath.Join(workflow.Root, StateSubDir, statebasename)
	log.Tracef("state file %s", statefile)

	if opts.Force {
		os.Remove(statefile)
	}

	if st, err := os.Stat(statefile); err == nil && st.IsDir() == false {
		state.Load(statefile)
	}

	for _, step := range steps {
		log.Tracef("run task %s for %s", step.Name, workflow.Fullpath)
		taskErr := step.Task.Run(state)

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

	return nil
}
