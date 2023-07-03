package cleanup

import (
	"encoding/json"
	"fmt"
)

func (me *Cleanup) ProcessFile(srcdir string, statefile string, opts *Options, stats *Stats) error {

	result := RunCleanup(srcdir, statefile, opts, stats)

	if !result.HasError() {
		return nil
	}

	buf, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", buf)

	return nil
}
