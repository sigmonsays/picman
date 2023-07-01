package cleanup

import (
	"github.com/sigmonsays/picman/core"
)

func RunCleanup(statefile string, opts *Options, stats *Stats) error {
	log.Tracef("")
	log.Tracef("start %s", statefile)
	stats.Processed++

	state := core.NewState()
	err := state.Load(statefile)
	if err != nil {
		return err
	}

	return nil
}
