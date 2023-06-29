package autosort

import (
	"path/filepath"

	"github.com/sigmonsays/picman/autosort/task"
	"github.com/sigmonsays/picman/core"
)

func RunWorkflow(workflow *core.Workflow) error {
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
	statefile := filepath.Join(workflow.Root, ".picman/state")

	for _, step := range steps {
		log.Tracef("run task %s for %s", step.Name, workflow.Fullpath)
		err := step.Task.Run(state)
		if err != nil {
			log.Warnf("step %s on %s failed: %s", step.Name, workflow.Fullpath, err)
			break
		}

		state.Save(statepath)
	}
	return nil
}
