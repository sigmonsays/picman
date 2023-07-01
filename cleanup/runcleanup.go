package cleanup

import (
	"errors"
	"os"

	"github.com/sigmonsays/picman/core"
)

var (
	DestinationEmpty     = errors.New("destination field is empty")
	NoDestination        = errors.New("No destination file exists")
	DestinationIsNotFile = errors.New("destination is not a regular file")
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

	if state.DestinationFilename == "" {
		return DestinationEmpty
	}

	// make sure the destination exists and is a regular file
	destInfo, err := os.Stat(state.DestinationFilename)
	if err != nil {
		return NoDestination
	}
	if destInfo.Mode().IsRegular() == false {
		return DestinationIsNotFile
	}

	// if DoNotProcess is set we just delete the metadata file

	return nil
}
