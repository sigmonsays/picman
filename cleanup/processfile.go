package cleanup

import (
	"fmt"
	"strings"
)

func (me *Cleanup) ProcessFile(srcdir string, statefile string, opts *Options, stats *Stats) error {

	result := RunCleanup(srcdir, statefile, opts, stats)
	result.Finish()

	fmt.Printf("ROW %s\n", strings.Join(result.Row, "\t"))

	return nil
}
