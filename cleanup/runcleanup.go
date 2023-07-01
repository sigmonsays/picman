package cleanup

import (
	"fmt"
	"os"

	"github.com/sigmonsays/picman/core"
)

func NewResult(symbol string, msg string) *Result {
	ret := &Result{
		Symbol:  symbol,
		Message: msg,
	}
	return ret
}

type Result struct {
	Symbol  string
	Message string
}

func (me *Result) Error() string {
	return fmt.Sprintf("%s: %s", me.Symbol, me.Message)
}

var (
	DestinationEmpty     = NewResult("DestinationEmpty", "destination field is empty")
	NoDestination        = NewResult("NoDestination", "No destination file exists")
	DestinationIsNotFile = NewResult("DestinationIsNotFile", "destination is not a regular file")
	DoNotProcess         = NewResult("DoNotProcess", "do not process")
)

func RunCleanup(srcdir string, statefile string, opts *Options, stats *Stats) error {
	log.Tracef("")
	log.Tracef("start %s", statefile)
	stats.Processed++

	state := core.NewState()
	err := state.Load(statefile)
	if err != nil {
		return err
	}

	if state.DestinationFilename == "" {
		return DestinationEmpty
	}

	var result *Result

	// make sure the destination exists and is a regular file
	destExists := false
	destInfo, err := os.Stat(state.DestinationFilename)
	if err != nil {
		result = NoDestination
		destExists = false
	}
	if destInfo.Mode().IsRegular() == false {
		result = DestinationIsNotFile
		destExists = false
	}

	// check if our type is supported

	// buf, _ := json.MarshalIndent(state, "", "  ")

	// todo: if DoNotProcess is set we just delete the metadata file

	// fmt.Printf("%s\n", buf)

	return nil
}
