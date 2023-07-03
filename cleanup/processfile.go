package cleanup

import (
	"fmt"
	"strings"
)

func (me *Cleanup) ProcessFile(srcdir string, statefile string, opts *Options, stats *Stats) error {

	result := RunCleanup(srcdir, statefile, opts, stats)
	result.Finish()

	if result.Row[0] == "OK" {
		return nil
	}

	fmt.Printf("%s\n", strings.Join(result.Row, "\t"))

	return nil
}
