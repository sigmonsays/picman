package autosort

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/sigmonsays/picman/autosort/task"
	"github.com/sigmonsays/picman/core"
)

var StateSubDir = ".picman/state"

func RunWorkflow(workflow *core.Workflow) error {
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
	}

	// start a state file
	state := core.NewState()

	// determine the state path
	cs := sha256.New()
	fmt.Fprintf(cs, workflow.Fullpath)
	sha := cs.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	basename := filepath.Base(workflow.Fullpath)
	statebasename := basename + "-" + shaStr[:6] + ".json"
	statefile := filepath.Join(workflow.Root, StateSubDir, statebasename)
	log.Tracef("state file %s", statefile)

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
			if err == core.StopProcessing {
				return core.StopProcessing

			} else if err == core.SkipStep {
				continue
			}

			log.Warnf("step %s on %s failed: %s", step.Name, workflow.Fullpath, err)
			break
		}

	}
	return nil
}
