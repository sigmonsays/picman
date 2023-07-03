package cleanup

import (
	"os"

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
	MissingDate          = NewReason("MissingDate", "no date found")
	SourceMissing        = NewReason("SourceMissing", "source missing")
)

type Result struct {
	OriginalFile    string
	DestinationFile string

	LoadError            bool `json:"load_error,omitempty"`
	DestinationEmpty     bool `json:"destination_empty,omitempty"`
	NoDestination        bool `json:"no_destination,omitempty"`
	DestinationIsNotFile bool `json:"destination_is_not_file,omitempty"`
	DoNotProcess         bool `json:"do_not_process,omitempty"`
	FileTypeNotSupported bool `json:"file_type_not_supported,omitempty"`
	MissingDate          bool `json:"missing_date,omitempty"`
	SourceMissing        bool `json:"source_missing,omitempty"`
}

func (me *Result) HasError() bool {
	if me.LoadError ||
		me.DestinationEmpty ||
		me.DestinationEmpty ||
		me.NoDestination ||
		me.DestinationIsNotFile ||
		me.DoNotProcess ||
		me.FileTypeNotSupported ||
		me.MissingDate ||
		me.SourceMissing {
		return true
	}
	return false
}

func RunCleanup(srcdir string, statefile string, opts *Options, stats *Stats) *Result {
	log.Tracef("")
	log.Tracef("start %s", statefile)
	stats.Processed++

	ret := &Result{}

	state := core.NewState()
	err := state.Load(statefile)
	if err != nil {
		ret.LoadError = true
	}
	if state == nil {
		return nil
	}

	ret.OriginalFile = state.OriginalFilename
	ret.DestinationFile = state.DestinationFilename

	if state.DestinationFilename == "" {
		ret.DestinationEmpty = true
	}

	// make sure the source exists and is a regular file
	sourceExists := PathExistsRegularFile(state.OriginalFilename)
	if !sourceExists {
		ret.SourceMissing = true
	}

	// make sure the destination exists and is a regular file
	destExists := PathExistsRegularFile(state.OriginalFilename)
	if !destExists {
		ret.NoDestination = true
	}

	// check if our type is supported
	if core.IsFileExtSupported(state.Ext) != nil {
		ret.FileTypeNotSupported = true
	}

	// make sure we have a date
	if state.Date == nil || (state.Date.Year == 0 || state.Date.Month == 0) {
		ret.MissingDate = true
	}

	// todo: if DoNotProcess is set we just delete the metadata file
	// maybe we need a --delete flag

	return ret
}

// return true if path exists and is a regular file
func PathExistsRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.Mode().IsRegular() == false {
		return false
	}
	return true
}
