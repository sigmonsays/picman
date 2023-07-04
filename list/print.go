package list

import (
	"encoding/json"
	"fmt"

	"github.com/sigmonsays/picman/core"
)

func PrintState(srcdir string, statefile string, opts *Options) error {
	state := core.NewState()
	err := state.Load(statefile)
	if err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("state is nil")
	}

	buf, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", buf)
	return nil
}
