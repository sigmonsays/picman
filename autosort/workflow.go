package autosort

import (
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

	for _, step := range steps {
		err := step.Task.Run(state)
		if err != nil {
			log.Warnf("step %s on %s failed: %s", step.Name, workflow.Fullpath, err)
			break
		}
	}
	return nil
}
