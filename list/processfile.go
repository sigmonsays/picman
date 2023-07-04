package list

import (
	"encoding/json"
	"fmt"
)

func (me *List) ProcessFile(srcdir string, statefile string, opts *Options) error {

	result := PrintState(srcdir, statefile, opts)

	buf, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", buf)

	return nil
}
