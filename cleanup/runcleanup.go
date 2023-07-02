package cleanup

import (
	"os"
	"path/filepath"

	"github.com/sigmonsays/picman/core"
)

func NewReason(symbol string, msg string) *Reason {
	ret := &Reason{
		Symbol:  symbol,
		Message: msg,
	}
	return ret
}

type Reason struct {
	Symbol  string
	Message string
}

var (
	LoadError            = NewReason("LoadError", "load state error")
	DestinationEmpty     = NewReason("DestinationEmpty", "destination field is empty")
	NoDestination        = NewReason("NoDestination", "No destination file exists")
	DestinationIsNotFile = NewReason("DestinationIsNotFile", "destination is not a regular file")
	DoNotProcess         = NewReason("DoNotProcess", "do not process")
	FileTypeNotSupported = NewReason("FileTypeNotSupported", "file type not supported")
)

type Result struct {
	Row []string
}

func (me *Result) Finish() {
	expected := 5
	if len(me.Row) == expected {
		return
	}
	for i := 0; i < expected-len(me.Row); i++ {
		me.Row = append(me.Row, "\t")
	}
}
func (me *Result) WithReason(r *Reason) *Result {
	me.Row = append(me.Row, r.Symbol)
	return me
}
func (me *Result) Print() error {
	return nil
}

func RunCleanup(srcdir string, statefile string, opts *Options, stats *Stats) *Result {
	log.Tracef("")
	log.Tracef("start %s", statefile)
	stats.Processed++

	ret := &Result{}

	state := core.NewState()
	err := state.Load(statefile)
	if err != nil {
		return ret.WithReason(LoadError)
	}

	if state.DestinationFilename == "" {
		ret.WithReason(DestinationEmpty)
	}

	var reason *Reason

	// make sure the destination exists and is a regular file
	destExists := false
	destInfo, err := os.Stat(state.DestinationFilename)
	if err != nil {
		reason = NoDestination
		destExists = false
	}

	if destExists && destInfo.Mode().IsRegular() == false {
		reason = DestinationIsNotFile
		destExists = false
	}

	if destExists {
	}

	// check if our type is supported
	if core.IsFileExtSupported(state.Ext) != nil {
		reason = FileTypeNotSupported
	}

	row := []string{}

	if reason == nil {
		row = append(row, "OK")
	} else {
		row = append(row, reason.Symbol)
	}
	row = append(row, state.Ext)

	// src path
	srcrel, _ := filepath.Rel(srcdir, state.OriginalFilename)
	row = append(row, srcrel)
	row = append(row, state.DestinationFilename)
	ret.Row = row

	// todo: if DoNotProcess is set we just delete the metadata file
	// maybe we need a --delete flag

	return ret
}
