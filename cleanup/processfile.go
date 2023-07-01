package cleanup

import (
	"encoding/json"
	"fmt"

	"github.com/sigmonsays/picman/core"
)

func (me *Cleanup) ProcessFile(statefile string, opts *Options, stats *Stats) error {

	state := core.NewState()

	err := RunCleanup(statefile, opts, stats)

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
