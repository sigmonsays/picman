package task

import (
	"github.com/sigmonsays/picman/core"
)

func NewStartWorkflow(w *core.Workflow) *StartWorkflow {
	ret := &StartWorkflow{}
	ret.Workflow = w
	return ret
}

type StartWorkflow struct {
	Workflow *core.Workflow
}

func (me *StartWorkflow) Run(state *core.State) error {

	state.Source = me.Workflow.Source
	state.Stat.Size = int(me.Workflow.Info.Size())
	state.Stat.MTime = me.Workflow.Info.ModTime()
	state.OriginalFilename = me.Workflow.Fullpath

	log.Tracef("started workflow for %s size:%d mtime:%d",
		state.OriginalFilename, state.Stat.Size, state.Stat.MTime.Unix())

	return nil
}
